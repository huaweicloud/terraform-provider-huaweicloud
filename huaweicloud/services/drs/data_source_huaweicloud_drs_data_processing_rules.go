package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/data-processing-rules
func DataSourceDrsDataProcessingRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsDataProcessingRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_process_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataProcessInfoSchema(),
			},
		},
	}
}

func dataProcessInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter_conditions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     filterConditionsSchema(),
			},
			"is_batch_process": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"add_columns": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     addColumnsSchema(),
			},
			"ddl_operation": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dml_operation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_object_column_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dbObjectColumnInfoSchema(),
			},
			"db_or_table_rename_rule": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dbOrTableRenameRuleSchema(),
			},
			"db_object": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dbObjectSchema(),
			},
			"is_synchronized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process_rule_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func filterConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filtering_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func addColumnsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"column_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"column_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"column_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dbObjectColumnInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"column_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     columnInfosSchema(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func columnInfosSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"column_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"column_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_key_or_unique_index": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"column_mapped_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_filtered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_partition_key": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dbOrTableRenameRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"prefix_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"suffix_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dbObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_scope": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_root_db": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     targetRootDbSchema(),
			},
			// Convert this field to JSON string value from API response.
			"object_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func targetRootDbSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_encoding": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDataProcessingRulesQueryParams(offset int) string {
	if offset == 0 {
		return ""
	}

	return fmt.Sprintf("?offset=%d", offset)
}

func dataSourceDrsDataProcessingRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		jobId   = d.Get("job_id").(string)
		httpUrl = "v5/{project_id}/jobs/{job_id}/data-processing-rules"
		offset  = 0
		result  []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithParams := requestPath + buildDataProcessingRulesQueryParams(offset)
		resp, err := client.Request("GET", requestPathWithParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS data processing rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataList := utils.PathSearch("data_process_info", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data_process_info", flattenDataProcessingRules(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataProcessingRules(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"filter_conditions":       flattenFilterConditions(utils.PathSearch("filter_conditions", item, nil)),
			"is_batch_process":        utils.PathSearch("is_batch_process", item, nil),
			"add_columns":             flattenAddColumns(utils.PathSearch("add_columns", item, nil)),
			"ddl_operation":           flattenMapValues(utils.PathSearch("ddl_operation", item, nil)),
			"dml_operation":           utils.PathSearch("dml_operation", item, nil),
			"db_object_column_info":   flattenDbObjectColumnInfo(utils.PathSearch("db_object_column_info", item, nil)),
			"db_or_table_rename_rule": flattenDbOrTableRenameRule(utils.PathSearch("db_or_table_rename_rule", item, nil)),
			"db_object":               flattenDbObject(utils.PathSearch("db_object", item, nil)),
			"is_synchronized":         utils.PathSearch("is_synchronized", item, nil),
			"source":                  utils.PathSearch("source", item, nil),
			"process_rule_level":      utils.PathSearch("process_rule_level", item, nil),
		})
	}
	return result
}

func flattenMapValues(rawMap interface{}) map[string]string {
	mapRaw, ok := rawMap.(map[string]interface{})
	if !ok || len(mapRaw) == 0 {
		return nil
	}

	return utils.ExpandToStringMap(mapRaw)
}

func flattenFilterConditions(conditions interface{}) []interface{} {
	conditionsRaw, ok := conditions.([]interface{})
	if !ok || len(conditionsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(conditionsRaw))
	for _, item := range conditionsRaw {
		result = append(result, map[string]interface{}{
			"value":          utils.PathSearch("value", item, nil),
			"filtering_type": utils.PathSearch("filtering_type", item, nil),
		})
	}
	return result
}

func flattenAddColumns(columns interface{}) []interface{} {
	columnsRaw, ok := columns.([]interface{})
	if !ok || len(columnsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(columnsRaw))
	for _, item := range columnsRaw {
		result = append(result, map[string]interface{}{
			"column_type":  utils.PathSearch("column_type", item, nil),
			"column_name":  utils.PathSearch("column_name", item, nil),
			"column_value": utils.PathSearch("column_value", item, nil),
			"data_type":    utils.PathSearch("data_type", item, nil),
		})
	}
	return result
}

func flattenDbObjectColumnInfo(info interface{}) []interface{} {
	if info == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"db_name":      utils.PathSearch("db_name", info, nil),
			"schema_name":  utils.PathSearch("schema_name", info, nil),
			"table_name":   utils.PathSearch("table_name", info, nil),
			"column_infos": flattenColumnInfos(utils.PathSearch("column_infos", info, nil)),
			"total_count":  utils.PathSearch("total_count", info, nil),
		},
	}
}

func flattenColumnInfos(infos interface{}) []interface{} {
	infosRaw, ok := infos.([]interface{})
	if !ok || len(infosRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(infosRaw))
	for _, item := range infosRaw {
		result = append(result, map[string]interface{}{
			"column_name":                 utils.PathSearch("column_name", item, nil),
			"column_type":                 utils.PathSearch("column_type", item, nil),
			"primary_key_or_unique_index": utils.PathSearch("primary_key_or_unique_index", item, nil),
			"column_mapped_name":          utils.PathSearch("column_mapped_name", item, nil),
			"is_filtered":                 utils.PathSearch("is_filtered", item, nil),
			"is_partition_key":            utils.PathSearch("is_partition_key", item, nil),
		})
	}
	return result
}

func flattenDbOrTableRenameRule(rule interface{}) []interface{} {
	if rule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"prefix_name": utils.PathSearch("prefix_name", rule, nil),
			"suffix_name": utils.PathSearch("suffix_name", rule, nil),
			"type":        utils.PathSearch("type", rule, nil),
		},
	}
}

func flattenDbObject(dbObject interface{}) []interface{} {
	if dbObject == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"object_scope":   utils.PathSearch("object_scope", dbObject, nil),
			"target_root_db": flattenTargetRootDb(utils.PathSearch("target_root_db", dbObject, nil)),
			"object_info":    utils.JsonToString(utils.PathSearch("object_info", dbObject, nil)),
		},
	}
}

func flattenTargetRootDb(targetRootDb interface{}) []interface{} {
	if targetRootDb == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"db_name":     utils.PathSearch("db_name", targetRootDb, nil),
			"db_encoding": utils.PathSearch("db_encoding", targetRootDb, nil),
		},
	}
}
