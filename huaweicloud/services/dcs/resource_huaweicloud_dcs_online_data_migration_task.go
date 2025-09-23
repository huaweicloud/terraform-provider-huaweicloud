package dcs

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var onlineDataMigrationTaskNonUpdatableParams = []string{"task_name", "vpc_id", "subnet_id", "security_group_id",
	"description"}

// @API DCS POST /v2/{project_id}/migration/instance
// @API DCS GET /v2/{project_id}/jobs/{job_id}
// @API DCS POST /v2/{project_id}/migration/{task_id}/task
// @API DCS GET /v2/{project_id}/migration-task/{task_id}
// @API DCS GET /v2/{project_id}/migration-tasks
// @API DCS POST /v2/{project_id}/migration-task/{task_id}/stop
// @API DCS DELETE /v2/{project_id}/migration-tasks/delete
func ResourceDcsOnlineDataMigrationTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsOnlineDataMigrationTaskCreate,
		ReadContext:   resourceDcsOnlineDataMigrationTaskRead,
		UpdateContext: resourceDcsOnlineDataMigrationTaskUpdate,
		DeleteContext: resourceDcsOnlineDataMigrationTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(onlineDataMigrationTaskNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"migration_method": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resume_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"migration_method"},
			},
			"bandwidth_limit_mb": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"migration_method"},
			},
			"source_instance": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				Elem:         onlineDataMigrationInstanceConfig(),
				RequiredWith: []string{"migration_method"},
			},
			"target_instance": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				Elem:         onlineDataMigrationInstanceConfig(),
				RequiredWith: []string{"migration_method"},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"ecs_tenant_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"supported_features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"released_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func onlineDataMigrationInstanceConfig() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"addrs": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDcsOnlineDataMigrationTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/migration/instance"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildCreateOnlineDataMigrationTaskBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating online data migration task: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("instance_id", createRespBody, "").(string)
	if taskId == "" {
		return diag.Errorf("error creating online data migration task: instance_id is not found in API response")
	}

	d.SetId(taskId)

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating online data migration task: job_id is not found in API response")
	}
	err = checkDcsInstanceJobFinish(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := d.GetOk("migration_method"); ok {
		err = updateOnlineDataMigrationTask(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
		targetStatus := "SUCCESS"
		if v.(string) == "incremental_migration" {
			targetStatus = "INCRMIGEATING"
		}
		err = checkMigrationTaskFinish(ctx, client, taskId, []string{targetStatus}, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDcsOnlineDataMigrationTaskRead(ctx, d, meta)
}

func buildCreateOnlineDataMigrationTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":              d.Get("task_name"),
		"vpc_id":            d.Get("vpc_id"),
		"subnet_id":         d.Get("subnet_id"),
		"security_group_id": d.Get("security_group_id"),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func updateOnlineDataMigrationTask(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v2/{project_id}/migration/{task_id}/task"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{task_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	updateOpt.JSONBody = utils.RemoveNil(buildConfigOnlineDataMigrationTaskBodyParams(d))

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating online data migration task: %s", err)
	}

	return nil
}

func buildConfigOnlineDataMigrationTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"migration_method":   d.Get("migration_method"),
		"resume_mode":        d.Get("resume_mode"),
		"source_instance":    buildOnlineDataMigrationInstanceConfigBodyParams(d.Get("source_instance")),
		"target_instance":    buildOnlineDataMigrationInstanceConfigBodyParams(d.Get("target_instance")),
		"bandwidth_limit_mb": utils.ValueIgnoreEmpty(d.Get("bandwidth_limit_mb")),
	}
	return bodyParams
}

func buildOnlineDataMigrationInstanceConfigBodyParams(resp interface{}) map[string]interface{} {
	configRaw := resp.([]interface{})
	if len(configRaw) == 0 {
		return nil
	}

	bodyParams := map[string]interface{}{
		"id":       utils.ValueIgnoreEmpty(configRaw[0].(map[string]interface{})["id"]),
		"addrs":    utils.ValueIgnoreEmpty(configRaw[0].(map[string]interface{})["addrs"]),
		"password": utils.ValueIgnoreEmpty(configRaw[0].(map[string]interface{})["password"]),
	}
	return bodyParams
}

func resourceDcsOnlineDataMigrationTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	getRespBody, err := getMigrationTask(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DCS.4133"),
			"error getting DCS online data migration task")
	}

	status := utils.PathSearch("status", getRespBody, "").(string)
	if status == "DELETED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting DCS online data migration task")
	} else if status == "SUCCESS" {
		// when the migration task is deleted, the value of the task status queried by the task detail API may be SUCCESS,
		// but it can not be queried by the list API
		getListRespBody, err := getMigrationTaskList(client)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error getting DCS online data migration task")
		}

		task := utils.PathSearch(fmt.Sprintf("migration_tasks[?task_id=='%s']|[0]", d.Id()), getListRespBody, nil)
		if task == nil {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting DCS online data migration task")
		}
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("task_name", utils.PathSearch("task_name", getRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("tenant_vpc_id", getRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("tenant_subnet_id", getRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("tenant_security_group_id", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("migration_method", utils.PathSearch("migration_method", getRespBody, nil)),
		d.Set("resume_mode", utils.PathSearch("resume_mode", getRespBody, nil)),
		d.Set("bandwidth_limit_mb", utils.PathSearch("bandwidth_limit_mb", getRespBody, nil)),
		d.Set("source_instance", flattenMigrationInstanceConfig(d, getRespBody, "source_instance")),
		d.Set("target_instance", flattenMigrationInstanceConfig(d, getRespBody, "target_instance")),
		d.Set("ecs_tenant_private_ip", utils.PathSearch("ecs_tenant_private_ip", getRespBody, nil)),
		d.Set("network_type", utils.PathSearch("network_type", getRespBody, nil)),
		d.Set("status", status),
		d.Set("supported_features", utils.PathSearch("supported_features", getRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
		d.Set("released_at", utils.PathSearch("released_at", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMigrationInstanceConfig(d *schema.ResourceData, getRespBody interface{}, expression string) []interface{} {
	resp := utils.PathSearch(expression, getRespBody, nil)
	if resp == nil {
		return nil
	}

	instanceConfig := map[string]interface{}{
		"id":    utils.PathSearch("id", resp, nil),
		"addrs": utils.PathSearch("addrs", resp, nil),
		"name":  utils.PathSearch("name", resp, nil),
	}
	rawInstanceConfig := d.Get(expression).([]interface{})
	if len(rawInstanceConfig) > 0 {
		instanceConfig["password"] = rawInstanceConfig[0].(map[string]interface{})["password"]
	}

	return []interface{}{instanceConfig}
}

func resourceDcsOnlineDataMigrationTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	err = updateOnlineDataMigrationTask(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	targetStatus := "SUCCESS"
	if d.Get("migration_method").(string) == "incremental_migration" {
		targetStatus = "INCRMIGEATING"
	}
	err = checkMigrationTaskFinish(ctx, client, d.Id(), []string{targetStatus}, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDcsOnlineDataMigrationTaskRead(ctx, d, meta)
}

func resourceDcsOnlineDataMigrationTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	if d.Get("migration_method") == "incremental_migration" {
		getRespBody, err := getMigrationTask(client, d.Id())
		if err != nil {
			return diag.Errorf("error getting DCS online data migration task: %s", err)
		}
		status := utils.PathSearch("status", getRespBody, "").(string)
		if status == "INCRMIGEATING" {
			err = stopOnlineDataMigrationTask(d, client)
			if err != nil {
				return diag.FromErr(err)
			}

			err = checkMigrationTaskFinish(ctx, client, d.Id(), []string{"TERMINATED"}, d.Timeout(schema.TimeoutDelete))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	err = deleteMigrationTask(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	err = checkMigrationTaskDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func stopOnlineDataMigrationTask(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	httpUrl := "v2/{project_id}/migration-task/{task_id}/stop"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{task_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating online data migration task: %s", err)
	}

	return nil
}
