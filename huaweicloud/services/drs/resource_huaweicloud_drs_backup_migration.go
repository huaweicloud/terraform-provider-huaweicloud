package drs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var backupMigrationNonUpdatableParams = []string{
	"base_info.*.engine_type",
	"base_info.*.tags",
	"base_info.*.tags.*.key",
	"base_info.*.tags.*.value",
	"base_info.*.enterprise_project_id",
	"target_db_info",
	"target_db_info.*.target_instance_id",
	"target_db_info.*.ms_file_stream_status",
	"target_db_info.*.file_id",
	"backup_info",
	"backup_info.*.file_source",
	"backup_info.*.bucket_name",
	"backup_info.*.files",
	"backup_info.*.files.*.name",
	"backup_info.*.files.*.obs_path",
	"backup_info.*.files.*.rds_version",
	"backup_info.*.files.*.rds_source_instance_id",
	"options",
	"options.*.is_cover",
	"options.*.is_default_restore",
	"options.*.is_last_backup",
	"options.*.is_precheck",
	"options.*.recovery_mode",
	"options.*.db_names",
	"options.*.reset_db_name_map",
	"options.*.is_delete_backup_file",
}

// @API DRS POST /v5/{project_id}/backup-migration-jobs
// @API DRS GET /v5/{project_id}/backup-migration-jobs/{job_id}
// @API DRS PUT /v5/{project_id}/backup-migration-jobs/{job_id}
// @API DRS DELETE /v5/{project_id}/backup-migration-jobs/{job_id}
func ResourceBackupMigration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupMigrationCreate,
		ReadContext:   resourceBackupMigrationRead,
		UpdateContext: resourceBackupMigrationUpdate,
		DeleteContext: resourceBackupMigrationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(backupMigrationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"base_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     backupMigrationBaseInfoSchema(),
			},
			"target_db_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     backupMigrationTargetDbInfoSchema(),
			},
			"backup_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     backupMigrationBackupInfoSchema(),
			},
			"options": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     backupMigrationOptionsSchema(),
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
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finish_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_db_names": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func backupMigrationBaseInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Field `name` can be updated.
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Field `description` can be updated.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     backupMigrationResourceTagSchema(),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func backupMigrationResourceTagSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func backupMigrationTargetDbInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ms_file_stream_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"file_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func backupMigrationBackupInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"files": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     backupMigrationBackupFileInfoSchema(),
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func backupMigrationBackupFileInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"obs_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rds_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rds_source_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func backupMigrationOptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_last_backup": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_precheck": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"recovery_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			// This field `is_cover` is bool type in API, but config it as string type here.
			"is_cover": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// This field `is_default_restore` is bool type in API, but config it as string type here.
			"is_default_restore": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// This field `is_delete_backup_file` is bool type in API, but config it as string type here.
			"is_delete_backup_file": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"db_names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"reset_db_name_map": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildResourceTagParams(rawArray []interface{}) []map[string]interface{} {
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
			"key":   utils.ValueIgnoreEmpty(rawMap["key"]),
			"value": utils.ValueIgnoreEmpty(rawMap["value"]),
		})
	}

	return rst
}

func buildBaseInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"name":                  rawMap["name"],
		"engine_type":           rawMap["engine_type"],
		"description":           utils.ValueIgnoreEmpty(rawMap["description"]),
		"tags":                  buildResourceTagParams(rawMap["tags"].([]interface{})),
		"enterprise_project_id": utils.ValueIgnoreEmpty(rawMap["enterprise_project_id"]),
	}
}

func buildTargetDbInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"target_instance_id":    rawMap["target_instance_id"],
		"ms_file_stream_status": utils.ValueIgnoreEmpty(rawMap["ms_file_stream_status"]),
		"file_id":               utils.ValueIgnoreEmpty(rawMap["file_id"]),
	}
}

func buildBackupFileInfoParams(rawArray []interface{}) []map[string]interface{} {
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
			"name":                   rawMap["name"],
			"obs_path":               rawMap["obs_path"], // The API has limitations; by default, an empty string is required.
			"rds_version":            utils.ValueIgnoreEmpty(rawMap["rds_version"]),
			"rds_source_instance_id": utils.ValueIgnoreEmpty(rawMap["rds_source_instance_id"]),
		})
	}

	return rst
}

func buildBackupInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"file_source": rawMap["file_source"],
		"files":       buildBackupFileInfoParams(rawMap["files"].([]interface{})),
		"bucket_name": utils.ValueIgnoreEmpty(rawMap["bucket_name"]),
	}
}

func convertStringToBoolParam(rawString string) interface{} {
	if rawString == "" {
		return nil
	}

	return rawString == "true"
}

func buildStringListParams(rawArray []interface{}) interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	return utils.ExpandToStringList(rawArray)
}

func buildStringMapParams(rawMap map[string]interface{}) interface{} {
	if len(rawMap) == 0 {
		return nil
	}

	rst := make(map[string]string)
	// The API restricts the key to be an empty string.
	for k, v := range rawMap {
		strVal, ok := v.(string)
		if !ok {
			continue
		}
		rst[k] = strVal
	}

	return rst
}

func buildOptionsParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"is_last_backup":        rawMap["is_last_backup"],
		"is_precheck":           rawMap["is_precheck"],
		"recovery_mode":         rawMap["recovery_mode"],
		"is_cover":              convertStringToBoolParam(rawMap["is_cover"].(string)),
		"is_default_restore":    convertStringToBoolParam(rawMap["is_default_restore"].(string)),
		"is_delete_backup_file": convertStringToBoolParam(rawMap["is_delete_backup_file"].(string)),
		"db_names":              buildStringListParams(rawMap["db_names"].([]interface{})),
		"reset_db_name_map":     buildStringMapParams(rawMap["reset_db_name_map"].(map[string]interface{})),
	}
}

func buildBackupMigrationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"base_info":      buildBaseInfoParams(d.Get("base_info").([]interface{})),
		"target_db_info": buildTargetDbInfoParams(d.Get("target_db_info").([]interface{})),
		"backup_info":    buildBackupInfoParams(d.Get("backup_info").([]interface{})),
		"options":        buildOptionsParams(d.Get("options").([]interface{})),
	}
	return bodyParams
}

func QueryMigrationDetail(client *golangsdk.ServiceClient, jobId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/backup-migration-jobs/{job_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForMigrationJobSuccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			jobDetail, err := QueryMigrationDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", jobDetail, "").(string)
			if status == "" {
				return nil, "ERROR", errors.New("status is not found in job detail API response")
			}

			if status == "SUCCESS" {
				return jobDetail, "COMPLETED", nil
			}

			// FAILED, PRECHECK FAILED
			if strings.Contains(status, "FAILED") {
				errorLog := utils.PathSearch("error_log", jobDetail, "").(string)
				return jobDetail, "FAILED", fmt.Errorf("error migrate backup: %s", errorLog)
			}

			return jobDetail, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceBackupMigrationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/backup-migration-jobs"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildBackupMigrationBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating DRS backup migration job: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job_id from the API response")
	}
	d.SetId(jobId)

	if err := waitingForMigrationJobSuccess(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for migration job (%s) to success: %s", d.Id(), err)
	}

	return resourceBackupMigrationRead(ctx, d, meta)
}

func flattenBaseInfoTags(respBody interface{}) []interface{} {
	respArray := utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return rst
}

func flattenBaseInfo(respBody interface{}) []interface{} {
	respMap := utils.PathSearch("base_info", respBody, nil)
	if respMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"name":                  utils.PathSearch("name", respMap, nil),
			"engine_type":           utils.PathSearch("engine_type", respMap, nil),
			"description":           utils.PathSearch("description", respMap, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", respMap, nil),
			"tags":                  flattenBaseInfoTags(respBody),
		},
	}
}

func flattenTargetDBInfo(respBody interface{}) []interface{} {
	respMap := utils.PathSearch("target_db_info", respBody, nil)
	if respMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"target_instance_id":    utils.PathSearch("target_instance_id", respMap, nil),
			"ms_file_stream_status": utils.PathSearch("ms_file_stream_status", respMap, nil),
			"file_id":               utils.PathSearch("file_id", respMap, nil),
		},
	}
}

