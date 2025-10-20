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

var scheduledEventAcceptNonUpdatableParams = []string{"event_id", "not_before"}

// @API ECS POST /v3/{project_id}/instance-scheduled-events/{id}/actions/accept
func ResourceComputeScheduledEventAccept() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeScheduledEventAcceptCreate,
		ReadContext:   resourceComputeScheduledEventAcceptRead,
		UpdateContext: resourceComputeScheduledEventAcceptUpdate,
		DeleteContext: resourceComputeScheduledEventAcceptDelete,

		CustomizeDiff: config.FlexibleForceNew(scheduledEventAcceptNonUpdatableParams),

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

func resourceComputeScheduledEventAcceptCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instance-scheduled-events/{id}/actions/accept"
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
	createOpt.JSONBody = utils.RemoveNil(buildScheduledEventAcceptBodyParams(d))

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ECS scheduled event accept: %s", err)
	}

	d.SetId(d.Get("event_id").(string))

	return nil
}

func buildScheduledEventAcceptBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"not_before": utils.ValueIgnoreEmpty(d.Get("not_before")),
	}
	return bodyParams
}

func resourceComputeScheduledEventAcceptRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeScheduledEventAcceptUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComputeScheduledEventAcceptDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ECS scheduled event accept resource is not supported. The resource is only removed from the" +
		" state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
