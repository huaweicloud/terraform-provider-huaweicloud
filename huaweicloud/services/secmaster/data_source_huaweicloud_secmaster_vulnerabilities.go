package secmaster

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Secmaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/vulnerability/search
func DataSourceVulnerabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVulnerabilitiesRead,

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
			"sort_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"from_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to_date": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"condition": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     dataConditionSchema(),
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"format_version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workspace_id": {
							Type:     schema.TypeString,
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
						"dataclass_ref": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataDataClassRefSchema(),
						},
						"data_object": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataDataObjectSchema(),
						},
					},
				},
			},
		},
	}
}

func dataDataObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"vul_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"first_observed_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"batch_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remediation": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDataObjectRemediationSchema(),
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extend_properties": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDataObjectExtendPropertiesSchema(),
			},
			"region_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vulnerability_type": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDataObjectVulnerabilityTypeSchema(),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_observed_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDataObjectResourceSchema(),
			},
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vulnerability": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDataObjectVulnerabilitySchema(),
			},
			"dataclass_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDataSourceSchema(),
			},
			"arrive_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"environment": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataEnvironmentSchema(),
			},
			"trigger_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"handled": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataEnvironmentSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func dataDataSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain_id": {
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
			"company_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"product_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_feature": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataDataObjectVulnerabilitySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"solution": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"repair_severity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"related": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataDataObjectResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ep_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataDataObjectVulnerabilityTypeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_en": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_zh": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vulnerability_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vulnerability_type_en": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vulnerability_type_zh": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataDataObjectExtendPropertiesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"operations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataDataObjectExtendPropertiesOperationsSchema(),
			},
		},
	}
}

func dataDataObjectExtendPropertiesOperationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_build_in": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataDataObjectRemediationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recommendation": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataDataClassRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     dataConditionConditionsSchema(),
			},
			"logics": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataConditionConditionsSchema() *schema.Resource {
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

func buildVulnerabilitiesBodyParams(d *schema.ResourceData, offset int) map[string]interface{} {
	rst := map[string]interface{}{
		"sort_by":   utils.ValueIgnoreEmpty(d.Get("sort_by")),
		"order":     utils.ValueIgnoreEmpty(d.Get("order")),
		"from_date": utils.ValueIgnoreEmpty(d.Get("from_date")),
		"to_date":   utils.ValueIgnoreEmpty(d.Get("to_date")),
		"condition": buildConditionBodyParams(d.Get("condition").([]interface{})),
		"offset":    utils.ValueIgnoreEmpty(offset),
	}

	return rst
}

func buildConditionBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"conditions": buildConditionConditionsBodyParams(rawMap["conditions"].([]interface{})),
		"logics":     utils.ValueIgnoreEmpty(rawMap["logics"]),
	}
}

func buildConditionConditionsBodyParams(rawArray []interface{}) []map[string]interface{} {
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

func mapToReader(data map[string]interface{}) (io.Reader, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBytes), nil
}

func dataSourceVulnerabilitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/vulnerability/search"
		offset  = 0
		allData = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		// `golangsdk` adds a newline character when converting the `JSONBody` data sequence.
		// This newline character will cause the API to report an error `400`, so it is replaced here with `RawBody`.
		rawBody, err := mapToReader(utils.RemoveNil(buildVulnerabilitiesBodyParams(d, offset)))
		if err != nil {
			return diag.Errorf("error converting map to reader: %s", err)
		}

		requestOpt.RawBody = rawBody
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster vulnerability list: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		data := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(data) == 0 {
			break
		}

		allData = append(allData, data...)
		totalCount := utils.PathSearch("total", respBody, 0).(int)
		if len(allData) >= totalCount {
			break
		}

		offset += len(data)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenVulnerabilities(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVulnerabilities(respArray []interface{}) []map[string]interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"format_version": utils.PathSearch("format_version", v, nil),
			"version":        utils.PathSearch("version", v, nil),
			"project_id":     utils.PathSearch("project_id", v, nil),
			"workspace_id":   utils.PathSearch("workspace_id", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
			"dataclass_ref":  flattenDataClassRefAttribute(utils.PathSearch("dataclass_ref", v, nil)),
			"data_object":    flattenDataObjectAttribute(utils.PathSearch("data_object", v, nil)),
		})
	}

	return rst
}

func flattenDataClassRefAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":   utils.PathSearch("id", respBody, nil),
			"name": utils.PathSearch("name", respBody, nil),
		},
	}
}

func flattenDataObjectAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"vul_name":            utils.PathSearch("vul_name", respBody, nil),
			"first_observed_time": utils.PathSearch("first_observed_time", respBody, nil),
			"batch_number":        utils.PathSearch("batch_number", respBody, nil),
			"description":         utils.PathSearch("description", respBody, nil),
			"resource_num":        utils.PathSearch("resource_num", respBody, nil),
			"domain_id":           utils.PathSearch("domain_id", respBody, nil),
			"workspace_id":        utils.PathSearch("workspace_id", respBody, nil),
			"remediation":         flattenRemediationAttribute(utils.PathSearch("remediation", respBody, nil)),
			"domain_name":         utils.PathSearch("domain_name", respBody, nil),
			"update_time":         utils.PathSearch("update_time", respBody, nil),
			"is_deleted":          utils.PathSearch("is_deleted", respBody, nil),
			"project_id":          utils.PathSearch("project_id", respBody, nil),
			"extend_properties":   flattenExtendPropertiesAttribute(utils.PathSearch("extend_properties", respBody, nil)),
			"region_name":         utils.PathSearch("region_name", respBody, nil),
			"id":                  utils.PathSearch("id", respBody, nil),
			"vulnerability_type":  flattenVulnerabilityTypeAttribute(utils.PathSearch("vulnerability_type", respBody, nil)),
			"create_time":         utils.PathSearch("create_time", respBody, nil),
			"last_observed_time":  utils.PathSearch("last_observed_time", respBody, nil),
			"resource":            flattenResourceAttribute(utils.PathSearch("resource", respBody, nil)),
			"count":               utils.PathSearch("count", respBody, nil),
			"region_id":           utils.PathSearch("region_id", respBody, nil),
			"vulnerability":       flattenVulnerabilityAttribute(utils.PathSearch("vulnerability", respBody, nil)),
			"dataclass_id":        utils.PathSearch("dataclass_id", respBody, nil),
			"version":             utils.PathSearch("version", respBody, nil),
			"data_source":         flattenDataSourceAttribute(utils.PathSearch("data_source", respBody, nil)),
			"arrive_time":         utils.PathSearch("arrive_time", respBody, nil),
			"environment":         flattenEnvironmentAttribute(utils.PathSearch("environment", respBody, nil)),
			"trigger_flag":        utils.PathSearch("trigger_flag", respBody, nil),
			"handled":             utils.PathSearch("handled", respBody, nil),
		},
	}
}

func flattenEnvironmentAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"domain_id":   utils.PathSearch("domain_id", respBody, nil),
			"domain_name": utils.PathSearch("domain_name", respBody, nil),
			"project_id":  utils.PathSearch("project_id", respBody, nil),
			"region_id":   utils.PathSearch("region_id", respBody, nil),
			"region_name": utils.PathSearch("region_name", respBody, nil),
			"vendor_type": utils.PathSearch("vendor_type", respBody, nil),
		},
	}
}

func flattenDataSourceAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"domain_id":       utils.PathSearch("domain_id", respBody, nil),
			"project_id":      utils.PathSearch("project_id", respBody, nil),
			"region_id":       utils.PathSearch("region_id", respBody, nil),
			"company_name":    utils.PathSearch("company_name", respBody, nil),
			"source_type":     utils.PathSearch("source_type", respBody, nil),
			"product_name":    utils.PathSearch("product_name", respBody, nil),
			"product_feature": utils.PathSearch("product_feature", respBody, nil),
		},
	}
}

func flattenVulnerabilityAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":              utils.PathSearch("id", respBody, nil),
			"type":            utils.PathSearch("type", respBody, nil),
			"url":             utils.PathSearch("url", respBody, nil),
			"status":          utils.PathSearch("status", respBody, nil),
			"level":           utils.PathSearch("level", respBody, nil),
			"reason":          utils.PathSearch("reason", respBody, nil),
			"solution":        utils.PathSearch("solution", respBody, nil),
			"repair_severity": utils.PathSearch("repair_severity", respBody, nil),
			"related":         utils.PathSearch("related", respBody, nil),
			"tags":            utils.PathSearch("tags", respBody, nil),
		},
	}
}

func flattenResourceAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":         utils.PathSearch("id", respBody, nil),
			"name":       utils.PathSearch("name", respBody, nil),
			"type":       utils.PathSearch("type", respBody, nil),
			"provider":   utils.PathSearch("provider", respBody, nil),
			"region_id":  utils.PathSearch("region_id", respBody, nil),
			"domain_id":  utils.PathSearch("domain_id", respBody, nil),
			"project_id": utils.PathSearch("project_id", respBody, nil),
			"ep_id":      utils.PathSearch("ep_id", respBody, nil),
			"tags":       utils.PathSearch("tags", respBody, nil),
		},
	}
}

func flattenVulnerabilityTypeAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":                    utils.PathSearch("id", respBody, nil),
			"category":              utils.PathSearch("category", respBody, nil),
			"category_en":           utils.PathSearch("category_en", respBody, nil),
			"category_zh":           utils.PathSearch("category_zh", respBody, nil),
			"vulnerability_type":    utils.PathSearch("vulnerability_type", respBody, nil),
			"vulnerability_type_en": utils.PathSearch("vulnerability_type_en", respBody, nil),
			"vulnerability_type_zh": utils.PathSearch("vulnerability_type_zh", respBody, nil),
		},
	}
}

func flattenExtendPropertiesAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"operations": flattenExtendPropertiesOperationsAttribute(utils.PathSearch("operations", respBody, nil)),
		},
	}
}

func flattenExtendPropertiesOperationsAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"is_build_in": utils.PathSearch("is_build_in", respBody, nil),
		},
	}
}

func flattenRemediationAttribute(respBody interface{}) []map[string]interface{} {
	if respBody == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"recommendation": utils.PathSearch("recommendation", respBody, nil),
		},
	}
}
