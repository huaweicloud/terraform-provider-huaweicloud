package coc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC POST /v1/{ticket_type}/list-tickets
func DataSourceCocIssueTickets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocIssueTicketsRead,

		Schema: map[string]*schema.Schema{
			"ticket_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"string_filters": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     DataSourceCocIssueTicketsFilters(),
			},
			"sort_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     DataSourceCocIssueTicketsFilters(),
			},
			"contain_sub_ticket": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"ticket_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tickets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issue_correlation_sla": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_cause_cloud_service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_cause_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"current_cloud_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issue_contact_person": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"issue_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"commit_upload_attachment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"virtual_schedule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"regions": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"root_cause_comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"solution": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"regions_search": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level_approve_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"suspension_approve_config": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"found_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_common_issue": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_need_change": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_enable_suspension": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_start_process_async": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_update_null": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_return_full_info": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_start_process": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignee": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"participator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"work_flow_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"baseline_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ticket_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_tickets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     ticketsSubTicketsElem(),
						},
						"enum_data_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     ticketsEnumDataListElem(),
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"meta_data_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ticket_type_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"form_info": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceCocIssueTicketsFilters() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func ticketsEnumDataListElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_zh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name_en": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"biz_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prop_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"model_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ticketsSubTicketsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"change_ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"change_ticket_id_sub": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"whether_to_change": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"main_ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ticket_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCocIssueTicketsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	instanceBatches, err := queryIssueTickets(client, d)
	if err != nil {
		return diag.Errorf("error querying issue tickets: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("tickets", flattenCocIssueTickets(
			utils.PathSearch("data.tickets", instanceBatches, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func queryIssueTickets(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/{ticket_type}/list-tickets"
	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{ticket_type}", d.Get("ticket_type").(string))
	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildQueryIssueTicketsBodyParams(d)),
	}

	queryResp, err := client.Request("POST", queryPath, &queryOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying COC issue tickets: %s", err)
	}

	queryRespBody, err := utils.FlattenResponse(queryResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening COC issue tickets: %s", err)
	}

	return queryRespBody, nil
}

func buildQueryIssueTicketsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"string_filters":     buildQueryIssueTicketsStringFiltersBodyParams(d.Get("string_filters")),
		"sort_filter":        buildQueryIssueTicketsSortFiltersBodyParams(d.Get("sort_filter")),
		"contain_all":        true,
		"contain_total":      true,
		"contain_sub_ticket": utils.ValueIgnoreEmpty(d.Get("contain_sub_ticket")),
		"ticket_types":       utils.ValueIgnoreEmpty(d.Get("ticket_types")),
	}

	return bodyParams
}

func buildQueryIssueTicketsStringFiltersBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"operator": raw["operator"],
				"field":    raw["field"],
				"values":   raw["values"],
			}
		}
		return params
	}

	return nil
}

func buildQueryIssueTicketsSortFiltersBodyParams(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"operator": raw["operator"],
			"field":    raw["field"],
			"values":   raw["values"],
		}
		return param
	}

	return nil
}

