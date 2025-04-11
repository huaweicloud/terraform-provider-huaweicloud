package elb

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v3/{project_id}/elb/ipgroups
// @API ELB GET /v3/{project_id}/elb/ipgroups/{ipgroup_id}
// @API ELB PUT /v3/{project_id}/elb/ipgroups/{ipgroup_id}
// @API ELB DELETE /v3/{project_id}/elb/ipgroups/{ipgroup_id}
func ResourceIpGroupV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpGroupV3Create,
		ReadContext:   resourceIpGroupV3Read,
		UpdateContext: resourceIpGroupV3Update,
		DeleteContext: resourceIpGroupV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"listener_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIpGroupV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/ipgroups"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateIpGroupBodyParams(d, cfg))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB IP group: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB IP group: %s", err)
	}
	ipGroupId := utils.PathSearch("ipgroup.id", createRespBody, "").(string)
	if ipGroupId == "" {
		return diag.Errorf("error creating ELB IP group: ID is not found in API response")
	}

	d.SetId(ipGroupId)

	return resourceIpGroupV3Read(ctx, d, meta)
}

func buildCreateIpGroupBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"ip_list":               buildCreateIpGroupIpList(d.Get("ip_list").(*schema.Set).List()),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
	}
	bodyParams := map[string]interface{}{
		"ipgroup": params,
	}
	return bodyParams
}

func buildCreateIpGroupIpList(rawIpList []interface{}) []map[string]interface{} {
	if len(rawIpList) == 0 {
		return nil
	}

	ipList := make([]map[string]interface{}, 0, len(rawIpList))
	for _, rawIp := range rawIpList {
		if v, ok := rawIp.(map[string]interface{}); ok {
			ipList = append(ipList, map[string]interface{}{
				"ip":          v["ip"],
				"description": utils.ValueIgnoreEmpty(v["description"]),
			})
		}
	}
	return ipList
}

func resourceIpGroupV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/elb/ipgroups/{ipgroup_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{ipgroup_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB IP group")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("ipgroup.name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("ipgroup.description", getRespBody, nil)),
		d.Set("ip_list", flattenIpGroupIpList(getRespBody)),
		d.Set("listener_ids", flattenIpGroupListenerIds(getRespBody)),
		d.Set("enterprise_project_id", utils.PathSearch("ipgroup.enterprise_project_id", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("ipgroup.created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("ipgroup.updated_at", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIpGroupIpList(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("ipgroup.ip_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"ip":          utils.PathSearch("ip", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func flattenIpGroupListenerIds(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("ipgroup.listener_ids", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, utils.PathSearch("id", v, nil))
	}
	return rst
}

func resourceIpGroupV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/ipgroups/{ipgroup_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{ipgroup_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateIpGroupBodyParams(d))
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ELB IP group: %s", err)
	}

	return resourceIpGroupV3Read(ctx, d, meta)
}

func buildUpdateIpGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"ip_list":     buildCreateIpGroupIpList(d.Get("ip_list").(*schema.Set).List()),
	}
	bodyParams := map[string]interface{}{
		"ipgroup": params,
	}
	return bodyParams
}

func resourceIpGroupV3Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/ipgroups/{ipgroup_id}"
		product = "elb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{ipgroup_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting ELB IP group: %s", err)
	}

	return nil
}
