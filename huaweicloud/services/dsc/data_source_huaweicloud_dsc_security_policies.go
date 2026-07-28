package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/security-policies
func DataSourceSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesPolicySchema(),
			},
		},
	}
}

func securityPoliciesPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dbss_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesDbssPolicySchema(),
			},
			"ddm_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesDdmConfigSchema(),
			},
			"ddm_policy_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesDdmPolicySchema(),
			},
			"dom_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesDomConfigSchema(),
			},
			"dom_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesDomPolicySchema(),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gde_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesGdeConfigSchema(),
			},
			"gde_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesGdePolicySchema(),
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_datasource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_datasource_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_datasource_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"related_instance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesResourceInfoSchema(),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func securityPoliciesDdmConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"proxy_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zk_election_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"zk_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func securityPoliciesGdeConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enc_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"proxy_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func securityPoliciesDbssPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_mask": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"show_result": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func securityPoliciesDomConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"deploy_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func securityPoliciesDomPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"custom_policy": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"data_audit": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"default_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"intelligent_protection_baseline": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"virtual_patch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func securityPoliciesDdmPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesColumnSchema(),
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func securityPoliciesGdePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityPoliciesColumnSchema(),
			},
			"table": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func securityPoliciesColumnSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"mask": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func securityPoliciesResourceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"case_sensitive": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extra_params": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"res_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"res_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"res_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildSecurityPoliciesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceSecurityPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/security-policies"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildSecurityPoliciesQueryParams(d, limit, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC security policies: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		policyList := utils.PathSearch("policy_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policyList...)

		if len(policyList) < limit {
			break
		}

		offset += len(policyList)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policy_list", flattenSecurityPolicies(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"dbss_policy": flattenSecurityPoliciesDbssPolicy(v),
			"ddm_config":  flattenSecurityPoliciesDdmConfig(v),
			"ddm_policy_list": flattenSecurityPoliciesDdmPolicyList(
				utils.PathSearch("ddm_policy_list", v, make([]interface{}, 0)).([]interface{})),
			"dom_config":              flattenSecurityPoliciesDomConfig(v),
			"dom_policy":              flattenSecurityPoliciesDomPolicy(v),
			"enabled":                 utils.PathSearch("enabled", v, nil),
			"gde_config":              flattenSecurityPoliciesGdeConfig(v),
			"gde_policy":              flattenSecurityPoliciesGdePolicy(v),
			"id":                      utils.PathSearch("id", v, nil),
			"name":                    utils.PathSearch("name", v, nil),
			"related_datasource_id":   utils.PathSearch("related_datasource_id", v, nil),
			"related_datasource_name": utils.PathSearch("related_datasource_name", v, nil),
			"related_datasource_type": utils.PathSearch("related_datasource_type", v, nil),
			"related_instance_id":     utils.PathSearch("related_instance_id", v, nil),
			"related_instance_name":   utils.PathSearch("related_instance_name", v, nil),
			"related_instance_type":   utils.PathSearch("related_instance_type", v, nil),
			"resource":                flattenSecurityPoliciesResourceInfo(v),
			"status":                  utils.PathSearch("status", v, nil),
			"type":                    utils.PathSearch("type", v, nil),
			"update_time":             utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}

func flattenSecurityPoliciesDdmConfig(respBody interface{}) []interface{} {
	ddmConfig := utils.PathSearch("ddm_config", respBody, nil)
	if ddmConfig == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"proxy_port":       utils.PathSearch("proxy_port", ddmConfig, nil),
			"zk_election_port": utils.PathSearch("zk_election_port", ddmConfig, nil),
			"zk_port":          utils.PathSearch("zk_port", ddmConfig, nil),
		},
	}
}

func flattenSecurityPoliciesGdeConfig(respBody interface{}) []interface{} {
	gdeConfig := utils.PathSearch("gde_config", respBody, nil)
	if gdeConfig == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"enc_mode":   utils.PathSearch("enc_mode", gdeConfig, nil),
			"proxy_port": utils.PathSearch("proxy_port", gdeConfig, nil),
		},
	}
}

func flattenSecurityPoliciesDbssPolicy(respBody interface{}) []interface{} {
	dbssPolicy := utils.PathSearch("dbss_policy", respBody, nil)
	if dbssPolicy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"data_mask":   utils.PathSearch("data_mask", dbssPolicy, nil),
			"show_result": utils.PathSearch("show_result", dbssPolicy, nil),
		},
	}
}

func flattenSecurityPoliciesDomConfig(respBody interface{}) []interface{} {
	domConfig := utils.PathSearch("dom_config", respBody, nil)
	if domConfig == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"deploy_mode": utils.PathSearch("deploy_mode", domConfig, nil),
		},
	}
}

func flattenSecurityPoliciesDomPolicy(respBody interface{}) []interface{} {
	domPolicy := utils.PathSearch("dom_policy", respBody, nil)
	if domPolicy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"custom_policy":  utils.PathSearch("custom_policy", domPolicy, nil),
			"data_audit":     utils.PathSearch("data_audit", domPolicy, nil),
			"default_action": utils.PathSearch("default_action", domPolicy, nil),
			"intelligent_protection_baseline": utils.PathSearch(
				"intelligent_protection_baseline", domPolicy, nil),
			"virtual_patch": utils.PathSearch("virtual_patch", domPolicy, nil),
		},
	}
}

func flattenSecurityPoliciesDdmPolicyList(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"columns": flattenSecurityPoliciesColumns(
				utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
			"namespace": utils.PathSearch("namespace", v, nil),
			"table":     utils.PathSearch("table", v, nil),
		})
	}

	return rst
}

func flattenSecurityPoliciesGdePolicy(respBody interface{}) []interface{} {
	gdePolicy := utils.PathSearch("gde_policy", respBody, nil)
	if gdePolicy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"action": utils.PathSearch("action", gdePolicy, nil),
			"alg":    utils.PathSearch("alg", gdePolicy, nil),
			"columns": flattenSecurityPoliciesColumns(
				utils.PathSearch("columns", gdePolicy, make([]interface{}, 0)).([]interface{})),
			"table": utils.PathSearch("table", gdePolicy, nil),
		},
	}
}

func flattenSecurityPoliciesColumns(columns []interface{}) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(columns))
	for _, v := range columns {
		rst = append(rst, map[string]interface{}{
			"mask": utils.PathSearch("mask", v, nil),
			"name": utils.PathSearch("name", v, nil),
		})
	}

	return rst
}

func flattenSecurityPoliciesResourceInfo(respBody interface{}) []interface{} {
	resourceInfo := utils.PathSearch("resource", respBody, nil)
	if resourceInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"account":        utils.PathSearch("account", resourceInfo, nil),
			"address":        utils.PathSearch("address", resourceInfo, nil),
			"address_type":   utils.PathSearch("address_type", resourceInfo, nil),
			"case_sensitive": utils.PathSearch("case_sensitive", resourceInfo, nil),
			"database_name":  utils.PathSearch("database_name", resourceInfo, nil),
			"extra_params":   utils.PathSearch("extra_params", resourceInfo, nil),
			"password":       utils.PathSearch("password", resourceInfo, nil),
			"port":           utils.PathSearch("port", resourceInfo, nil),
			"res_id":         utils.PathSearch("res_id", resourceInfo, nil),
			"res_type":       utils.PathSearch("res_type", resourceInfo, nil),
			"res_version":    utils.PathSearch("res_version", resourceInfo, nil),
		},
	}
}
