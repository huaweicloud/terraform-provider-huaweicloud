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

var drsSubscriptionInfoNonUpdatableParams = []string{
	"job_id",
	"name",
	"description",
	"consume_time",
}

// @API DRS PUT /v5/{project_id}/subscriptions/{job_id}/info
func ResourceDrsSubscriptionInfo() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsSubscriptionInfoCreate,
		ReadContext:   resourceDrsSubscriptionInfoRead,
		UpdateContext: resourceDrsSubscriptionInfoUpdate,
		DeleteContext: resourceDrsSubscriptionInfoDelete,

		CustomizeDiff: config.FlexibleForceNew(drsSubscriptionInfoNonUpdatableParams),

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
				Description: "Specifies the ID of the DRS subscription task.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the name of the subscription task.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the description of the subscription task.",
			},
			"consume_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Specifies the consumption time point, in timestamp format.",
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

func resourceDrsSubscriptionInfoCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/subscriptions/{job_id}/info"
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
	createOpt.JSONBody = utils.RemoveNil(buildCreateDrsSubscriptionInfoBodyParams(d))

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating DRS subscription info: %s", err)
	}

	_, err = utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The API returns an empty response body, use job_id as the resource ID
	d.SetId(d.Get("job_id").(string))

	return resourceDrsSubscriptionInfoRead(ctx, d, meta)
}

func buildCreateDrsSubscriptionInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         utils.ValueIgnoreEmpty(d.Get("name")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"consume_time": utils.ValueIgnoreEmpty(d.Get("consume_time")),
	}
	return bodyParams
}

func resourceDrsSubscriptionInfoRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsSubscriptionInfoUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsSubscriptionInfoDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS subscription info resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
