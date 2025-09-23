package ga

import (
	"context"
	"fmt"
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

// @API GA POST /v1/logtanks
// @API GA GET /v1/logtanks/{logtank_id}
// @API GA PUT /v1/logtanks/{logtank_id}
// @API GA DELETE /v1/logtanks/{logtank_id}
func ResourceAccessLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessLogCreate,
		UpdateContext: resourceAccessLogUpdate,
		ReadContext:   resourceAccessLogRead,
		DeleteContext: resourceAccessLogDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the resource to which the access log belongs.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the resource to which the access log belongs.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the log group to which the access log belongs.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the log stream to which the access log belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the access log.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the access log.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the access log.`,
			},
		},
	}
}

func resourceAccessLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/logtanks"
	)

	client, err := cfg.NewServiceClient("ga", region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateAccessLogBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating access log: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	logId := utils.PathSearch("logtank.id", respBody, "").(string)
	if logId == "" {
		return diag.Errorf("error creating access log: ID is not found in API response")
	}

	d.SetId(logId)

	err = waitingForAccessLogStateCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), d.Id())
	if err != nil {
		return diag.Errorf("error waiting for the access log (%s) creation to complete: %s", d.Id(), err)
	}

	return resourceAccessLogRead(ctx, d, meta)
}

func buildCreateAccessLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"logtank": map[string]interface{}{
			"resource_type": d.Get("resource_type"),
			"resource_id":   d.Get("resource_id"),
			"log_group_id":  d.Get("log_group_id"),
			"log_stream_id": d.Get("log_stream_id"),
		},
	}

	return bodyParams
}

func GetAccessLogInfo(client *golangsdk.ServiceClient, logId string) (interface{}, error) {
	httpUrl := "v1/logtanks/{logtank_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{logtank_id}", logId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return utils.FlattenResponse(resp)
}

func resourceAccessLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("ga", region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	respBody, err := GetAccessLogInfo(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving access log")
	}

	mErr := multierror.Append(
		d.Set("resource_type", utils.PathSearch("logtank.resource_type", respBody, nil)),
		d.Set("resource_id", utils.PathSearch("logtank.resource_id", respBody, nil)),
		d.Set("log_group_id", utils.PathSearch("logtank.log_group_id", respBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("logtank.log_stream_id", respBody, nil)),
		d.Set("status", utils.PathSearch("logtank.status", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("logtank.created_at", respBody, "").(string))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
			utils.PathSearch("logtank.updated_at", respBody, "").(string))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAccessLogUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/logtanks/{logtank_id}"
	)

	client, err := cfg.NewServiceClient("ga", region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	if d.HasChanges("log_group_id", "log_stream_id") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{logtank_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         buildUpdateAccessLogBodyParams(d),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating access log: %s", err)
		}

		err = waitingForAccessLogStateCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), d.Id())
		if err != nil {
			return diag.Errorf("error waiting for the access log (%s) update to complete: %s", d.Id(), err)
		}
	}

	return resourceAccessLogRead(ctx, d, meta)
}

func buildUpdateAccessLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"logtank": map[string]interface{}{
			"log_group_id":  d.Get("log_group_id"),
			"log_stream_id": d.Get("log_stream_id"),
		},
	}

	return bodyParams
}

func resourceAccessLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/logtanks/{logtank_id}"
	)

	client, err := cfg.NewServiceClient("ga", region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{logtank_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting access log")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      waitAccessLogStatusRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	return diag.FromErr(err)
}

func waitingForAccessLogStateCompleted(ctx context.Context, client *golangsdk.ServiceClient, t time.Duration, logId string) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      waitAccessLogStatusRefreshFunc(client, logId),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitAccessLogStatusRefreshFunc(client *golangsdk.ServiceClient, logId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetAccessLogInfo(client, logId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "Resource Not Found", "DELETED", nil
			}

			return nil, "ERROR", err
		}

		status := utils.PathSearch("logtank.status", respBody, "").(string)

		if utils.StrSliceContains([]string{"ERROR"}, status) {
			return respBody, "ERROR", fmt.Errorf("unexpected status: '%s'", status)
		}

		if utils.StrSliceContains([]string{"ACTIVE"}, status) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}
