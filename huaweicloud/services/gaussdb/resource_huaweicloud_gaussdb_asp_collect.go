package gaussdb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var gaussDbAspCollectNonUpdatableParams = []string{
	"instance_id",
	"start_time",
	"end_time",
}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/asp/collect
// @API GaussDB GET /v3/{project_id}/jobs
func ResourceGaussDbAspCollect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbAspCollectCreate,
		UpdateContext: resourceGaussDbAspCollectUpdate,
		ReadContext:   resourceGaussDbAspCollectRead,
		DeleteContext: resourceGaussDbAspCollectDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussDbAspCollectNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
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
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceGaussDbAspCollectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/asp/collect"
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
	createOpt.JSONBody = utils.RemoveNil(buildCreateAspCollectBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB ASP collect: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, "").(string)
	if jobId != "" {
		if err = checkGaussDBOpenGaussJobFinish(ctx, client, jobId, 10, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error creating GaussDB ASP collect: %s", err)
		}
	}

	d.SetId(d.Get("instance_id").(string))

	return nil
}

func buildCreateAspCollectBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time": d.Get("start_time"),
		"end_time":   d.Get("end_time"),
	}
	return bodyParams
}

func resourceGaussDbAspCollectRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbAspCollectUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbAspCollectDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB ASP collect resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
