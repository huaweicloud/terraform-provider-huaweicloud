package drs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var updateDataProgressRulesNonUpdatableParams = []string{
	"job_id",
	"data_process_info",
	"data_process_info.*.filter_conditions",
	"data_process_info.*.filter_conditions.*.value",
	"data_process_info.*.filter_conditions.*.filtering_type",
	"data_process_info.*.is_batch_process",
	"data_process_info.*.add_columns",
	"data_process_info.*.add_columns.*.column_type",
	"data_process_info.*.add_columns.*.column_name",
	"data_process_info.*.add_columns.*.column_value",
	"data_process_info.*.add_columns.*.data_type",
	"data_process_info.*.ddl_operation",
	"data_process_info.*.dml_operation",
	"data_process_info.*.db_object_column_info",
	"data_process_info.*.db_object_column_info.*.db_name",
	"data_process_info.*.db_object_column_info.*.schema_name",
	"data_process_info.*.db_object_column_info.*.table_name",
	"data_process_info.*.db_object_column_info.*.column_infos",
	"data_process_info.*.db_object_column_info.*.column_infos.*.column_name",
	"data_process_info.*.db_object_column_info.*.column_infos.*.column_type",
	"data_process_info.*.db_object_column_info.*.column_infos.*.primary_key_or_unique_index",
	"data_process_info.*.db_object_column_info.*.column_infos.*.column_mapped_name",
	"data_process_info.*.db_object_column_info.*.column_infos.*.is_filtered",
	"data_process_info.*.db_object_column_info.*.column_infos.*.is_partition_key",
	"data_process_info.*.db_object_column_info.*.total_count",
	"data_process_info.*.db_or_table_rename_rule",
	"data_process_info.*.db_or_table_rename_rule.*.prefix_name",
	"data_process_info.*.db_or_table_rename_rule.*.suffix_name",
	"data_process_info.*.db_or_table_rename_rule.*.type",
	"data_process_info.*.db_object",
	"data_process_info.*.db_object.*.object_scope",
	"data_process_info.*.db_object.*.target_root_db",
	"data_process_info.*.db_object.*.target_root_db.*.db_name",
	"data_process_info.*.db_object.*.target_root_db.*.db_encoding",
	"data_process_info.*.db_object.*.object_info",
	"data_process_info.*.is_synchronized",
	"data_process_info.*.source",
	"data_process_info.*.process_rule_level",
}

// @API DRS PUT /v5/{project_id}/jobs/{job_id}/data-processing-rules
// @API DRS GET /v5/{project_id}/jobs/{job_id}/data-processing-rules/result
func ResourceDrsUpdateDataProgressRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUpdateDataProgressRulesCreate,
		ReadContext:   resourceUpdateDataProgressRulesRead,
		UpdateContext: resourceUpdateDataProgressRulesUpdate,
		DeleteContext: resourceUpdateDataProgressRulesDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(updateDataProgressRulesNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_process_info": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     updateDataProgressRulesDataProcessInfoSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func updateDataProgressRulesDataProcessInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter_conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     updateDataProgressRulesFilterConditionsSchema(),
			},
			"is_batch_process": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"add_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     updateDataProgressRulesAddColumnsSchema(),
			},
			"ddl_operation": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"dml_operation": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_object_column_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     updateDataProgressRulesDbObjectColumnInfoSchema(),
			},
			"db_or_table_rename_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     updateDataProgressRulesDbOrTableRenameRuleSchema(),
			},
			"db_object": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     updateDataProgressRulesDbObjectSchema(),
			},
			"is_synchronized": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"process_rule_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func updateDataProgressRulesFilterConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filtering_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func updateDataProgressRulesAddColumnsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"column_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"column_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"column_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func updateDataProgressRulesDbObjectColumnInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"table_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"column_infos": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     updateDataProgressRulesColumnInfosSchema(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func updateDataProgressRulesColumnInfosSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"column_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"column_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_key_or_unique_index": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"column_mapped_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_filtered": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_partition_key": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func updateDataProgressRulesDbOrTableRenameRuleSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"prefix_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"suffix_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func updateDataProgressRulesDbObjectSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"object_scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_root_db": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     updateDataProgressRulesTargetRootDbSchema(),
			},
			// Convert this field from struct to JSON string value.
			"object_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func updateDataProgressRulesTargetRootDbSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildUpdateDataProgressRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"data_process_info": buildUpdateDataProgressRulesDataProcessInfoParams(
			d.Get("data_process_info").([]interface{})),
	})
}

func buildUpdateDataProgressRulesDataProcessInfoParams(rawArray []interface{}) []map[string]interface{} {
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
			"filter_conditions": buildUpdateDataProgressRulesFilterConditionsParams(
				rawMap["filter_conditions"].([]interface{})),
			"is_batch_process": utils.ValueIgnoreEmpty(rawMap["is_batch_process"]),
			"add_columns": buildUpdateDataProgressRulesAddColumnsParams(
				rawMap["add_columns"].([]interface{})),
			"ddl_operation": buildUpdateDataProgressRulesMapValues(rawMap["ddl_operation"]),
			"dml_operation": utils.ValueIgnoreEmpty(rawMap["dml_operation"]),
			"db_object_column_info": buildUpdateDataProgressRulesDbObjectColumnInfoParams(
				rawMap["db_object_column_info"].([]interface{})),
			"db_or_table_rename_rule": buildUpdateDataProgressRulesDbOrTableRenameRuleParams(
				rawMap["db_or_table_rename_rule"].([]interface{})),
			"db_object": buildUpdateDataProgressRulesDbObjectParams(
				rawMap["db_object"].([]interface{})),
			"is_synchronized":    utils.ValueIgnoreEmpty(rawMap["is_synchronized"]),
			"source":             utils.ValueIgnoreEmpty(rawMap["source"]),
			"process_rule_level": utils.ValueIgnoreEmpty(rawMap["process_rule_level"]),
		})
	}

	return rst
}

