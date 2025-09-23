package iotda

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	// First, determine whether to upload batch task files.
	var targetFileId string
	if targetsFile, ok := d.GetOk("targets_file"); ok && targetsFile.(string) != "" {
		fileId, err := uploadBatchTaskFile(client, targetsFile.(string))
		if err != nil {
			return diag.FromErr(err)
		}

		targetFileId = fileId
	}

	createPath := client.Endpoint + "v5/iot/{project_id}/batchtasks"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildBatchTaskCreateParams(d, targetFileId)),
	}

	createResp, createTaskErr := client.Request("POST", createPath, &createOpt)
	if createTaskErr != nil {
		// When creating the batch task fails, it is necessary to delete the uploaded task file.
		var deleteTaskFileErr error
		if targetFileId != "" {
			deleteTaskFileErr = deleteBatchTaskFile(client, targetFileId)
		}

		return diag.Errorf("error creating IoTDA batch task: %s, %s", createTaskErr, deleteTaskFileErr)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", createRespBody, "").(string)
	if taskId == "" {
		var deleteTaskFileErr error
		if targetFileId != "" {
			deleteTaskFileErr = deleteBatchTaskFile(client, targetFileId)
		}

		return diag.Errorf("error creating IoTDA batch task: ID is not found in API response, %s", deleteTaskFileErr)
	}

	d.SetId(taskId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      batchTaskStateRefreshFunc(client, taskId),
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

func uploadBatchTaskFile(client *golangsdk.ServiceClient, targetsFile string) (string, error) {
	file, openErr := os.Open(targetsFile)
	if openErr != nil {
		return "", fmt.Errorf("error opening batch task file: %s", openErr)
	}

	defer file.Close()

	var requestBody bytes.Buffer
	multiPartWriter := multipart.NewWriter(&requestBody)

	formFile, err := multiPartWriter.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(formFile, file)
	if err != nil {
		return "", err
	}

	err = multiPartWriter.Close()
	if err != nil {
		return "", err
	}

	uploadBatchTaskFilePath := client.Endpoint + "v5/iot/{project_id}/batchtask-files"
	uploadBatchTaskFilePath = strings.ReplaceAll(uploadBatchTaskFilePath, "{project_id}", client.ProjectID)
	uploadBatchTaskFileOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":         multiPartWriter.FormDataContentType(),
			"X-Sdk-Content-Sha256": "UNSIGNED-PAYLOAD",
		},
		RawBody: &requestBody,
	}

	uploadResp, uploadErr := client.Request("POST", uploadBatchTaskFilePath, &uploadBatchTaskFileOpt)
	if uploadErr != nil {
		return "", fmt.Errorf("error uploading IoTDA batch task file: %s", uploadErr)
	}

	uploadRespBody, err := utils.FlattenResponse(uploadResp)
	if err != nil {
		return "", err
	}

	fileId := utils.PathSearch("file_id", uploadRespBody, "").(string)
	if fileId == "" {
		return "", errors.New("error uploading IoTDA batch task file: ID is not found in API response")
	}

	return fileId, nil
}

func buildBatchTaskCreateParams(d *schema.ResourceData, fileId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"task_name":       d.Get("name"),
		"task_type":       d.Get("type"),
		"app_id":          utils.ValueIgnoreEmpty(d.Get("space_id").(string)),
		"targets":         utils.ExpandToStringList(d.Get("targets").([]interface{})),
		"document_source": fileId,
	}

	targetsFilter := d.Get("targets_filter").([]interface{})
	if len(targetsFilter) > 0 {
		groupIDs := targetsFilter[0].(map[string]interface{})
		groupIDList := utils.ExpandToStringList(groupIDs["group_ids"].([]interface{}))

		targetGroupIDs := make(map[string]interface{})
		targetGroupIDs["group_ids"] = groupIDList

		bodyParams["targets_filter"] = targetGroupIDs
	}

	return bodyParams
}

