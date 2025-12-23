package workspace

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var desktopPoolExpandNonUpdatableParams = []string{
	"pool_id",
	"size",
}

// @API Workspace POST /v2/{project_id}/desktop-pools/{pool_id}/expand
func ResourceDesktopPoolExpand() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopPoolExpandCreate,
		ReadContext:   resourceDesktopPoolExpandRead,
		UpdateContext: resourceDesktopPoolExpandUpdate,
		DeleteContext: resourceDesktopPoolExpandDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(desktopPoolExpandNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the desktop pool is located.`,
			},

			// Required parameters.
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the desktop pool to be expanded.`,
			},
			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The number of desktops to be added to the pool.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDesktopPoolExpandBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"size": d.Get("size"),
	}
}

func resourceDesktopPoolExpandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		poolID  = d.Get("pool_id").(string)
		httpUrl = "v2/{project_id}/desktop-pools/{pool_id}/expand"
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{pool_id}", poolID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDesktopPoolExpandBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error expanding desktop pool: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "")
	if jobId.(string) != "" {
		_, err = waitForJobCompleted(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceDesktopPoolExpandRead(ctx, d, meta)
}

func resourceDesktopPoolExpandRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This is a one-time action resource, so there's no need to read the resource state.
	return nil
}

func resourceDesktopPoolExpandUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This is a one-time action resource, so there's no need to update the resource.
	return nil
}

func resourceDesktopPoolExpandDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for expanding a desktop pool. Deleting this
    resource will not clear the corresponding request record, but will only remove the resource information
    from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
