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

var pgTableRestoreNonUpdatableParams = []string{
	"instance_id",
	"restore_time",
	"databases",
	"databases.*.database",
	"databases.*.schemas",
	"databases.*.schemas.*.schema",
	"databases.*.schemas.*.tables",
	"databases.*.schemas.*.tables.*.old_name",
	"databases.*.schemas.*.tables.*.new_name",
}

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
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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
				Type:     schema.TypeList,
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

	url := "v3/{project_id}/instances/batch/restore/tables"
	path := client.Endpoint + url
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
		JSONBody: utils.RemoveNil(buildPgTablesRestoreBodyParams(d)),
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", path, &opt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}

	instanceID := d.Get("instance_id").(string)

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
		return diag.Errorf("error creating RDS PostgreSQL tables restore: %s", err)
	}

	res, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("restore_result[0].job_id", res, nil)
	if jobId == nil {
		return diag.Errorf("error creating RDS PostgreSQL table restore: job_id not found in response")
	}

	d.SetId(jobId.(string))

	if err := checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error creating RDS PostgreSQL table restore: %s", err)
	}

	return nil
}

func buildPgTablesRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instances": []map[string]interface{}{
			{
				"restore_time": d.Get("restore_time"),
				"instance_id":  d.Get("instance_id"),
				"databases":    buildPgRestoreDatabases(d.Get("databases")),
			},
		},
	}
	return bodyParams
}

func buildPgRestoreDatabases(databasesRaw interface{}) []map[string]interface{} {
	rawDatabases := databasesRaw.([]interface{})

	if len(rawDatabases) == 0 {
		return nil
	}

	databases := make([]map[string]interface{}, len(rawDatabases))
	for i, db := range rawDatabases {
		dbMap := db.(map[string]interface{})
		databases[i] = map[string]interface{}{
			"database": dbMap["database"],
			"schemas":  buildPgRestoreSchemas(convertTypeListToTypeSet(dbMap["schemas"])),
		}
	}
	return databases
}

func convertTypeListToTypeSet(schemasRaw interface{}) interface{} {
	return schemasRaw
}

func buildPgRestoreSchemas(schemasRaw interface{}) []map[string]interface{} {
	rawSchemas := schemasRaw.([]interface{})

	if len(rawSchemas) == 0 {
		return nil
	}

	schemas := make([]map[string]interface{}, len(rawSchemas))
	for i, sc := range rawSchemas {
		scMap := sc.(map[string]interface{})
		schemas[i] = map[string]interface{}{
			"schema": scMap["schema"],
			"tables": buildPgRestoreTables(scMap["tables"]),
		}
	}
	return schemas
}

func buildPgRestoreTables(tablesRaw interface{}) []map[string]interface{} {
	rawTables := tablesRaw.([]interface{})

	if len(rawTables) == 0 {
		return nil
	}

	tables := make([]map[string]interface{}, len(rawTables))
	for i, t := range rawTables {
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
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  "Deleting restoration record is not supported. The record is removed from state only.",
		},
	}
}
