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

var jobsBatchStopNonUpdatableParams = []string{
	"jobs",
	"jobs.*.job_id",
	"jobs.*.is_force_stop",
}

// @API DRS POST /v5/{project_id}/jobs/batch-stop
func ResourceJobsBatchStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobsBatchStopCreate,
		ReadContext:   resourceJobsBatchStopRead,
		UpdateContext: resourceJobsBatchStopUpdate,
		DeleteContext: resourceJobsBatchStopDelete,

		CustomizeDiff: config.FlexibleForceNew(jobsBatchStopNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Elem:     jobsBatchStopSchema(),
				Required: true,
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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func jobsBatchStopSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_force_stop": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func buildJobsBatchStopBodyParams(d *schema.ResourceData) map[string]interface{} {
	jobsRaw := d.Get("jobs").([]interface{})
	jobs := make([]map[string]interface{}, 0, len(jobsRaw))

	for _, jobRaw := range jobsRaw {
		job, ok := jobRaw.(map[string]interface{})
		if !ok {
			continue
		}

		jobMap := map[string]interface{}{
			"job_id": job["job_id"],
		}

		if v, ok := job["is_force_stop"]; ok {
			jobMap["is_force_stop"] = v
		}

		jobs = append(jobs, jobMap)
	}

	return map[string]interface{}{
		"jobs": jobs,
	}
}

func resourceJobsBatchStopCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/batch-stop"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v5 client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildJobsBatchStopBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch stopping DRS jobs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	results := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
	if len(results) == 0 {
		return diag.Errorf("unable to find the results from the API response")
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	mErr := multierror.Append(nil,
		d.Set("results", flattenJobsBatchStopResults(results)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DRS batch stop job fields: %s", mErr)
	}

	return nil
}

func flattenJobsBatchStopResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, result := range results {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", result, nil),
			"name":       utils.PathSearch("name", result, nil),
			"status":     utils.PathSearch("status", result, nil),
			"error_code": utils.PathSearch("error_code", result, nil),
			"error_msg":  utils.PathSearch("error_msg", result, nil),
		})
	}
	return rst
}

func resourceJobsBatchStopRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceJobsBatchStopUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceJobsBatchStopDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch stop DRS jobs. Deleting this 
resource will not restore the stopped jobs or undo the stop action, but will only remove the resource information 
from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
