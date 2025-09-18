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

// @API COC POST /v2/incidents/{incident_id}/histories
func DataSourceCocIncidentActionHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocIncidentActionHistoriesRead,

		Schema: map[string]*schema.Schema{
			"incident_id": {
				Type:     schema.TypeString,
				Required: true,
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
							Elem:     DataSourceCocIncidentsFilters(),
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
				Elem:     DataSourceCocIncidentsFilters(),
			},
			"int_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     DataSourceCocIncidentsFilters(),
			},
			"string_filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     DataSourceCocIncidentsFilters(),
			},
			"sort_filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     DataSourceCocIncidentsFilters(),
			},
			"condition": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
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
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stop_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enum_data_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     DataSourceCocIncidentsHistoriesEnumDataListElem(),
						},
					},
				},
			},
		},
	}
}

func DataSourceCocIncidentsFilters() *schema.Resource {
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

func DataSourceCocIncidentsHistoriesEnumDataListElem() *schema.Resource {
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

func dataSourceCocIncidentActionHistoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/incidents/{incident_id}/histories"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	incidentID := d.Get("incident_id").(string)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{incident_id}", incidentID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildIncidentActionHistoriesQueryOpts(d)),
	}
	getResp, err := client.Request("POST", getPath, &getOpt)

	if err != nil {
		return diag.Errorf("error retrieving COC incident action histories: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	incidentHistories := flattenCocIncidentActionHistories(getRespBody)
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
		d.Set("data", incidentHistories),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildIncidentActionHistoriesQueryOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"incident_id":     d.Get("incident_id"),
		"count_filters":   buildIncidentActionHistoriesCountFiltersQueryOpts(d.Get("count_filters")),
		"fields":          utils.ValueIgnoreEmpty(d.Get("fields")),
		"group_by_filter": buildIncidentActionHistoriesSingleFiltersQueryOpts(d.Get("group_by_filter")),
		"int_filters":     buildIncidentActionHistoriesFiltersQueryOpts(d.Get("int_filters")),
		"string_filters":  buildIncidentActionHistoriesFiltersQueryOpts(d.Get("string_filters")),
		"sort_filter":     buildIncidentActionHistoriesSingleFiltersQueryOpts(d.Get("sort_filter")),
		"condition":       utils.ValueIgnoreEmpty(d.Get("condition")),
	}

	return bodyParams
}

func buildIncidentActionHistoriesCountFiltersQueryOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":    utils.ValueIgnoreEmpty(raw["name"]),
				"filters": buildIncidentActionHistoriesFiltersQueryOpts(utils.ValueIgnoreEmpty(raw["filters"])),
			}
		}
		return params
	}

	return nil
}

func buildIncidentActionHistoriesSingleFiltersQueryOpts(rawParam interface{}) map[string]interface{} {
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

func buildIncidentActionHistoriesFiltersQueryOpts(rawParams interface{}) []map[string]interface{} {
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

func flattenCocIncidentActionHistories(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	dataJson := utils.PathSearch("data.info", resp, make([]interface{}, 0))
	dataArray := dataJson.([]interface{})
	if len(dataArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataArray))
	for _, data := range dataArray {
		result = append(result, map[string]interface{}{
			"action":         utils.PathSearch("action", data, nil),
			"action_name_zh": utils.PathSearch("action_name_zh", data, nil),
			"action_name_en": utils.PathSearch("action_name_en", data, nil),
			"operator":       utils.PathSearch("operator", data, nil),
			"status":         utils.PathSearch("status", data, nil),
			"start_time":     utils.PathSearch("start_time", data, nil),
			"stop_time":      utils.PathSearch("stop_time", data, nil),
			"comment":        utils.PathSearch("comment", data, nil),
			"enum_data_list": flattenCocIncidentActionHistoriesEnumDataList(
				utils.PathSearch("enum_data_list", data, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenCocIncidentActionHistoriesEnumDataList(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"prop_id": utils.PathSearch("prop_id", raw, nil),
				"biz_id":  utils.PathSearch("biz_id", raw, nil),
				"name_zh": utils.PathSearch("name_zh", raw, nil),
				"name_en": utils.PathSearch("name_en", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}
