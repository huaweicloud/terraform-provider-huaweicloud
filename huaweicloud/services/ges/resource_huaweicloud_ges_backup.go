// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GES
// ---------------------------------------------------------------

package ges

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

// @API GES GET /v2/{project_id}/graphs/{graph_id}/backups
// @API GES POST /v2/{project_id}/graphs/{graph_id}/backups
// @API GES DELETE /v2/{project_id}/graphs/{graph_id}/backups/{id}
// @API GES GET /v2/{project_id}/graphs/{graph_id}/jobs/{job_id}/status
func ResourceGesBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGesBackupCreate,
		ReadContext:   resourceGesBackupRead,
		DeleteContext: resourceGesBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGesBackupImportState,
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
			"graph_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The backup name.`,
			},
			"backup_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup method. The value can be **auto** or **manual**. `,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup status.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Start time of a backup job.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `End time of a backup job.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Backup file size (MB).`,
			},
			"duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Backup duration (seconds).`,
			},
			"encrypted": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the data is encrypted.`,
			},
		},
	}
}

func resourceGesBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createBackup: create a GES backup.
	var (
		createBackupHttpUrl = "v2/{project_id}/graphs/{graph_id}/backups"
		createBackupProduct = "ges"
	)
	createBackupClient, err := cfg.NewServiceClient(createBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	createBackupPath := createBackupClient.Endpoint + createBackupHttpUrl
	createBackupPath = strings.ReplaceAll(createBackupPath, "{project_id}", createBackupClient.ProjectID)
	createBackupPath = strings.ReplaceAll(createBackupPath, "{graph_id}", fmt.Sprintf("%v", d.Get("graph_id")))

	createBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	createBackupResp, err := createBackupClient.Request("POST", createBackupPath, &createBackupOpt)
	if err != nil {
		return diag.Errorf("error creating GesBackup: %s", err)
	}

	createBackupRespBody, err := utils.FlattenResponse(createBackupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	backupId := utils.PathSearch("backup_id", createBackupRespBody, "").(string)
	if backupId == "" {
		return diag.Errorf("unable to find the GES backup ID from the API response")
	}
	d.SetId(backupId)

	err = createBackupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate),
		utils.PathSearch("job_id", createBackupRespBody, "").(string))
	if err != nil {
		return diag.Errorf("error waiting for the Create of GesBackup (%s) to complete: %s", d.Id(), err)
	}
	return resourceGesBackupRead(ctx, d, meta)
}

func createBackupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{},
	t time.Duration, jobId string) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createBackupWaiting: waiting backup is available
			var (
				createBackupWaitingHttpUrl = "v2/{project_id}/graphs/{graph_id}/jobs/{job_id}/status"
				createBackupWaitingProduct = "ges"
			)
			createBackupWaitingClient, err := cfg.NewServiceClient(createBackupWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating GES Client: %s", err)
			}

			createBackupWaitingPath := createBackupWaitingClient.Endpoint + createBackupWaitingHttpUrl
			createBackupWaitingPath = strings.ReplaceAll(createBackupWaitingPath, "{project_id}", createBackupWaitingClient.ProjectID)
			createBackupWaitingPath = strings.ReplaceAll(createBackupWaitingPath, "{graph_id}", fmt.Sprintf("%v", d.Get("graph_id")))
			createBackupWaitingPath = strings.ReplaceAll(createBackupWaitingPath, "{job_id}", fmt.Sprintf("%v", jobId))

			createBackupWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
			}

			createBackupWaitingResp, err := createBackupWaitingClient.Request("GET", createBackupWaitingPath, &createBackupWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createBackupWaitingRespBody, err := utils.FlattenResponse(createBackupWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`status`, createBackupWaitingRespBody, "").(string)

			targetStatus := []string{
				"success",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createBackupWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"failed",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createBackupWaitingRespBody, status, nil
			}

			return createBackupWaitingRespBody, "PENDING", nil
		},
		Timeout:                   t,
		Delay:                     30 * time.Second,
		PollInterval:              30 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceGesBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getBackup: Query the GES backup.
	var (
		getBackupHttpUrl = "v2/{project_id}/graphs/{graph_id}/backups"
		getBackupProduct = "ges"
	)
	getBackupClient, err := cfg.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.ReplaceAll(getBackupPath, "{project_id}", getBackupClient.ProjectID)
	getBackupPath = strings.ReplaceAll(getBackupPath, "{graph_id}", fmt.Sprintf("%v", d.Get("graph_id")))

	getBackupPath += "?limit=120"

	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GesBackup")
	}

	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("backup_list[?id == '%s']|[0]", d.Id())
	getBackupRespBody = utils.PathSearch(jsonPath, getBackupRespBody, nil)
	if getBackupRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getBackupRespBody, nil)),
		d.Set("backup_method", utils.PathSearch("backup_method", getBackupRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getBackupRespBody, nil)),
		d.Set("start_time", utils.PathSearch("start_time", getBackupRespBody, nil)),
		d.Set("end_time", utils.PathSearch("end_time", getBackupRespBody, nil)),
		d.Set("size", utils.PathSearch("size", getBackupRespBody, nil)),
		d.Set("duration", utils.PathSearch("duration", getBackupRespBody, nil)),
		d.Set("encrypted", utils.PathSearch("encrypted", getBackupRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGesBackupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteBackup: delete GES backup
	var (
		deleteBackupHttpUrl = "v2/{project_id}/graphs/{graph_id}/backups/{id}"
		deleteBackupProduct = "ges"
	)
	deleteBackupClient, err := cfg.NewServiceClient(deleteBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	deleteBackupPath := deleteBackupClient.Endpoint + deleteBackupHttpUrl
	deleteBackupPath = strings.ReplaceAll(deleteBackupPath, "{project_id}", deleteBackupClient.ProjectID)
	deleteBackupPath = strings.ReplaceAll(deleteBackupPath, "{graph_id}", fmt.Sprintf("%v", d.Get("graph_id")))
	deleteBackupPath = strings.ReplaceAll(deleteBackupPath, "{id}", d.Id())

	deleteBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = deleteBackupClient.Request("DELETE", deleteBackupPath, &deleteBackupOpt)
	if err != nil {
		return diag.Errorf("error deleting GesBackup: %s", err)
	}

	return nil
}

func resourceGesBackupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <graph_id>/<id>")
	}

	d.Set("graph_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
