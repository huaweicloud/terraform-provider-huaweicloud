package dsc

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

// @API DSC GET /v1/{project_id}/sdg/watermark/extract-logs
func DataSourceWatermarkExtractLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWatermarkExtractLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"water_mark_log": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     watermarkExtractLogsSchema(),
			},
		},
	}
}

func watermarkExtractLogsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"blind_watermark": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dest_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"doc_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"download_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_exist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"file_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finish_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_file_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"visible_watermark": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildWatermarkExtractLogsQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func dataSourceWatermarkExtractLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v1/{project_id}/sdg/watermark/extract-logs"
		product         = "dsc"
		limit           = 1000
		offset          = 0
		allWaterMarkLog = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildWatermarkExtractLogsQueryParams(limit, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC watermark extract logs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		logResp := utils.PathSearch("water_mark_log", respBody, make([]interface{}, 0)).([]interface{})
		allWaterMarkLog = append(allWaterMarkLog, logResp...)
		if len(logResp) < limit {
			break
		}

		offset += len(logResp)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("water_mark_log", flattenWatermarkExtractLogs(allWaterMarkLog)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWatermarkExtractLogs(waterMarkLogList []interface{}) []interface{} {
	if len(waterMarkLogList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(waterMarkLogList))
	for _, v := range waterMarkLogList {
		rst = append(rst, map[string]interface{}{
			"blind_watermark":   utils.PathSearch("blind_watermark", v, nil),
			"dest_url":          utils.PathSearch("dest_url", v, nil),
			"doc_path":          utils.PathSearch("doc_path", v, nil),
			"download_url":      utils.PathSearch("download_url", v, nil),
			"file_exist":        utils.PathSearch("file_exist", v, nil),
			"file_url":          utils.PathSearch("file_url", v, nil),
			"finish_num":        utils.PathSearch("finish_num", v, nil),
			"project_id":        utils.PathSearch("project_id", v, nil),
			"remark":            utils.PathSearch("remark", v, nil),
			"state":             utils.PathSearch("state", v, nil),
			"task_id":           utils.PathSearch("task_id", v, nil),
			"task_name":         utils.PathSearch("task_name", v, nil),
			"total_file_num":    utils.PathSearch("total_file_num", v, nil),
			"visible_watermark": utils.PathSearch("visible_watermark", v, nil),
		})
	}

	return rst
}