func deleteBatchTaskFile(client *golangsdk.ServiceClient, fileId string) error {
	deletePath := client.Endpoint + "v5/iot/{project_id}/batchtask-files/{file_id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{file_id}", fileId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return fmt.Errorf("error deleting IoTDA batch task file after failed creation of IoTDA batch task: %s", err)
	}

	return nil
}

func batchTaskStateRefreshFunc(client *golangsdk.ServiceClient, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// There is no need to handle pagination related parameters, as it only applies to the structure of subtasks.
		resp, err := GetBatchTaskById(client, taskId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("batchtask.status", resp, "").(string)
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
			return resp, status, nil
		}

		return resp, "PENDING", nil
	}
}

func GetBatchTaskById(client *golangsdk.ServiceClient, taskId string) (interface{}, error) {
	getPath := client.Endpoint + "v5/iot/{project_id}/batchtasks/{task_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving IoTDA batch task: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceBatchTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		isDerived      = WithDerivedAuth(cfg, region)
		product        = "iotda"
		httpUrl        = "v5/iot/{project_id}/batchtasks/{task_id}"
		offset         = 0
		targetResp     interface{}
		allTaskDetails = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", d.Id())
	getPath += "?limit=50"
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// The purpose of pagination is to get all subtasks under the batch task.
	// The pagination parameter only takes effect on the structure of subtask details(task_details).
	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			// If the resource does not exist, the API will return a `404` status code.
			return common.CheckDeletedDiag(d, err, "error retrieving IoTDA batch task")
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		if targetResp == nil {
			targetResp = getRespBody
		}

		taskDetailsResp := utils.PathSearch("task_details", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(taskDetailsResp) == 0 {
			break
		}

		allTaskDetails = append(allTaskDetails, taskDetailsResp...)
		offset += len(taskDetailsResp)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("batchtask.task_name", targetResp, nil)),
		d.Set("type", utils.PathSearch("batchtask.task_type", targetResp, nil)),
		d.Set("status", utils.PathSearch("batchtask.status", targetResp, nil)),
		d.Set("status_desc", utils.PathSearch("batchtask.status_desc", targetResp, nil)),
		d.Set("created_at", utils.PathSearch("batchtask.create_time", targetResp, nil)),
		d.Set("task_progress", flattenTaskProgress(utils.PathSearch("batchtask.task_progress", targetResp, nil))),
		d.Set("task_details", flattenTaskDetails(allTaskDetails)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTaskProgress(taskProgressResp interface{}) []interface{} {
	if taskProgressResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"total":   utils.PathSearch("total", taskProgressResp, nil),
		"success": utils.PathSearch("success", taskProgressResp, nil),
		"fail":    utils.PathSearch("fail", taskProgressResp, nil),
	}

	return []interface{}{result}
}

func flattenTaskDetails(taskDetailsResp []interface{}) []interface{} {
	if len(taskDetailsResp) == 0 {
		return nil
	}

	rst := make([]interface{}, len(taskDetailsResp))
	for i, v := range taskDetailsResp {
		rst[i] = map[string]interface{}{
			"target": utils.PathSearch("target", v, nil),
			"status": utils.PathSearch("status", v, nil),
			"output": utils.PathSearch("output", v, nil),
			"error":  flattenErrorInfo(utils.PathSearch("error", v, nil)),
		}
	}

	return rst
}

func flattenErrorInfo(errorInfoResp interface{}) []interface{} {
	if errorInfoResp == nil {
		return nil
	}

	result := map[string]interface{}{
		"error_code": utils.PathSearch("error_code", errorInfoResp, nil),
		"error_msg":  utils.PathSearch("error_msg", errorInfoResp, nil),
	}

	return []interface{}{result}
}

func resourceBatchTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + "v5/iot/{project_id}/batchtasks/{task_id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	// When deleting non-existent resource, the API returns a `404` status code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA batch task")
	}

	return nil
}