func buildUpdateDataProgressRulesMapValues(rawMap interface{}) map[string]string {
	mapRaw, ok := rawMap.(map[string]interface{})
	if !ok || len(mapRaw) == 0 {
		return nil
	}

	return utils.ExpandToStringMap(mapRaw)
}

func buildUpdateDataProgressRulesFilterConditionsParams(rawArray []interface{}) []map[string]interface{} {
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
			"value":          utils.ValueIgnoreEmpty(rawMap["value"]),
			"filtering_type": utils.ValueIgnoreEmpty(rawMap["filtering_type"]),
		})
	}

	return rst
}

func buildUpdateDataProgressRulesAddColumnsParams(rawArray []interface{}) []map[string]interface{} {
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
			"column_type":  utils.ValueIgnoreEmpty(rawMap["column_type"]),
			"column_name":  utils.ValueIgnoreEmpty(rawMap["column_name"]),
			"column_value": utils.ValueIgnoreEmpty(rawMap["column_value"]),
			"data_type":    utils.ValueIgnoreEmpty(rawMap["data_type"]),
		})
	}

	return rst
}

func buildUpdateDataProgressRulesDbObjectColumnInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"db_name":      utils.ValueIgnoreEmpty(rawMap["db_name"]),
		"schema_name":  utils.ValueIgnoreEmpty(rawMap["schema_name"]),
		"table_name":   utils.ValueIgnoreEmpty(rawMap["table_name"]),
		"column_infos": buildUpdateDataProgressRulesColumnInfosParams(rawMap["column_infos"].([]interface{})),
		"total_count":  utils.ValueIgnoreEmpty(rawMap["total_count"]),
	}
}

func buildUpdateDataProgressRulesColumnInfosParams(rawArray []interface{}) []map[string]interface{} {
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
			"column_name":                 utils.ValueIgnoreEmpty(rawMap["column_name"]),
			"column_type":                 utils.ValueIgnoreEmpty(rawMap["column_type"]),
			"primary_key_or_unique_index": utils.ValueIgnoreEmpty(rawMap["primary_key_or_unique_index"]),
			"column_mapped_name":          utils.ValueIgnoreEmpty(rawMap["column_mapped_name"]),
			"is_filtered":                 utils.ValueIgnoreEmpty(rawMap["is_filtered"]),
			"is_partition_key":            utils.ValueIgnoreEmpty(rawMap["is_partition_key"]),
		})
	}

	return rst
}

func buildUpdateDataProgressRulesDbOrTableRenameRuleParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"prefix_name": utils.ValueIgnoreEmpty(rawMap["prefix_name"]),
		"suffix_name": utils.ValueIgnoreEmpty(rawMap["suffix_name"]),
		"type":        utils.ValueIgnoreEmpty(rawMap["type"]),
	}
}

func buildUpdateDataProgressRulesDbObjectParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"object_scope": utils.ValueIgnoreEmpty(rawMap["object_scope"]),
		"target_root_db": buildUpdateDataProgressRulesTargetRootDbParams(
			rawMap["target_root_db"].([]interface{})),
		"object_info": utils.StringToJson(rawMap["object_info"].(string)),
	}
}

func buildUpdateDataProgressRulesTargetRootDbParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"db_name":     utils.ValueIgnoreEmpty(rawMap["db_name"]),
		"db_encoding": utils.ValueIgnoreEmpty(rawMap["db_encoding"]),
	}
}

func resourceUpdateDataProgressRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		jobId   = d.Get("job_id").(string)
		httpUrl = "v5/{project_id}/jobs/{job_id}/data-processing-rules"
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
		JSONBody: utils.RemoveNil(buildUpdateDataProgressRulesBodyParams(d)),
	}

	resp, err := client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DRS data processing rules: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	queryId := utils.PathSearch("id", respBody, "").(string)
	if queryId == "" {
		return diag.Errorf("unable to find the ID from the API response")
	}
	d.SetId(queryId)

	if err := waitForUpdateDataProcessingRulesCompleted(ctx, client, jobId, queryId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForUpdateDataProcessingRulesCompleted(ctx context.Context, client *golangsdk.ServiceClient,
	jobId, queryId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"pending"},
		Target:  []string{"success"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getUpdateDataProcessingRulesResult(client, jobId, queryId)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "failed" {
				return respBody, "failed", errors.New("the data processing rules update task failed")
			}

			return respBody, status, nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DRS data processing rules update (%s) to complete: %s",
			queryId, err)
	}

	return nil
}

func getUpdateDataProcessingRulesResult(client *golangsdk.ServiceClient, jobId, queryId string) (interface{}, error) {
	httpUrl := "v5/{project_id}/jobs/{job_id}/data-processing-rules/result"
	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestPath += fmt.Sprintf("?query_id=%s", queryId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 202},
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DRS update data processing rules result: %s", err)
	}

	return utils.FlattenResponse(resp)
}

func resourceUpdateDataProgressRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUpdateDataProgressRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUpdateDataProgressRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to update DRS data processing rules. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information
    from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
