package rds

import (
	"context"
	"net/http"
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

var drInstanceDrCapabilityNonUpdatableParams = []string{
	"instance_id", "target_instance_id", "target_project_id", "target_region", "target_ip",
}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/action
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
// @API RDS POST /v3/{project_id}/instances/disaster-recovery-infos
// @API IAM GET /v3/projects
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/delete-disaster-recovery
func ResourceDrInstanceDrCapability() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrInstanceDrCapabilityCreate,
		ReadContext:   resourceDrInstanceDrCapabilityRead,
		UpdateContext: resourceDrInstanceDrCapabilityUpdate,
		DeleteContext: resourceDrInstanceDrCapabilityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(drInstanceDrCapabilityNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"target_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_ip": {
				Type:     schema.TypeString,
				Required: true,
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
			"replica_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_process": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wal_receive_replay_delay_in_ms": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wal_write_receive_delay_in_mb": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wal_write_replay_delay_in_mb": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDrInstanceDrCapabilityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/action"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDrInstanceDrCapabilityBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: 2 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS DR instance DR capability: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating RDS DR instance DR capability: job_id is not found in API response")
	}

	res, err := getInstanceDrCapability(client, "", d.Get("target_instance_id").(string), instanceId)
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("instance_dr_infos|[0].id", res, "").(string)
	if id == "" {
		return diag.Errorf("error creating RDS DR instance DR capability: ID is not found")
	}
	d.SetId(id)

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDrInstanceDrCapabilityRead(ctx, d, meta)
}

func buildCreateDrInstanceDrCapabilityBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"target_instance_id": d.Get("target_instance_id"),
		"target_project_id":  d.Get("target_project_id"),
		"target_region":      d.Get("target_region"),
		"target_ip":          d.Get("target_ip"),
	}
	bodyParams := map[string]interface{}{
		"build_slave_dr_relation": params,
	}
	return bodyParams
}

func resourceDrInstanceDrCapabilityRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getRespBody, err := getInstanceDrCapability(client, d.Id(), "", "")
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS DR instance DR capability")
	}

	instanceDrInfo := utils.PathSearch("instance_dr_infos|[0]", getRespBody, nil)
	if instanceDrInfo == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving RDS DR instance DR capability")
	}

	targetInstanceID := utils.PathSearch("master_instance_id", instanceDrInfo, "").(string)
	if targetInstanceID == "" {
		return diag.Errorf("error getting master instance ID")
	}

	targetRegion := utils.PathSearch("master_region", instanceDrInfo, "").(string)
	if targetRegion == "" {
		return diag.Errorf("error getting master region")
	}

	instance, err := getTargetInstanceById(cfg, targetInstanceID, targetRegion)
	if err != nil {
		return diag.Errorf("error getting RDS instance: %s", err)
	}

	project, err := getProjectIdByRegion(d, cfg, targetRegion)
	if err != nil {
		return diag.Errorf("error getting IAM project: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("slave_instance_id", instanceDrInfo, nil)),
		d.Set("target_instance_id", targetInstanceID),
		d.Set("target_project_id", utils.PathSearch("projects|[0].id", project, nil)),
		d.Set("target_region", targetRegion),
		d.Set("target_ip", utils.PathSearch("instances|[0].private_ips[0]", instance, nil)),
		d.Set("status", utils.PathSearch("status", instanceDrInfo, nil)),
		d.Set("replica_state", utils.PathSearch("replica_state", instanceDrInfo, nil)),
		d.Set("time", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("time", instanceDrInfo, float64(0)).(float64))/1000, false)),
		d.Set("build_process", utils.PathSearch("build_process", instanceDrInfo, nil)),
		d.Set("wal_receive_replay_delay_in_ms", utils.PathSearch("wal_receive_replay_delay_in_ms",
			instanceDrInfo, nil)),
		d.Set("wal_write_receive_delay_in_mb", utils.PathSearch("wal_write_receive_delay_in_mb",
			instanceDrInfo, nil)),
		d.Set("wal_write_replay_delay_in_mb", utils.PathSearch("wal_write_replay_delay_in_mb",
			instanceDrInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDrInstanceDrCapabilityUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrInstanceDrCapabilityDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/delete-disaster-recovery"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = buildDeleteDrInstanceDrCapabilityBodyParams(d)

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.01010244"),
			"error deleting RDS DR instance DR capability")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error deleting RDS dr instance dr capability: job_id is not found in API response")
	}

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteDrInstanceDrCapabilityBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"target_instance_id": d.Get("target_instance_id"),
		"target_project_id":  d.Get("target_project_id"),
		"target_region":      d.Get("target_region"),
		"target_ip":          d.Get("target_ip"),
		"is_master":          false,
	}
	return bodyParams
}
