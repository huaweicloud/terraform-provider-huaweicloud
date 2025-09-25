package coc

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var alarmLinkedIncidentNonUpdatableParams = []string{"alarm_ids", "current_cloud_service_id", "description",
	"is_service_interrupt", "level_id", "mtm_type", "title", "enterprise_project_id", "assignee", "assignee_role",
	"assignee_scene", "attachment", "is_change_event", "mtm_region", "source_id"}

// @API COC POST /v1/alarm-mgmt/alarms-linked-incident
func ResourceAlarmLinkedIncident() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmLinkedIncidentCreate,
		ReadContext:   resourceAlarmLinkedIncidentRead,
		UpdateContext: resourceAlarmLinkedIncidentUpdate,
		DeleteContext: resourceAlarmLinkedIncidentDelete,

		CustomizeDiff: config.FlexibleForceNew(alarmLinkedIncidentNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"alarm_ids": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"current_cloud_service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_service_interrupt": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"level_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mtm_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assignee": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assignee_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"assignee_scene": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"attachment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_change_event": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"mtm_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
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

func buildAlarmLinkedIncidentCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_ids":                d.Get("alarm_ids"),
		"current_cloud_service_id": d.Get("current_cloud_service_id"),
		"description":              d.Get("description"),
		"is_service_interrupt":     d.Get("is_service_interrupt"),
		"level_id":                 d.Get("level_id"),
		"mtm_type":                 d.Get("mtm_type"),
		"title":                    d.Get("title"),
		"enterprise_project_id":    utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"assignee":                 utils.ValueIgnoreEmpty(d.Get("assignee")),
		"assignee_role":            utils.ValueIgnoreEmpty(d.Get("assignee_role")),
		"assignee_scene":           utils.ValueIgnoreEmpty(d.Get("assignee_scene")),
		"attachment":               utils.ValueIgnoreEmpty(d.Get("attachment")),
		"is_change_event":          utils.ValueIgnoreEmpty(d.Get("is_change_event")),
		"mtm_region":               utils.ValueIgnoreEmpty(d.Get("mtm_region")),
		"source_id":                utils.ValueIgnoreEmpty(d.Get("source_id")),
	}

	return bodyParams
}

func resourceAlarmLinkedIncidentCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/alarm-mgmt/alarms-linked-incident"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAlarmLinkedIncidentCreateOpts(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC alarm linked incident: %s", err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	return nil
}

func resourceAlarmLinkedIncidentRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlarmLinkedIncidentUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlarmLinkedIncidentDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting alarm linked incident resource is not supported. The alarm linked incident resource is only" +
		" removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
