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

var batchSetDefinerNonUpdatableParams = []string{
	"jobs",
	"jobs.*.job_id",
	"jobs.*.replace_definer",
}

// @API DRS POST /v3/{project_id}/jobs/batch-replace-definer
func ResourceBatchSetDefiner() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchSetDefinerCreate,
		ReadContext:   resourceBatchSetDefinerRead,
		UpdateContext: resourceBatchSetDefinerUpdate,
		DeleteContext: resourceBatchSetDefinerDelete,

		CustomizeDiff: config.FlexibleForceNew(batchSetDefinerNonUpdatableParams),

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
				Elem:     batchSetDefinerJobsSchema(),
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
				Elem:     batchSetDefinerResultsSchema(),
			},
		},
	}
}

func batchSetDefinerJobsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"replace_definer": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func batchSetDefinerResultsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
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
	}
}

func buildBatchSetDefinerBodyParams(d *schema.ResourceData) map[string]interface{} {
	jobsRaw := d.Get("jobs").([]interface{})
	jobs := make([]map[string]interface{}, 0, len(jobsRaw))

	for _, jobRaw := range jobsRaw {
		job, ok := jobRaw.(map[string]interface{})
		if !ok {
			continue
		}

		jobs = append(jobs, map[string]interface{}{
			"job_id":          job["job_id"],
			"replace_definer": job["replace_definer"],
		})
	}

	return map[string]interface{}{
		"jobs": jobs,
	}
}

func resourceBatchSetDefinerCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/jobs/batch-replace-definer"
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
		JSONBody:         buildBatchSetDefinerBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch setting DRS job definer: %s", err)
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
		d.Set("results", flattenBatchSetDefinerResults(results)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting batch set DRS job definer fields: %s", mErr)
	}

	return nil
}

func flattenBatchSetDefinerResults(results []interface{}) []interface{} {
	if len(results) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(results))
	for _, result := range results {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", result, nil),
			"status":     utils.PathSearch("status", result, nil),
			"error_code": utils.PathSearch("error_code", result, nil),
			"error_msg":  utils.PathSearch("error_msg", result, nil),
		})
	}

	return rst
}

func resourceBatchSetDefinerRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchSetDefinerUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchSetDefinerDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch set DRS job definer. Deleting this resource
    will not restore the definer setting or undo the set action, but will only remove the resource information from the
    tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
