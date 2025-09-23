package taurusdb

import (
	"context"
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

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/backups/restore/tables
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
func ResourceGaussDBMysqlTableRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBMysqlTableRestoreCreate,
		ReadContext:   resourceGaussDBMysqlTableRestoreRead,
		DeleteContext: resourceGaussDBMysqlTableRestoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"restore_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"restore_tables": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     restoreTablesSchema(),
			},
			"last_table_info": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func restoreTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem:     restoreTablesTablesSchema(),
			},
		},
	}
}

func restoreTablesTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"old_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"new_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceGaussDBMysqlTableRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/backups/restore/tables"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreateTableRestoreBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	r, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     GaussDBInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error restoring tables to GaussDB MySQL instance (%s): %s", instanceId, err)
	}

	createRespBody, err := utils.FlattenResponse(r.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("unable to find the job_id from the response: %s", err)
	}

	d.SetId(jobId.(string))

	err = checkGaussDBMySQLJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for tables to GaussDB MySQL instance(%s) to complete: %s", instanceId, err)
	}

	return nil
}

func buildCreateTableRestoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"restore_time":    d.Get("restore_time"),
		"restore_tables":  buildCreateTableRestoreRestoreTablesBodyParams(d.Get("restore_tables")),
		"last_table_info": utils.ValueIgnoreEmpty(d.Get("last_table_info")),
	}
	return bodyParams
}

func buildCreateTableRestoreRestoreTablesBodyParams(restoreTables interface{}) []map[string]interface{} {
	restoreTablesRaw := restoreTables.([]interface{})
	res := make([]map[string]interface{}, 0)
	for _, restoreTableRaw := range restoreTablesRaw {
		if v, ok := restoreTableRaw.(map[string]interface{}); ok {
			res = append(res, map[string]interface{}{
				"database": v["database"],
				"tables":   buildCreateTableRestoreRestoreTablesTablesBodyParams(v["tables"]),
			})
		}
	}
	return res
}

func buildCreateTableRestoreRestoreTablesTablesBodyParams(tables interface{}) []map[string]interface{} {
	tablesRaw := tables.([]interface{})
	res := make([]map[string]interface{}, 0)
	for _, tableRaw := range tablesRaw {
		if v, ok := tableRaw.(map[string]interface{}); ok {
			res = append(res, map[string]interface{}{
				"old_name": v["old_name"],
				"new_name": v["new_name"],
			})
		}
	}
	return res
}

func resourceGaussDBMysqlTableRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBMysqlTableRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting restoration tables is not supported. The restoration tables is only removed from the state," +
		" but it remains in the cloud. And the instance doesn't return to the state before restoration."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
