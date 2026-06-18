package taurusdb

import (
	"context"
	"errors"
	"fmt"
	"log"
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

var htapStarrocksReplicationNoneUpdatableParams = []string{
	"instance_id", "task_name", "source_instance_id",
}

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/database-config-check
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/table-config-check
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/configuration
// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication
// @API TaurusDB PUT /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication
// @API TaurusDB DELETE /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/pause
// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/resume
// @API TaurusDB GET /v3/{project_id}/jobs
func ResourceTaurusDBHtapStarrocksReplication() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBHtapStarrocksReplicationCreate,
		ReadContext:   resourceTaurusDBHtapStarrocksReplicationRead,
		UpdateContext: resourceTaurusDBHtapStarrocksReplicationUpdate,
		DeleteContext: resourceTaurusDBHtapStarrocksReplicationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceTaurusDBHtapStarrocksReplicationImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(htapStarrocksReplicationNoneUpdatableParams),

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
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_node_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_instance_level_sync": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"database_repl_scope": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"source_database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_database_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"db_configs": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     starrocksReplicationDbConfigSchema(),
			},
			"tables_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     starrocksReplicationTablesConfigSchema(),
			},
			"table_repl_config": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     starrocksReplicationTableReplConfigSchema(),
			},
			"enable_sync": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
			},
			"sync_action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			// Computed fields in configuration response
			"database_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationDatabaseInfoSchema(),
			},
			"table_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationTableConfigCheckResultSchema(),
			},
			"new_table_repl_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationTableReplConfigSchema(),
			},
			"is_tables_change": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_error_of_alter_table": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_support_reg_exp": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Computed fields in status response
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"percentage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_need_repair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_main_task": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func starrocksReplicationDbConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"param_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func starrocksReplicationTablesConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"table_config": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func starrocksReplicationTableReplConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"repl_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repl_scope": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func starrocksReplicationDatabaseInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_config_check_results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationDbConfigCheckResultSchema(),
			},
		},
	}
}

func starrocksReplicationDbConfigCheckResultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"param_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func starrocksReplicationTableConfigCheckResultSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"table_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"table_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTaurusDBHtapStarrocksReplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	if v, ok := d.GetOk("tables_configs"); ok {
		tablesConfigs := v.([]interface{})
		if len(tablesConfigs) > 0 {
			err := createStarrocksDatabaseReplicationWithTablesConfigs(client, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	} else {
		err := createStarrocksDatabaseReplication(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Waiting for replication task creation completed
	instanceId := d.Get("instance_id").(string)
	taskName := d.Get("task_name").(string)
	details, err := waitForReplicationTaskCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), instanceId, taskName)
	if err != nil {
		return diag.Errorf("error waiting for creating HTAP instance (%s) replication task (%s): %s",
			instanceId, taskName, err)
	}

	// Check the replication task status, if the status is abnormal, delete the task and return error
	errMsg := utils.PathSearch("last_error_of_alter_table", details, "").(string)
	if errMsg != "" {
		log.Printf("[WARN] created HTAP instance (%s) replication task (%s) is abnormal: %s", instanceId, taskName, errMsg)
		err = deleteReplicationTask(client, instanceId, taskName)
		if err != nil {
			return diag.Errorf("error deleting abnormal HTAP instance (%s) replication task (%s): %s", instanceId, taskName, err)
		}
		return diag.Errorf("error creating HTAP instance (%s) replication task (%s): %s", instanceId, taskName, errMsg)
	}

	// Start the replication task
	if v, ok := d.Get("enable_sync").(string); ok && v == "true" {
		err = startSyncInReplicationTask(client, instanceId, taskName, details)
		if err != nil {
			return diag.Errorf("error starting HTAP insatnce (%s) replication task (%s): %s", instanceId, taskName, err)
		}
		_, err = waitForReplicationTaskCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), instanceId, taskName)
		if err != nil {
			return diag.Errorf("error waiting for starting HTAP instance (%s) replication task (%s) completed: %s",
				instanceId, taskName, err)
		}
	}

	// operate action the replication task
	if v, ok := d.Get("sync_action").(string); ok {
		if v == "pause" {
			err = pauseSyncInReplicationTask(client, instanceId, taskName)
			if err != nil {
				return diag.Errorf("error pausing HTAP insatnce (%s) replication task (%s): %s", instanceId, taskName, err)
			}
			_, err = waitForReplicationTaskCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), instanceId, taskName)
			if err != nil {
				return diag.Errorf("error waiting for pausing HTAP instance (%s) replication task (%s) completed: %s",
					instanceId, taskName, err)
			}
		}
		if v == "resume" {
			err = resumeSyncInReplicationTask(client, instanceId, taskName)
			if err != nil {
				return diag.Errorf("error resuming HTAP insatnce (%s) replication task (%s): %s", instanceId, taskName, err)
			}
			_, err = waitForReplicationTaskCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), instanceId, taskName)
			if err != nil {
				return diag.Errorf("error waiting for resuming HTAP instance (%s) replication task (%s) completed: %s",
					instanceId, taskName, err)
			}
		}
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceId, taskName))
	return resourceTaurusDBHtapStarrocksReplicationRead(ctx, d, meta)
}

