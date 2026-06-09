package gaussdb

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

var gaussDbDrDrillNonUpdatableParams = []string{
	"instance_id",
	"disaster_type",
	"xlog_keep_ratio",
}

// @API GaussDB POST /v3.5/{project_id}/instances/{instance_id}/disaster-recovery/simulation-start
// @API GaussDB POST /v3.5/{project_id}/instances/{instance_id}/disaster-recovery/simulation-stop
// @API GaussDB GET /v3/{project_id}/instances
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceGaussDbDrDrill() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbDrDrillCreate,
		ReadContext:   resourceGaussDbDrDrillRead,
		UpdateContext: resourceGaussDbDrDrillUpdate,
		DeleteContext: resourceGaussDbDrDrillDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussDbDrDrillNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			"xlog_keep_ratio": {
				Type:     schema.TypeInt,
				Optional: true,
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

func resourceGaussDbDrDrillCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/instances/{instance_id}/disaster-recovery/simulation-start"
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
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDbDrDrillBodyParams(d))

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
		return diag.Errorf("error creating GaussDB DR drill: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error creating GaussDB DR drill: job_id is not found in API response")
	}

	d.SetId(d.Get("instance_id").(string))

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error creating GaussDB DR drill: %s", err)
	}

	return resourceGaussDbDrDrillRead(ctx, d, meta)
}

func buildCreateGaussDbDrDrillBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"disaster_type":   d.Get("disaster_type"),
		"xlog_keep_ratio": utils.ValueIgnoreEmpty(d.Get("xlog_keep_ratio")),
	}

	return bodyParams
}

func resourceGaussDbDrDrillRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbDrDrillUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbDrDrillDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/instances/{instance_id}/disaster-recovery/simulation-stop"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDbDrDrillBodyParams(d))

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
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting GaussDB DR drill: %s", err)
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, "").(string)
	if jobId == "" {
		return diag.Errorf("error deleting GaussDB DR drill: job_id is not found in API response")
	}

	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 2, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error deleting GaussDB DR drill: %s", err)
	}

	return nil
}

func buildDeleteGaussDbDrDrillBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"disaster_type": d.Get("disaster_type"),
	}

	return bodyParams
}
