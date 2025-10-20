package ecs

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

var scheduledEventUpdateNonUpdatableParams = []string{"event_id", "not_before"}

// @API ECS PUT /v3/{project_id}/instance-scheduled-events/{id}
func ResourceComputeScheduledEventUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeScheduledEventUpdateCreate,
		ReadContext:   resourceComputeScheduledEventUpdateRead,
		UpdateContext: resourceComputeScheduledEventUpdateUpdate,
		DeleteContext: resourceComputeScheduledEventUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(scheduledEventUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"event_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"not_before": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceComputeScheduledEventUpdateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instance-scheduled-events/{id}"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{id}", d.Get("event_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildScheduledEventUpdateBodyParams(d))

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ECS scheduled event update: %s", err)
	}

	d.SetId(d.Get("event_id").(string))

	return nil
}

func buildScheduledEventUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"not_before": d.Get("not_before"),
	}
	return bodyParams
}

func resourceComputeScheduledEventUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeScheduledEventUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeScheduledEventUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ECS scheduled event update resource is not supported. The resource is only removed from the" +
		" state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