func createStarrocksDatabaseReplication(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/database-config-check"
		instanceId = d.Get("instance_id").(string)
		taskName   = d.Get("task_name").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = buildStarrocksReplicationConfigBodyParams(d)

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error creating  HTAP instance (%s) replication task (%s) task: %s", instanceId, taskName, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}
	taskNameResp := utils.PathSearch("task_name", respBody, "")
	if taskNameResp == "" {
		return fmt.Errorf("error creating HTAP instance (%s) replication task (%s): task name is empty", instanceId, taskName)
	}
	return nil
}

func createStarrocksDatabaseReplicationWithTablesConfigs(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/table-config-check"
		instanceId = d.Get("instance_id").(string)
		taskName   = d.Get("task_name").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	bodyParams := buildStarrocksReplicationConfigBodyParams(d)
	bodyParams["is_create_task"] = "true"
	createOpt.JSONBody = bodyParams

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error creating HTAP instance (%s) replication task (%s) with tables configs: %s",
			instanceId, taskName, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}
	taskNameResp := utils.PathSearch("task_name", respBody, "")
	if taskNameResp == "" {
		return fmt.Errorf("error creating HTAP instance (%s) replication task (%s) with tables configs: "+
			"task name is empty", instanceId, taskName)
	}

	return nil
}

func getReplicationDetails(client *golangsdk.ServiceClient, instanceId, taskName string) (interface{}, error) {
	taskConfig, err := getReplicationConfig(client, instanceId, taskName)
	if err != nil {
		return nil, err
	}
	taskStatus, err := getReplicationStatus(client, instanceId, taskName)
	if err != nil {
		return nil, err
	}

	// Merge config and status into a single object
	details := make(map[string]interface{})
	// Add all fields from taskConfig
	if configMap, ok := taskConfig.(map[string]interface{}); ok {
		for k, v := range configMap {
			details[k] = v
		}
	}

	// Add all fields from taskStatus
	if statusMap, ok := taskStatus.(map[string]interface{}); ok {
		for k, v := range statusMap {
			details[k] = v
		}
	}

	return details, nil
}

func getReplicationStatus(client *golangsdk.ServiceClient, instanceId, taskName string) (interface{}, error) {
	replications, err := getHtapStarrocksReplications(client, instanceId)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HTAP instance(%s) replications: %s", instanceId, err)
	}

	// Filter by task_name
	for _, replication := range replications {
		if utils.PathSearch("task_name", replication, "").(string) == taskName {
			return replication, nil
		}
	}

	return nil, golangsdk.ErrDefault404{}
}

