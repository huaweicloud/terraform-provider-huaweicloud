package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/policy-records/{policy_id}/search
func DataSourcePolicyRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     policyRecordsConditionSchema(),
			},
			"sort": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     policyRecordsSortSchema(),
			},
			"group_by": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     policyRecordsGroupBySchema(),
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_task_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_task_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_automation": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"block_target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"trigger_flag": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dataclass_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workflow_instance": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"workflow_instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"workflow_instance_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"workflow_instance_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"block_age": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_block_ageing": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"block_ageing": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"policyrecord_type": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"environment": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"domain_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vendor_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"defense_policy_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"defense_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_policy_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_policy_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_connection_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_connection_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_connection_region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_connection_region_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_block_status": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_block_action": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_failure_description": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_creator_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"defense_creator_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_project_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_enterprise_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"target_enterprise_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"description": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func policyRecordsConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     policyRecordsConditionsSchema(),
			},
			"logics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func policyRecordsConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func policyRecordsSortSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sort_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func policyRecordsGroupBySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_by_fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group_by_hit": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     policyRecordsGroupByHitSchema(),
			},
		},
	}
}

func policyRecordsGroupByHitSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dest": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourcePolicyRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/policy-records/{policy_id}/search"
		result  = make([]interface{}, 0)
		limit   = 1000
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", d.Get("policy_id").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type":        "application/json;charset=UTF-8",
			"X-Secmaster-Version": "25.5.0",
		},
	}

	for {
		currentBodyParams := buildPolicyRecordsBodyParams(d, limit, offset)
		reqOpt.JSONBody = currentBodyParams

		resp, err := client.Request("POST", requestPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster policy records: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataRaw := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataRaw) == 0 {
			break
		}

		result = append(result, dataRaw...)

		total := utils.PathSearch("total", respBody, float64(0)).(float64)
		if int(total) <= offset+len(dataRaw) {
			break
		}

		offset += len(dataRaw)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("records", flattenPolicyRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPolicyRecordsBodyParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	return map[string]interface{}{
		"limit":     limit,
		"offset":    offset,
		"condition": buildPolicyRecordsConditionBodyParams(d.Get("condition").([]interface{})),
		"sort":      buildPolicyRecordsSortBodyParams(d.Get("sort").([]interface{})),
		"group_by":  buildPolicyRecordsGroupByBodyParams(d.Get("group_by").([]interface{})),
	}
}

func buildPolicyRecordsConditionBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"conditions": buildPolicyRecordsConditionsBodyParams(rawMap["conditions"].([]interface{})),
		"logics":     utils.ValueIgnoreEmpty(rawMap["logics"]),
	}
}

func buildPolicyRecordsConditionsBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(rawMap["name"]),
			"data": utils.ValueIgnoreEmpty(rawMap["data"]),
		})
	}

	return rst
}

func buildPolicyRecordsSortBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"sort_by": utils.ValueIgnoreEmpty(rawMap["sort_by"]),
			"order":   utils.ValueIgnoreEmpty(rawMap["order"]),
		})
	}

	return rst
}

func buildPolicyRecordsGroupByBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"group_by_fields": utils.ValueIgnoreEmpty(rawMap["group_by_fields"]),
		"group_by_hit":    buildPolicyRecordsGroupByHitBodyParams(rawMap["group_by_hit"].([]interface{})),
	}
}

func buildPolicyRecordsGroupByHitBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"source": utils.ValueIgnoreEmpty(rawMap["source"]),
		"dest":   utils.ValueIgnoreEmpty(rawMap["dest"]),
	}
}

func flattenPolicyRecords(records []interface{}) []map[string]interface{} {
	if len(records) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(records))
	for _, record := range records {
		dataObject := utils.PathSearch("data_object", record, nil)
		if dataObject == nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"id":                  utils.PathSearch("id", dataObject, nil),
			"policy_id":           utils.PathSearch("policy_id", dataObject, nil),
			"policy_task_id":      utils.PathSearch("policy_task_id", dataObject, nil),
			"policy_task_name":    utils.PathSearch("policy_task_name", dataObject, nil),
			"policy_category":     utils.PathSearch("policy_category", dataObject, nil),
			"policy_direction":    utils.PathSearch("policy_direction", dataObject, nil),
			"policy_automation":   utils.PathSearch("policy_automation", dataObject, nil),
			"block_target":        utils.PathSearch("block_target", dataObject, nil),
			"description":         utils.PathSearch("description", dataObject, nil),
			"is_deleted":          utils.PathSearch("is_deleted", dataObject, nil),
			"trigger_flag":        utils.PathSearch("trigger_flag", dataObject, nil),
			"create_time":         utils.PathSearch("create_time", dataObject, nil),
			"update_time":         utils.PathSearch("update_time", dataObject, nil),
			"dataclass_id":        utils.PathSearch("dataclass_id", dataObject, nil),
			"domain_id":           utils.PathSearch("domain_id", dataObject, nil),
			"domain_name":         utils.PathSearch("domain_name", dataObject, nil),
			"project_id":          utils.PathSearch("project_id", dataObject, nil),
			"region_id":           utils.PathSearch("region_id", dataObject, nil),
			"region_name":         utils.PathSearch("region_name", dataObject, nil),
			"workflow_instance":   flattenWorkflowInstance(dataObject),
			"block_age":           flattenBlockAge(dataObject),
			"policyrecord_type":   flattenPolicyrecordType(dataObject),
			"environment":         flattenEnvironment(dataObject),
			"defense_policy_list": flattenDefensePolicyList(dataObject),
		})
	}

	return result
}

