package oms

import (
	"context"
	"fmt"
	"strconv"
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

// @API OMS POST /v2/{project_id}/sync-tasks
// @API OMS GET /v2/{project_id}/sync-tasks/{sync_task_id}
// @API OMS DELETE /v2/{project_id}/sync-tasks/{sync_task_id}
// @API OMS POST /v2/{project_id}/sync-tasks/{sync_task_id}/stop
// @API OMS POST /v2/{project_id}/sync-tasks/{sync_task_id}/start

func ResourceMigrationSyncTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrationSyncTaskCreate,
		ReadContext:   resourceMigrationSyncTaskRead,
		UpdateContext: resourceMigrationSyncTaskUpdate,
		DeleteContext: resourceMigrationSyncTaskDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"src_region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"src_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"src_ak": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"src_sk": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
			},

			"dst_bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dst_ak": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dst_sk": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
				ForceNew:  true,
			},
			"src_cloud_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_kms": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"enable_restore": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"enable_metadata_migration": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"source_cdn": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"http", "https",
							}, false),
						},
						"authentication_type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "NONE",
						},
						"authentication_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"consistency_check": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"start", "stop",
				}, false),
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_start_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"object_overwrite_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dst_storage_policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"monthly_acceptance_request": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"monthly_success_object": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"monthly_failure_object": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"monthly_skip_object": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"monthly_size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func resourceMigrationSyncTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createSyncTaskHttpUrl = "v2/{project_id}/sync-tasks"
		createSyncTaskProduct = "oms"
	)
	createSyncTaskClient, err := cfg.NewServiceClient(createSyncTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	createSyncTaskPath := createSyncTaskClient.Endpoint + createSyncTaskHttpUrl
	createSyncTaskPath = strings.ReplaceAll(createSyncTaskPath, "{project_id}", createSyncTaskClient.ProjectID)

	createSyncTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createSyncTaskOpt.JSONBody = utils.RemoveNil(buildcreateSyncTaskBodyParams(d, region))
	createSyncTaskResp, err := createSyncTaskClient.Request("POST", createSyncTaskPath, &createSyncTaskOpt)
	if err != nil {
		return diag.Errorf("error creating OMS migration sync task: %s", err)
	}

	createSyncTaskRespBody, err := utils.FlattenResponse(createSyncTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("sync_task_id", createSyncTaskRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating OMS migration sync task: ID is not found in API response")
	}

	d.SetId(id.(string))

	if d.Get("action").(string) == "stop" {
		err := stopSyncTask(createSyncTaskClient, d)
		if err != nil {
			return diag.Errorf("error stopping OMS migration sync task: %s", err)
		}
	}

	return resourceMigrationSyncTaskRead(ctx, d, meta)
}

func buildcreateSyncTaskBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	createOpts := map[string]interface{}{
		"src_cloud_type":            utils.ValueIgnoreEmpty(d.Get("src_cloud_type")),
		"src_region":                d.Get("src_region"),
		"src_bucket":                d.Get("src_bucket"),
		"src_ak":                    d.Get("src_ak"),
		"src_sk":                    d.Get("src_sk"),
		"source_cdn":                buildSourceCdnOpts(d.Get("source_cdn").([]interface{})),
		"dst_region":                region,
		"dst_bucket":                d.Get("dst_bucket"),
		"dst_ak":                    d.Get("dst_ak"),
		"dst_sk":                    d.Get("dst_sk"),
		"description":               utils.ValueIgnoreEmpty(d.Get("description")),
		"enable_kms":                d.Get("enable_kms"),
		"enable_restore":            d.Get("enable_restore"),
		"enable_metadata_migration": d.Get("enable_metadata_migration"),
		"app_id":                    utils.ValueIgnoreEmpty(d.Get("app_id")),
		"consistency_check":         utils.ValueIgnoreEmpty(d.Get("consistency_check")),
	}

	return createOpts
}

func resourceMigrationSyncTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	queryTime := strconv.FormatInt(time.Now().UnixMilli(), 10)

	var (
		getSyncTaskHttpUrl = "v2/{project_id}/sync-tasks/{sync_task_id}?query_time=" + queryTime
		getSyncTaskProduct = "oms"
	)
	getSyncTaskClient, err := cfg.NewServiceClient(getSyncTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	getSyncTaskPath := getSyncTaskClient.Endpoint + getSyncTaskHttpUrl
	getSyncTaskPath = strings.ReplaceAll(getSyncTaskPath, "{project_id}", getSyncTaskClient.ProjectID)
	getSyncTaskPath = strings.ReplaceAll(getSyncTaskPath, "{sync_task_id}", d.Id())

	getSyncTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getSyncTaskResp, err := getSyncTaskClient.Request("GET", getSyncTaskPath, &getSyncTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OMS migration sync task")
	}

	getSyncTaskRespBody, err := utils.FlattenResponse(getSyncTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("status", getSyncTaskRespBody, "").(string)
	actionValue := "start"
	if status == "STOPPED" {
		actionValue = "stop"
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("src_cloud_type", utils.PathSearch("src_cloud_type", getSyncTaskRespBody, nil)),
		d.Set("src_region", utils.PathSearch("src_region", getSyncTaskRespBody, nil)),
		d.Set("src_bucket", utils.PathSearch("src_bucket", getSyncTaskRespBody, nil)),
		d.Set("dst_bucket", utils.PathSearch("dst_bucket", getSyncTaskRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getSyncTaskRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getSyncTaskRespBody, nil)),
		d.Set("enable_kms", utils.PathSearch("enable_kms", getSyncTaskRespBody, nil)),
		d.Set("enable_metadata_migration", utils.PathSearch("enable_metadata_migration", getSyncTaskRespBody, nil)),
		d.Set("enable_restore", utils.PathSearch("enable_restore", getSyncTaskRespBody, nil)),
		d.Set("app_id", utils.PathSearch("app_id", getSyncTaskRespBody, nil)),
		d.Set("object_overwrite_mode", utils.PathSearch("object_overwrite_mode", getSyncTaskRespBody, nil)),
		d.Set("dst_storage_policy", utils.PathSearch("dst_storage_policy", getSyncTaskRespBody, nil)),
		d.Set("consistency_check", utils.PathSearch("consistency_check", getSyncTaskRespBody, nil)),
		d.Set("monthly_acceptance_request", utils.PathSearch("monthly_acceptance_request", getSyncTaskRespBody, nil)),
		d.Set("monthly_success_object", utils.PathSearch("monthly_success_object", getSyncTaskRespBody, nil)),
		d.Set("monthly_failure_object", utils.PathSearch("monthly_failure_object", getSyncTaskRespBody, nil)),
		d.Set("monthly_skip_object", utils.PathSearch("monthly_skip_object", getSyncTaskRespBody, nil)),
		d.Set("monthly_size", utils.PathSearch("monthly_size", getSyncTaskRespBody, nil)),
		d.Set("action", actionValue),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getSyncTaskRespBody, float64(0)).(float64))/1000, false)),
		d.Set("last_start_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("last_start_time", getSyncTaskRespBody, float64(0)).(float64))/1000, false)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OMS migration sync task fields: %s", err)
	}

	return nil
}

func resourceMigrationSyncTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateSyncTaskProduct = "oms"
	)
	updateSyncTaskClient, err := cfg.NewServiceClient(updateSyncTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		if action == "start" {
			err := startSyncTask(updateSyncTaskClient, d)
			if err != nil {
				return diag.Errorf("error starting OMS migration sync task: %s", err)
			}
		} else {
			if d.Get("action").(string) == "stop" {
				err := stopSyncTask(updateSyncTaskClient, d)
				if err != nil {
					return diag.Errorf("error stopping OMS migration sync task: %s", err)
				}
			}
		}
	}

	return resourceMigrationSyncTaskRead(ctx, d, meta)
}

func buildstartSyncTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	startActionOpts := map[string]interface{}{
		"src_ak": d.Get("src_ak").(string),
		"src_sk": d.Get("src_sk").(string),
		"dst_ak": d.Get("dst_ak").(string),
		"dst_sk": d.Get("dst_sk").(string),
	}
	if sourceCDNs := d.Get("source_cdn").([]interface{}); len(sourceCDNs) > 0 {
		sourceCdn := sourceCDNs[0].(map[string]interface{})
		startActionOpts["source_cdn_authentication_key"] = utils.String(sourceCdn["authentication_key"].(string))
	}

	return startActionOpts
}

func startSyncTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	startSyncTaskHttpUrl := "v2/{project_id}/sync-tasks/{sync_task_id}/start"
	startSyncTaskPath := client.Endpoint + startSyncTaskHttpUrl
	startSyncTaskPath = strings.ReplaceAll(startSyncTaskPath, "{project_id}", client.ProjectID)
	startSyncTaskPath = strings.ReplaceAll(startSyncTaskPath, "{sync_task_id}", d.Id())

	startSyncTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	startSyncTaskOpt.JSONBody = utils.RemoveNil(buildstartSyncTaskBodyParams(d))
	_, err := client.Request("POST", startSyncTaskPath, &startSyncTaskOpt)
	if err != nil {
		return fmt.Errorf("error starting OMS migration sync task: %s", err)
	}

	return nil
}

func stopSyncTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stopSyncTaskHttpUrl := "v2/{project_id}/sync-tasks/{sync_task_id}/stop"
	stopSyncTaskPath := client.Endpoint + stopSyncTaskHttpUrl
	stopSyncTaskPath = strings.ReplaceAll(stopSyncTaskPath, "{project_id}", client.ProjectID)
	stopSyncTaskPath = strings.ReplaceAll(stopSyncTaskPath, "{sync_task_id}", d.Id())

	stopSyncTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err := client.Request("POST", stopSyncTaskPath, &stopSyncTaskOpt)
	if err != nil {
		return fmt.Errorf("error stopping OMS migration sync task: %s", err)
	}

	return nil
}

func resourceMigrationSyncTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	queryTime := strconv.FormatInt(time.Now().UnixMilli(), 10)

	var (
		deleteSyncTaskHttpUrl = "v2/{project_id}/sync-tasks/{sync_task_id}"
		deleteSyncTaskProduct = "oms"
	)
	deleteSyncTaskClient, err := cfg.NewServiceClient(deleteSyncTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	deleteSyncTaskPath := deleteSyncTaskClient.Endpoint + deleteSyncTaskHttpUrl
	deleteSyncTaskPath = strings.ReplaceAll(deleteSyncTaskPath, "{project_id}", deleteSyncTaskClient.ProjectID)
	deleteSyncTaskPath = strings.ReplaceAll(deleteSyncTaskPath, "{sync_task_id}", d.Id())

	deleteSyncTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getSyncTaskResp, err := deleteSyncTaskClient.Request("GET", fmt.Sprintf("%s?query_time=%s", deleteSyncTaskPath, queryTime), &deleteSyncTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OMS migration sync task")
	}

	getSyncTaskRespBody, err := utils.FlattenResponse(getSyncTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	status := utils.PathSearch("status", getSyncTaskRespBody, "").(string)

	// Unable to delete running synchronization task, please stop it first.
	if status == "SYNCHRONIZING" {
		err := stopSyncTask(deleteSyncTaskClient, d)
		if err != nil {
			return diag.Errorf("error stopping OMS migration sync task: %s", err)
		}
	}

	_, err = deleteSyncTaskClient.Request("DELETE", deleteSyncTaskPath, &deleteSyncTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting OMS migration sync task")
	}

	return nil
}
