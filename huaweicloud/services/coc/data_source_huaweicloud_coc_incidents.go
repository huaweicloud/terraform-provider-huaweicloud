package coc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC POST /v2/incidents/list
func DataSourceCocIncidents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocIncidentsRead,

		Schema: map[string]*schema.Schema{
			"contain_sub_ticket": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"string_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     DataSourceCocIncidentsFiltersElem(),
			},
			"sort_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     DataSourceCocIncidentsFiltersElem(),
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"count_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"filters": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     DataSourceCocIncidentsFiltersElem(),
						},
					},
				},
			},
			"fields": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group_by_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     DataSourceCocIncidentsFiltersElem(),
			},
			"int_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     DataSourceCocIncidentsFiltersElem(),
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
						"current_cloud_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mtm_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mtm_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_service_interrupt": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"work_flow_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"phase": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"assignee": {
							Type:     schema.TypeString,
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
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"handle_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"incident_ownership": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enum_data_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     DataSourceCocIncidentsTicketsEnumDataListElem(),
						},
					},
				},
			},
		},
	}
}

func DataSourceCocIncidentsFiltersElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field": {
				Type:     schema.TypeString,
				Required: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"match_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"priority_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func DataSourceCocIncidentsTicketsEnumDataListElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"prop_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"biz_id": {
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
		},
	}
}

func dataSourceCocIncidentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/incidents/list"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	getPath := client.Endpoint + httpUrl

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildIncidentsQueryOpts(d)),
	}
	getResp, err := client.Request("POST", getPath, &getOpt)

	if err != nil {
		return diag.Errorf("error retrieving COC incidents: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	incidents := flattenCocIncidents(getRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("tickets", incidents),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildIncidentsQueryOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"contain_all":        true,
		"contain_total":      true,
		"contain_sub_ticket": utils.ValueIgnoreEmpty(d.Get("contain_sub_ticket")),
		"string_filters":     buildIncidentsFiltersQueryOpts(d.Get("string_filters")),
		"sort_filter":        buildIncidentsSingleFiltersQueryOpts(d.Get("sort_filter")),
		"condition":          utils.ValueIgnoreEmpty(d.Get("condition")),
		"count_filters":      buildIncidentsCountFiltersQueryOpts(d.Get("count_filters")),
		"fields":             utils.ValueIgnoreEmpty(d.Get("fields")),
		"group_by_filter":    buildIncidentsFiltersQueryOpts(d.Get("group_by_filter")),
		"int_filters":        buildIncidentsFiltersQueryOpts(d.Get("int_filters")),
		"ticket_types":       utils.ValueIgnoreEmpty(d.Get("ticket_types")),
	}

	return bodyParams
}

func buildIncidentsCountFiltersQueryOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":    utils.ValueIgnoreEmpty(raw["name"]),
				"filters": buildIncidentsFiltersQueryOpts(utils.ValueIgnoreEmpty(raw["filters"])),
			}
		}
		return params
	}

	return nil
}

func buildIncidentsFiltersQueryOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"field":         raw["field"],
				"operator":      raw["operator"],
				"values":        raw["values"],
				"name":          utils.ValueIgnoreEmpty(raw["name"]),
				"group":         utils.ValueIgnoreEmpty(raw["group"]),
				"match_type":    utils.ValueIgnoreEmpty(raw["match_type"]),
				"priority_type": utils.ValueIgnoreEmpty(raw["priority_type"]),
			}
		}
		return params
	}

	return nil
}

func buildIncidentsSingleFiltersQueryOpts(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"field":         raw["field"],
			"operator":      raw["operator"],
			"values":        raw["values"],
			"name":          utils.ValueIgnoreEmpty(raw["name"]),
			"group":         utils.ValueIgnoreEmpty(raw["group"]),
			"match_type":    utils.ValueIgnoreEmpty(raw["match_type"]),
			"priority_type": utils.ValueIgnoreEmpty(raw["priority_type"]),
		}

		return param
	}

	return nil
}

func flattenCocIncidents(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	dataJson := utils.PathSearch("data.tickets", resp, make([]interface{}, 0))
	dataArray := dataJson.([]interface{})
	if len(dataArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataArray))
	for _, data := range dataArray {
		result = append(result, map[string]interface{}{
			"current_cloud_service_id": utils.PathSearch("current_cloud_service_id", data, nil),
			"level_id":                 utils.PathSearch("level_id", data, nil),
			"mtm_region":               utils.PathSearch("mtm_region", data, nil),
			"source_id":                utils.PathSearch("source_id", data, nil),
			"forward_rule_id":          utils.PathSearch("forward_rule_id", data, nil),
			"enterprise_project_id":    utils.PathSearch("enterprise_project_id", data, nil),
			"mtm_type":                 utils.PathSearch("mtm_type", data, nil),
			"title":                    utils.PathSearch("title", data, nil),
			"description":              utils.PathSearch("description", data, nil),
			"ticket_id":                utils.PathSearch("ticket_id", data, nil),
			"is_service_interrupt":     utils.PathSearch("is_service_interrupt", data, nil),
			"work_flow_status":         utils.PathSearch("work_flow_status", data, nil),
			"phase":                    utils.PathSearch("phase", data, nil),
			"assignee":                 utils.PathSearch("assignee", data, nil),
			"creator":                  utils.PathSearch("creator", data, nil),
			"operator":                 utils.PathSearch("operator", data, nil),
			"update_time":              utils.PathSearch("update_time", data, nil),
			"create_time":              utils.PathSearch("create_time", data, nil),
			"start_time":               utils.PathSearch("start_time", data, nil),
			"handle_time":              utils.PathSearch("handle_time", data, nil),
			"incident_ownership":       utils.PathSearch("incident_ownership", data, nil),
			"enum_data_list": flattenCocIncidentsEnumDataList(
				utils.PathSearch("enum_data_list", data, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenCocIncidentsEnumDataList(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			m := map[string]interface{}{
				"prop_id": utils.PathSearch("prop_id", params, nil),
				"biz_id":  utils.PathSearch("biz_id", params, nil),
				"name_zh": utils.PathSearch("name_zh", params, nil),
				"name_en": utils.PathSearch("name_en", params, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}
