package er

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

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/flow-logs
// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/flow-logs/{flow_log_id}
// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/flow-logs/{flow_log_id}/enable
// @API ER POST /v3/{project_id}/enterprise-router/{er_id}/flow-logs/{flow_log_id}/disable
// @API ER PUT /v3/{project_id}/enterprise-router/{er_id}/flow-logs/{flow_log_id}
// @API ER DELETE /v3/{project_id}/enterprise-router/{er_id}/flow-logs/{flow_log_id}
func ResourceFlowLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlowLogCreate,
		UpdateContext: resourceFlowLogUpdate,
		ReadContext:   resourceFlowLogRead,
		DeleteContext: resourceFlowLogDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceFlowLogImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_store_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"state": {
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
		},
	}
}

func buildCreateFlowLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"flow_log": map[string]interface{}{
			"log_store_type": d.Get("log_store_type"),
			"log_group_id":   d.Get("log_group_id"),
			"log_stream_id":  d.Get("log_stream_id"),
			"resource_type":  d.Get("resource_type"),
			"resource_id":    d.Get("resource_id"),
			"name":           d.Get("name"),
			"description":    utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func resourceFlowLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	createFlowLogHttpUrl := "enterprise-router/{er_id}/flow-logs"
	createFlowLogPath := client.ResourceBaseURL() + createFlowLogHttpUrl
	createFlowLogPath = strings.ReplaceAll(createFlowLogPath, "{er_id}", d.Get("instance_id").(string))

	createFlowLogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	createFlowLogOpt.JSONBody = utils.RemoveNil(buildCreateFlowLogBodyParams(d))
	createFlowLogResp, err := client.Request("POST", createFlowLogPath, &createFlowLogOpt)
	if err != nil {
		return diag.Errorf("error creating flow log: %s", err)
	}

	createInstanceRespBody, err := utils.FlattenResponse(createFlowLogResp)
	if err != nil {
		return diag.FromErr(err)
	}

	flowLogId := utils.PathSearch("flow_log.id", createInstanceRespBody, "").(string)
	if flowLogId == "" {
		return diag.Errorf("unable to find the ER flow log ID from the API response")
	}

	d.SetId(flowLogId)

	err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the create of flow log (%s) to complete: %s", d.Id(), err)
	}

	enabled := d.Get("enabled").(bool)
	if !enabled {
		err = updateFlowLogState(client, d, "disable")
		if err != nil {
			return diag.FromErr(err)
		}

		err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of flow log (%s) to complete: %s", d.Id(), err)
		}
	}

	return resourceFlowLogRead(ctx, d, meta)
}

