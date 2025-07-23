package sdrs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SDRS DELETE /v1/{project_id}/task-center/{server_group_id}/failure-jobs/batch
func ResourceDeleteSpecifiedGroupFailureJobs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeleteSpecifiedGroupFailureJobsCreate,
		ReadContext:   resourceDeleteSpecifiedGroupFailureJobsRead,
		UpdateContext: resourceDeleteSpecifiedGroupFailureJobsUpdate,
		DeleteContext: resourceDeleteSpecifiedGroupFailureJobsDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"server_group_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"server_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the protected group to delete all failure jobs from.`,
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

func resourceDeleteSpecifiedGroupFailureJobsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/task-center/{server_group_id}/failure-jobs/batch"
		product = "sdrs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SDRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{server_group_id}", d.Get("server_group_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting all failure jobs from specified group: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	return resourceDeleteSpecifiedGroupFailureJobsRead(ctx, d, meta)
}

func resourceDeleteSpecifiedGroupFailureJobsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceDeleteSpecifiedGroupFailureJobsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceDeleteSpecifiedGroupFailureJobsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to delete all failure jobs from a specified group. Deleting this 
resource will not change the current configuration, but will only remove the resource information from the 
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
