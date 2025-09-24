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

var ticketActionNonUpdatableParams = []string{"ticket_type", "user_id", "ticket_id", "task_id", "action", "params"}

// @API COC POST /v1/{ticket_type}/actions
func ResourceTicketAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTicketActionCreate,
		ReadContext:   resourceTicketActionRead,
		UpdateContext: resourceTicketActionUpdate,
		DeleteContext: resourceTicketActionDelete,

		CustomizeDiff: config.FlexibleForceNew(ticketActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"ticket_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ticket_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"params": {
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

func buildTicketActionCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ticket_id": d.Get("ticket_id"),
		"user_id":   utils.ValueIgnoreEmpty(d.Get("user_id")),
		"task_id":   utils.ValueIgnoreEmpty(d.Get("task_id")),
		"action":    utils.ValueIgnoreEmpty(d.Get("action")),
		"params":    parseJson(d.Get("params").(string)),
	}

	return bodyParams
}

func resourceTicketActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/{ticket_type}/actions"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{ticket_type}", d.Get("ticket_type").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildTicketActionCreateOpts(d),
	}

	ticketID := d.Get("ticket_id").(string)

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error requesting COC ticket action (%s): %s", ticketID, err)
	}

	d.SetId(ticketID)

	return nil
}

func resourceTicketActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTicketActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTicketActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ticket action resource is not supported. The ticket action resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
