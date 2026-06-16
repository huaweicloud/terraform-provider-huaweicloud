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

var batchDeleteJobsNonUpdatableParams = []string{"jobs"}

// @API DRS DELETE /v5/{project_id}/jobs
func ResourceBatchDeleteJobs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchDeleteJobsCreate,
		ReadContext:   resourceBatchDeleteJobsRead,
		UpdateContext: resourceBatchDeleteJobsUpdate,
		DeleteContext: resourceBatchDeleteJobsDelete,

		CustomizeDiff: config.FlexibleForceNew(batchDeleteJobsNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
					},
				},
			},
		},
	}
}

func buildBatchDeleteJobsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"jobs": utils.ExpandToStringList(d.Get("jobs").([]interface{})),
	}
}

func resourceBatchDeleteJobsCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs"
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
		JSONBody:         buildBatchDeleteJobsBodyParams(d),
	}

	resp, err := client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch deleting DRS jobs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId.String())

	mErr := multierror.Append(
		nil,
		d.Set("results", flattenBatchDeleteJobsResults(respBody)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting DRS batch delete jobs fields: %s", mErr)
	}

	return nil
}

func flattenBatchDeleteJobsResults(respBody interface{}) []interface{} {
	jobsResp := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
	if len(jobsResp) == 0 {
		return nil
	}

	results := make([]interface{}, 0, len(jobsResp))
	for _, job := range jobsResp {
		results = append(results, map[string]interface{}{
			"id":     utils.PathSearch("id", job, nil),
			"name":   utils.PathSearch("name", job, nil),
			"status": utils.PathSearch("status", job, nil),
		})
	}

	return results
}

func resourceBatchDeleteJobsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteJobsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteJobsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch delete DRS jobs. Deleting this resource
    will not restore the deleted jobs or undo the delete action, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
