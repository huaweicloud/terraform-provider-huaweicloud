package drs

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/drs/v3/jobs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var startJobNonUpdatableParams = []string{"job_id", "start_time"}

// @API DRS POST /v3/{project_id}/jobs/batch-precheck
// @API DRS POST /v3/{project_id}/jobs/batch-precheck-result
// @API DRS POST /v3/{project_id}/jobs/batch-starting
// @API DRS POST /v3/{project_id}/jobs/batch-status
func ResourceStartJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStartJobCreate,
		ReadContext:   resourceStartJobRead,
		UpdateContext: resourceStartJobUpdate,
		DeleteContext: resourceStartJobDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(startJobNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
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

func resourceStartJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		jobId  = d.Get("job_id").(string)
	)

	client, err := cfg.DrsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v3 client: %s", err)
	}

	err = preCheck(ctx, client, jobId, d.Timeout(schema.TimeoutCreate), "forStartJob")
	if err != nil {
		return diag.FromErr(err)
	}

	startTime := d.Get("start_time").(string)
	startMode := "start"
	if startTime != "" && startTime != "0" {
		startMode = "start_later"
	}

	startReq := jobs.StartJobReq{
		Jobs: []jobs.StartInfo{
			{
				JobId:     jobId,
				StartTime: startTime,
			},
		},
	}
	_, err = jobs.Start(client, startReq)
	if err != nil {
		return diag.Errorf("Start DRS job failed: %s", err)
	}

	err = waitingforJobStatus(ctx, client, jobId, startMode, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(jobId)

	return nil
}

func resourceStartJobRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceStartJobUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceStartJobDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to start a DRS job. Deleting this resource will not
    stop the job or undo the start action, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