func getReplicationConfig(client *golangsdk.ServiceClient, instanceId, taskName string) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/configuration"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = fmt.Sprintf("%s?task_name=%s", getPath, taskName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func resourceTaurusDBHtapStarrocksReplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	taskName := d.Get("task_name").(string)
	details, err := getReplicationDetails(client, instanceId, taskName)
	if err != nil {
		// task not found
		err = common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.270022")
		return common.CheckDeletedDiag(d, err, "error retrieving HTAP instance replication task")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("source_instance_id", utils.PathSearch("source_instance_id", details, nil)),
		d.Set("source_node_id", utils.PathSearch("source_node_id", details, nil)),
		d.Set("source_database_name", utils.PathSearch("source_database", details, nil)),
		d.Set("target_database_name", utils.PathSearch("target_database", details, nil)),
		d.Set("task_name", utils.PathSearch("task_name", details, nil)),
		d.Set("is_instance_level_sync", utils.PathSearch("is_instance_level_sync", details, nil)),
		d.Set("database_repl_scope", utils.PathSearch("database_repl_scope", details, nil)),
		d.Set("database_info", flattenStarrocksReplicationDatabaseInfo(details)),
		d.Set("table_infos", flattenStarrocksReplicationTableInfos(details)),
		d.Set("table_repl_config", flattenStarrocksReplicationTableReplConfig(utils.PathSearch("table_repl_config", details, nil))),
		d.Set("new_table_repl_config", flattenStarrocksReplicationTableReplConfig(utils.PathSearch("new_table_repl_config", details, nil))),
		d.Set("is_support_reg_exp", utils.PathSearch("is_support_reg_exp", details, nil)),
		d.Set("is_tables_change", utils.PathSearch("is_tables_change", details, false).(bool)),
		d.Set("error_msg", utils.PathSearch("error_msg", details, nil)),
		d.Set("last_error_of_alter_table", utils.PathSearch("last_error_of_alter_table", details, nil)),
		// Set fields from taskStatus
		d.Set("enable_sync", d.Get("enable_sync").(string)),
		d.Set("sync_action", d.Get("sync_action").(string)),
		d.Set("status", utils.PathSearch("status", details, nil)),
		d.Set("stage", utils.PathSearch("stage", details, nil)),
		d.Set("percentage", utils.PathSearch("percentage", details, nil)),
		d.Set("is_need_repair", utils.PathSearch("is_need_repair", details, nil)),
		d.Set("is_main_task", utils.PathSearch("is_main_task", details, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTaurusDBHtapStarrocksReplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		taskName   = d.Get("task_name").(string)
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}
	details, err := getReplicationDetails(client, instanceId, taskName)
	if err != nil {
		return diag.Errorf("error retrieving HTAP instance (%s) replication task (%s): %s", instanceId, taskName, err)
	}

	updateConfigFields := []string{
		"is_instance_level_sync",
		"database_repl_scope",
		"source_node_id",
		"source_database_name",
		"target_database_name",
		"db_configs",
		"table_repl_config",
		"tables_configs",
	}
	// update config
	if d.HasChanges(updateConfigFields...) {
		err = updateHtapStarrocksReplication(ctx, client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// update task status
	if d.HasChanges("enable_sync") {
		if v, ok := d.Get("enable_sync").(string); ok && v == "true" {
			err = startSyncInReplicationTask(client, instanceId, taskName, details)
			if err != nil {
				return diag.Errorf("error starting HTAP insatnce (%s) replication task (%s): %s", instanceId, taskName, err)
			}
			_, err = waitForReplicationTaskCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId, taskName)
			if err != nil {
				return diag.Errorf("error waiting for starting HTAP instance (%s) replication task (%s) completed: %s",
					instanceId, taskName, err)
			}
		} else if v == "false" {
			return diag.FromErr(errors.New("the modification changing enable_sync to value 'false' is not allowed"))
		}
		// If v is empty string (""), do nothing, just ignore the change
	}

	if d.HasChanges("sync_action") {
		syncAction := d.Get("sync_action").(string)
		if syncAction == "pause" {
			err = pauseSyncInReplicationTask(client, instanceId, taskName)
			if err != nil {
				return diag.Errorf("error pausing HTAP insatnce (%s) replication task (%s): %s", instanceId, taskName, err)
			}
			_, err = waitForReplicationTaskCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId, taskName)
			if err != nil {
				return diag.Errorf("error waiting for pausing HTAP instance (%s) replication task (%s) completed: %s",
					instanceId, taskName, err)
			}
		}
		if syncAction == "resume" {
			err = resumeSyncInReplicationTask(client, instanceId, taskName)
			if err != nil {
				return diag.Errorf("error resuming HTAP insatnce (%s) replication task (%s): %s", instanceId, taskName, err)
			}
			_, err = waitForReplicationTaskCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId, taskName)
			if err != nil {
				return diag.Errorf("error waiting for resuming HTAP instance (%s) replication task (%s) completed: %s",
					instanceId, taskName, err)
			}
		}
	}

	return resourceTaurusDBHtapStarrocksReplicationRead(ctx, d, meta)
}

