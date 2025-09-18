package coc

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

var incidentActionNonUpdatableParams = []string{"incident_id", "task_id", "action", "params"}

// @API COC POST /v2/incidents/{incident_id}/actions
func ResourceIncidentAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentActionCreate,
		ReadContext:   resourceIncidentActionRead,
		UpdateContext: resourceIncidentActionUpdate,
		DeleteContext: resourceIncidentActionDelete,

		CustomizeDiff: config.FlexibleForceNew(incidentActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"incident_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
			},
			"params": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

func buildIncidentActionCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"incident_id": d.Get("incident_id"),
		"task_id":     d.Get("task_id"),
		"action":      d.Get("action"),
		"params":      utils.ValueIgnoreEmpty(d.Get("params")),
	}

	return bodyParams
}

func resourceIncidentActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	incidentID := d.Get("incident_id").(string)
	httpUrl := "v2/incidents/{incident_id}/actions"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{incident_id}", incidentID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildIncidentActionCreateOpts(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC incident action: %s", err)
	}

	d.SetId(incidentID)

	return nil
}

func resourceIncidentActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIncidentActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIncidentActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting incident action resource is not supported. The incident action resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
