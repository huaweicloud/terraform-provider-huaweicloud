package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS PUT /v5/{project_id}/jobs/{job_id}/filter-rules
// @API DRS DELETE /v5/{project_id}/jobs/{job_id}/filter-rules
func ResourceFilterRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFilterRulesCreate,
		ReadContext:   resourceFilterRulesRead,
		UpdateContext: resourceFilterRulesUpdate,
		DeleteContext: resourceFilterRulesDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_filter_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "common",
				ValidateFunc: validation.StringInSlice([]string{"common", "config"}, false),
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "job",
				ValidateFunc: validation.StringInSlice([]string{"job", "compare"}, false),
			},
			"general_filtering_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_object_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"parent_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_alias_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"object_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"source_filter_condition": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"advanced_setting_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"col_names": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prim_key_or_indexes": {
							Type:     schema.TypeString,
							Required: true,
						},
						"indexes": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func buildFilterRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"general_filtering_list": buildGeneralFilteringListBodyParams(d.Get("general_filtering_list")),
		"advanced_setting_list":  buildAdvancedSettingListBodyParams(d.Get("advanced_setting_list")),
		"source":                 utils.ValueIgnoreEmpty(d.Get("source")),
		"data_filter_type":       utils.ValueIgnoreEmpty(d.Get("data_filter_type")),
	}
	return bodyParams
}

func buildGeneralFilteringListBodyParams(rawParams interface{}) []interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, rawItem := range rawArray {
		item := rawItem.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"filter_object_list":      buildFilterObjectListBodyParams(item["filter_object_list"]),
			"source_filter_condition": utils.ValueIgnoreEmpty(item["source_filter_condition"]),
		})
	}
	return result
}

func buildFilterObjectListBodyParams(rawParams interface{}) []interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, rawItem := range rawArray {
		item := rawItem.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"id":                item["id"],
			"parent_id":         item["parent_id"],
			"object_name":       item["object_name"],
			"object_alias_name": utils.ValueIgnoreEmpty(item["object_alias_name"]),
			"object_type":       utils.ValueIgnoreEmpty(item["object_type"]),
		})
	}
	return result
}

func buildAdvancedSettingListBodyParams(rawParams interface{}) []interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, rawItem := range rawArray {
		item := rawItem.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"db_name":             item["db_name"],
			"table_name":          item["table_name"],
			"col_names":           item["col_names"],
			"prim_key_or_indexes": item["prim_key_or_indexes"],
			"indexes":             item["indexes"],
			"values":              item["values"],
		})
	}
	return result
}

func buildDeleteFilterRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"rule_ids": utils.ValueIgnoreEmpty(d.Get("rule_ids")),
		"type":     d.Get("data_filter_type"),
	}
}

func resourceFilterRulesCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFilterRulesBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error saving DRS data filter rules: %s", err)
	}

	d.SetId(d.Get("job_id").(string))

	return nil
}

func resourceFilterRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// The query API for data filter rules is an asynchronous collection task and does not directly return
	// the current resource state, so no processing is performed in the 'Read()' method.
	return nil
}

func resourceFilterRulesUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFilterRulesBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DRS data filter rules: %s", err)
	}

	return nil
}

func resourceFilterRulesDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// The delete API requires the rule IDs generated by the backend. If the rule IDs are not specified,
	// the resource is only removed from the state.
	ruleIds := d.Get("rule_ids").([]interface{})
	if len(ruleIds) == 0 {
		return nil
	}

	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDeleteFilterRulesBodyParams(d)),
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DRS data filter rules: %s", err)
	}

	return nil
}
