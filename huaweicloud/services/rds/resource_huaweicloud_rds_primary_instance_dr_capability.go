package rds

import (
	"context"
	"fmt"
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

var primaryInstanceDrCapabilityNonUpdatableParams = []string{
	"instance_id", "target_instance_id", "target_project_id", "target_region", "target_ip", "target_subnet",
}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/action
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
// @API RDS POST /v3/{project_id}/instances/disaster-recovery-infos
// @API VPC GET /v1/{project_id}/subnets/{subnet_id}
// @API IAM GET /v3/projects
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/delete-disaster-recovery
func ResourcePrimaryInstanceDrCapability() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrimaryInstanceDrCapabilityCreate,
		ReadContext:   resourcePrimaryInstanceDrCapabilityRead,
		UpdateContext: resourcePrimaryInstanceDrCapabilityUpdate,
		DeleteContext: resourcePrimaryInstanceDrCapabilityDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(primaryInstanceDrCapabilityNonUpdatableParams),

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
			"target_subnet": {
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
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"build_process": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePrimaryInstanceDrCapabilityCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	createOpt.JSONBody = utils.RemoveNil(buildCreatePrimaryInstanceDrCapabilityBodyParams(d))
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
		return diag.Errorf("error creating RDS primary instance DR capability: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating RDS primary instance DR capability: job_id is not found in API response")
	}

	res, err := getInstanceDrCapability(client, "", instanceId, d.Get("target_instance_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	id := utils.PathSearch("instance_dr_infos|[0].id", res, "").(string)
	if id == "" {
		return diag.Errorf("error creating RDS primary instance DR capability: ID is not found")
	}
	d.SetId(id)

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePrimaryInstanceDrCapabilityRead(ctx, d, meta)
}

func buildCreatePrimaryInstanceDrCapabilityBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"target_instance_id": d.Get("target_instance_id"),
		"target_project_id":  d.Get("target_project_id"),
		"target_region":      d.Get("target_region"),
		"target_ip":          d.Get("target_ip"),
		"target_subnet":      d.Get("target_subnet"),
	}
	bodyParams := map[string]interface{}{
		"build_master_dr_relation": params,
	}
	return bodyParams
}

func resourcePrimaryInstanceDrCapabilityRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, "error retrieving RDS primary instance DR capability")
	}

	instanceDrInfo := utils.PathSearch("instance_dr_infos|[0]", getRespBody, nil)
	if instanceDrInfo == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving RDS primary instance DR capability")
	}

	targetInstanceID := utils.PathSearch("slave_instance_id", instanceDrInfo, "").(string)
	if targetInstanceID == "" {
		return diag.Errorf("error getting slave instance ID")
	}

	targetRegion := utils.PathSearch("slave_region", instanceDrInfo, "").(string)
	if targetRegion == "" {
		return diag.Errorf("error getting slave region")
	}

	instance, err := getTargetInstanceById(cfg, targetInstanceID, targetRegion)
	if err != nil {
		return diag.Errorf("error getting RDS instance: %s", err)
	}

	subnetId := utils.PathSearch("instances|[0].subnet_id", instance, "").(string)
	if subnetId == "" {
		return diag.Errorf("error getting subnet ID form instance(%s): %s", targetInstanceID, err)
	}
	subnet, err := getTargetSubnetById(cfg, subnetId, targetRegion)
	if err != nil {
		return diag.Errorf("error getting VPC subnet: %s", err)
	}

	project, err := getProjectIdByRegion(d, cfg, targetRegion)
	if err != nil {
		return diag.Errorf("error getting IAM project: %s", err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("master_instance_id", instanceDrInfo, nil)),
		d.Set("target_instance_id", targetInstanceID),
		d.Set("target_project_id", utils.PathSearch("projects|[0].id", project, nil)),
		d.Set("target_region", targetRegion),
		d.Set("target_ip", utils.PathSearch("instances|[0].private_ips[0]", instance, nil)),
		d.Set("target_subnet", utils.PathSearch("subnet.cidr", subnet, nil)),
		d.Set("status", utils.PathSearch("status", instanceDrInfo, nil)),
		d.Set("time", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("time", instanceDrInfo, float64(0)).(float64))/1000, false)),
		d.Set("build_process", utils.PathSearch("build_process", instanceDrInfo, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getInstanceDrCapability(client *golangsdk.ServiceClient, id, masterInstanceId, slaveInstanceId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/disaster-recovery-infos"
	)
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = utils.RemoveNil(buildGetInstanceDrCapabilityQueryParams(id, masterInstanceId, slaveInstanceId))

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func buildGetInstanceDrCapabilityQueryParams(id, masterInstanceId, slaveInstanceId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":                 utils.ValueIgnoreEmpty(id),
		"master_instance_id": utils.ValueIgnoreEmpty(masterInstanceId),
		"slave_instance_id":  utils.ValueIgnoreEmpty(slaveInstanceId),
	}
	return bodyParams
}

func getTargetInstanceById(cfg *config.Config, id, region string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances?id={id}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", id)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func getTargetSubnetById(cfg *config.Config, id, region string) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/subnets/{subnet_id}"
		product = "vpc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{subnet_id}", id)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func getProjectIdByRegion(d *schema.ResourceData, cfg *config.Config, name string) (interface{}, error) {
	var (
		httpUrl = "v3/projects?name={name}"
		product = "identity"
	)
	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating IAM client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{name}", name)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourcePrimaryInstanceDrCapabilityUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePrimaryInstanceDrCapabilityDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	deleteOpt.JSONBody = buildDeletePrimaryInstanceDrCapabilityBodyParams(d)

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.01010244"),
			"error deleting RDS primary instance DR capability")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error deleting RDS primary instance DR capability: job_id is not found in API response")
	}

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeletePrimaryInstanceDrCapabilityBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"target_instance_id": d.Get("target_instance_id"),
		"target_project_id":  d.Get("target_project_id"),
		"target_region":      d.Get("target_region"),
		"target_ip":          d.Get("target_ip"),
		"is_master":          true,
	}
	return bodyParams
}