func updateHtapStarrocksReplication(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication"
		instanceId = d.Get("instance_id").(string)
		taskName   = d.Get("task_name").(string)
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	updateOpt.JSONBody = buildStarrocksReplicationConfigBodyParams(d)

	res, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating HTAP instance (%s) replication task (%s): %s", instanceId, taskName, err)
	}

	updateRespBody, err := utils.FlattenResponse(res)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error updating HTAP instance (%s) replication task (%s): job_id is empty in the response",
			instanceId, taskName)
	}
	err = waitForJobComplete(ctx, client, d.Timeout(schema.TimeoutUpdate), instanceId, jobId)
	if err != nil {
		return fmt.Errorf("error waitting for updating HTAP instance (%s) replication task (%s) to completed: %s",
			instanceId, taskName, err)
	}
	return nil
}

func buildStarrocksReplicationConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"task_name":            d.Get("task_name"),
		"source_instance_id":   d.Get("source_instance_id"),
		"source_database_name": d.Get("source_database_name"),
		"target_database_name": d.Get("target_database_name"),
		"db_configs":           buildStarrocksReplicationDbConfigs(d.Get("db_configs").([]interface{})),
		"table_repl_config":    buildStarrocksReplicationTableReplConfig(d.Get("table_repl_config").([]interface{})),
	}
	if v, ok := d.GetOk("is_instance_level_sync"); ok {
		bodyParams["is_instance_level_sync"] = v
	}
	if v, ok := d.GetOk("database_repl_scope"); ok {
		bodyParams["database_repl_scope"] = v
	}
	if v, ok := d.GetOk("source_node_id"); ok {
		bodyParams["source_node_id"] = v
	}
	if v, ok := d.GetOk("tables_configs"); ok {
		bodyParams["tables_configs"] = buildStarrocksReplicationTablesConfigs(v.([]interface{}))
	}
	return bodyParams
}

func buildStarrocksReplicationDbConfigs(dbConfigs []interface{}) []interface{} {
	if len(dbConfigs) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(dbConfigs))
	for _, v := range dbConfigs {
		configMap := v.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"param_name": configMap["param_name"],
			"value":      configMap["value"],
		})
	}
	return result
}

func buildStarrocksReplicationTableReplConfig(configs []interface{}) interface{} {
	if len(configs) == 0 {
		return nil
	}
	configMap := configs[0].(map[string]interface{})
	result := map[string]interface{}{
		"repl_type":  configMap["repl_type"],
		"repl_scope": configMap["repl_scope"],
	}
	if v, ok := configMap["tables"]; ok && v != nil {
		result["tables"] = v.([]interface{})
	}
	return result
}

func buildStarrocksReplicationTablesConfigs(configs []interface{}) []interface{} {
	if len(configs) == 0 {
		return nil
	}
	result := make([]interface{}, 0, len(configs))
	for _, v := range configs {
		configMap := v.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"table_name":   configMap["table_name"],
			"table_config": configMap["table_config"],
		})
	}
	return result
}

func deleteReplicationTask(client *golangsdk.ServiceClient, instanceId, taskName string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"task_name": taskName,
		},
	}
	resp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return err
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error deleting HTAP instance (%s) replication task (%s): job_id is empty in the response",
			instanceId, taskName)
	}
	return nil
}

func resourceTaurusDBHtapStarrocksReplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}
	instanceId := d.Get("instance_id").(string)
	taskName := d.Get("task_name").(string)
	err = deleteReplicationTask(client, instanceId, taskName)
	if err != nil {
		// task not found
		err = common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.05000091")
		return common.CheckDeletedDiag(d, err, "error deleting HTAP instance replication task")
	}
	return nil
}

func flattenStarrocksReplicationDatabaseInfo(resp interface{}) []interface{} {
	databaseInfo := utils.PathSearch("database_info", resp, nil)
	if databaseInfo == nil {
		return nil
	}

	dbConfigCheckResults := utils.PathSearch("dbCfgCheckResults", databaseInfo, make([]interface{}, 0)).([]interface{})
	checkResults := make([]interface{}, 0, len(dbConfigCheckResults))
	for _, v := range dbConfigCheckResults {
		checkResults = append(checkResults, map[string]interface{}{
			"param_name":   utils.PathSearch("param_name", v, nil),
			"value":        utils.PathSearch("value", v, nil),
			"check_result": utils.PathSearch("check_result", v, nil),
		})
	}

	return []interface{}{
		map[string]interface{}{
			"database_name":           utils.PathSearch("databaseName", databaseInfo, nil),
			"db_config_check_results": checkResults,
		},
	}
}

