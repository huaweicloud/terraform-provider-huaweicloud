package dcs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var backupImportTaskNonUpdatableParams = []string{"task_name", "migration_type", "migration_method", "description",
	"backup_files", "backup_files.*.file_source", "backup_files.*.bucket_name", "backup_files.*.backup_id",
	"backup_files.*.files", "backup_files.*.files.*.file_name", "backup_files.*.files.*.size", "backup_files.*.files.*.update_at",
	"target_instance", "target_instance.*.id", "target_instance.*.password"}

// @API DCS POST /v2/{project_id}/migration-task
// @API DCS GET /v2/{project_id}/migration-task/{task_id}
// @API DCS DELETE /v2/{project_id}/migration-tasks/delete
// @API DCS GET /v2/{project_id}/migration-tasks
func ResourceDcsBackupImportTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsBackupImportTaskCreate,
		ReadContext:   resourceDcsBackupImportTaskRead,
		UpdateContext: resourceDcsBackupImportTaskUpdate,
		DeleteContext: resourceDcsBackupImportTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: config.FlexibleForceNew(backupImportTaskNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"migration_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"migration_method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"backup_files": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     backupImportTaskBackupFilesSchema(),
			},
			"target_instance": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     backupImportTaskTargetInstanceSchema(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"released_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func backupImportTaskBackupFilesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"files": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     backupImportTaskBackupFilesFilesSchema(),
			},
			"backup_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func backupImportTaskBackupFilesFilesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_name": {
				Type:     schema.TypeString,
				Required: true,
			}, "size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"update_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func backupImportTaskTargetInstanceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDcsBackupImportTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/migration-task"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateBackupImportTasBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating backup import task: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating backup import task: id is not found in API response")
	}

	d.SetId(id)

	err = checkMigrationTaskFinish(ctx, client, id, []string{"SUCCESS"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcsBackupImportTaskRead(ctx, d, meta)
}

func buildCreateBackupImportTasBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"task_name":        d.Get("task_name"),
		"migration_type":   d.Get("migration_type"),
		"migration_method": d.Get("migration_method"),
		"backup_files":     buildCreateBackupImportTaskBackupFilesBodyParams(d),
		"target_instance":  buildCreateBackupImportTaskTargetInstanceBodyParams(d),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func buildCreateBackupImportTaskBackupFilesBodyParams(d *schema.ResourceData) map[string]interface{} {
	backupFiles := d.Get("backup_files").([]interface{})[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"file_source": backupFiles["file_source"],
		"bucket_name": utils.ValueIgnoreEmpty(backupFiles["bucket_name"]),
		"backup_id":   utils.ValueIgnoreEmpty(backupFiles["backup_id"]),
		"files":       utils.ValueIgnoreEmpty(buildCreateBackupImportTaskBackupFilesFilesBodyParams(backupFiles["files"])),
	}
	return bodyParams
}

func buildCreateBackupImportTaskBackupFilesFilesBodyParams(backupFiles interface{}) []map[string]interface{} {
	files := backupFiles.([]interface{})
	if len(files) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, 0, len(files))
	for _, v := range files {
		file := v.(map[string]interface{})
		bodyParams = append(bodyParams, map[string]interface{}{
			"file_name": file["file_name"],
			"size":      utils.ValueIgnoreEmpty(file["size"]),
			"update_at": utils.ValueIgnoreEmpty(file["update_at"]),
		})
	}
	return bodyParams
}

func buildCreateBackupImportTaskTargetInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	targetInstance := d.Get("target_instance").([]interface{})[0].(map[string]interface{})
	bodyParams := map[string]interface{}{
		"id":       targetInstance["id"],
		"password": utils.ValueIgnoreEmpty(targetInstance["password"]),
	}
	return bodyParams
}

func resourceDcsBackupImportTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getRespBody, err := getMigrationTask(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DCS.4133"),
			"error getting DCS backup import task")
	}

	status := utils.PathSearch("status", getRespBody, "").(string)
	if status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting DCS backup import task")
	}

	mErr = multierror.Append(nil,
		d.Set("region", region),
		d.Set("task_name", utils.PathSearch("task_name", getRespBody, nil)),
		d.Set("migration_type", utils.PathSearch("migration_type", getRespBody, nil)),
		d.Set("migration_method", utils.PathSearch("migration_method", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("backup_files", flattenBackupImportTaskBackupFiles(getRespBody)),
		d.Set("target_instance", flattenBackupImportTaskTargetInstance(getRespBody)),
		d.Set("status", status),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
		d.Set("released_at", utils.PathSearch("released_at", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBackupImportTaskBackupFiles(resp interface{}) []interface{} {
	curJson := utils.PathSearch("backup_files", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"file_source": utils.PathSearch("file_source", curJson, nil),
			"bucket_name": utils.PathSearch("bucket_name", curJson, nil),
			"files":       flattenBackupImportTaskBackupFilesFiles(curJson),
			"backup_id":   utils.PathSearch("backup_record.backup_id", curJson, nil),
		},
	}
	return rst
}

func flattenBackupImportTaskBackupFilesFiles(resp interface{}) []interface{} {
	curJson := utils.PathSearch("files", resp, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"file_name": utils.PathSearch("file_name", v, nil),
			"size":      utils.PathSearch("size", v, nil),
			"update_at": utils.PathSearch("update_at", v, nil),
		})
	}
	return rst
}

func flattenBackupImportTaskTargetInstance(resp interface{}) []interface{} {
	curJson := utils.PathSearch("target_instance", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"id":   utils.PathSearch("id", curJson, nil),
			"name": utils.PathSearch("name", curJson, nil),
		},
	}
	return rst
}

func resourceDcsBackupImportTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsBackupImportTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	err = deleteMigrationTask(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = checkMigrationTaskDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func deleteMigrationTask(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v2/{project_id}/migration-tasks/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = buildDeleteMigrationTaskBodyParams(d.Id())

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return err
	}
	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return err
	}

	id := utils.PathSearch("task_id_list[0]", deleteRespBody, "").(string)
	if id == "" {
		return errors.New("error deleting backup import task, id is not found in the response")
	}

	return nil
}

func buildDeleteMigrationTaskBodyParams(taskId string) interface{} {
	bodyParams := map[string]interface{}{
		"task_id_list": []string{taskId},
	}
	return bodyParams
}

func checkMigrationTaskFinish(ctx context.Context, client *golangsdk.ServiceClient, taskId string, target []string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       target,
		Refresh:      migrationTaskRefreshFunc(client, taskId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for migration task(%s) to be completed: %s ", taskId, err)
	}
	return nil
}

func migrationTaskRefreshFunc(client *golangsdk.ServiceClient, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getMigrationTask(client, taskId)
		if err != nil {
			if errCode, ok := err.(golangsdk.ErrDefault400); ok {
				var response interface{}
				if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
					errorCode := utils.PathSearch("error_code", response, "").(string)
					if errorCode == "DCS.4133" {
						return "", "DELETED", nil
					}
				}
			}
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		directReturnStatus := []string{"SUCCESS", "FAILED", "TERMINATED", "INCRMIGEATING", "MIGRATION_FAILED",
			"RELEASED", "DELETED"}
		if utils.StrSliceContains(directReturnStatus, status) {
			return getRespBody, status, nil
		}

		return getRespBody, "PENDING", nil
	}
}

func getMigrationTask(client *golangsdk.ServiceClient, taskId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/migration-task/{task_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}

func checkMigrationTaskDeleted(ctx context.Context, client *golangsdk.ServiceClient, taskId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      migrationTaskDeleteRefreshFunc(client, taskId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for migration task(%s) to be deleted: %s ", taskId, err)
	}
	return nil
}

func migrationTaskDeleteRefreshFunc(client *golangsdk.ServiceClient, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRespBody, err := getMigrationTaskList(client)
		if err != nil {
			return nil, "ERROR", err
		}

		task := utils.PathSearch(fmt.Sprintf("migration_tasks[?task_id=='%s']|[0]", taskId), getRespBody, nil)
		if task == nil {
			return getRespBody, "DELETED", nil
		}

		return getRespBody, "PENDING", nil
	}
}

// when the migration task is deleted, the value of the task status queried by the task detail API may be SUCCESS,
// but it can not be queried by the list API
func getMigrationTaskList(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v2/{project_id}/migration-tasks"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getResp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return nil, err
	}

	getRespJson, err := json.Marshal(getResp)
	if err != nil {
		return nil, err
	}
	var getRespBody interface{}
	err = json.Unmarshal(getRespJson, &getRespBody)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}
