package fgs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var FunctionTriggerStatusActionNonUpdatableParams = []string{
	"function_urn",
	"trigger_type_code",
	"trigger_id",
	"trigger_status",
}

// @API FunctionGraph PUT /v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}
func ResourceFunctionTriggerStatusAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFunctionTriggerStatusActionCreate,
		ReadContext:   resourceFunctionTriggerStatusActionRead,
		UpdateContext: resourceFunctionTriggerStatusActionUpdate,
		DeleteContext: resourceFunctionTriggerStatusActionDelete,

		CustomizeDiff: config.FlexibleForceNew(FunctionTriggerStatusActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the function trigger is located.`,
			},
			"function_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The URN of the function.`,
			},
			"trigger_type_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The trigger type code.`,
			},
			"trigger_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The trigger ID.`,
			},
			"trigger_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The status of the trigger. Valid values are ACTIVE and DISABLED.`,
			},
			"event_data": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The event data of the trigger.`,
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

func resourceFunctionTriggerStatusActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		httpUrl         = "v2/{project_id}/fgs/triggers/{function_urn}/{trigger_type_code}/{trigger_id}"
		functionUrn     = d.Get("function_urn").(string)
		triggerTypeCode = d.Get("trigger_type_code").(string)
		triggerId       = d.Get("trigger_id").(string)
		triggerStatus   = d.Get("trigger_status").(string)
	)

	client, err := cfg.NewServiceClient("fgs", region)
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{function_urn}", functionUrn)
	updatePath = strings.ReplaceAll(updatePath, "{trigger_type_code}", triggerTypeCode)
	updatePath = strings.ReplaceAll(updatePath, "{trigger_id}", triggerId)

	parsedEventData, err := parseEventDataAndDecryptSentisiveParams(ctx, meta, d, utils.StringToJson(d.Get("event_data").(string)))
	if err != nil {
		return diag.Errorf("error parsing event data: %s", err)
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"trigger_status": triggerStatus,
			"event_data":     parsedEventData,
		}),
	}

	_, err = client.Request("PUT", updatePath, &opt)
	if err != nil {
		return diag.Errorf("error updating function trigger status: %s", err)
	}

	d.SetId(triggerId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      functionTriggerStatusRefreshFunc(client, functionUrn, triggerTypeCode, triggerId, []string{"ACTIVE", "DISABLED"}),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		diag.Errorf("error waiting for the function trigger (%s) status to become expected value (%s): %s",
			triggerId, triggerStatus, err)
	}

	return resourceFunctionTriggerStatusActionRead(ctx, d, meta)
}

func resourceFunctionTriggerStatusActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFunctionTriggerStatusActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFunctionTriggerStatusActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for updating function trigger status. Deleting this resource will
not revert the trigger status, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
