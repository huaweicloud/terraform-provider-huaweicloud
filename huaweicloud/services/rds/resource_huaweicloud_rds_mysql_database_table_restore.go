package rds

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS POST /v3.1/{project_id}/instances/{instance_id}/restore/tables
// @API RDS POST /v3/{project_id}/instances/batch/restore/databases
// @API RDS GET /v3/{project_id}/instances
func ResourceMysqlDatabaseTableRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlDatabaseTableRestoreCreate,
		ReadContext:   resourceMysqlDatabaseTableRestoreRead,
		DeleteContext: resourceMysqlDatabaseTableRestoreDelete,

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
				Description: `Specifies the ID of RDS MySQL instance.`,
			},
			"restore_time": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the restoration time point.`,
			},
			"is_fast_restore": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to use fast restoration.`,
			},
			"databases": {
				Type:         schema.TypeList,
				Elem:         mysqlDatabasesRestoreSchema(),
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"restore_tables"},
				Description:  `Specifies the databases that will be restored.`,
			},
			"restore_tables": {
				Type:         schema.TypeList,
				Elem:         mysqlTablesRestoreSchema(),
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"databases"},
				Description:  ` Specifies the tables that will be restored.`,
			},
		},
	}
}

func mysqlDatabasesRestoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"old_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the database before restoration.`,
			},
			"new_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the database after restoration.`,
			},
		},
	}
	return &sc
}

func mysqlTablesRestoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"database": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the database name.`,
			},
			"tables": {
				Type:        schema.TypeList,
				Elem:        mysqlTableNameRestoreSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the tables.`,
			},
		},
	}
	return &sc
}

func mysqlTableNameRestoreSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"old_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the table before restoration.`,
			},
			"new_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the table after restoration.`,
			},
		},
	}
	return &sc
}

func resourceMysqlDatabaseTableRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if _, ok := d.GetOk("databases"); ok {
		bodyParams := buildCreateDatabasesRestoreBodyParams(d)
		err = databaseTableRestore(ctx, d, client, "databases", bodyParams)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		bodyParams := buildCreateTablesRestoreBodyParams(d)
		err = databaseTableRestore(ctx, d, client, "tables", bodyParams)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func databaseTableRestore(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient, restoreType string,
	bodyParams map[string]interface{}) error {
	var (
		databaseHttpUrl = "v3/{project_id}/instances/batch/restore/databases"
		tablesHttpUrl   = "v3.1/{project_id}/instances/{instance_id}/restore/tables"
	)
	var httpUrl string
	if restoreType == "databases" {
		httpUrl = databaseHttpUrl
	} else {
		httpUrl = tablesHttpUrl
	}
	instanceID := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(bodyParams)

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
		return fmt.Errorf("error creating RDS Mysql %s restore: %s", restoreType, err)
	}

	res, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return err
	}
	jobId := utils.PathSearch("job_id||restore_result[0].job_id", res, nil)
	if jobId == nil {
		return fmt.Errorf("error creating RDS Mysql %s restore: job_id is not found in API response", restoreType)
	}

	d.SetId(jobId.(string))

	if err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate)); err != nil {
		return fmt.Errorf("error creating RDS Mysql %s restore: %s", restoreType, err)
	}
	return nil
}

func buildCreateDatabasesRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"instances": []map[string]interface{}{
			{
				"restore_time":    d.Get("restore_time"),
				"instance_id":     d.Get("instance_id"),
				"databases":       buildCreateNamesRestoreBody(d.Get("databases")),
				"is_fast_restore": d.Get("is_fast_restore"),
			},
		},
	}
	return bodyParams
}

func buildCreateTablesRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"restore_time":    d.Get("restore_time"),
		"restore_tables":  buildCreateTablesRestoreBody(d.Get("restore_tables")),
		"is_fast_restore": d.Get("is_fast_restore"),
	}
	return bodyParams
}

func buildCreateTablesRestoreBody(tablesRaw interface{}) []map[string]interface{} {
	rawParams := tablesRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, len(rawParams))
	for i, v := range rawParams {
		raw := v.(map[string]interface{})
		params[i] = map[string]interface{}{
			"database": raw["database"],
			"tables":   buildCreateNamesRestoreBody(raw["tables"]),
		}
	}

	return params
}

func buildCreateNamesRestoreBody(databasesRaw interface{}) []map[string]interface{} {
	rawParams := databasesRaw.([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, len(rawParams))
	for i, v := range rawParams {
		raw := v.(map[string]interface{})
		params[i] = map[string]interface{}{
			"old_name": raw["old_name"],
			"new_name": raw["new_name"],
		}
	}

	return params
}

func resourceMysqlDatabaseTableRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMysqlDatabaseTableRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restoration record is not supported. The restoration record is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before restoration."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
