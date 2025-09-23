package sfsturbo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/{feature}/tasks
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/{feature}/tasks/{task_id}
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/fs/{feature}/tasks/{task_id}
func ResourceDuTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDuTaskCreate,
		ReadContext:   resourceDuTaskRead,
		DeleteContext: resourceDuTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDuTaskImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dir_usage": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dirUsageSchema(),
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dirUsageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_capacity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_count": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     fileCountSchema(),
			},
		},
	}
	return &sc
}

func fileCountSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"dir": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"regular": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pipe": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"char": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"block": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"socket": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"symlink": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildCreateDuTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"path": d.Get("path"),
	}
	return bodyParams
}

func resourceDuTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	createDuTaskHttpUrl := "sfs-turbo/shares/{share_id}/fs/{feature}/tasks"
	createDuTaskPath := client.ResourceBaseURL() + createDuTaskHttpUrl
	createDuTaskPath = strings.ReplaceAll(createDuTaskPath, "{share_id}", d.Get("share_id").(string))
	createDuTaskPath = strings.ReplaceAll(createDuTaskPath, "{feature}", "dir-usage")

	createDuTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createDuTaskOpt.JSONBody = utils.RemoveNil(buildCreateDuTaskBodyParams(d))
	createDuTaskResp, err := client.Request("POST", createDuTaskPath, &createDuTaskOpt)
	if err != nil {
		return diag.Errorf("error creating DU task: %s", err)
	}

	createDuTaskRespBody, err := utils.FlattenResponse(createDuTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", createDuTaskRespBody, "").(string)
	if taskId == "" {
		return diag.Errorf("unable to find the DU task ID from the API response")
	}

	d.SetId(taskId)

	return resourceDuTaskRead(ctx, d, meta)
}

func resourceDuTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	getDuTaskHttpUrl := "sfs-turbo/shares/{share_id}/fs/{feature}/tasks/{task_id}"
	getDuTaskPath := client.ResourceBaseURL() + getDuTaskHttpUrl
	getDuTaskPath = strings.ReplaceAll(getDuTaskPath, "{share_id}", d.Get("share_id").(string))
	getDuTaskPath = strings.ReplaceAll(getDuTaskPath, "{feature}", "dir-usage")
	getDuTaskPath = strings.ReplaceAll(getDuTaskPath, "{task_id}", d.Id())
	getDuTaskOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getDuTaskResp, err := client.Request("GET", getDuTaskPath, &getDuTaskOpts)
	if err != nil {
		if hasSpecifyErrorCode403(err, "SFS.TURBO.9000") {
			err = golangsdk.ErrDefault404{}
		}
		return common.CheckDeletedDiag(d, err, "error retrieving DU task")
	}

	getDuTaskRespBody, err := utils.FlattenResponse(getDuTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	beginTime := utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("begin_time", getDuTaskRespBody, "").(string), "2006-01-02 15:04:05")
	endTime := utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("end_time", getDuTaskRespBody, "").(string), "2006-01-02 15:04:05")

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("path", utils.PathSearch("dir_usage.path", getDuTaskRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getDuTaskRespBody, nil)),
		d.Set("dir_usage", flattenDuTaskResponse(utils.PathSearch("dir_usage", getDuTaskRespBody, make([]interface{}, 0)))),
		d.Set("begin_time", utils.FormatTimeStampRFC3339(beginTime/1000, false)),
		d.Set("end_time", utils.FormatTimeStampRFC3339(endTime/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDuTaskResponse(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	params := map[string]interface{}{
		"path":          utils.PathSearch("path", resp, nil),
		"used_capacity": utils.PathSearch("used_capacity", resp, nil),
		"message":       utils.PathSearch("message", resp, nil),
		"file_count":    flattenFileCount(utils.PathSearch("file_count", resp, nil)),
	}

	return []map[string]interface{}{params}
}

func flattenFileCount(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	params := map[string]interface{}{
		"dir":     utils.PathSearch("dir", resp, nil),
		"regular": utils.PathSearch("regular", resp, nil),
		"pipe":    utils.PathSearch("pipe", resp, nil),
		"char":    utils.PathSearch("char", resp, nil),
		"block":   utils.PathSearch("block", resp, nil),
		"socket":  utils.PathSearch("socket", resp, nil),
		"symlink": utils.PathSearch("symlink", resp, nil),
	}

	return []map[string]interface{}{params}
}

func resourceDuTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	deleteDuTaskHttpUrl := "sfs-turbo/shares/{share_id}/fs/{feature}/tasks/{task_id}"
	deleteDuTaskPath := client.ResourceBaseURL() + deleteDuTaskHttpUrl
	deleteDuTaskPath = strings.ReplaceAll(deleteDuTaskPath, "{share_id}", d.Get("share_id").(string))
	deleteDuTaskPath = strings.ReplaceAll(deleteDuTaskPath, "{feature}", "dir-usage")
	deleteDuTaskPath = strings.ReplaceAll(deleteDuTaskPath, "{task_id}", d.Id())

	deleteDuTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	_, err = client.Request("DELETE", deleteDuTaskPath, &deleteDuTaskOpt)
	if err != nil {
		if hasSpecifyErrorCode403(err, "SFS.TURBO.9000") {
			err = golangsdk.ErrDefault404{}
		}
		return common.CheckDeletedDiag(d, err, "error deleting DU task")
	}

	return nil
}

func resourceDuTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<share_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("share_id", parts[0]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

// When the SFS Turbo does not exist, the response body example of the details interface is as follows:
// {"errCode":"SFS.TURBO.9000","errMsg":"no privileges to operate"}
func hasSpecifyErrorCode403(err error, specCode string) bool {
	if errCode, ok := err.(golangsdk.ErrDefault403); ok {
		var response interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &response); jsonErr == nil {
			errorCode, parseErr := jmespath.Search("errCode", response)
			if parseErr != nil {
				log.Printf("[WARN] failed to parse errCode from response body: %s", parseErr)
			}

			if errorCode == specCode {
				return true
			}
		}
	}

	return false
}
