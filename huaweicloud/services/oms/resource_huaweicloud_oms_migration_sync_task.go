package oms

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	omsmodel "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/oms/v2/model"

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
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.HcOmsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	createOpts, err := buildCreateOpts(d, region)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateSyncTask(&omsmodel.CreateSyncTaskRequest{Body: &createOpts})
	if err != nil {
		return diag.Errorf("error creating OMS migration sync task: %s", err)
	}
	if resp.SyncTaskId == nil {
		return diag.Errorf("unable to find the task ID")
	}

	syncTaskId := *resp.SyncTaskId
	d.SetId(syncTaskId)

	if d.Get("action").(string) == "stop" {
		stopOpts := omsmodel.StopSyncTaskRequest{
			SyncTaskId: d.Id(),
		}
		_, err := client.StopSyncTask(&stopOpts)
		if err != nil {
			return diag.Errorf("error stopping OMS migration sync task: %s", err)
		}
	}

	return resourceMigrationSyncTaskRead(ctx, d, meta)
}

func buildCreateOpts(d *schema.ResourceData, region string) (omsmodel.CreateSyncTaskReq, error) {
	createOpts := omsmodel.CreateSyncTaskReq{
		SrcCloudType:            utils.StringIgnoreEmpty(d.Get("src_cloud_type").(string)),
		SrcRegion:               d.Get("src_region").(string),
		SrcBucket:               d.Get("src_bucket").(string),
		SrcAk:                   d.Get("src_ak").(string),
		SrcSk:                   d.Get("src_sk").(string),
		DstRegion:               region,
		DstBucket:               d.Get("dst_bucket").(string),
		DstAk:                   d.Get("dst_ak").(string),
		DstSk:                   d.Get("dst_sk").(string),
		Description:             utils.StringIgnoreEmpty(d.Get("description").(string)),
		EnableKms:               utils.Bool(d.Get("enable_kms").(bool)),
		EnableRestore:           utils.Bool(d.Get("enable_restore").(bool)),
		EnableMetadataMigration: utils.Bool(d.Get("enable_metadata_migration").(bool)),
		AppId:                   utils.StringIgnoreEmpty(d.Get("app_id").(string)),
	}

	sourceCdn, err := buildSourceCdnOpts(d.Get("source_cdn").([]interface{}))
	if err != nil {
		return omsmodel.CreateSyncTaskReq{}, err
	}
	createOpts.SourceCdn = sourceCdn

	consistencyCheck := d.Get("consistency_check").(string)
	if consistencyCheck != "" {
		var consistencyCheckOpt omsmodel.CreateSyncTaskReqConsistencyCheck
		if err := consistencyCheckOpt.UnmarshalJSON([]byte(consistencyCheck)); err != nil {
			return omsmodel.CreateSyncTaskReq{}, fmt.Errorf("error parsing the argument consistency_check: %s", err)
		}
		createOpts.ConsistencyCheck = &consistencyCheckOpt
	}
	return createOpts, nil
}

func resourceMigrationSyncTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.HcOmsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}
	queryTime := strconv.FormatInt(time.Now().UnixMilli(), 10)
	resp, err := client.ShowSyncTask(&omsmodel.ShowSyncTaskRequest{SyncTaskId: d.Id(), QueryTime: queryTime})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving OMS migration sync task")
	}
	createTime := *resp.CreateTime
	lastStartTime := *resp.LastStartTime
	dstStoragePolicyValue := ""
	if resp.DstStoragePolicy != nil {
		dstStoragePolicyValue = resp.DstStoragePolicy.Value()
	}

	status := resp.Status.Value()
	actionValue := "start"
	if status == "STOPPED" {
		actionValue = "stop"
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("src_cloud_type", resp.SrcCloudType.Value()),
		d.Set("src_region", resp.SrcRegion),
		d.Set("src_bucket", resp.SrcBucket),
		d.Set("dst_bucket", resp.DstBucket),

		d.Set("description", resp.Description),
		d.Set("status", status),
		d.Set("enable_kms", resp.EnableKms),
		d.Set("enable_metadata_migration", resp.EnableMetadataMigration),
		d.Set("enable_restore", resp.EnableRestore),

		d.Set("created_at", utils.FormatTimeStampRFC3339(createTime/1000, false)),
		d.Set("last_start_at", utils.FormatTimeStampRFC3339(lastStartTime/1000, false)),
		d.Set("app_id", resp.AppId),
		d.Set("object_overwrite_mode", resp.ObjectOverwriteMode.Value()),
		d.Set("dst_storage_policy", dstStoragePolicyValue),
		d.Set("consistency_check", resp.ConsistencyCheck.Value()),

		d.Set("monthly_acceptance_request", resp.MonthlyAcceptanceRequest),
		d.Set("monthly_success_object", resp.MonthlySuccessObject),
		d.Set("monthly_failure_object", resp.MonthlyFailureObject),
		d.Set("monthly_skip_object", resp.MonthlySkipObject),
		d.Set("monthly_size", resp.MonthlySize),
		d.Set("action", actionValue),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting OMS migration sync task fields: %s", err)
	}

	return nil
}

func resourceMigrationSyncTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.HcOmsV2Client(region)

	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	if d.HasChange("action") {
		action := d.Get("action").(string)
		if action == "start" {
			startActionOpts := omsmodel.StartSyncTaskReq{
				SrcAk: d.Get("src_ak").(string),
				SrcSk: d.Get("src_sk").(string),
				DstAk: d.Get("dst_ak").(string),
				DstSk: d.Get("dst_sk").(string),
			}
			if sourceCDNs := d.Get("source_cdn").([]interface{}); len(sourceCDNs) > 0 {
				sourceCdn := sourceCDNs[0].(map[string]interface{})
				startActionOpts.SourceCdnAuthenticationKey = utils.String(sourceCdn["authentication_key"].(string))
			}
			_, err := client.StartSyncTask(&omsmodel.StartSyncTaskRequest{SyncTaskId: d.Id(), Body: &startActionOpts})
			if err != nil {
				return diag.Errorf("error starting OMS migration sync task: %s", err)
			}
		} else {
			stopOpts := omsmodel.StopSyncTaskRequest{
				SyncTaskId: d.Id(),
			}
			_, err := client.StopSyncTask(&stopOpts)
			if err != nil {
				return diag.Errorf("error stopping OMS migration sync task: %s", err)
			}
		}
	}

	return resourceMigrationSyncTaskRead(ctx, d, meta)
}

func resourceMigrationSyncTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.HcOmsV2Client(region)

	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}
	taskId := d.Id()
	queryTime := strconv.FormatInt(time.Now().UnixMilli(), 10)
	resp, err := client.ShowSyncTask(&omsmodel.ShowSyncTaskRequest{SyncTaskId: taskId, QueryTime: queryTime})
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return nil
		}
		return diag.Errorf("error retrieving OMS migration sync task: %s", err)
	}

	if resp.Status == nil {
		return diag.Errorf("unable to find the status OMS migration sync task: %s", taskId)
	}

	taskStatus := resp.Status.Value()
	// Unable to delete running synchronization task, please stop it first.
	if taskStatus == "SYNCHRONIZING" {
		stopOpts := omsmodel.StopSyncTaskRequest{
			SyncTaskId: d.Id(),
		}
		_, err := client.StopSyncTask(&stopOpts)
		if err != nil {
			return diag.Errorf("error stopping OMS migration sync task: %s", err)
		}
	}

	_, err = client.DeleteSyncTask(&omsmodel.DeleteSyncTaskRequest{SyncTaskId: taskId})
	if err != nil {
		return diag.Errorf("error deleting OMS migration sync task: %s", err)
	}

	return nil
}
