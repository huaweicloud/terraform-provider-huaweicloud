// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/contentgateway/url-tasks
func DataSourceCacheUrlTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCacheUrlTasksRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the start timestamp, in milliseconds.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the end timestamp, in milliseconds.`,
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the refresh or preheat URL.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the URL status.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the file type.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Elem:        cacheUrlTasksSchema(),
				Computed:    true,
				Description: `The list of URL task information.`,
			},
		},
	}
}

func cacheUrlTasksSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The URL ID.`,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL status.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task type.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory refresh mode.`,
			},
			"task_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The task ID.`,
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The modification time.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file type.`,
			},
		},
	}
	return &sc
}

func resourceCacheUrlTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		limit    = int32(100)
		offset   int32
		respUrls []model.Urls
		mErr     *multierror.Error
	)

	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	request := &model.ShowUrlTaskInfoRequest{
		StartTime: utils.Int64IgnoreEmpty(int64(d.Get("start_time").(int))),
		EndTime:   utils.Int64IgnoreEmpty(int64(d.Get("end_time").(int))),
		Limit:     utils.Int32(limit),
		Url:       utils.StringIgnoreEmpty(d.Get("url").(string)),
		TaskType:  utils.StringIgnoreEmpty(d.Get("task_type").(string)),
		Status:    utils.StringIgnoreEmpty(d.Get("status").(string)),
		FileType:  utils.StringIgnoreEmpty(d.Get("file_type").(string)),
	}

	for {
		request.Offset = utils.Int32(offset)
		resp, err := hcCdnClient.ShowUrlTaskInfo(request)
		if err != nil {
			return diag.Errorf("error retrieving CDN cache url tasks: %s", err)
		}

		if resp == nil || resp.Result == nil || len(*resp.Result) == 0 {
			break
		}
		respUrls = append(respUrls, *resp.Result...)
		offset += limit
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("tasks", flattenCacheUrlTasks(respUrls)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCacheUrlTasks(respUrls []model.Urls) []interface{} {
	if len(respUrls) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respUrls))
	for _, v := range respUrls {
		rst = append(rst, map[string]interface{}{
			"id":          v.Id,
			"url":         v.Url,
			"status":      v.Status,
			"task_type":   v.Type,
			"mode":        v.Mode,
			"task_id":     v.TaskId,
			"modify_time": flattenCreatedAt(v.ModifyTime),
			"created_at":  flattenCreatedAt(v.CreateTime),
			"file_type":   v.FileType,
		})
	}
	return rst
}
