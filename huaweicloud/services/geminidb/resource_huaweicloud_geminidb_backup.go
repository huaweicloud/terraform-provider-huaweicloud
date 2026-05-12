package geminidb

import (
	"context"
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

var geminiDBBackupNonUpdatableParams = []string{
	"instance_id",
	"name",
	"description",
	"database_tables",
	"database_tables.*.database_name",
	"database_tables.*.table_names",
}

// @API GaussDBforNoSQL GET /v3/{project_id}/instances
// @API GaussDBforNoSQL GET /v3/{project_id}/jobs
// @API GaussDBforNoSQL GET /v4/{project_id}/backups
// @API GaussDBforNoSQL POST /v3/{project_id}/instances/{instance_id}/backups
// @API GaussDBforNoSQL DELETE /v3/{project_id}/backups/{backup_id}
func ResourceGeminiDBBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBBackupCreate,
		UpdateContext: resourceGeminiDBBackupUpdate,
		ReadContext:   resourceGeminiDBBackupRead,
		DeleteContext: resourceGeminiDBBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(geminiDBBackupNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database_tables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     geminiDBBackupDatabaseTablesSchema(),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBBackupDatastoreSchema(),
			},
		},
	}
}

func geminiDBBackupDatabaseTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_names": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func geminiDBBackupDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGeminiDBBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/backups"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGeminiDBBackupBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GeminiDB backup: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	backupID := utils.PathSearch("backup_id", createRespBody, "").(string)
	if backupID == "" {
		return diag.Errorf("error creating GeminiDB backup: backup_id is not found in API response")
	}
	d.SetId(backupID)

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating GeminiDB backup: job_id is not found in API response")
	}
	err = checkGeminiDbJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGeminiDBBackupRead(ctx, d, meta)
}

func resourceGeminiDBBackupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	backupID := d.Id()

	backup, err := GetBackup(client, backupID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB backup")
	}
	if backup == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GeminiDB backup")
	}
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", backup, nil)),
		d.Set("name", utils.PathSearch("name", backup, nil)),
		d.Set("description", utils.PathSearch("description", backup, nil)),
		d.Set("status", utils.PathSearch("status", backup, nil)),
		d.Set("type", utils.PathSearch("type", backup, nil)),
		d.Set("size", utils.PathSearch("size", backup, nil)),
		d.Set("begin_time", utils.PathSearch("begin_time", backup, nil)),
		d.Set("end_time", utils.PathSearch("end_time", backup, nil)),
		d.Set("datastore", flattenGeminiDBBackupDatastore(backup)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting GeminiDB backup fields: %s", err)
	}

	return nil
}

func GetBackup(client *golangsdk.ServiceClient, backupID string) (interface{}, error) {
	backup, err := getBackupInfo(client, backupID, "Instance")
	if err != nil {
		return nil, err
	}
	if backup == nil {
		backup, err = getBackupInfo(client, backupID, "DatabaseTable")
		if err != nil {
			return nil, err
		}
	}

	return backup, nil
}

func getBackupInfo(client *golangsdk.ServiceClient, backupID, backupType string) (interface{}, error) {
	httpUrl := "v4/{project_id}/backups"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetGeminiDBBackupQueryParams(backupID, backupType)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	backup := utils.PathSearch("backups[0]", getRespBody, nil)

	return backup, nil
}

func buildGetGeminiDBBackupQueryParams(backupID, backupType string) string {
	return fmt.Sprintf("?backup_id=%s&type=%s", backupID, backupType)
}

func resourceGeminiDBBackupUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGeminiDBBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/backups/{backup_id}"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{backup_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.201214"),
			"error retrieving Geminidb backup")
	}

	_, err = utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	err = waitingForBackupDeleteCompleted(ctx, d, d.Timeout(schema.TimeoutDelete), client)
	if err != nil {
		return diag.Errorf("error waiting for the Create of Backup (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func buildCreateGeminiDBBackupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	if databaseTables, ok := d.GetOk("database_tables"); ok {
		bodyParams["database_tables"] = buildGeminiDBBackupDatabaseTables(databaseTables.([]interface{}))
	}

	return bodyParams
}

func buildGeminiDBBackupDatabaseTables(databaseTables []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(databaseTables))
	for _, table := range databaseTables {
		tableMap := table.(map[string]interface{})
		item := map[string]interface{}{
			"database_name": tableMap["database_name"],
		}
		if tableNames, ok := tableMap["table_names"]; ok {
			item["table_names"] = tableNames
		}
		result = append(result, item)
	}
	return result
}

func flattenGeminiDBBackupDatastore(backup interface{}) []map[string]interface{} {
	datastoreRaw := utils.PathSearch("datastore", backup, nil)
	if datastoreRaw == nil {
		return nil
	}

	datastore := map[string]interface{}{
		"type":    utils.PathSearch("type", datastoreRaw, nil),
		"version": utils.PathSearch("version", datastoreRaw, nil),
	}

	return []map[string]interface{}{datastore}
}

func waitingForBackupDeleteCompleted(ctx context.Context, d *schema.ResourceData, t time.Duration, client *golangsdk.ServiceClient) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"BUILDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			backup, err := GetBackup(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			if backup == nil {
				return "", "COMPLETED", nil
			}
			return "", "PENDING", nil
		},
		Timeout:      t,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
