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

var redistributionParametersModifyNonUpdatableParams = []string{
	"instance_id",
	"redis_join_tables",
	"redis_parallel_jobs",
	"redis_resource_level",
}

// @API GaussDB PUT /v3/{project_id}/instances/{instance_id}/redistribution-parameters
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceGaussDBRedistributionParametersModify() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBRedistributionParametersModifyCreate,
		ReadContext:   resourceGaussDBRedistributionParametersModifyRead,
		UpdateContext: resourceGaussDBRedistributionParametersModifyUpdate,
		DeleteContext: resourceGaussDBRedistributionParametersModifyDelete,

		CustomizeDiff: config.FlexibleForceNew(redistributionParametersModifyNonUpdatableParams),

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
			"redis_join_tables": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			"redis_parallel_jobs": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"redis_resource_level": {
				Type:     schema.TypeString,
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

func resourceGaussDBRedistributionParametersModifyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/redistribution-parameters"
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
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBRedistributionParametersModifyBodyParams(d))

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("PUT", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	res, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error modifying GaussDB redistribution parameters: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(res.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("instance_id").(string))

	jobId := utils.PathSearch("result", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error modifying GaussDB redistribution parameters, result is not found in the response")
	}
	err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId.(string), 2, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildCreateGaussDBRedistributionParametersModifyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"redis_join_tables":    utils.ValueIgnoreEmpty(d.Get("redis_join_tables")),
		"redis_parallel_jobs":  utils.ValueIgnoreEmpty(d.Get("redis_parallel_jobs")),
		"redis_resource_level": utils.ValueIgnoreEmpty(d.Get("redis_resource_level")),
	}
	return bodyParams
}

func resourceGaussDBRedistributionParametersModifyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBRedistributionParametersModifyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDBRedistributionParametersModifyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB redistribution parameters modify resource is not supported. The resource is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