func flattenCocIssueTickets(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			m := map[string]interface{}{
				"issue_correlation_sla":     utils.PathSearch("issue_correlation_sla", params, nil),
				"level":                     utils.PathSearch("level", params, nil),
				"root_cause_cloud_service":  utils.PathSearch("root_cause_cloud_service", params, nil),
				"root_cause_type":           utils.PathSearch("root_cause_type", params, nil),
				"current_cloud_service_id":  utils.PathSearch("current_cloud_service_id", params, nil),
				"issue_contact_person":      utils.PathSearch("issue_contact_person", params, nil),
				"issue_version":             utils.PathSearch("issue_version", params, nil),
				"source":                    utils.PathSearch("source", params, nil),
				"commit_upload_attachment":  utils.PathSearch("commit_upload_attachment", params, nil),
				"source_id":                 utils.PathSearch("source_id", params, nil),
				"enterprise_project_id":     utils.PathSearch("enterprise_project", params, nil),
				"virtual_schedule_type":     utils.PathSearch("virtual_schedule_type", params, nil),
				"title":                     utils.PathSearch("title", params, nil),
				"regions":                   utils.PathSearch("regions", params, nil),
				"description":               utils.PathSearch("description", params, nil),
				"root_cause_comment":        utils.PathSearch("root_cause_comment", params, nil),
				"solution":                  utils.PathSearch("solution", params, nil),
				"regions_search":            utils.PathSearch("regions_serch", params, nil),
				"level_approve_config":      utils.PathSearch("level_approve_config", params, nil),
				"suspension_approve_config": utils.PathSearch("suspension_approve_config", params, nil),
				"handle_time":               utils.PathSearch("handle_time", params, nil),
				"found_time":                utils.PathSearch("fount_time", params, nil),
				"is_common_issue":           utils.PathSearch("is_common_issue", params, nil),
				"is_need_change":            utils.PathSearch("is_need_change", params, nil),
				"is_enable_suspension":      utils.PathSearch("is_enable_suspension", params, nil),
				"is_start_process_async":    utils.PathSearch("is_start_process_async", params, nil),
				"is_update_null":            utils.PathSearch("is_update_null", params, nil),
				"creator":                   utils.PathSearch("creator", params, nil),
				"operator":                  utils.PathSearch("operator", params, nil),
				"is_return_full_info":       utils.PathSearch("is_return_full_info", params, nil),
				"is_start_process":          utils.PathSearch("is_start_process", params, nil),
				"ticket_id":                 utils.PathSearch("ticket_id", params, nil),
				"real_ticket_id":            utils.PathSearch("real_ticket_id", params, nil),
				"assignee":                  utils.PathSearch("assignee", params, nil),
				"participator":              utils.PathSearch("participator", params, nil),
				"work_flow_status":          utils.PathSearch("work_flow_status", params, nil),
				"engine_error_msg":          utils.PathSearch("engine_error_msg", params, nil),
				"baseline_status":           utils.PathSearch("baseline_status", params, nil),
				"ticket_type":               utils.PathSearch("ticket_type", params, nil),
				"phase":                     utils.PathSearch("phase", params, nil),
				"sub_tickets":               flattenCocIssueTicketsSubTickets(utils.PathSearch("sub_tickets", params, nil)),
				"enum_data_list":            flattenCocIssueTicketsEnumDataList(utils.PathSearch("enum_data_list", params, nil)),
				"id":                        utils.PathSearch("id", params, nil),
				"meta_data_version":         utils.PathSearch("meta_data_version", params, nil),
				"update_time":               utils.PathSearch("update_time", params, nil),
				"create_time":               utils.PathSearch("create_time", params, nil),
				"is_deleted":                utils.PathSearch("is_deleted", params, nil),
				"ticket_type_id":            utils.PathSearch("ticket_type_id", params, nil),
				"form_info":                 utils.JsonToString(utils.PathSearch("_form_info", params, nil)),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func flattenCocIssueTicketsSubTickets(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"change_ticket_id":     utils.PathSearch("change_ticket_id", raw, nil),
				"change_ticket_id_sub": utils.PathSearch("change_ticket_id_sub", raw, nil),
				"whether_to_change":    utils.PathSearch("whether_to_change", raw, nil),
				"is_deleted":           utils.PathSearch("is_deleted", raw, nil),
				"id":                   utils.PathSearch("id", raw, nil),
				"main_ticket_id":       utils.PathSearch("main_ticket_id", raw, nil),
				"parent_ticket_id":     utils.PathSearch("parent_ticket_id", raw, nil),
				"ticket_id":            utils.PathSearch("ticket_id", raw, nil),
				"real_ticket_id":       utils.PathSearch("real_ticket_id", raw, nil),
				"ticket_path":          utils.PathSearch("ticket_path", raw, nil),
				"target_value":         utils.PathSearch("target_value", raw, nil),
				"target_type":          utils.PathSearch("target_type", raw, nil),
				"create_time":          utils.PathSearch("create_time", raw, nil),
				"update_time":          utils.PathSearch("update_time", raw, nil),
				"creator":              utils.PathSearch("creator", raw, nil),
				"operator":             utils.PathSearch("operator", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func flattenCocIssueTicketsEnumDataList(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"is_deleted":     utils.PathSearch("is_deleted", raw, nil),
				"match_type":     utils.PathSearch("match_type", raw, nil),
				"ticket_id":      utils.PathSearch("ticket_id", raw, nil),
				"real_ticket_id": utils.PathSearch("real_ticket_id", raw, nil),
				"name_zh":        utils.PathSearch("name_zh", raw, nil),
				"name_en":        utils.PathSearch("name_en", raw, nil),
				"biz_id":         utils.PathSearch("biz_id", raw, nil),
				"prop_id":        utils.PathSearch("prop_id", raw, nil),
				"model_id":       utils.PathSearch("model_id", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}
