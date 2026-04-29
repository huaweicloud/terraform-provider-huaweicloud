package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var stopJobNonUpdatableParams = []string{"job_id", "is_force_stop"}

// @API DRS POST /v5/{project_id}/jobs/{job_id}/stop
func ResourceStopJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStopJobCreate,
		ReadContext:   resourceStopJobRead,
		UpdateContext: resourceStopJobUpdate,
		DeleteContext: resourceStopJobDelete,

		CustomizeDiff: config.FlexibleForceNew(stopJobNonUpdatableParams),

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
			"is_force_stop": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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
	}
}

func buildStopJobBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"is_force_stop": utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "is_force_stop"),
	}
}

func resourceStopJobCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/stop"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS v5 client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildStopJobBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error stopping DRS job: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("job_id").(string))

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting stop DRS job fields: %s", mErr)
	}

	return nil
}

func resourceStopJobRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceStopJobUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceStopJobDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to stop a DRS job. Deleting this resource will not
    restore the job or undo the stop action, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
