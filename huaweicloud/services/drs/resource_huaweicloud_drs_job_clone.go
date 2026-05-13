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

var cloneJobNonUpdatableParams = []string{"job_id", "name", "task_version"}

// @API DRS POST /v5/{project_id}/jobs/clone
func ResourceDrsJobClone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceJobCloneCreate,
		ReadContext:   resourceJobCloneRead,
		UpdateContext: resourceJobCloneUpdate,
		DeleteContext: resourceJobCloneDelete,

		CustomizeDiff: config.FlexibleForceNew(cloneJobNonUpdatableParams),

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_clone_job": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCloneJobBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"job_id":       d.Get("job_id"),
		"name":         d.Get("name"),
		"task_version": d.Get("task_version"),
	}
	return utils.RemoveNil(bodyParams)
}

func resourceJobCloneCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/clone"
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
		JSONBody:         utils.RemoveNil(buildCloneJobBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error cloning DRS job: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	clonedJobId := utils.PathSearch("id", respBody, "").(string)
	if clonedJobId == "" {
		return diag.Errorf("unable to find the ID from the API response")
	}
	d.SetId(clonedJobId)

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("is_clone_job", utils.PathSearch("is_clone_job", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceJobCloneRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceJobCloneUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceJobCloneDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to clone a DRS job. Deleting this resource will not 
delete the cloned job from the cloud, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