func flattenStarrocksReplicationTableInfos(resp interface{}) []interface{} {
	tableInfos := utils.PathSearch("table_infos", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, 0, len(tableInfos))
	for _, v := range tableInfos {
		result = append(result, map[string]interface{}{
			"table_name":   utils.PathSearch("table_name", v, nil),
			"table_config": utils.PathSearch("table_config", v, nil),
			"check_result": utils.PathSearch("check_result", v, nil),
		})
	}
	return result
}

func flattenStarrocksReplicationTableReplConfig(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	tables := utils.PathSearch("tables", resp, make([]interface{}, 0)).([]interface{})

	return []interface{}{
		map[string]interface{}{
			"repl_type":  utils.PathSearch("repl_type", resp, nil),
			"repl_scope": utils.PathSearch("repl_scope", resp, nil),
			"tables":     utils.ExpandToStringList(tables),
		},
	}
}

func waitForReplicationTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	instanceId, taskName string) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Waiting"},
		Target:       []string{"Completed"},
		Refresh:      replicationTaskStatusRefreshFunc(client, instanceId, taskName),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func replicationTaskStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, taskName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		details, err := getReplicationDetails(client, instanceId, taskName)
		if err != nil {
			// If the task is not found yet, return PENDING state
			log.Printf("[DEBUG] The replication task (%s) is not found yet, continue waiting", taskName)
			return details, "Waiting", nil
		}
		status := utils.PathSearch("status", details, "").(string)
		if status == "Wait" || status == "" {
			return details, "Waiting", nil
		}
		return details, "Completed", nil
	}
}

func startSyncInReplicationTask(client *golangsdk.ServiceClient, instanceId, taskName string, details interface{}) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication"
	startPath := client.Endpoint + httpUrl
	startPath = strings.ReplaceAll(startPath, "{project_id}", client.ProjectID)
	startPath = strings.ReplaceAll(startPath, "{instance_id}", instanceId)

	startOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	startOpt.JSONBody = map[string]interface{}{
		"source_instance_id": utils.PathSearch("source_instance_id", details, ""),
		"source_database":    utils.PathSearch("source_database", details, ""),
		"target_database":    utils.PathSearch("target_database_name", details, ""),
		"task_name":          taskName,
	}

	if v := utils.PathSearch("source_node_id", details, ""); v.(string) != "" {
		startOpt.JSONBody.(map[string]interface{})["source_node_id"] = v.(string)
	}

	resp, err := client.Request("POST", startPath, &startOpt)
	if err != nil {
		return fmt.Errorf("error starting HTAP instance (%s) replication task (%s): %s", instanceId, taskName, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error starting HTAP instance (%s) replication task (%s): job_id not found in the response",
			instanceId, taskName)
	}
	return nil
}

func pauseSyncInReplicationTask(client *golangsdk.ServiceClient, instanceId, taskName string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/pause"
	pausePath := client.Endpoint + httpUrl
	pausePath = strings.ReplaceAll(pausePath, "{project_id}", client.ProjectID)
	pausePath = strings.ReplaceAll(pausePath, "{instance_id}", instanceId)

	pauseOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	pauseOpt.JSONBody = map[string]interface{}{
		"task_name": taskName,
	}

	resp, err := client.Request("POST", pausePath, &pauseOpt)
	if err != nil {
		return fmt.Errorf("error pausing HTAP instance (%s) replication task (%s):%s", instanceId, taskName, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error pausing HTAP instance (%s) replication task (%s): job_id not found in the response",
			instanceId, taskName)
	}

	return nil
}

func resumeSyncInReplicationTask(client *golangsdk.ServiceClient, instanceId, taskName string) error {
	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication/resume"
	resumePath := client.Endpoint + httpUrl
	resumePath = strings.ReplaceAll(resumePath, "{project_id}", client.ProjectID)
	resumePath = strings.ReplaceAll(resumePath, "{instance_id}", instanceId)

	resumeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	resumeOpt.JSONBody = map[string]interface{}{
		"task_name": taskName,
	}

	resp, err := client.Request("POST", resumePath, &resumeOpt)
	if err != nil {
		return fmt.Errorf("error resuming HTAP instance (%s) replication task (%s): %s", instanceId, taskName, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return fmt.Errorf("error resuming HTAP instance (%s) replication task (%s): job_id not found in the response",
			instanceId, taskName)
	}

	return nil
}

func resourceTaurusDBHtapStarrocksReplicationImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<task_name>")
	}

	instanceId := parts[0]
	taskName := parts[1]

	mErr := multierror.Append(nil,
		d.Set("instance_id", instanceId),
		d.Set("task_name", taskName),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
