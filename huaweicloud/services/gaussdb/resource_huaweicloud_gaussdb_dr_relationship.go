package gaussdb

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

var gaussDbDrRelationshipNonUpdatableParams = []string{
	"instance_id",
	"disaster_type",
	"dr_ip",
	"dr_user_name",
	"dr_user_password",
	"dr_task_name",
}

// @API GaussDB POST /v3.5/{project_id}/instances/{instance_id}/disaster-recovery/construct
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3.5/{project_id}/disaster-recovery/relations
// @API GaussDB POST /v3.5/{project_id}/instances/{instance_id}/disaster-recovery/release
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceGaussDbDrRelationship() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbDrRelationshipCreate,
		ReadContext:   resourceGaussDbDrRelationshipRead,
		UpdateContext: resourceGaussDbDrRelationshipUpdate,
		DeleteContext: resourceGaussDbDrRelationshipDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussDbDrRelationshipNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"disaster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dr_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dr_user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dr_user_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"dr_task_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"synchronization_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dr_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"precheck_failed_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disaster_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slave_region_instance_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbDrRelationshipRegionInstanceInfoSchema(),
			},
			"master_region_instance_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbDrRelationshipRegionInstanceInfoSchema(),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func gaussDbDrRelationshipRegionInstanceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"region_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGaussDbDrRelationshipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/instances/{instance_id}/disaster-recovery/construct"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDbDrRelationshipBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating GaussDB DR relationship: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating GaussDB DR relationship: job_id is not found in API response")
	}

	d.SetId(d.Get("instance_id").(string))

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error creating GaussDB DR relationship: %s", err)
	}

	return resourceGaussDbDrRelationshipRead(ctx, d, meta)
}

func buildCreateGaussDbDrRelationshipBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"disaster_type":    d.Get("disaster_type"),
		"dr_ip":            d.Get("dr_ip"),
		"dr_user_name":     d.Get("dr_user_name"),
		"dr_user_password": d.Get("dr_user_password"),
		"dr_task_name":     utils.ValueIgnoreEmpty(d.Get("dr_task_name")),
	}
	return bodyParams
}

func resourceGaussDbDrRelationshipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/disaster-recovery/relations"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildGetGaussDbDrRelationshipQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB DR relationship")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	relationship := utils.PathSearch(fmt.Sprintf("relations[?instance_id=='%s']|[0]", d.Id()), getRespBody, nil)
	if relationship == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB DR relationship")
	}

	status := utils.PathSearch("status", relationship, "").(string)
	if status == "completed" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving GaussDB DR relationship")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", relationship, nil)),
		d.Set("dr_task_name", utils.PathSearch("name", relationship, nil)),
		d.Set("synchronization_id", utils.PathSearch("synchronization_id", relationship, nil)),
		d.Set("dr_id", utils.PathSearch("id", relationship, nil)),
		d.Set("status", status),
		d.Set("precheck_failed_reason", utils.PathSearch("precheck_failed_reason", relationship, nil)),
		d.Set("disaster_type", utils.PathSearch("disaster_type", relationship, nil)),
		d.Set("disaster_role", utils.PathSearch("disaster_role", relationship, nil)),
		d.Set("created", utils.PathSearch("created", relationship, nil)),
		d.Set("updated", utils.PathSearch("updated", relationship, nil)),
		d.Set("instance_name", utils.PathSearch("instance_name", relationship, nil)),
		d.Set("instance_status", utils.PathSearch("instance_status", relationship, nil)),
		d.Set("actions", utils.PathSearch("actions", relationship, nil)),
		d.Set("slave_region_instance_info", flattenGaussDbDrRelationshipRegionInstanceInfo(
			utils.PathSearch("slave_region_instance_info", relationship, nil))),
		d.Set("master_region_instance_info", flattenGaussDbDrRelationshipRegionInstanceInfo(
			utils.PathSearch("master_region_instance_info", relationship, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetGaussDbDrRelationshipQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?instance_id=%s", d.Id())
}

func flattenGaussDbDrRelationshipRegionInstanceInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"region_code":  utils.PathSearch("region_code", resp, nil),
			"instance_id":  utils.PathSearch("instance_id", resp, nil),
			"project_id":   utils.PathSearch("project_id", resp, nil),
			"project_name": utils.PathSearch("project_name", resp, nil),
			"ip_address":   utils.PathSearch("ip_address", resp, nil),
		},
	}
}

func resourceGaussDbDrRelationshipUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbDrRelationshipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/instances/{instance_id}/disaster-recovery/release"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDbDrRelationshipBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", deletePath, &deleteOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	deleteResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceStateRefreshFunc(client, d.Get("instance_id").(string)),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.06020105"),
			"error deleting GaussDB DR relationship")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting GaussDB DR relationship: job_id is not found in API response")
	}

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error deleting GaussDB DR relationship: %s", err)
	}

	return nil
}

func buildDeleteGaussDbDrRelationshipBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"disaster_type": d.Get("disaster_type"),
	}
	return bodyParams
}
