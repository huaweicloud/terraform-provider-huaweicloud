package rds

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var pgTableRestoreNonUpdatableParams = []string{"instance_id", "restore_time", "databases"}

// @API RDS POST /v3/{project_id}/instances/batch/restore/tables
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
func ResourceRdsPgTableRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePgTableRestoreCreate,
		UpdateContext: resourcePgTableRestoreUpdate,
		ReadContext:   resourcePgTableRestoreRead,
		DeleteContext: resourcePgTableRestoreDelete,

		CustomizeDiff: config.FlexibleForceNew(pgTableRestoreNonUpdatableParams),

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
				Type:     schema.TypeString,
				Required: true,
			},
			"restore_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"databases": {
				Type:     schema.TypeSet,
				Elem:     pgTableRestoreDatabaseSchema(),
				Required: true,
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

func pgTableRestoreDatabaseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schemas": {
				Type:     schema.TypeSet,
				Elem:     pgTableRestoreSchemasSchema(),
				Required: true,
			},
		},
	}
}

func pgTableRestoreSchemasSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schema": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tables": {
				Type:     schema.TypeSet,
				Elem:     pgTableRestoreTableSchema(),
				Required: true,
			},
		},
	}
}

func pgTableRestoreTableSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"old_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"new_name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourcePgTableRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/batch/restore/tables"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildPgTablesRestoreBodyParams(d)),
	}

	instanceID := d.Get("instance_id").(string)
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS PostgreSQL table restore: %s", err)
	}

	res, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("restore_result[0].job_id", res, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating RDS PostgreSQL table restore: job_id not found in response")
	}

	d.SetId(jobId)

	if err = checkRDSInstanceJobFinish(client, jobId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildPgTablesRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instances": []map[string]interface{}{
			{
				"restore_time": d.Get("restore_time"),
				"instance_id":  d.Get("instance_id"),
				"databases":    buildPgRestoreDatabases(d.Get("databases").(*schema.Set).List()),
			},
		},
	}
	return bodyParams
}

func buildPgRestoreDatabases(databasesRaw []interface{}) []map[string]interface{} {
	if len(databasesRaw) == 0 {
		return nil
	}

	databases := make([]map[string]interface{}, len(databasesRaw))
	for i, db := range databasesRaw {
		dbMap := db.(map[string]interface{})
		databases[i] = map[string]interface{}{
			"database": dbMap["database"],
			"schemas":  buildPgRestoreSchemas(dbMap["schemas"].(*schema.Set).List()),
		}
	}
	return databases
}

func buildPgRestoreSchemas(schemasRaw []interface{}) []map[string]interface{} {
	if len(schemasRaw) == 0 {
		return nil
	}

	schemas := make([]map[string]interface{}, len(schemasRaw))
	for i, sc := range schemasRaw {
		scMap := sc.(map[string]interface{})
		schemas[i] = map[string]interface{}{
			"schema": scMap["schema"],
			"tables": buildPgRestoreTables(scMap["tables"].(*schema.Set).List()),
		}
	}
	return schemas
}

func buildPgRestoreTables(tablesRaw []interface{}) []map[string]interface{} {
	if len(tablesRaw) == 0 {
		return nil
	}

	tables := make([]map[string]interface{}, len(tablesRaw))
	for i, t := range tablesRaw {
		tMap := t.(map[string]interface{})
		tables[i] = map[string]interface{}{
			"old_name": tMap["old_name"],
			"new_name": tMap["new_name"],
		}
	}
	return tables
}

func resourcePgTableRestoreUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePgTableRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePgTableRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting RDS PostgreSQL table restore resource is not supported. The resource is only removed from " +
		"the state."
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
