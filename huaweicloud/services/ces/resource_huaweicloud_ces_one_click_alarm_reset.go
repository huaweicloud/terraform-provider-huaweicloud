package ces

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

var oneClickAlarmResetNonUpdatableFields = []string{
	"one_click_alarm_id",
}

// @API CES POST /v2/{project_id}/one-click-alarms
func ResourceOneClickAlarmReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOneClickAlarmResetCreate,
		UpdateContext: resourceOneClickAlarmResetUpdate,
		ReadContext:   resourceOneClickAlarmResetRead,
		DeleteContext: resourceOneClickAlarmResetDelete,

		CustomizeDiff: config.FlexibleForceNew(oneClickAlarmResetNonUpdatableFields),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"one_click_alarm_id": {
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

func resourceOneClickAlarmResetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/one-click-alarms"
		product = "ces"
	)
	createClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createPath := createClient.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildOneClickAlarmResetBodyParams(d)),
	}

	_, err = createClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CES reset alarm rules: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	return nil
}

func buildOneClickAlarmResetBodyParams(d *schema.ResourceData) map[string]interface{} {
	param := map[string]interface{}{
		"one_click_alarm_id":   d.Get("one_click_alarm_id"),
		"is_reset":             true,
		"notification_enabled": false,
		"dimension_names": map[string]interface{}{
			"metric": [1]string{"default"},
			"event":  []string{},
		},
	}

	return param
}

func resourceOneClickAlarmResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOneClickAlarmResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOneClickAlarmResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting the reset alarm rules for one service in one-click monitoring resource is not supported." +
		" The reset alarm rules for one service in one-click monitoring resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