func getFlowLogInfo(d *schema.ResourceData, meta interface{}) (*http.Response, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating ER v3 client: %s", err)
	}

	getFlowLogHttpUrl := "enterprise-router/{er_id}/flow-logs/{flow_log_id}"
	getFlowLogPath := client.ResourceBaseURL() + getFlowLogHttpUrl
	getFlowLogPath = strings.ReplaceAll(getFlowLogPath, "{er_id}", d.Get("instance_id").(string))
	getFlowLogPath = strings.ReplaceAll(getFlowLogPath, "{flow_log_id}", d.Id())
	getFlowLogOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	resp, err := client.Request("GET", getFlowLogPath, &getFlowLogOpts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func resourceFlowLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	getFlowLogResp, err := getFlowLogInfo(d, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ER flow log")
	}

	getFlowLogRespBody, err := utils.FlattenResponse(getFlowLogResp)
	if err != nil {
		return diag.FromErr(err)
	}

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("log_store_type", utils.PathSearch("flow_log.log_store_type", getFlowLogRespBody, nil)),
		d.Set("log_group_id", utils.PathSearch("flow_log.log_group_id", getFlowLogRespBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("flow_log.log_stream_id", getFlowLogRespBody, nil)),
		d.Set("name", utils.PathSearch("flow_log.name", getFlowLogRespBody, nil)),
		d.Set("description", utils.PathSearch("flow_log.description", getFlowLogRespBody, nil)),
		d.Set("resource_type", utils.PathSearch("flow_log.resource_type", getFlowLogRespBody, nil)),
		d.Set("resource_id", utils.PathSearch("flow_log.resource_id", getFlowLogRespBody, nil)),
		d.Set("state", utils.PathSearch("flow_log.state", getFlowLogRespBody, nil)),
		// The time results are not the time in RF3339 format without milliseconds.
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("flow_log.created_at",
			getFlowLogRespBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("flow_log.updated_at",
			getFlowLogRespBody, "").(string))/1000, false)),
		d.Set("enabled", utils.PathSearch("flow_log.enabled", getFlowLogRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateFlowLogState(client *golangsdk.ServiceClient, d *schema.ResourceData, operateType string) error {
	updateFlowLogStateHttpUrl := "enterprise-router/{er_id}/flow-logs/{flow_log_id}/" + operateType
	updateFlowLogStatePath := client.ResourceBaseURL() + updateFlowLogStateHttpUrl
	updateFlowLogStatePath = strings.ReplaceAll(updateFlowLogStatePath, "{er_id}", d.Get("instance_id").(string))
	updateFlowLogStatePath = strings.ReplaceAll(updateFlowLogStatePath, "{flow_log_id}", d.Id())

	updateFlowLogStateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err := client.Request("POST", updateFlowLogStatePath, &updateFlowLogStateOpts)
	if err != nil {
		return fmt.Errorf("error %s flow log: %s", operateType, err)
	}

	return nil
}

func buildUpdateFlowLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourceFlowLogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 Client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateFlowLogHttpUrl := "enterprise-router/{er_id}/flow-logs/{flow_log_id}"
		updateFlowLogPath := client.ResourceBaseURL() + updateFlowLogHttpUrl
		updateFlowLogPath = strings.ReplaceAll(updateFlowLogPath, "{er_id}", d.Get("instance_id").(string))
		updateFlowLogPath = strings.ReplaceAll(updateFlowLogPath, "{flow_log_id}", d.Id())

		updateFlowLogOpts := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateFlowLogOpts.JSONBody = utils.RemoveNil(buildUpdateFlowLogBodyParams(d))
		_, err = client.Request("PUT", updateFlowLogPath, &updateFlowLogOpts)
		if err != nil {
			return diag.Errorf("error updating ER flow log: %s", err)
		}

		err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of flow log (%s) to complete: %s", d.Id(), err)
		}
	}

	if d.HasChanges("enabled") {
		enabled := d.Get("enabled").(bool)
		if enabled {
			err = updateFlowLogState(client, d, "enable")
		} else {
			err = updateFlowLogState(client, d, "disable")
		}

		if err != nil {
			return diag.FromErr(err)
		}
		err = flowLogWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the update of flow log (%s) to complete: %s", d.Id(), err)
		}
	}

	return resourceFlowLogRead(ctx, d, meta)
}

func resourceFlowLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 Client: %s", err)
	}

	deleteFlowLogHttpUrl := "enterprise-router/{er_id}/flow-logs/{flow_log_id}"
	deleteFlowLogPath := client.ResourceBaseURL() + deleteFlowLogHttpUrl
	deleteFlowLogPath = strings.ReplaceAll(deleteFlowLogPath, "{er_id}", d.Get("instance_id").(string))
	deleteFlowLogPath = strings.ReplaceAll(deleteFlowLogPath, "{flow_log_id}", d.Id())

	deleteFlowLogOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			202,
		},
	}
	_, err = client.Request("DELETE", deleteFlowLogPath, &deleteFlowLogOpt)
	if err != nil {
		return diag.Errorf("error deleting flow log: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      flowLogStatusRefreshFunc(d, meta, true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flowLogStatusRefreshFunc(d *schema.ResourceData, meta interface{}, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getFlowLogInfo(d, meta)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
				return "Resource Not Found", "DELETED", nil
			}

			return nil, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("flow_log.state", respBody, "").(string)

		if utils.StrSliceContains([]string{"fail"}, status) {
			return respBody, "", fmt.Errorf("unexpected status: '%s'", status)
		}

		if utils.StrSliceContains([]string{"available"}, status) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func flowLogWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      flowLogStatusRefreshFunc(d, meta, false),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceFlowLogImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<instance_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
