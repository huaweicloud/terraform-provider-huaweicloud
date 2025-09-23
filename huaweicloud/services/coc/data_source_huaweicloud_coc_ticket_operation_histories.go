package coc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC POST /v1/{ticket_type}/list-histories
func DataSourceCocTicketOperationHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocTicketOperationHistoriesRead,

		Schema: map[string]*schema.Schema{
			"ticket_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"string_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     ticketOperationHistoriesFiltersSchema(),
			},
			"sort_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     ticketOperationHistoriesFiltersSchema(),
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"stop_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"action_name_zh": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_name_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_template_zh": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_template_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"final_sub_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enum_data_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
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
									"user_name": {
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
							},
						},
					},
				},
			},
		},
	}
}

func ticketOperationHistoriesFiltersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"operator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"values": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	return &sc
}

func buildTicketOperationHistoriesCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"string_filters": buildTicketOperationHistoriesStringFilterCreateOpts(d.Get("string_filters")),
		"sort_filter":    buildTicketOperationHistoriesSortFilterCreateOpts(d.Get("sort_filter")),
	}

	return bodyParams
}

func buildTicketOperationHistoriesStringFilterCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				params[i] = map[string]interface{}{
					"operator": utils.ValueIgnoreEmpty(raw["operator"]),
					"field":    utils.ValueIgnoreEmpty(raw["field"]),
					"name":     utils.ValueIgnoreEmpty(raw["name"]),
					"values":   utils.ValueIgnoreEmpty(raw["values"]),
				}
			}
		}
		return params
	}

	return nil
}

func buildTicketOperationHistoriesSortFilterCreateOpts(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		if raw, ok := rawArray[0].(map[string]interface{}); ok {
			param := map[string]interface{}{
				"operator": utils.ValueIgnoreEmpty(raw["operator"]),
				"field":    utils.ValueIgnoreEmpty(raw["field"]),
				"name":     utils.ValueIgnoreEmpty(raw["name"]),
				"values":   utils.ValueIgnoreEmpty(raw["values"]),
			}
			return param
		}
	}

	return nil
}

func dataSourceCocTicketOperationHistoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	ticketType := d.Get("ticket_type").(string)
	httpUrl := "v1/{ticket_type}/list-histories"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{ticket_type}", ticketType)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildTicketOperationHistoriesCreateOpts(d)),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving COC ticket operation histories: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("data", flattenCocTicketOperationHistoriesData(
			utils.PathSearch("data.info", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCocTicketOperationHistoriesData(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			m := map[string]interface{}{
				"action_id":          utils.PathSearch("action_id", params, nil),
				"action":             utils.PathSearch("action", params, nil),
				"sub_action":         utils.PathSearch("sub_action", params, nil),
				"operator":           utils.PathSearch("operator", params, nil),
				"comment":            utils.PathSearch("comment", params, nil),
				"id":                 utils.PathSearch("id", params, nil),
				"ticket_id":          utils.PathSearch("ticket_id", params, nil),
				"start_time":         utils.PathSearch("start_time", params, nil),
				"stop_time":          utils.PathSearch("stop_time", params, nil),
				"target_type":        utils.PathSearch("target_type", params, nil),
				"target_value":       utils.PathSearch("target_value", params, nil),
				"is_deleted":         utils.PathSearch("is_deteted", params, nil),
				"update_time":        utils.PathSearch("update_time", params, nil),
				"action_name_zh":     utils.PathSearch("action_name_zh", params, nil),
				"action_name_en":     utils.PathSearch("action_name_en", params, nil),
				"action_template_zh": utils.PathSearch("action_template_zh", params, nil),
				"action_template_en": utils.PathSearch("action_template_en", params, nil),
				"status":             utils.PathSearch("status", params, nil),
				"final_sub_action":   utils.PathSearch("final_sub_action", params, nil),
				"enum_data_list": flattenCocTicketOperationHistoriesDataEnumDataList(
					utils.PathSearch("enum_data_list", params, make([]interface{}, 0)).([]interface{})),
			}
			rst = append(rst, m)
		}
		return rst
	}
	return nil
}

func flattenCocTicketOperationHistoriesDataEnumDataList(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			m := map[string]interface{}{
				"is_deleted":     utils.PathSearch("is_deleted", params, nil),
				"match_type":     utils.PathSearch("match_type", params, nil),
				"ticket_id":      utils.PathSearch("ticket_id", params, nil),
				"real_ticket_id": utils.PathSearch("real_ticket_id", params, nil),
				"name_zh":        utils.PathSearch("name_zh", params, nil),
				"name_en":        utils.PathSearch("name_en", params, nil),
				"user_name":      utils.PathSearch("user_name", params, nil),
				"biz_id":         utils.PathSearch("biz_id", params, nil),
				"prop_id":        utils.PathSearch("prop_id", params, nil),
				"model_id":       utils.PathSearch("model_id", params, nil),
			}
			rst = append(rst, m)
		}
		return rst
	}
	return nil
}
