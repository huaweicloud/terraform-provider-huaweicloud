// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RDS
// ---------------------------------------------------------------

package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var sqlServerDatabaseNonUpdatableParams = []string{"instance_id", "name"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/database
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/instances/{instance_id}/database/detail
// @API RDS DELETE /v3.1/{project_id}/instances/{instance_id}/database/{db_name}
func ResourceSQLServerDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLServerDatabaseCreate,
		ReadContext:   resourceSQLServerDatabaseRead,
		UpdateContext: resourceSQLServerDatabaseUpdate,
		DeleteContext: resourceSQLServerDatabaseDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(sqlServerDatabaseNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				Description: `Specifies the ID of the RDS SQLServer instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the database name.`,
			},
			"character_set": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the character set used by the database.`,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the database status.`,
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

func resourceSQLServerDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSQLServerDatabase: create RDS SQLServer database.
	var (
		createSQLServerDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database"
		createSQLServerDatabaseProduct = "rds"
	)
	createSQLServerDatabaseClient, err := cfg.NewServiceClient(createSQLServerDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createSQLServerDatabasePath := createSQLServerDatabaseClient.Endpoint + createSQLServerDatabaseHttpUrl
	createSQLServerDatabasePath = strings.ReplaceAll(createSQLServerDatabasePath, "{project_id}",
		createSQLServerDatabaseClient.ProjectID)
	createSQLServerDatabasePath = strings.ReplaceAll(createSQLServerDatabasePath, "{instance_id}", instanceId)

	createSQLServerDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createSQLServerDatabaseOpt.JSONBody = utils.RemoveNil(buildCreateSQLServerDatabaseBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		_, err = createSQLServerDatabaseClient.Request("POST", createSQLServerDatabasePath, &createSQLServerDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(createSQLServerDatabaseClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS SQLServer database: %s", err)
	}

	dbName := d.Get("name").(string)
	d.SetId(instanceId + "/" + dbName)

	return resourceSQLServerDatabaseRead(ctx, d, meta)
}

func buildCreateSQLServerDatabaseBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceSQLServerDatabaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSQLServerDatabase: query RDS SQLServer database
	var (
		getSQLServerDatabaseHttpUrl = "v3/{project_id}/instances/{instance_id}/database/detail?page=1&limit=100"
		getSQLServerDatabaseProduct = "rds"
	)
	getSQLServerDatabaseClient, err := cfg.NewServiceClient(getSQLServerDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database name from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	getSQLServerDatabasePath := getSQLServerDatabaseClient.Endpoint + getSQLServerDatabaseHttpUrl
	getSQLServerDatabasePath = strings.ReplaceAll(getSQLServerDatabasePath, "{project_id}",
		getSQLServerDatabaseClient.ProjectID)
	getSQLServerDatabasePath = strings.ReplaceAll(getSQLServerDatabasePath, "{instance_id}", instanceId)

	getSQLServerDatabaseResp, err := pagination.ListAllItems(
		getSQLServerDatabaseClient,
		"page",
		getSQLServerDatabasePath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS SQLServer database")
	}

	getSQLServerDatabaseRespJson, err := json.Marshal(getSQLServerDatabaseResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getSQLServerDatabaseRespBody interface{}
	err = json.Unmarshal(getSQLServerDatabaseRespJson, &getSQLServerDatabaseRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	database := utils.PathSearch(fmt.Sprintf("databases[?name=='%s']|[0]", dbName), getSQLServerDatabaseRespBody, nil)
	if database == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", database, nil)),
		d.Set("character_set", utils.PathSearch("character_set", database, nil)),
		d.Set("state", utils.PathSearch("state", database, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSQLServerDatabaseUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSQLServerDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSQLServerDatabase: delete RDS SQLServer database
	var (
		deleteSQLServerDatabaseHttpUrl = "v3.1/{project_id}/instances/{instance_id}/database/{db_name}"
		deleteSQLServerDatabaseProduct = "rds"
	)
	deleteSQLServerDatabaseClient, err := cfg.NewServiceClient(deleteSQLServerDatabaseProduct, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	deleteSQLServerDatabasePath := deleteSQLServerDatabaseClient.Endpoint + deleteSQLServerDatabaseHttpUrl
	deleteSQLServerDatabasePath = strings.ReplaceAll(deleteSQLServerDatabasePath, "{project_id}",
		deleteSQLServerDatabaseClient.ProjectID)
	deleteSQLServerDatabasePath = strings.ReplaceAll(deleteSQLServerDatabasePath, "{instance_id}", instanceId)
	deleteSQLServerDatabasePath = strings.ReplaceAll(deleteSQLServerDatabasePath, "{db_name}",
		d.Get("name").(string))

	deleteSQLServerDatabaseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := deleteSQLServerDatabaseClient.Request("DELETE", deleteSQLServerDatabasePath, &deleteSQLServerDatabaseOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	deleteResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(deleteSQLServerDatabaseClient, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting RDS SQLServer database: %s", err)
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}
	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error deleting RDS SQL server database: job_id is not found in API response")
	}

	err = checkRDSInstanceJobFinish(deleteSQLServerDatabaseClient, jobId.(string), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
