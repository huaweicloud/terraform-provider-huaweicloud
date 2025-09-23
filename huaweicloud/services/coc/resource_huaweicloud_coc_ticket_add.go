package coc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var ticketAddNonUpdatableParams = []string{"ticket_type", "title", "change_notes", "description", "enterprise_project_id",
	"change_type", "level", "issue_ticket_type", "change_scheme", "change_guides", "commit_upload_attachment",
	"regions", "change_scene_id", "current_cloud_service_id", "root_cause_cloud_service", "source", "source_id",
	"found_time", "virtual_schedule_type", "issue_contact_person", "schedule_scenes", "schedule_roles",
	"schedule_roles_name", "approve_type", "virtual_schedule_role", "location_id", "plan_task_sub_type", "plan_task_id",
	"plan_task_name", "plan_task_params", "is_start_process", "sub_tickets", "sub_tickets.*.app_name",
	"sub_tickets.*.region", "sub_tickets.*.target_type", "sub_tickets.*.target_value", "sub_tickets.*.expected_start_time",
	"sub_tickets.*.expected_end_time", "sub_tickets.*.executors", "sub_tickets.*.cooperators"}

// @API COC POST /v1/{ticket_type}/tickets
func ResourceTicketAdd() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTicketAddCreate,
		ReadContext:   resourceTicketAddRead,
		UpdateContext: resourceTicketAddUpdate,
		DeleteContext: resourceTicketAddDelete,

		CustomizeDiff: config.FlexibleForceNew(ticketAddNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"ticket_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"change_notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"change_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"issue_ticket_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"change_scheme": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"change_guides": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"commit_upload_attachment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"change_scene_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"current_cloud_service_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"root_cause_cloud_service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"found_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"virtual_schedule_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"issue_contact_person": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_scenes": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"schedule_roles": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"schedule_roles_name": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"approve_type": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"virtual_schedule_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"location_id": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
			},
			"plan_task_sub_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_task_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"plan_task_params": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_start_process": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"sub_tickets": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"region": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"target_value": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expected_start_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"expected_end_time": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"executors": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"cooperators": {
							Type:     schema.TypeString,
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
			"ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_start_process_async": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_update_null": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_return_full_info": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceTicketAddCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createRespBody, err := createTicketAddTicket(client, d)
	if err != nil {
		return diag.Errorf("error flattening creating the COC ticket response: %s", err)
	}

	id := utils.PathSearch("data.id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating the COC ticket: can not find id in return")
	}

	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("ticket_id", utils.PathSearch("data.ticket_id", createRespBody, nil)),
		d.Set("is_start_process_async", utils.PathSearch("data.is_start_process_async", createRespBody, nil)),
		d.Set("is_update_null", utils.PathSearch("data.is_update_null", createRespBody, nil)),
		d.Set("is_return_full_info", utils.PathSearch("data.is_return_full_info", createRespBody, nil)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error creating the COC ticket return fields: %s", mErr)
	}

	return nil
}

func createTicketAddTicket(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/{ticket_type}/tickets"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{ticket_type}", d.Get("ticket_type").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildTicketAddCreateOpts(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating the COC ticket: %s", err)
	}

	return utils.FlattenResponse(createResp)
}

func buildTicketAddCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"title":                    d.Get("title"),
		"change_notes":             utils.ValueIgnoreEmpty(d.Get("change_notes")),
		"description":              utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project":       utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"change_type":              utils.ValueIgnoreEmpty(d.Get("change_type")),
		"level":                    utils.ValueIgnoreEmpty(d.Get("level")),
		"ticket_type":              utils.ValueIgnoreEmpty(d.Get("issue_ticket_type")),
		"change_scheme":            utils.ValueIgnoreEmpty(d.Get("change_scheme")),
		"change_guides":            utils.ValueIgnoreEmpty(d.Get("change_guides")),
		"commit_upload_attachment": utils.ValueIgnoreEmpty(d.Get("commit_upload_attachment")),
		"regions":                  utils.ValueIgnoreEmpty(d.Get("regions")),
		"change_scene_id":          utils.ValueIgnoreEmpty(d.Get("change_scene_id")),
		"current_cloud_service_id": utils.ValueIgnoreEmpty(d.Get("current_cloud_service_id")),
		"root_cause_cloud_service": utils.ValueIgnoreEmpty(d.Get("root_cause_cloud_service")),
		"source":                   utils.ValueIgnoreEmpty(d.Get("source")),
		"source_id":                utils.ValueIgnoreEmpty(d.Get("source_id")),
		"fount_time":               utils.ValueIgnoreEmpty(d.Get("found_time")),
		"virtual_schedule_type":    utils.ValueIgnoreEmpty(d.Get("virtual_schedule_type")),
		"issue_contact_person":     utils.ValueIgnoreEmpty(d.Get("issue_contact_person")),
		"schedule_scenes":          utils.ValueIgnoreEmpty(d.Get("schedule_scenes")),
		"schedule_roles":           utils.ValueIgnoreEmpty(d.Get("schedule_roles")),
		"schedule_roles_name":      utils.ValueIgnoreEmpty(d.Get("schedule_roles_name")),
		"approve_type":             utils.ValueIgnoreEmpty(d.Get("approve_type")),
		"virtual_schedule_role":    utils.ValueIgnoreEmpty(d.Get("virtual_schedule_role")),
		"location_id":              utils.ValueIgnoreEmpty(d.Get("location_id")),
		"plan_task_sub_type":       utils.ValueIgnoreEmpty(d.Get("plan_task_sub_type")),
		"plan_task_id":             utils.ValueIgnoreEmpty(d.Get("plan_task_id")),
		"plan_task_name":           utils.ValueIgnoreEmpty(d.Get("plan_task_name")),
		"plan_task_params":         utils.ValueIgnoreEmpty(d.Get("plan_task_params")),
		"is_start_process":         utils.ValueIgnoreEmpty(d.Get("is_start_process")),
		"sub_tickets":              buildTicketAddSubTicketsCreateOpts(d.Get("sub_tickets")),
	}

	return bodyParams
}

func buildTicketAddSubTicketsCreateOpts(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				params[i] = map[string]interface{}{
					"target_type":         utils.ValueIgnoreEmpty(raw["target_type"]),
					"app_name":            utils.ValueIgnoreEmpty(raw["app_name"]),
					"region":              utils.ValueIgnoreEmpty(raw["region"]),
					"target_value":        utils.ValueIgnoreEmpty(raw["target_value"]),
					"expected_start_time": utils.ValueIgnoreEmpty(raw["expected_start_time"]),
					"expected_end_time":   utils.ValueIgnoreEmpty(raw["expected_end_time"]),
					"executors":           utils.ValueIgnoreEmpty(raw["executors"]),
					"cooperators":         utils.ValueIgnoreEmpty(raw["cooperators"]),
				}
			}
		}
		return params
	}

	return nil
}

func resourceTicketAddRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTicketAddUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTicketAddDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ticket add resource is not supported. The ticket add resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
