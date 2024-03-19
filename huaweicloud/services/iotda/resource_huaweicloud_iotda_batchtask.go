package iotda

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	iotdav5 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	batchTaskInitializing   = "Initializing"
	batchTaskWaiting        = "Waitting"
	batchTaskProcessing     = "Processing"
	batchTaskSuccess        = "Success"
	batchTaskFail           = "Fail"
	batchTaskPartialSuccess = "PartialSuccess"
)

// @API IoTDA POST /v5/iot/{project_id}/batchtask-files
// @API IoTDA DELETE /v5/iot/{project_id}/batchtask-files/{file_id}
// @API IoTDA POST /v5/iot/{project_id}/batchtasks
// @API IoTDA GET /v5/iot/{project_id}/batchtasks/{task_id}
// @API IoTDA DELETE /v5/iot/{project_id}/batchtasks/{task_id}
func ResourceBatchTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchTaskCreate,
		ReadContext:   resourceBatchTaskRead,
		DeleteContext: resourceBatchTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"targets": {
				Type:         schema.TypeList,
				Optional:     true,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				ExactlyOneOf: []string{"targets", "targets_filter", "targets_file"},
			},
			"targets_filter": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_ids": {
							Type:     schema.TypeList,
							Required: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"targets_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_desc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_progress": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"success": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"fail": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"task_details": {
				Type:      schema.TypeList,
				Sensitive: true,
				Computed:  true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"output": {
							Type:      schema.TypeString,
							Sensitive: true,
							Computed:  true,
						},
						"error": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"error_code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_msg": {
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

func resourceBatchTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	// First, determine whether to upload batch task files.
	var targetFileId *string
	if v, ok := d.GetOk("targets_file"); ok && v.(string) != "" {
		file, openErr := os.Open(v.(string))
		if openErr != nil {
			return diag.Errorf("error opening batch task file: %s", openErr)
		}
		defer file.Close()

		uploadOpts := buildBatchTaskFileUploadParams(file)
		uploadTaskFileResp, uploadErr := client.UploadBatchTaskFile(uploadOpts)
		if uploadErr != nil {
			return diag.Errorf("error uploading IoTDA batch task file: %s", uploadErr)
		}

		if uploadTaskFileResp == nil || uploadTaskFileResp.FileId == nil {
			return diag.Errorf("error uploading IoTDA batch task file: ID is not found in API response")
		}

		targetFileId = uploadTaskFileResp.FileId
	}

	createOpts := buildBatchTaskCreateParams(d, targetFileId)
	createTaskResp, createTaskErr := client.CreateBatchTask(createOpts)
	if createTaskErr != nil {
		// When creating the batch task fails, it is necessary to delete the uploaded task file.
		var deleteTaskFileErr error
		if targetFileId != nil {
			deleteTaskFileErr = deleteBatchTaskFile(client, *targetFileId)
		}

		return diag.Errorf("error creating IoTDA batch task: %s, %s", createTaskErr, deleteTaskFileErr)
	}

	if createTaskResp == nil || createTaskResp.TaskId == nil {
		var deleteTaskFileErr error
		if targetFileId != nil {
			deleteTaskFileErr = deleteBatchTaskFile(client, *targetFileId)
		}

		return diag.Errorf("error creating IoTDA batch task: ID is not found in API response, %s", deleteTaskFileErr)
	}

	d.SetId(*createTaskResp.TaskId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      batchTaskStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for IoTDA batch task to complete: %s", err)
	}

	return resourceBatchTaskRead(ctx, d, meta)
}

func buildBatchTaskFileUploadParams(file *os.File) *model.UploadBatchTaskFileRequest {
	req := model.UploadBatchTaskFileRequest{
		Body: &model.UploadBatchTaskFileRequestBody{
			File: &def.FilePart{
				Content: file,
			},
		},
	}
	return &req
}

func buildBatchTaskCreateParams(d *schema.ResourceData, fileId *string) *model.CreateBatchTaskRequest {
	req := model.CreateBatchTaskRequest{
		Body: &model.CreateBatchTask{
			TaskName:       d.Get("name").(string),
			TaskType:       d.Get("type").(string),
			AppId:          utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			Targets:        utils.ExpandToStringListPointer(d.Get("targets").([]interface{})),
			DocumentSource: fileId,
		},
	}

	targetsFilter := d.Get("targets_filter").([]interface{})
	if len(targetsFilter) > 0 {
		groupIDs := targetsFilter[0].(map[string]interface{})
		groupIDList := utils.ExpandToStringList(groupIDs["group_ids"].([]interface{}))

		targetGroupIDs := make(map[string]interface{})
		targetGroupIDs["group_ids"] = groupIDList

		req.Body.TargetsFilter = targetGroupIDs
	}

	return &req
}

func deleteBatchTaskFile(client *iotdav5.IoTDAClient, fileId string) error {
	deleteOpts := &model.DeleteBatchTaskFileRequest{
		FileId: fileId,
	}
	_, err := client.DeleteBatchTaskFile(deleteOpts)
	if err != nil {
		return fmt.Errorf("error deleting IoTDA batch task file after failed creation of IoTDA batch task: %s", err)
	}

	return nil
}

func batchTaskStateRefreshFunc(client *iotdav5.IoTDAClient, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// There is no need to handle pagination related parameters, as it only applies to the structure of subtasks.
		getOpts := &model.ShowBatchTaskRequest{
			TaskId: taskId,
		}
		resp, err := client.ShowBatchTask(getOpts)
		if err != nil {
			return nil, "ERROR", err
		}

		status := *resp.Batchtask.Status

		targetStatus := []string{
			batchTaskSuccess, batchTaskFail, batchTaskPartialSuccess,
		}
		if utils.StrSliceContains(targetStatus, status) {
			return resp, "COMPLETED", nil
		}

		pendingStatus := []string{
			batchTaskInitializing, batchTaskWaiting, batchTaskProcessing,
		}
		if !utils.StrSliceContains(pendingStatus, status) {
			return resp, *resp.Batchtask.Status, nil
		}

		return resp, "PENDING", nil
	}
}

func resourceBatchTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allTaskDetails []model.TaskDetail
		targetResp     *model.ShowBatchTaskResponse
		limit          = int32(50)
		offset         int32
	)

	// The purpose of pagination is to get all subtasks under the batch task.
	// The pagination parameter only takes effect on the structure of subtask details(task_details).
	for {
		getOpts := &model.ShowBatchTaskRequest{
			TaskId: d.Id(),
			Limit:  utils.Int32(limit),
			Offset: &offset,
		}
		// If the resource does not exist, the API will return a 404 status code.
		resp, getErr := client.ShowBatchTask(getOpts)
		if getErr != nil {
			return common.CheckDeletedDiag(d, getErr, "error querying IoTDA batch task")
		}

		if targetResp == nil {
			targetResp = resp
		}

		if len(*resp.TaskDetails) == 0 {
			break
		}
		allTaskDetails = append(allTaskDetails, *resp.TaskDetails...)
		offset += limit
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", targetResp.Batchtask.TaskName),
		d.Set("type", targetResp.Batchtask.TaskType),
		d.Set("status", targetResp.Batchtask.Status),
		d.Set("status_desc", targetResp.Batchtask.StatusDesc),
		d.Set("created_at", targetResp.Batchtask.CreateTime),
		d.Set("task_progress", flattenTaskProgress(targetResp.Batchtask.TaskProgress)),
		d.Set("task_details", flattenTaskDetails(allTaskDetails)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTaskProgress(taskProgress *model.TaskProgress) []interface{} {
	if taskProgress == nil {
		return nil
	}

	rst := map[string]interface{}{
		"total":   taskProgress.Total,
		"success": taskProgress.Success,
		"fail":    taskProgress.Fail,
	}

	return []interface{}{rst}
}

func flattenTaskDetails(taskDetails []model.TaskDetail) []interface{} {
	if len(taskDetails) == 0 {
		return nil
	}

	rst := make([]interface{}, len(taskDetails))
	for i, v := range taskDetails {
		rst[i] = map[string]interface{}{
			"target": v.Target,
			"status": v.Status,
			"output": v.Output,
			"error":  flattenErrorInfo(v.Error),
		}
	}

	return rst
}

func flattenErrorInfo(errorInfo *model.ErrorInfo) []interface{} {
	if errorInfo == nil {
		return nil
	}

	rst := map[string]interface{}{
		"error_code": errorInfo.ErrorCode,
		"error_msg":  errorInfo.ErrorMsg,
	}

	return []interface{}{rst}
}

func resourceBatchTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := model.DeleteBatchTaskRequest{
		TaskId: d.Id(),
	}
	_, err = client.DeleteBatchTask(&deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA batch task: %s", err)
	}

	return nil
}
