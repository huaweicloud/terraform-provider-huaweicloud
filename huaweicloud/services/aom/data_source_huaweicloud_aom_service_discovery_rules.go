package aom

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM GET /v1/{project_id}/inv/servicediscoveryrules
func DataSourceServiceDiscoveryRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceDiscoveryRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"discovery_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"check_content": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"check_mode": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"check_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service_name_rule":     &dataSourceSchemaBasicNameRule,
									"application_name_rule": &dataSourceSchemaBasicNameRule,
								},
							},
						},
						"log_file_suffix": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"discovery_rule_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_default_rule": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"detect_log_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"log_path_rules": &dataSourceSchemaBasicNameRule,
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

var dataSourceSchemaBasicNameRule = schema.Schema{
	Type:     schema.TypeList,
	Computed: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"args": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

func dataSourceServiceDiscoveryRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := listServiceDiscoveryRules(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenServiceDiscoveryRules(results.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListServiceDiscoveryRulesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("rule_id"); ok {
		res = fmt.Sprintf("?&id=%v", v)
	}

	return res
}

func listServiceDiscoveryRules(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v1/{project_id}/inv/servicediscoveryrules"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListServiceDiscoveryRulesQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving service discovery rules: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening service discovery rules: %s", err)
	}

	return utils.PathSearch("appRules", listRespBody, make([]interface{}, 0)), nil
}

func flattenServiceDiscoveryRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		isDefaultRule, _ := strconv.ParseBool(utils.PathSearch("spec.isDefaultRule", rule, "false").(string))
		detectLogEnabled, _ := strconv.ParseBool(utils.PathSearch("spec.detectLog", rule, "false").(string))
		createdAt, _ := strconv.ParseInt(utils.PathSearch("createTime", rule, "").(string), 10, 64)

		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", rule, nil),
			"name":         utils.PathSearch("name", rule, nil),
			"service_type": utils.PathSearch("spec.appType", rule, nil),
			"discovery_rules": flattenDateSourceDiscoveryRules(
				utils.PathSearch("spec.discoveryRule", rule, make([]interface{}, 0)).([]interface{})),
			"name_rules":             flattenNameRules((utils.PathSearch("spec.nameRule", rule, nil))),
			"log_file_suffix":        utils.PathSearch("spec.logFileFix", rule, nil),
			"discovery_rule_enabled": utils.PathSearch("enable", rule, nil),
			"is_default_rule":        isDefaultRule,
			"detect_log_enabled":     detectLogEnabled,
			"priority":               utils.PathSearch("spec.priority", rule, nil),
			"log_path_rules": flattenBasicNameRules(
				utils.PathSearch("spec.logPathRule", rule, make([]interface{}, 0)).([]interface{})),
			"description": utils.PathSearch("desc", rule, nil),
			"created_at":  utils.FormatTimeStampRFC3339(createdAt/1000, false),
		})
	}
	return result
}

func flattenDateSourceDiscoveryRules(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"check_content": utils.PathSearch("checkContent", params, nil),
			"check_mode":    utils.PathSearch("checkMode", params, nil),
			"check_type":    utils.PathSearch("checkType", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenNameRules(paramsList interface{}) []interface{} {
	if paramsList == nil {
		return nil
	}
	rst := map[string]interface{}{
		"service_name_rule": flattenBasicNameRules(
			utils.PathSearch("appNameRule", paramsList, make([]interface{}, 0)).([]interface{})),
		"application_name_rule": flattenBasicNameRules(
			utils.PathSearch("applicationNameRule", paramsList, make([]interface{}, 0)).([]interface{})),
	}

	return []interface{}{rst}
}

func flattenBasicNameRules(paramsList []interface{}) []map[string]interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name_type": utils.PathSearch("nameType", params, nil),
			"args":      utils.PathSearch("args", params, nil),
			"value":     utils.PathSearch("value", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}