func flattenBackupInfo(respBody interface{}, d *schema.ResourceData) []interface{} {
	respMap := utils.PathSearch("backup_info", respBody, nil)
	if respMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"file_source": utils.PathSearch("file_source", respMap, nil),
			"bucket_name": utils.PathSearch("bucket_name", respMap, nil),
			"files":       buildBackupFileInfoParams(d.Get("backup_info.0.files").([]interface{})),
		},
	}
}

func flattenBoolToString(respBool interface{}) string {
	if respBool == nil {
		return ""
	}

	boolValue, ok := respBool.(bool)
	if !ok {
		return ""
	}

	if boolValue {
		return "true"
	}

	return "false"
}

func flattenStringList(respBody interface{}) []string {
	respArray := utils.PathSearch("db_names", respBody, make([]interface{}, 0)).([]interface{})
	if len(respArray) == 0 {
		return nil
	}

	return utils.ExpandToStringList(respArray)
}

func flattenStringMap(respBody interface{}) map[string]string {
	respMap := utils.PathSearch("reset_db_name_map", respBody, make(map[string]interface{})).(map[string]interface{})
	if len(respMap) == 0 {
		return nil
	}

	rst := make(map[string]string)
	// The API restricts the key to be an empty string.
	for k, v := range respMap {
		strVal, ok := v.(string)
		if !ok {
			continue
		}
		rst[k] = strVal
	}

	return rst
}

func flattenRecoveryMode(respValue interface{}) string {
	strValue, ok := respValue.(string)
	if !ok {
		return ""
	}

	rst := ""
	switch strValue {
	case "1":
		rst = "full"
	case "2":
		rst = "incre"
	default:
		rst = strValue
	}

	return rst
}

func flattenOptions(respBody interface{}) []interface{} {
	respMap := utils.PathSearch("options", respBody, nil)
	if respMap == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"is_last_backup":        utils.PathSearch("is_last_backup", respMap, nil),
			"is_precheck":           utils.PathSearch("is_precheck", respMap, nil),
			"recovery_mode":         flattenRecoveryMode(utils.PathSearch("recovery_mode", respMap, nil)),
			"is_cover":              flattenBoolToString(utils.PathSearch("is_cover", respMap, nil)),
			"is_default_restore":    flattenBoolToString(utils.PathSearch("is_default_restore", respMap, nil)),
			"is_delete_backup_file": flattenBoolToString(utils.PathSearch("is_delete_backup_file", respMap, nil)),
			"db_names":              flattenStringList(respMap),
			"reset_db_name_map":     flattenStringMap(respMap),
		},
	}
}

func resourceBackupMigrationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	respBody, err := QueryMigrationDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DRS.10000010"),
			"error retrieving DRS backup migration job detail",
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("base_info", flattenBaseInfo(respBody)),
		d.Set("target_db_info", flattenTargetDBInfo(respBody)),
		d.Set("backup_info", flattenBackupInfo(respBody, d)),
		d.Set("options", flattenOptions(respBody)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("finish_time", utils.PathSearch("finish_time", respBody, nil)),
		d.Set("new_db_names", utils.PathSearch("new_db_names", respBody, nil)),
		d.Set("instance_name", utils.PathSearch("instance_name", respBody, nil)),
		d.Set("error_log", utils.PathSearch("error_log", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBackupMigrationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	if d.HasChanges("base_info.0.name", "base_info.0.description") {
		requestPath := client.Endpoint + "v5/{project_id}/backup-migration-jobs/{job_id}"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
			JSONBody: map[string]interface{}{
				"name":        d.Get("base_info.0.name"),
				"description": d.Get("base_info.0.description"),
			},
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating DRS backup migration job: %s", err)
		}
	}

	return resourceBackupMigrationRead(ctx, d, meta)
}

func resourceBackupMigrationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/backup-migration-jobs/{job_id}"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DRS backup migration job: %s", err)
	}

	return nil
}
