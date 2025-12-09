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

var raspProtectionPolicyNonUpdatableParams = []string{"os_type", "enterprise_project_id"}

// @API HSS POST /v5/{project_id}/rasp/policy
// @API HSS GET /v5/{project_id}/rasp/policies
// @API HSS GET /v5/{project_id}/rasp/policy/detail
// @API HSS PUT /v5/{project_id}/rasp/policy
// @API HSS DELETE /v5/{project_id}/rasp/policy
func ResourceRaspProtectionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRaspProtectionPolicyCreate,
		ReadContext:   resourceRaspProtectionPolicyRead,
		UpdateContext: resourceRaspProtectionPolicyUpdate,
		DeleteContext: resourceRaspProtectionPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRaspProtectionPolicyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(raspProtectionPolicyNonUpdatableParams),

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
			// The field `feature_list` does not exist in API response body.
			// Instead of `rule_list`.
			"feature_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chk_feature_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"protective_action": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"feature_configure": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The field `enterprise_project_id` does not exist in API response body.
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
			"rule_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chk_feature_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"chk_feature_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"chk_feature_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"feature_configure": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protective_action": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"optional_protective_action": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"editable": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildRaspProtectionPolicyCreateQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?policy_name=%v", d.Get("policy_name"))

	if v, ok := d.GetOk("os_type"); ok {
		queryParams = fmt.Sprintf("%s&os_type=%v", queryParams, v)
	}

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func buildRaspProtectionPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"feature_list": buildFeatureListBodyParams(d.Get("feature_list").([]interface{})),
	}

	return bodyParams
}

func buildFeatureListBodyParams(featureInfo []interface{}) []map[string]interface{} {
	if len(featureInfo) == 0 {
		return nil
	}

	checkRuleInfo := make([]map[string]interface{}, 0, len(featureInfo))
	for _, v := range featureInfo {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"chk_feature_id":    raw["chk_feature_id"],
			"protective_action": raw["protective_action"],
			"enabled":           raw["enabled"],
			"feature_configure": raw["feature_configure"],
		}

		checkRuleInfo = append(checkRuleInfo, params)
	}

	return checkRuleInfo
}

func resourceRaspProtectionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		epsId      = cfg.GetEnterpriseProjectID(d)
		policyName = d.Get("policy_name").(string)
		httpUrl    = "v5/{project_id}/rasp/policy"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildRaspProtectionPolicyCreateQueryParams(d, epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildRaspProtectionPolicyBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating protection policy: %s", err)
	}

	policy, err := getRaspProtectionPolicyId(client, region, policyName, epsId)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy_id", policy, "").(string)
	if policyId == "" {
		return diag.Errorf("error creating protection policy: unable to find policy ID")
	}

	d.SetId(policyId)

	return resourceRaspProtectionPolicyRead(ctx, d, meta)
}

func getRaspProtectionPolicyId(client *golangsdk.ServiceClient, region, policyName, epsId string) (interface{}, error) {
	httpUrl := "v5/{project_id}/rasp/policies"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?policy_name=%s", listPath, policyName)
	if epsId != "" {
		listPath = fmt.Sprintf("%s&enterprise_project_id=%s", listPath, epsId)
	}

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	resp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		return nil, fmt.Errorf("error retrieving protection policy list: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	policy := utils.PathSearch("data_list|[0]", respBody, nil)
	if policy != nil {
		return policy, nil
	}

	return nil, errors.New("error creating protection policy: unable to find policy in list API")
}

func buildRaspProtectionPolicyQueryParams(policyId, epsId string) string {
	queryParams := fmt.Sprintf("?policy_id=%s", policyId)

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func resourceRaspProtectionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	policy, err := GetRaspProtectionPolicy(client, d.Id(), epsId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving protection policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policy_name", utils.PathSearch("policy_name", policy, nil)),
		d.Set("os_type", utils.PathSearch("os_type", policy, nil)),
		d.Set("rule_list", flattenRaspProtectionPolicyRules(
			utils.PathSearch("rule_list", policy, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetRaspProtectionPolicy(client *golangsdk.ServiceClient, policyId, epsId string) (interface{}, error) {
	httpUrl := "v5/{project_id}/rasp/policy/detail"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildRaspProtectionPolicyQueryParams(policyId, epsId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	ruleList := utils.PathSearch("rule_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(ruleList) > 0 {
		return respBody, nil
	}

	return nil, golangsdk.ErrDefault404{}
}

func flattenRaspProtectionPolicyRules(ruleListInfo []interface{}) []interface{} {
	if len(ruleListInfo) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(ruleListInfo))
	for _, v := range ruleListInfo {
		result = append(result, map[string]interface{}{
			"chk_feature_id":             utils.PathSearch("chk_feature_id", v, nil),
			"chk_feature_name":           utils.PathSearch("chk_feature_name", v, nil),
			"chk_feature_desc":           utils.PathSearch("chk_feature_desc", v, nil),
			"feature_configure":          utils.PathSearch("feature_configure", v, nil),
			"protective_action":          utils.PathSearch("protective_action", v, nil),
			"optional_protective_action": utils.PathSearch("optional_protective_action", v, nil),
			"enabled":                    utils.PathSearch("enabled", v, nil),
			"editable":                   utils.PathSearch("editable", v, nil),
		})
	}

	return result
}

func buildRaspProtectionPolicyUpdateQueryParams(policyId, policyName, epsId string) string {
	queryParams := fmt.Sprintf("?policy_id=%s&policy_name=%s", policyId, policyName)

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func resourceRaspProtectionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		policyName = d.Get("policy_name").(string)
		epsId      = cfg.GetEnterpriseProjectID(d)
		httpUrl    = "v5/{project_id}/rasp/policy"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath += buildRaspProtectionPolicyUpdateQueryParams(d.Id(), policyName, epsId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildRaspProtectionPolicyBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating protection policy: %s", err)
	}

	return resourceRaspProtectionPolicyRead(ctx, d, meta)
}

func resourceRaspProtectionPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/rasp/policy"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath += buildRaspProtectionPolicyQueryParams(d.Id(), epsId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When the policy does not exist, the delete API will return `200`.
	// So we need to check if the policy exists before deleting it.
	_, err = GetRaspProtectionPolicy(client, d.Id(), epsId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving protection policy")
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting protection policy")
	}

	return nil
}

func resourceRaspProtectionPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<enterprise_project_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[0])
}
