package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// There are too many subfield layers to list here.
var checkDataFilterNonUpdatableParams = []string{
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

// @API DRS POST /v5/{project_id}/job/{job_id}/data-filtering/check
func ResourceDrsCheckDataFilter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCheckDataFilterCreate,
		ReadContext:   resourceCheckDataFilterRead,
		UpdateContext: resourceCheckDataFilterUpdate,
		DeleteContext: resourceCheckDataFilterDelete,

		CustomizeDiff: config.FlexibleForceNew(checkDataFilterNonUpdatableParams),

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
				Required: true,
				Elem:     checkDataFilterDataProcessInfoSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func checkDataFilterDataProcessInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"filter_conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     checkDataFilterFilterConditionsSchema(),
			},
			"is_batch_process": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"add_columns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     checkDataFilterAddColumnsSchema(),
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
				Elem:     checkDataFilterDbObjectColumnInfoSchema(),
			},
			"db_or_table_rename_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     checkDataFilterDbOrTableRenameRuleSchema(),
			},
			"db_object": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     checkDataFilterDbObjectSchema(),
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

func checkDataFilterFilterConditionsSchema() *schema.Resource {
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

func checkDataFilterAddColumnsSchema() *schema.Resource {
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

func checkDataFilterDbObjectColumnInfoSchema() *schema.Resource {
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
				Elem:     checkDataFilterColumnInfosSchema(),
			},
			"total_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func checkDataFilterColumnInfosSchema() *schema.Resource {
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

func checkDataFilterDbOrTableRenameRuleSchema() *schema.Resource {
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

func checkDataFilterDbObjectSchema() *schema.Resource {
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
				Elem:     checkDataFilterTargetRootDbSchema(),
			},
			// Convert this field from struct to JSON string value.
			"object_info": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func checkDataFilterTargetRootDbSchema() *schema.Resource {
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

func buildCheckDataFilterBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"data_process_info": buildDataProcessInfoParams(d.Get("data_process_info").([]interface{})),
	}
	return bodyParams
}

func buildDataProcessInfoParams(rawArray []interface{}) []map[string]interface{} {
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
			"filter_conditions":       buildFilterConditionsParams(rawMap["filter_conditions"].([]interface{})),
			"is_batch_process":        utils.ValueIgnoreEmpty(rawMap["is_batch_process"]),
			"add_columns":             buildAddColumnsParams(rawMap["add_columns"].([]interface{})),
			"ddl_operation":           buildMapValues(rawMap["ddl_operation"]),
			"dml_operation":           utils.ValueIgnoreEmpty(rawMap["dml_operation"]),
			"db_object_column_info":   buildDbObjectColumnInfoParams(rawMap["db_object_column_info"].([]interface{})),
			"db_or_table_rename_rule": buildDbOrTableRenameRuleParams(rawMap["db_or_table_rename_rule"].([]interface{})),
			"db_object":               buildDbObjectParams(rawMap["db_object"].([]interface{})),
			"is_synchronized":         utils.ValueIgnoreEmpty(rawMap["is_synchronized"]),
			"source":                  utils.ValueIgnoreEmpty(rawMap["source"]),
			"process_rule_level":      utils.ValueIgnoreEmpty(rawMap["process_rule_level"]),
		})
	}

	return rst
}

func buildMapValues(rawMap interface{}) map[string]string {
	mapRaw, ok := rawMap.(map[string]interface{})
	if !ok || len(mapRaw) == 0 {
		return nil
	}

	return utils.ExpandToStringMap(mapRaw)
}

func buildFilterConditionsParams(rawArray []interface{}) []map[string]interface{} {
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

func buildAddColumnsParams(rawArray []interface{}) []map[string]interface{} {
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

func buildDbObjectColumnInfoParams(rawArray []interface{}) map[string]interface{} {
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
		"column_infos": buildColumnInfosParams(rawMap["column_infos"].([]interface{})),
		"total_count":  utils.ValueIgnoreEmpty(rawMap["total_count"]),
	}
}

func buildColumnInfosParams(rawArray []interface{}) []map[string]interface{} {
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

func buildDbOrTableRenameRuleParams(rawArray []interface{}) map[string]interface{} {
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

func buildDbObjectParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"object_scope":   utils.ValueIgnoreEmpty(rawMap["object_scope"]),
		"target_root_db": buildTargetRootDbParams(rawMap["target_root_db"].([]interface{})),
		"object_info":    utils.StringToJson(rawMap["object_info"].(string)),
	}
}

func buildTargetRootDbParams(rawArray []interface{}) map[string]interface{} {
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

func resourceCheckDataFilterCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/job/{job_id}/data-filtering/check"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCheckDataFilterBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error checking DRS data filter: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	checkId := utils.PathSearch("id", respBody, "").(string)
	if checkId == "" {
		return diag.Errorf("unable to find the ID from the API response")
	}
	d.SetId(checkId)

	mErr := multierror.Append(nil,
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCheckDataFilterRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCheckDataFilterUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCheckDataFilterDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to check DRS data filter. Deleting this resource will not 
delete the check result from the cloud, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
