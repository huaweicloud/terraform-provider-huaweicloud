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

var alarmClearNonUpdatableParams = []string{"alarm_ids", "remarks", "is_service_interrupt", "start_time",
	"fault_recovery_time"}

// @API COC POST /v1/alarm-mgmt/alarms/cancel
func ResourceAlarmClear() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmClearCreate,
		ReadContext:   resourceAlarmClearRead,
		UpdateContext: resourceAlarmClearUpdate,
		DeleteContext: resourceAlarmClearDelete,

		CustomizeDiff: config.FlexibleForceNew(alarmClearNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"alarm_ids": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remarks": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_service_interrupt": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"fault_recovery_time": {
				Type:     schema.TypeInt,
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

func buildAlarmClearCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_ids":            d.Get("alarm_ids"),
		"remarks":              utils.ValueIgnoreEmpty(d.Get("remarks")),
		"is_service_interrupt": utils.ValueIgnoreEmpty(d.Get("is_service_interrupt")),
		"start_time":           utils.ValueIgnoreEmpty(d.Get("start_time")),
		"fault_recovery_time":  utils.ValueIgnoreEmpty(d.Get("fault_recovery_time")),
	}

	return bodyParams
}

func resourceAlarmClearCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/alarm-mgmt/alarms/cancel"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAlarmClearCreateOpts(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC alarm clear: %s", err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	return nil
}

func resourceAlarmClearRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlarmClearUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAlarmClearDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting alarm clear resource is not supported. The alarm clear resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
