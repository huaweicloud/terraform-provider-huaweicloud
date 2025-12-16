package bms

import (
	"context"
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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var volumeAttachNonUpdatableParams = []string{"server_id", "volume_id", "device"}

// @API BMS POST /v1/{project_id}/baremetalservers/{server_id}/attachvolume
// @API BMS GET /v1/{project_id}/baremetalservers/{server_id}/os-volume_attachments
// @API BMS DELETE /v1/{project_id}/baremetalservers/{server_id}/detachvolume/{attachment_id}
func ResourceVolumeAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVolumeAttachCreate,
		ReadContext:   resourceVolumeAttachRead,
		UpdateContext: resourceVolumeAttachUpdate,
		DeleteContext: resourceVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceVolumeAttachImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(volumeAttachNonUpdatableParams),

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
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"device": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceVolumeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/attachvolume"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", d.Get("server_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateVolumeAttachBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating BMS volume attach: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating BMS volume attach: job_id is not found in API response")
	}
	d.SetId(fmt.Sprintf("%s/%s", d.Get("server_id").(string), d.Get("volume_id").(string)))

	err = waitForJobComplete(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceVolumeAttachRead(ctx, d, meta)
}

func buildCreateVolumeAttachBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"volumeId": d.Get("volume_id"),
		"device":   utils.ValueIgnoreEmpty(d.Get("device")),
	}
	return map[string]interface{}{
		"volumeAttachment": bodyParams,
	}
}

func resourceVolumeAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/os-volume_attachments"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{server_id}", d.Get("server_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "BMS.1007"),
			"error retrieving BMS volume attach")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	searchPath := fmt.Sprintf("volumeAttachments[?volumeId=='%s']|[0]", d.Get("volume_id").(string))
	volumeAttach := utils.PathSearch(searchPath, getRespBody, nil)
	if volumeAttach == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("server_id", utils.PathSearch("serverId", volumeAttach, nil)),
		d.Set("volume_id", utils.PathSearch("volumeId", volumeAttach, nil)),
		d.Set("device", utils.PathSearch("device", volumeAttach, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVolumeAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceVolumeAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/detachvolume/{attachment_id}"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{server_id}", d.Get("server_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{attachment_id}", d.Get("volume_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error.code", "BMS.0114", "BMS.1006"),
			"error deleting BMS volume attach")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting BMS volume attach: job_id is not found in API response")
	}

	err = waitForJobComplete(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForJobComplete(ctx context.Context, client *golangsdk.ServiceClient, jobId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      bmsJobRefreshFunc(client, jobId),
		Timeout:      timeout,
		Delay:        2 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for BMS job(%s) to complete: %s ", jobId, err)
	}
	return nil
}

func bmsJobRefreshFunc(client *golangsdk.ServiceClient, jobId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl = "v1/{project_id}/jobs/{job_id}"
		)

		getPath := client.Endpoint + httpUrl
		getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
		getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)

		getOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return nil, "FAIL", err
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return nil, "FAIL", err
		}

		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "" {
			return nil, "FAIL", errors.New("status is not found")
		}

		if status == "FAIL" {
			return getRespBody, "FAIL", fmt.Errorf("the connection status is: %s", status)
		}
		if status == "SUCCESS" {
			return getRespBody, "SUCCESS", nil
		}

		return getRespBody, "PENDING", nil
	}
}

func resourceVolumeAttachImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <server_id>/<volume_id>")
	}

	mErr := multierror.Append(nil,
		d.Set("server_id", parts[0]),
		d.Set("volume_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
