// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DDS
// ---------------------------------------------------------------

package dds

import (
	"context"
	"encoding/json"
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

// @API DDS GET /v3/{project_id}/backups
// @API DDS POST /v3/{project_id}/backups
// @API DDS DELETE /v3/{project_id}/backups/{backup_id}
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/jobs
func ResourceDdsBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDdsBackupCreate,
		ReadContext:   resourceDdsBackupRead,
		DeleteContext: resourceDdsBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: backupImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a DDS instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the manual backup name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the manual backup description.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of a DDS instance.`,
			},
			"datastore": {
				Type:        schema.TypeList,
				Elem:        BackupDatastoreSchema(),
				Computed:    true,
				Description: `Indicates the database version.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup type.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the start time of the backup.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the end time of the backup.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup status.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the backup size in KB.`,
			},
		},
	}
}

func BackupDatastoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DB engine.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the database version.`,
			},
		},
	}
	return &sc
}

func resourceDdsBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createBackup: create DDS backup
	var (
		createBackupHttpUrl = "v3/{project_id}/backups"
		createBackupProduct = "dds"
	)
	createBackupClient, err := cfg.NewServiceClient(createBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	createBackupPath := createBackupClient.Endpoint + createBackupHttpUrl
	createBackupPath = strings.ReplaceAll(createBackupPath, "{project_id}", createBackupClient.ProjectID)

	createBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createBackupOpt.JSONBody = utils.RemoveNil(buildCreateBackupBodyParams(d))

	instanceId := d.Get("instance_id").(string)
	retryFunc := func() (interface{}, bool, error) {
		resp, err := createBackupClient.Request("POST", createBackupPath, &createBackupOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(createBackupClient, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating DDS backup: %s", err)
	}

	createBackupRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("backup_id", createBackupRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating DDS backup: backup_id is not found in API response")
	}

	d.SetId(id)

	jobId := utils.PathSearch("job_id", createBackupRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating DDS backup: job_id is not found in API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed"},
		Refresh:      ddsJobStatusRefreshFunc(jobId, region, cfg),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for job (%s) to complete: %s", jobId, err)
	}

	return resourceDdsBackupRead(ctx, d, meta)
}

func buildCreateBackupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_id": d.Get("instance_id"),
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	params := map[string]interface{}{
		"backup": bodyParams,
	}
	return params
}

func ddsJobStatusRefreshFunc(jobId, region string, cfg *config.Config) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getJobStatusHttpUrl = "v3/{project_id}/jobs"
			getJobStatusProduct = "dds"
		)
		getJobStatusClient, err := cfg.NewServiceClient(getJobStatusProduct, region)
		if err != nil {
			return nil, "", fmt.Errorf("error creating DDS client: %s", err)
		}

		getJobStatusPath := getJobStatusClient.Endpoint + getJobStatusHttpUrl
		getJobStatusPath = strings.ReplaceAll(getJobStatusPath, "{project_id}", getJobStatusClient.ProjectID)

		getJobStatusQueryParams := buildWaitJobQueryParams(jobId)
		getJobStatusPath += getJobStatusQueryParams

		getJobStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getJobStatusResp, err := getJobStatusClient.Request("GET", getJobStatusPath, &getJobStatusOpt)
		if err != nil {
			if errCode, ok := err.(golangsdk.ErrDefault400); ok {
				var apiError interface{}
				if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
					return nil, "", fmt.Errorf("error get DDS job: format error: %s", err)
				}
				errorCode := utils.PathSearch("error_code", apiError, "").(string)
				if errorCode == "" {
					return nil, "", fmt.Errorf("error parse errorCode from response body: error_code is not " +
						"found in the response")
				}
				// if the error_code is DBS.200543, it indicates that the deleting instance job has finished
				if errorCode == "DBS.200543" {
					return getJobStatusResp, "Deleted", nil
				}
			}
			return nil, "", err
		}

		getJobStatusRespBody, err := utils.FlattenResponse(getJobStatusResp)
		if err != nil {
			return nil, "", err
		}

		status := utils.PathSearch("job.status", getJobStatusRespBody, "").(string)
		if status == "" {
			return nil, "", fmt.Errorf("error get job status by job ID: %s", jobId)
		}
		return getJobStatusRespBody, status, nil
	}
}

func resourceDdsBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getBackup: Query DDS backup
	var (
		getBackupHttpUrl = "v3/{project_id}/backups"
		getBackupProduct = "dds"
	)
	getBackupClient, err := cfg.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.ReplaceAll(getBackupPath, "{project_id}", getBackupClient.ProjectID)

	getBackupQueryParams := buildGetBackupQueryParams(d.Id())
	getBackupPath += getBackupQueryParams

	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", []string{"DBS.201502", "DBS.201214"}...),
			"error retrieving DDS backup")
	}

	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return diag.FromErr(err)
	}
	backups := utils.PathSearch("backups", getBackupRespBody, make([]interface{}, 0)).([]interface{})
	if len(backups) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", backups[0], nil)),
		d.Set("instance_id", utils.PathSearch("instance_id", backups[0], nil)),
		d.Set("instance_name", utils.PathSearch("instance_name", backups[0], nil)),
		d.Set("datastore", flattenGetBackupResponseDatastore(backups[0])),
		d.Set("type", utils.PathSearch("type", backups[0], nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", backups[0], nil)),
		d.Set("end_time", utils.PathSearch("end_time", backups[0], nil)),
		d.Set("status", utils.PathSearch("status", backups[0], nil)),
		d.Set("size", utils.PathSearch("size", backups[0], 0)),
		d.Set("description", utils.PathSearch("description", backups[0], nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetBackupResponseDatastore(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("datastore", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"type":    utils.PathSearch("type", curJson, nil),
			"version": utils.PathSearch("version", curJson, nil),
		},
	}
	return rst
}

func buildGetBackupQueryParams(backupId string) string {
	return fmt.Sprintf("?backup_id=%s", backupId)
}

func buildWaitJobQueryParams(jobId string) string {
	return fmt.Sprintf("?id=%v", jobId)
}

func resourceDdsBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteBackup: Delete DDS backup
	var (
		deleteBackupHttpUrl = "v3/{project_id}/backups/{backup_id}"
		deleteBackupProduct = "dds"
	)
	deleteBackupClient, err := cfg.NewServiceClient(deleteBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	deleteBackupPath := deleteBackupClient.Endpoint + deleteBackupHttpUrl
	deleteBackupPath = strings.ReplaceAll(deleteBackupPath, "{project_id}", deleteBackupClient.ProjectID)
	deleteBackupPath = strings.ReplaceAll(deleteBackupPath, "{backup_id}", fmt.Sprintf("%v", d.Id()))

	deleteBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	instanceId := d.Get("instance_id").(string)
	retryFunc := func() (interface{}, bool, error) {
		resp, err := deleteBackupClient.Request("DELETE", deleteBackupPath, &deleteBackupOpt)
		retry, err := handleMultiOperationsError(err)
		return resp, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     ddsInstanceStateRefreshFunc(deleteBackupClient, instanceId),
		WaitTarget:   []string{"normal"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 5 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", []string{"DBS.201502", "DBS.201214"}...),
			"error deleting DDS backup")
	}

	deleteBackupRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteBackupRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting DDS backup: job_id is not found in API response")
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Running"},
		Target:       []string{"Completed", "Deleted"},
		Refresh:      ddsJobStatusRefreshFunc(jobId, region, cfg),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for job (%s) to be completed: %s", d.Id(), err)
	}

	return nil
}

func backupImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}
	instanceId := parts[0]
	backupId := parts[1]
	d.SetId(backupId)
	err := d.Set("instance_id", instanceId)
	if err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
