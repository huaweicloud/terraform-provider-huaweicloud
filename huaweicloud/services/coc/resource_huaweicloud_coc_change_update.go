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

var changeUpdateNonUpdatableParams = []string{"ticket_id", "phase", "work_flow_status", "action", "sub_tickets",
	"sub_tickets.*.ticket_id", "sub_tickets.*.change_result", "sub_tickets.*.is_verified_in_change_time",
	"sub_tickets.*.verified_docs", "sub_tickets.*.comment", "sub_tickets.*.change_fail_type",
	"sub_tickets.*.rollback_start_time", "sub_tickets.*.rollback_end_time",
	"sub_tickets.*.is_rollback_success", "sub_tickets.*.is_monitor_found"}

// @API COC PUT /v2/changes/{change_id}
func ResourceChangeUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChangeUpdateCreate,
		ReadContext:   resourceChangeUpdateRead,
		UpdateContext: resourceChangeUpdateUpdate,
		DeleteContext: resourceChangeUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(changeUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"ticket_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"phase": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"work_flow_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_tickets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ticket_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"change_result": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"is_verified_in_change_time": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"verified_docs": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"change_fail_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rollback_start_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"rollback_end_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"is_rollback_success": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"is_monitor_found": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
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

func buildChangeUpdateCreateOpts(d *schema.ResourceData) map[string]interface{} {
	ticketInfoParam := make(map[string]interface{})
	if v, ok := d.GetOk("phase"); ok {
		ticketInfoParam["phase"] = v
	}
	if v, ok := d.GetOk("work_flow_status"); ok {
		ticketInfoParam["work_flow_status"] = v
	}

	historyInfoParam := make(map[string]interface{})
	if v, ok := d.GetOk("action"); ok {
		historyInfoParam["action"] = v
	}

	bodyParams := map[string]interface{}{
		"sub_tickets": buildChangeUpdateSubTicketsCreateOpts(d.Get("sub_tickets")),
	}
	if len(ticketInfoParam) > 0 {
		bodyParams["ticket_info"] = ticketInfoParam
	}
	if len(historyInfoParam) > 0 {
		bodyParams["history_info"] = historyInfoParam
	}

	return bodyParams
}

func buildChangeUpdateSubTicketsCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				params[i] = map[string]interface{}{
					"ticket_id":                  utils.ValueIgnoreEmpty(raw["ticket_id"]),
					"change_result":              utils.ValueIgnoreEmpty(raw["change_result"]),
					"is_verified_in_change_time": utils.ValueIgnoreEmpty(raw["is_verified_in_change_time"]),
					"verified_docs":              utils.ValueIgnoreEmpty(raw["verified_docs"]),
					"comment":                    utils.ValueIgnoreEmpty(raw["comment"]),
					"change_fail_type":           utils.ValueIgnoreEmpty(raw["change_fail_type"]),
					"rollback_start_time":        utils.ValueIgnoreEmpty(raw["rollback_start_time"]),
					"rollback_end_time":          utils.ValueIgnoreEmpty(raw["rollback_end_time"]),
					"is_rollback_success":        utils.ValueIgnoreEmpty(raw["is_rollback_success"]),
					"is_monitor_found":           utils.ValueIgnoreEmpty(raw["is_monitor_found"]),
				}
			}
		}
		return params
	}

	return nil
}

func resourceChangeUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	ticketID := d.Get("ticket_id").(string)
	httpUrl := "v2/changes/{change_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{change_id}", ticketID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildChangeUpdateCreateOpts(d),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating COC change ticket: %s", err)
	}

	d.SetId(ticketID)

	return resourceChangeUpdateRead(ctx, d, meta)
}

func resourceChangeUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceChangeUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceChangeUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting change update resource is not supported. The change update resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
