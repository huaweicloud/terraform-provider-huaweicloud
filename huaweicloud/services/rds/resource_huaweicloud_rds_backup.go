// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var backupNonUpdatableParams = []string{"name", "instance_id", "description", "databases",
	"databases.*.name"}

// @API RDS DELETE /v3/{project_id}/backups/{id}
// @API RDS GET /v3/{project_id}/backups
// @API RDS POST /v3/{project_id}/backups
func ResourceBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBackupCreate,
		ReadContext:   resourceBackupRead,
		UpdateContext: resourceBackupUpdate,
		DeleteContext: resourceBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: backupImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(backupNonUpdatableParams),

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Backup name.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Instance ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description about the backup.`,
			},
			"databases": {
				Type:        schema.TypeList,
				Elem:        BackupBackupDatabaseSchema(),
				Optional:    true,
				Computed:    true,
				Description: `List of self-built Microsoft SQL Server databases that are partially backed up.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup start time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup end time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Backup status.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Backup size in KB.`,
			},
			"associated_with_ddm": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether a DDM instance has been associated.`,
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

func BackupBackupDatabaseSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Database to be backed up for Microsoft SQL Server.`,
			},
		},
	}
	return &sc
}

func resourceBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// createBackup: create a RDS backup.
	var (
		createBackupHttpUrl = "v3/{project_id}/backups"
		createBackupProduct = "rds"
	)
	createBackupClient, err := config.NewServiceClient(createBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating Backup Client: %s", err)
	}

	createBackupPath := createBackupClient.Endpoint + createBackupHttpUrl
	createBackupPath = strings.Replace(createBackupPath, "{project_id}", createBackupClient.ProjectID, -1)

	createBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createBackupOpt.JSONBody = utils.RemoveNil(buildCreateBackupBodyParams(d, config))
	var createBackupResp *http.Response
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		createBackupResp, err = createBackupClient.Request("POST", createBackupPath, &createBackupOpt)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error creating RDS database backup: %s", err)
	}

	createBackupRespBody, err := utils.FlattenResponse(createBackupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	backupId := utils.PathSearch("backup.id", createBackupRespBody, "").(string)
	if backupId == "" {
		return diag.Errorf("unable to find the RDS backup ID from the API response")
	}
	d.SetId(backupId)

	err = createBackupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of Backup (%s) to complete: %s", d.Id(), err)
	}
	return resourceBackupRead(ctx, d, meta)
}

func buildCreateBackupBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"instance_id": utils.ValueIgnoreEmpty(d.Get("instance_id")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"databases":   buildCreateBackupDatabasesChildBody(d),
	}
	return bodyParams
}

func buildCreateBackupDatabasesChildBody(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("databases").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, len(rawParams))
	for i, v := range rawParams {
		raw := v.(map[string]interface{})
		params[i] = map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(raw["name"]),
		}
	}

	return params
}

func buildCreateBackupWaitingQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&backup_id=%v", res, d.Id())

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func createBackupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// createBackupWaiting: missing operation notes
			var (
				createBackupWaitingHttpUrl = "v3/{project_id}/backups"
				createBackupWaitingProduct = "rds"
			)
			createBackupWaitingClient, err := config.NewServiceClient(createBackupWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Backup Client: %s", err)
			}

			createBackupWaitingPath := createBackupWaitingClient.Endpoint + createBackupWaitingHttpUrl
			createBackupWaitingPath = strings.Replace(createBackupWaitingPath, "{project_id}", createBackupWaitingClient.ProjectID, -1)

			createBackupWaitingqueryParams := buildCreateBackupWaitingQueryParams(d)
			createBackupWaitingPath = createBackupWaitingPath + createBackupWaitingqueryParams

			createBackupWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			createBackupWaitingResp, err := createBackupWaitingClient.Request("GET", createBackupWaitingPath, &createBackupWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createBackupWaitingRespBody, err := utils.FlattenResponse(createBackupWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`backups[0].status`, createBackupWaitingRespBody, "").(string)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `backups[0].status`)
			}

			if utils.StrSliceContains(strings.Split(`FAILED`, ","), status) {
				return createBackupWaitingRespBody, status, nil
			}

			if utils.StrSliceContains(strings.Split(`COMPLETED`, ","), status) {
				return createBackupWaitingRespBody, "COMPLETED", nil
			}

			return createBackupWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getBackup: Query the RDS manual backup
	var (
		getBackupHttpUrl = "v3/{project_id}/backups"
		getBackupProduct = "rds"
	)
	getBackupClient, err := config.NewServiceClient(getBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating Backup Client: %s", err)
	}

	getBackupPath := getBackupClient.Endpoint + getBackupHttpUrl
	getBackupPath = strings.Replace(getBackupPath, "{project_id}", getBackupClient.ProjectID, -1)

	getBackupqueryParams := buildGetBackupQueryParams(d)
	getBackupPath = getBackupPath + getBackupqueryParams

	getBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getBackupResp, err := getBackupClient.Request("GET", getBackupPath, &getBackupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Backup")
	}

	getBackupRespBody, err := utils.FlattenResponse(getBackupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("backups[0].name", getBackupRespBody, nil)),
		d.Set("instance_id", utils.PathSearch("backups[0].instance_id", getBackupRespBody, nil)),
		d.Set("description", utils.PathSearch("backups[0].description", getBackupRespBody, nil)),
		d.Set("begin_time", utils.PathSearch("backups[0].begin_time", getBackupRespBody, nil)),
		d.Set("end_time", utils.PathSearch("backups[0].end_time", getBackupRespBody, nil)),
		d.Set("status", utils.PathSearch("backups[0].status", getBackupRespBody, nil)),
		d.Set("size", utils.PathSearch("backups[0].size", getBackupRespBody, nil)),
		d.Set("associated_with_ddm", utils.PathSearch("backups[0].associated_with_ddm", getBackupRespBody, nil)),
		d.Set("databases", flattenGetBackupResponseBodyBackupDatabase(getBackupRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetBackupResponseBodyBackupDatabase(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("backups[0].databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
		})
	}
	return rst
}

func buildGetBackupQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&backup_id=%v", res, d.Id())

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func resourceBackupUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// deleteBackup: missing operation notes
	var (
		deleteBackupHttpUrl = "v3/{project_id}/backups/{id}"
		deleteBackupProduct = "rds"
	)
	deleteBackupClient, err := config.NewServiceClient(deleteBackupProduct, region)
	if err != nil {
		return diag.Errorf("error creating Backup Client: %s", err)
	}

	deleteBackupPath := deleteBackupClient.Endpoint + deleteBackupHttpUrl
	deleteBackupPath = strings.Replace(deleteBackupPath, "{project_id}", deleteBackupClient.ProjectID, -1)
	deleteBackupPath = strings.Replace(deleteBackupPath, "{id}", d.Id(), -1)

	deleteBackupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = deleteBackupClient.Request("DELETE", deleteBackupPath, &deleteBackupOpt)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting RDS database backup: %s", err)
	}

	err = deleteBackupWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of Backup (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func buildDeleteBackupWaitingQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&backup_id=%v", res, d.Id())

	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func deleteBackupWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			// deleteBackupWaiting: missing operation notes
			var (
				deleteBackupWaitingHttpUrl = "v3/{project_id}/backups"
				deleteBackupWaitingProduct = "rds"
			)
			deleteBackupWaitingClient, err := config.NewServiceClient(deleteBackupWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating Backup Client: %s", err)
			}

			deleteBackupWaitingPath := deleteBackupWaitingClient.Endpoint + deleteBackupWaitingHttpUrl
			deleteBackupWaitingPath = strings.Replace(deleteBackupWaitingPath, "{project_id}", deleteBackupWaitingClient.ProjectID, -1)

			deleteBackupWaitingqueryParams := buildDeleteBackupWaitingQueryParams(d)
			deleteBackupWaitingPath = deleteBackupWaitingPath + deleteBackupWaitingqueryParams

			deleteBackupWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}
			deleteBackupWaitingResp, err := deleteBackupWaitingClient.Request("GET", deleteBackupWaitingPath, &deleteBackupWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteBackupWaitingRespBody, err := utils.FlattenResponse(deleteBackupWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`total_count`, deleteBackupWaitingRespBody, nil)
			status := fmt.Sprintf("%v", statusRaw)
			if utils.StrSliceContains(strings.Split(`1`, ","), status) {
				return deleteBackupWaitingRespBody, "PENDING", nil
			}

			if utils.StrSliceContains(strings.Split(`0`, ","), status) {
				return deleteBackupWaitingRespBody, "COMPLETED", nil
			}

			return deleteBackupWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func backupImportState(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <instance_id>/<backup_id>")
	}
	instanceId := parts[0]
	backup_id := parts[1]
	d.SetId(backup_id)
	d.Set("instance_id", instanceId)
	return []*schema.ResourceData{d}, nil
}
