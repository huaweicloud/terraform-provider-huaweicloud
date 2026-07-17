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

// @API DSC GET /v1/{project_id}/devices/security-policies
func DataSourceDscDeviceSecurityPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscDeviceSecurityPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the policy name for filtering.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the policy type for filtering.",
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The security policy list.",
				Elem:        dscDevicePolicyDetailSchema(),
			},
		},
	}
}

func dscDevicePolicyDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The policy ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy name.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the policy is enabled.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy status.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The policy type.",
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The update time.",
			},
			"related_datasource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The related datasource ID.",
			},
			"related_datasource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The related datasource name.",
			},
			"related_datasource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The related datasource type.",
			},
			"related_device_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The related device ID.",
			},
			"related_device_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The related device name.",
			},
			"related_device_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The related device type.",
			},
			"target_datasource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target datasource ID.",
			},
			"target_datasource_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target datasource name.",
			},
			"target_datasource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target datasource type.",
			},
			"ddm_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The dynamic masking configuration.",
				Elem:        dscDdmConfigSchema(),
			},
			"ddm_policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The dynamic masking policy list.",
				Elem:        dscDdmPolicyDetailSchema(),
			},
			"gde_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The encryption configuration.",
				Elem:        dscGdeConfigSchema(),
			},
			"gde_policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The encryption policy.",
				Elem:        dscGdePolicyDetailSchema(),
			},
			"sdm_config": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The static masking configuration.",
				Elem:        dscSdmConfigSchema(),
			},
			"sdm_policy_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The static masking policy list.",
				Elem:        dscSdmPolicyDetailSchema(),
			},
			"resource": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The device resource information.",
				Elem:        dscDeviceResourceInfoSchema(),
			},
			"target_resource": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The target device resource information.",
				Elem:        dscDeviceResourceInfoSchema(),
			},
		},
	}
}

func dscDdmConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"proxy_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The proxy port.",
			},
			"zk_election_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ZK election port.",
			},
			"zk_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The ZK port.",
			},
		},
	}
}

func dscDdmPolicyDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The namespace name.",
			},
			"table": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The table name.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information list.",
				Elem:        dscPolicyColumnSchema(),
			},
		},
	}
}

func dscGdeConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enc_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The encryption mode. 1 means encrypt, 2 means decrypt.",
			},
			"proxy_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The proxy port.",
			},
		},
	}
}

func dscGdePolicyDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The action. 1 means encrypt, 2 means decrypt.",
			},
			"alg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The encryption algorithm.",
			},
			"table": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The table name.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information list.",
				Elem:        dscPolicyColumnSchema(),
			},
		},
	}
}

func dscSdmConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auto_rebuild_target": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to rebuild the target table.",
			},
			"clear_target": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to clear the target table.",
			},
			"select_param": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The extraction parameter value.",
			},
			"select_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The extraction type.",
			},
			"skip_dirty_data": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to skip dirty data.",
			},
		},
	}
}

func dscSdmPolicyDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The table name.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The namespace name.",
			},
			"do_mask": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to mask data.",
			},
			"do_move": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to move data.",
			},
			"columns": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The column information list.",
				Elem:        dscPolicyColumnSchema(),
			},
		},
	}
}

func dscPolicyColumnSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The column name.",
			},
			"mask": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The masking algorithm name or ID.",
			},
		},
	}
}

func dscDeviceResourceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The account name.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address.",
			},
			"address_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The address type.",
			},
			"case_sensitive": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether case sensitive.",
			},
			"database_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The database name.",
			},
			"extra_params": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The extra parameters.",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The password.",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port.",
			},
			"res_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource ID.",
			},
			"res_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"res_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource version.",
			},
		},
	}
}

func buildDscDeviceSecurityPoliciesQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDscDeviceSecurityPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/devices/security-policies"
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
		currentPath := requestPath + buildDscDeviceSecurityPoliciesQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC device security policies: %s", err)
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
		d.Set("policies", flattenDscDeviceSecurityPolicies(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscDeviceSecurityPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("id", v, nil),
			"name":                    utils.PathSearch("name", v, nil),
			"enabled":                 utils.PathSearch("enabled", v, nil),
			"status":                  utils.PathSearch("status", v, nil),
			"type":                    utils.PathSearch("type", v, nil),
			"update_time":             utils.PathSearch("update_time", v, nil),
			"related_datasource_id":   utils.PathSearch("related_datasource_id", v, nil),
			"related_datasource_name": utils.PathSearch("related_datasource_name", v, nil),
			"related_datasource_type": utils.PathSearch("related_datasource_type", v, nil),
			"related_device_id":       utils.PathSearch("related_device_id", v, nil),
			"related_device_name":     utils.PathSearch("related_device_name", v, nil),
			"related_device_type":     utils.PathSearch("related_device_type", v, nil),
			"target_datasource_id":    utils.PathSearch("target_datasource_id", v, nil),
			"target_datasource_name":  utils.PathSearch("target_datasource_name", v, nil),
			"target_datasource_type":  utils.PathSearch("target_datasource_type", v, nil),
			"ddm_config":              flattenDscDdmConfig(v),
			"ddm_policy_list":         flattenDscDdmPolicyList(utils.PathSearch("ddm_policy_list", v, make([]interface{}, 0)).([]interface{})),
			"gde_config":              flattenDscGdeConfig(v),
			"gde_policy":              flattenDscGdePolicy(v),
			"sdm_config":              flattenDscSdmConfig(v),
			"sdm_policy_list":         flattenDscSdmPolicyList(utils.PathSearch("sdm_policy_list", v, make([]interface{}, 0)).([]interface{})),
			"resource":                flattenDscDeviceResourceInfo(v, "resource"),
			"target_resource":         flattenDscDeviceResourceInfo(v, "target_resource"),
		})
	}

	return rst
}

func flattenDscDdmConfig(respBody interface{}) []interface{} {
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

func flattenDscDdmPolicyList(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"namespace": utils.PathSearch("namespace", v, nil),
			"table":     utils.PathSearch("table", v, nil),
			"columns":   flattenDscPolicyColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscGdeConfig(respBody interface{}) []interface{} {
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

func flattenDscGdePolicy(respBody interface{}) []interface{} {
	gdePolicy := utils.PathSearch("gde_policy", respBody, nil)
	if gdePolicy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"action":  utils.PathSearch("action", gdePolicy, nil),
			"alg":     utils.PathSearch("alg", gdePolicy, nil),
			"table":   utils.PathSearch("table", gdePolicy, nil),
			"columns": flattenDscPolicyColumns(utils.PathSearch("columns", gdePolicy, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenDscSdmConfig(respBody interface{}) []interface{} {
	sdmConfig := utils.PathSearch("sdm_config", respBody, nil)
	if sdmConfig == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"auto_rebuild_target": utils.PathSearch("auto_rebuild_target", sdmConfig, nil),
			"clear_target":        utils.PathSearch("clear_target", sdmConfig, nil),
			"select_param":        utils.PathSearch("select_param", sdmConfig, nil),
			"select_type":         utils.PathSearch("select_type", sdmConfig, nil),
			"skip_dirty_data":     utils.PathSearch("skip_dirty_data", sdmConfig, nil),
		},
	}
}

func flattenDscSdmPolicyList(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"table":     utils.PathSearch("table", v, nil),
			"namespace": utils.PathSearch("namespace", v, nil),
			"do_mask":   utils.PathSearch("do_mask", v, nil),
			"do_move":   utils.PathSearch("do_move", v, nil),
			"columns":   flattenDscPolicyColumns(utils.PathSearch("columns", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenDscPolicyColumns(columns []interface{}) []interface{} {
	if len(columns) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(columns))
	for _, v := range columns {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"mask": utils.PathSearch("mask", v, nil),
		})
	}

	return rst
}

func flattenDscDeviceResourceInfo(respBody interface{}, path string) []interface{} {
	resourceInfo := utils.PathSearch(path, respBody, nil)
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
