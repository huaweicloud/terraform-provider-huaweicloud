package taurusdb

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

// @API GaussDBforMySQL POST /v3/{project_id}/backups/create
// @API GaussDBforMySQL GET /v3/{project_id}/instances
// @API GaussDBforMySQL GET /v3/{project_id}/backups
// @API GaussDBforMySQL DELETE /v3/{project_id}/backups/{backup_id}
func ResourceGaussDBMysqlBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlBackupCreate,
		ReadContext:   resourceGaussDBMysqlBackupRead,
		DeleteContext: resourceGaussDBMysqlBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the GaussDB MySQL instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the backup.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the description of the backup.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup start time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the backup end time in the "yyyy-mm-ddThh:mm:ssZ" format.`,
			},
			"take_up_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the backup duration in minutes.`,
			},
			"size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the backup size in MB.`,
			},
			"datastore": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the database information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the database engine.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the database version.`,
						},
					},
				},
			},
		},
	}
}

func resourceGaussDBMysqlBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/backups/create"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBBackupBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL backup: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	backupID := utils.PathSearch("backup.id", createRespBody, nil)
	if backupID == nil {
		return diag.Errorf("error creating GaussDB MySQL backup: backup ID is not found in API response")
	}
	d.SetId(backupID.(string))

	err = waitForGaussDBBackupComplete(ctx, d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGaussDBMysqlBackupRead(ctx, d, meta)
}

func waitForGaussDBBackupComplete(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"BUILDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      gaussDBBackupStateRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for GaussDB MySQL backup (%s) to build complete: %s", d.Id(), err)
	}
	return nil
}

func gaussDBBackupStateRefreshFunc(client *golangsdk.ServiceClient, backupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		backup, err := getGaussDBBackup(client, backupID)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", backup, nil)
		if status == nil {
			return backup, "DELETED", fmt.Errorf("the GaussDB MySQL backup(%s) has been deleted", backupID)
		}
		return backup, status.(string), nil
	}
}

func buildCreateGaussDBBackupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instance_id": d.Get("instance_id"),
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceGaussDBMysqlBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	backup, err := getGaussDBBackup(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB MySQL backup")
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", backup, nil)),
		d.Set("name", utils.PathSearch("name", backup, nil)),
		d.Set("description", utils.PathSearch("description", backup, nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", backup, nil)),
		d.Set("end_time", utils.PathSearch("end_time", backup, nil)),
		d.Set("take_up_time", utils.PathSearch("take_up_time", backup, nil)),
		d.Set("size", utils.PathSearch("size", backup, nil)),
		d.Set("datastore", flattenGetGaussDBBackupResponseBodyDatastore(backup)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGaussDBBackupResponseBodyDatastore(backup interface{}) []interface{} {
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"type":    utils.PathSearch("datastore.type", backup, nil),
		"version": utils.PathSearch("datastore.version", backup, nil),
	})
	return rst
}

func getGaussDBBackup(client *golangsdk.ServiceClient, backupID string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/backups?backup_id={backup_id}"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{backup_id}", backupID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL backup: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}
	backup := utils.PathSearch("backups|[0]", getRespBody, nil)
	if backup == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return backup, nil
}

func resourceGaussDBMysqlBackupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/backups/{backup_id}"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{backup_id}", d.Id())

	deleteGaussDBDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteGaussDBDatabaseOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBDatabaseBodyParams(d))

	deleteGResp, err := client.Request("DELETE", deletePath, &deleteGaussDBDatabaseOpt)
	if err != nil {
		return diag.Errorf("error deleting GaussDB MySQL backup: %s", err)
	}

	_, err = utils.FlattenResponse(deleteGResp)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
