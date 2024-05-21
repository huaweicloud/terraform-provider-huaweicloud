// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/historytasks
func DataSourceCacheHistoryTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCacheHistoryTasksRead,
		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the enterprise project to which the resource belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task status.`,
			},
			"start_date": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the query start time.`,
			},
			"end_date": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the query end time.`,
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the field used for sorting.`,
			},
			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting type.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the file type.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task type.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Elem:        cacheHistoryTasksSchema(),
				Computed:    true,
				Description: `The history task list.`,
			},
		},
	}
}

func cacheHistoryTasksSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task result.`,
			},
			"processing": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of URLs that are being processed.`,
			},
			"succeed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of URLs processed.`,
			},
			"failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of URLs that failed to be processed.`,
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of URLs in the task.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task type.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the task was created.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the file type.`,
			},
		},
	}
	return &sc
}

func buildHistoryTaskRequestBodyStatusOpts(status string) *model.ShowHistoryTasksRequestStatus {
	if status == "" {
		return nil
	}

	statusToReq := new(model.ShowHistoryTasksRequestStatus)
	if err := statusToReq.UnmarshalJSON([]byte(status)); err != nil {
		log.Printf("[WARN] failed to parse status %s: %s", status, err)
		return nil
	}
	return statusToReq
}

func buildHistoryTaskRequestBodyFileTypeOpts(fileType string) *model.ShowHistoryTasksRequestFileType {
	if fileType == "" {
		return nil
	}

	fileTypeToReq := new(model.ShowHistoryTasksRequestFileType)
	if err := fileTypeToReq.UnmarshalJSON([]byte(fileType)); err != nil {
		log.Printf("[WARN] failed to parse file type %s: %s", fileType, err)
		return nil
	}
	return fileTypeToReq
}

func buildHistoryTaskRequestBodyTaskTypeOpts(taskType string) *model.ShowHistoryTasksRequestTaskType {
	if taskType == "" {
		return nil
	}

	taskTypeToReq := new(model.ShowHistoryTasksRequestTaskType)
	if err := taskTypeToReq.UnmarshalJSON([]byte(taskType)); err != nil {
		log.Printf("[WARN] failed to parse task type %s: %s", taskType, err)
		return nil
	}
	return taskTypeToReq
}

func resourceCacheHistoryTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		pageSize   = int32(10000)
		PageNumber = int32(1)
		respTasks  []model.TasksObject
		mErr       *multierror.Error
	)

	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	request := &model.ShowHistoryTasksRequest{
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		PageSize:            utils.Int32(pageSize),
		Status:              buildHistoryTaskRequestBodyStatusOpts(d.Get("status").(string)),
		StartDate:           utils.Int64IgnoreEmpty(int64(d.Get("start_date").(int))),
		EndDate:             utils.Int64IgnoreEmpty(int64(d.Get("end_date").(int))),
		OrderField:          utils.StringIgnoreEmpty(d.Get("order_field").(string)),
		OrderType:           utils.StringIgnoreEmpty(d.Get("order_type").(string)),
		FileType:            buildHistoryTaskRequestBodyFileTypeOpts(d.Get("file_type").(string)),
		TaskType:            buildHistoryTaskRequestBodyTaskTypeOpts(d.Get("task_type").(string)),
	}

	for {
		request.PageNumber = utils.Int32(PageNumber)
		resp, err := hcCdnClient.ShowHistoryTasks(request)
		if err != nil {
			return diag.Errorf("error retrieving CDN cache history tasks: %s", err)
		}

		if resp == nil || resp.Tasks == nil || len(*resp.Tasks) == 0 {
			break
		}
		respTasks = append(respTasks, *resp.Tasks...)
		PageNumber++
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("tasks", flattenCacheHistoryTasks(respTasks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCacheHistoryTasks(respTasks []model.TasksObject) []interface{} {
	if len(respTasks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respTasks))
	for _, v := range respTasks {
		rst = append(rst, map[string]interface{}{
			"id":         v.Id,
			"status":     v.Status,
			"processing": v.Processing,
			"succeed":    v.Succeed,
			"failed":     v.Failed,
			"total":      v.Total,
			"task_type":  flattenHistoryTaskTaskType(v.TaskType),
			"created_at": flattenCreatedAt(v.CreateTime),
			"file_type":  flattenHistoryTaskFileType(v.FileType),
		})
	}
	return rst
}

func flattenHistoryTaskTaskType(taskType *model.TasksObjectTaskType) string {
	if taskType == nil {
		return ""
	}
	return taskType.Value()
}

func flattenHistoryTaskFileType(fileType *model.TasksObjectFileType) string {
	if fileType == nil {
		return ""
	}
	return fileType.Value()
}
