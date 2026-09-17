package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsSubscriptionResetNonUpdatableParams = []string{
	"job_id",
}

// @API DRS POST /v5/{project_id}/subscription/{job_id}/reset
func ResourceDrsSubscriptionReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsSubscriptionResetCreate,
		ReadContext:   resourceDrsSubscriptionResetRead,
		UpdateContext: resourceDrsSubscriptionResetUpdate,
		DeleteContext: resourceDrsSubscriptionResetDelete,

		CustomizeDiff: config.FlexibleForceNew(drsSubscriptionResetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the DRS subscription task to reset.",
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

func resourceDrsSubscriptionResetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/subscription/{job_id}/reset"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", d.Get("job_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error resetting DRS subscription: %s", err)
	}

	_, err = utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API returns an empty response body, use job_id as the resource ID
	d.SetId(d.Get("job_id").(string))

	return resourceDrsSubscriptionResetRead(ctx, d, meta)
}

func resourceDrsSubscriptionResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsSubscriptionResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsSubscriptionResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS subscription reset resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
