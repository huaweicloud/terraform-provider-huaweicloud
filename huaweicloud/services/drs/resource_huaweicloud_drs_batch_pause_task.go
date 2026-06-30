package drs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchPauseTaskNonUpdatableParams = []string{
	"jobs",
	"jobs.*.job_id",
	"jobs.*.pause_mode",
}

// @API DRS POST /v3/{project_id}/jobs/batch-pause-task
func ResourceBatchPauseTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchPauseTaskCreate,
		ReadContext:   resourceBatchPauseTaskRead,
		UpdateContext: resourceBatchPauseTaskUpdate,
		DeleteContext: resourceBatchPauseTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(batchPauseTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"pause_mode": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildPauseTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	jobsRaw := d.Get("jobs").([]interface{})
	jobs := make([]map[string]interface{}, 0, len(jobsRaw))

	for _, jobRaw := range jobsRaw {
		job, ok := jobRaw.(map[string]interface{})
		if !ok {
			continue
		}
		jobs = append(jobs, map[string]interface{}{
			"job_id":     job["job_id"],
			"pause_mode": job["pause_mode"],
		})
	}

	return map[string]interface{}{
		"jobs": jobs,
	}
}

func resourceBatchPauseTaskCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/jobs/batch-pause-task"
	)

	client, err := cfg.NewServiceClient("drs", region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildPauseTaskBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error pausing DRS tasks: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	results := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
	if len(results) == 0 {
		return diag.Errorf("unable to find the results from the API response")
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	mErr := multierror.Append(nil,
		d.Set("results", flattenPauseTaskResults(results)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting pause DRS task fields: %s", mErr)
	}

	return nil
}

func flattenPauseTaskResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, result := range results {
		rst = append(rst, map[string]interface{}{
			"id":     utils.PathSearch("id", result, nil),
			"status": utils.PathSearch("status", result, nil),
		})
	}
	return rst
}

func resourceBatchPauseTaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchPauseTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchPauseTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to pause DRS tasks. Deleting this resource will not
undo the pause operation, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
