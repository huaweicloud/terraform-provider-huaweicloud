package sfsturbo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/{share_id}/hpc-cache/task
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/{share_id}/hpc-cache/task/{task_id}
func ResourceDataTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataTaskCreate,
		ReadContext:   resourceDataTaskRead,
		DeleteContext: resourceDataTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDataTaskImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
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
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"src_target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dest_target": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"src_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"dest_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
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

func buildCreateDataTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":        d.Get("type"),
		"src_target":  d.Get("src_target"),
		"dest_target": d.Get("dest_target"),
		"src_prefix":  utils.ValueIgnoreEmpty(d.Get("src_prefix")),
		"dest_prefix": utils.ValueIgnoreEmpty(d.Get("dest_prefix")),
	}
	return bodyParams
}

func resourceDataTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	createDataTaskHttpUrl := "sfs-turbo/{share_id}/hpc-cache/task"
	createDataTaskPath := client.ResourceBaseURL() + createDataTaskHttpUrl
	createDataTaskPath = strings.ReplaceAll(createDataTaskPath, "{share_id}", d.Get("share_id").(string))

	createDataTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createDataTaskOpt.JSONBody = utils.RemoveNil(buildCreateDataTaskBodyParams(d))
	createDataTaskResp, err := client.Request("POST", createDataTaskPath, &createDataTaskOpt)
	if err != nil {
		return diag.Errorf("error creating data task: %s", err)
	}

	createDataTaskRespBody, err := utils.FlattenResponse(createDataTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId, err := jmespath.Search("task_id", createDataTaskRespBody)
	if err != nil || taskId == nil {
		return diag.Errorf("error creating data task: ID is not found in API response")
	}

	d.SetId(taskId.(string))

	err = dataTaskWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the creation of data task (%s) to complete: %s", d.Id(), err)
	}

	return resourceDataTaskRead(ctx, d, meta)
}

func getDataTaskInfo(d *schema.ResourceData, meta interface{}) (*http.Response, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS v1 client: %s", err)
	}

	getDataTaskHttpUrl := "sfs-turbo/{share_id}/hpc-cache/task/{task_id}"
	getDataTaskPath := client.ResourceBaseURL() + getDataTaskHttpUrl
	getDataTaskPath = strings.ReplaceAll(getDataTaskPath, "{share_id}", d.Get("share_id").(string))
	getDataTaskPath = strings.ReplaceAll(getDataTaskPath, "{task_id}", d.Id())
	getDataTaskOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getDataTaskPath, &getDataTaskOpts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func resourceDataTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	getDataTaskResp, err := getDataTaskInfo(d, meta)
	if err != nil {
		// When the data task does not exist, the response body example of the details interface is as follows:
		// error message: {"errCode":"SFS.TURBO.0106","errMsg":"Invalid task id."}
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errCode", "SFS.TURBO.0106"),
			fmt.Sprintf("error retrieving data task, the error message: %s", err))
	}

	getDataTaskRespBody, err := utils.FlattenResponse(getDataTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	beginTime := utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("start_time", getDataTaskRespBody, "").(string), "2006-01-02T15:04:05")
	endTime := utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("end_time", getDataTaskRespBody, "").(string), "2006-01-02T15:04:05")

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("type", getDataTaskRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getDataTaskRespBody, nil)),
		d.Set("src_target", utils.PathSearch("src_target", getDataTaskRespBody, nil)),
		d.Set("dest_target", utils.PathSearch("dest_target", getDataTaskRespBody, nil)),
		d.Set("src_prefix", utils.PathSearch("src_prefix", getDataTaskRespBody, nil)),
		d.Set("dest_prefix", utils.PathSearch("dest_prefix", getDataTaskRespBody, nil)),
		d.Set("start_time", utils.FormatTimeStampRFC3339(beginTime/1000, false)),
		d.Set("end_time", utils.FormatTimeStampRFC3339(endTime/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// Destroying resources does not change the current state.
func resourceDataTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func dataTaskWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      dataTaskStatusRefreshFunc(d, meta),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func dataTaskStatusRefreshFunc(d *schema.ResourceData, meta interface{}) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getDataTaskInfo(d, meta)
		if err != nil {
			return nil, "ERROR", err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, "ERROR", err
		}

		state, err := jmespath.Search("status", respBody)
		if err != nil {
			return nil, "ERROR", fmt.Errorf("failed to obtain the status information, error message: %s", err)
		}

		statusRaw := fmt.Sprintf("%v", state)
		// Whether the `status` value is SUCCESS or FAIL, it indicates that the resource creation was successful.
		// The result can be obtained by querying the details API for detailed information.
		if utils.StrSliceContains([]string{"SUCCESS", "FAIL"}, statusRaw) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceDataTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
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
