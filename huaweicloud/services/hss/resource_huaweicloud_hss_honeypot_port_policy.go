package hss

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var honeypotPortPolicyNonUpdatableParams = []string{"enterprise_project_id"}

// @API HSS POST /v5/{project_id}/honeypot-port/policy
// @API HSS GET /v5/{project_id}/honeypot-port/policy-list
// @API HSS GET /v5/{project_id}/honeypot-port/policy/{policy_id}
// @API HSS PUT /v5/{project_id}/honeypot-port/policy/{policy_id}
// @API HSS DELETE /v5/{project_id}/honeypot-port/policy/{policy_id}
func ResourceHoneypotPortPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHoneypotPortPolicyCreate,
		ReadContext:   resourceHoneypotPortPolicyRead,
		UpdateContext: resourceHoneypotPortPolicyUpdate,
		DeleteContext: resourceHoneypotPortPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(honeypotPortPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ports_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"white_ip": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"host_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"host_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"port_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildHoneypotPortPolicyQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func buildHoneypotPortPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_name": d.Get("policy_name"),
		"os_type":     d.Get("os_type"),
		"ports_list":  buildPortsListBodyParams(d.Get("ports_list").(*schema.Set).List()),
		"white_ip":    utils.ExpandToStringList(d.Get("white_ip").([]interface{})),
		"host_id":     utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("host_id").([]interface{}))),
		"group_list":  utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("group_list").([]interface{}))),
	}

	return bodyParams
}

func buildPortsListBodyParams(portsList []interface{}) []map[string]interface{} {
	if len(portsList) == 0 {
		return nil
	}

	ports := make([]map[string]interface{}, 0, len(portsList))
	for _, v := range portsList {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"port":     raw["port"],
			"protocol": raw["protocol"],
		}
		ports = append(ports, params)
	}

	return ports
}

func resourceHoneypotPortPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		epsId      = cfg.GetEnterpriseProjectID(d)
		policyName = d.Get("policy_name").(string)
		httpUrl    = "v5/{project_id}/honeypot-port/policy"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildHoneypotPortPolicyQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildHoneypotPortPolicyBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating dynamic port honeypot policy: %s", err)
	}

	policy, err := getHoneypotPortPolicyId(client, policyName, epsId)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy_id", policy, "").(string)
	if policyId == "" {
		return diag.Errorf("error creating dynamic port honeypot policy: unable to find policy ID")
	}

	d.SetId(policyId)

	return resourceHoneypotPortPolicyRead(ctx, d, meta)
}

func getHoneypotPortPolicyId(client *golangsdk.ServiceClient, policyName, epsId string) (interface{}, error) {
	var (
		httpUrl = "v5/{project_id}/honeypot-port/policy-list?limit=200"
		offset  = 0
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		listPath = fmt.Sprintf("%s&enterprise_project_id=%s", listPath, epsId)
	}

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving dynamic port honeypot policy list: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		policies := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(policies) < 1 {
			break
		}

		policy := utils.PathSearch(fmt.Sprintf("[?policy_name=='%s']|[0]", policyName), policies, nil)
		if policy != nil {
			return policy, nil
		}

		offset += len(policies)
	}

	return nil, errors.New("error creating dynamic port honeypot policy: unable to find policy in list API")
}

func GetHoneypotPortPolicy(client *golangsdk.ServiceClient, policyId, epsId string) (interface{}, error) {
	httpUrl := "v5/{project_id}/honeypot-port/policy/{policy_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", policyId)
	getPath += buildHoneypotPortPolicyQueryParams(epsId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceHoneypotPortPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	policy, err := GetHoneypotPortPolicy(client, d.Id(), epsId)
	if err != nil {
		// When the policy does not exist, the response body example of the details interface is as follows:
		// error message: { "error_code": "HSS.1005","error_description": "无效的策略信息","error_msg": "无效的策略信息"}
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "HSS.1005"),
			"error retrieving dynamic port honeypot policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policy_name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("os_type", utils.PathSearch("os_type", policy, nil)),
		d.Set("port_list", flattenHoneypotPortPolicyPortList(utils.PathSearch("port_list", policy, make([]interface{}, 0)).([]interface{}))),
		d.Set("white_ip", utils.PathSearch("white_ip", policy, nil)),
		d.Set("host_list", utils.PathSearch("host_list", policy, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHoneypotPortPolicyPortList(rawPortList []interface{}) []interface{} {
	result := make([]interface{}, 0, len(rawPortList))
	for _, v := range rawPortList {
		result = append(result, map[string]interface{}{
			"port":     utils.PathSearch("port", v, nil),
			"protocol": utils.PathSearch("protocol", v, nil),
		})
	}

	return result
}

func resourceHoneypotPortPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/honeypot-port/policy/{policy_id}"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Id())
	updatePath += buildHoneypotPortPolicyQueryParams(epsId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildHoneypotPortPolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating dynamic port honeypot policy: %s", err)
	}

	return resourceHoneypotPortPolicyRead(ctx, d, meta)
}

func resourceHoneypotPortPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/honeypot-port/policy/{policy_id}"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", d.Id())
	deletePath += buildHoneypotPortPolicyQueryParams(epsId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the policy does not exist, the response HTTP status code of the deletion API is 400.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "HSS.1005"),
			fmt.Sprintf("error deleting dynamic port honeypot policy, the error message: %s", err))
	}

	return nil
}