func flattenWorkflowInstance(dataObject interface{}) []map[string]interface{} {
	instance := utils.PathSearch("workflow_instance", dataObject, nil)
	if instance == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"workflow_instance_id":     utils.PathSearch("workflow_instance_id", instance, nil),
			"workflow_instance_name":   utils.PathSearch("workflow_instance_name", instance, nil),
			"workflow_instance_status": utils.PathSearch("workflow_instance_status", instance, nil),
		},
	}
}

func flattenBlockAge(dataObject interface{}) []map[string]interface{} {
	blockAge := utils.PathSearch("block_age", dataObject, nil)
	if blockAge == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"is_block_ageing": utils.PathSearch("is_block_ageing", blockAge, nil),
			"block_ageing":    utils.PathSearch("block_ageing", blockAge, nil),
		},
	}
}

func flattenPolicyrecordType(dataObject interface{}) []map[string]interface{} {
	recordType := utils.PathSearch("policyrecord_type", dataObject, nil)
	if recordType == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"policy_type": utils.PathSearch("policy_type", recordType, nil),
			"id":          utils.PathSearch("id", recordType, nil),
			"category":    utils.PathSearch("category", recordType, nil),
		},
	}
}

func flattenEnvironment(dataObject interface{}) []map[string]interface{} {
	env := utils.PathSearch("environment", dataObject, nil)
	if env == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"domain_id":   utils.PathSearch("domain_id", env, nil),
			"domain_name": utils.PathSearch("domain_name", env, nil),
			"project_id":  utils.PathSearch("project_id", env, nil),
			"region_id":   utils.PathSearch("region_id", env, nil),
			"region_name": utils.PathSearch("region_name", env, nil),
			"vendor_type": utils.PathSearch("vendor_type", env, nil),
		},
	}
}

func flattenDefensePolicyList(dataObject interface{}) []map[string]interface{} {
	listRaw := utils.PathSearch("defense_policy_list", dataObject, make([]interface{}, 0))
	list, ok := listRaw.([]interface{})
	if !ok || len(list) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(list))
	for _, v := range list {
		result = append(result, map[string]interface{}{
			"defense_id":                     utils.PathSearch("defense_id", v, nil),
			"defense_type":                   utils.PathSearch("defense_type", v, nil),
			"defense_policy_id":              utils.PathSearch("defense_policy_id", v, nil),
			"defense_policy_name":            utils.PathSearch("defense_policy_name", v, nil),
			"defense_connection_id":          utils.PathSearch("defense_connection_id", v, nil),
			"defense_connection_name":        utils.PathSearch("defense_connection_name", v, nil),
			"defense_connection_region_id":   utils.PathSearch("defense_connection_region_id", v, nil),
			"defense_connection_region_name": utils.PathSearch("defense_connection_region_name", v, nil),
			"defense_block_status":           utils.PathSearch("defense_block_status", v, nil),
			"defense_block_action":           utils.PathSearch("defense_block_action", v, nil),
			"defense_failure_description":    utils.PathSearch("defense_failure_description", v, nil),
			"defense_creator_id":             utils.PathSearch("defense_creator_id", v, nil),
			"defense_creator_name":           utils.PathSearch("defense_creator_name", v, nil),
			"target_project_id":              utils.PathSearch("target_project_id", v, nil),
			"target_project_name":            utils.PathSearch("target_project_name", v, nil),
			"target_enterprise_id":           utils.PathSearch("target_enterprise_id", v, nil),
			"target_enterprise_name":         utils.PathSearch("target_enterprise_name", v, nil),
			"description":                    utils.PathSearch("description", v, nil),
		})
	}

	return result
}
