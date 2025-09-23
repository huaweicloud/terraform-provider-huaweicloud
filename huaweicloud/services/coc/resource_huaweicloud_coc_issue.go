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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var issueNonUpdatableParams = []string{"title", "description", "enterprise_project_id", "ticket_type",
	"virtual_schedule_type", "commit_upload_attachment", "regions", "level", "root_cause_cloud_service", "source",
	"source_id", "found_time", "issue_contact_person", "schedule_scenes", "virtual_schedule_role"}

// @API COC POST /v1/{ticket_type}/tickets
// @API COC GET /v1/{ticket_type}/tickets/{ticket_id}
func ResourceIssue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIssueCreate,
		ReadContext:   resourceIssueRead,
		UpdateContext: resourceIssueUpdate,
		DeleteContext: resourceIssueDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(issueNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ticket_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virtual_schedule_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"commit_upload_attachment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"regions": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"level": {
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
			"issue_contact_person": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_scenes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"virtual_schedule_role": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
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
			"is_start_process": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issue_correlation_sla": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_cause_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_cloud_service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issue_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"root_cause_comment": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"solution": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"regions_search": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"level_approve_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"suspension_approve_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"handle_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_common_issue": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_need_change": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_enable_suspension": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_ticket_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"assignee": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"participator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"work_flow_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_error_msg": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"baseline_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"phase": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_tickets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"main_ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parent_ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ticket_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"enum_data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_deleted": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"match_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"real_ticket_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_zh": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name_en": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"biz_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prop_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"model_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"meta_data_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ticket_type_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"form_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIssueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createRespBody, err := createIssueTicket(client, d)
	if err != nil {
		return diag.Errorf("error flattening creating the COC issue ticket response: %s", err)
	}

	executionID := utils.PathSearch("data.id", createRespBody, "").(string)
	if executionID == "" {
		return diag.Errorf("error creating the COC issue ticket: can not find id in return")
	}

	d.SetId(executionID)

	return resourceIssueRead(ctx, d, meta)
}

func createIssueTicket(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v1/issues_mgmt/tickets"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildIssueCreateOpts(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, fmt.Errorf("error creating the COC issue ticket: %s", err)
	}

	return utils.FlattenResponse(createResp)
}

func buildIssueCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"title":                    utils.ValueIgnoreEmpty(d.Get("title")),
		"description":              utils.ValueIgnoreEmpty(d.Get("description")),
		"enterprise_project":       utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"level":                    utils.ValueIgnoreEmpty(d.Get("level")),
		"ticket_type":              utils.ValueIgnoreEmpty(d.Get("ticket_type")),
		"commit_upload_attachment": utils.ValueIgnoreEmpty(d.Get("commit_upload_attachment")),
		"regions": strings.Join(
			utils.ExpandToStringList(d.Get("regions").(*schema.Set).List()), ","),
		"root_cause_cloud_service": utils.ValueIgnoreEmpty(d.Get("root_cause_cloud_service")),
		"source":                   utils.ValueIgnoreEmpty(d.Get("source")),
		"source_id":                utils.ValueIgnoreEmpty(d.Get("source_id")),
		"fount_time":               utils.ValueIgnoreEmpty(d.Get("found_time")),
		"virtual_schedule_type":    utils.ValueIgnoreEmpty(d.Get("virtual_schedule_type")),
		"issue_contact_person":     utils.ValueIgnoreEmpty(d.Get("issue_contact_person")),
		"schedule_scenes":          utils.ValueIgnoreEmpty(d.Get("schedule_scenes")),
		"virtual_schedule_role":    utils.ValueIgnoreEmpty(d.Get("virtual_schedule_role")),
	}

	return bodyParams
}

func resourceIssueRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	ticketID := d.Id()
	issue, err := GetIssueTicketDetail(client, ticketID)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"COC.00067003"), "error retrieving COC issue ticket")
	}

	mErr := multierror.Append(nil,
		d.Set("issue_correlation_sla", utils.PathSearch("data.issue_correlation_sla", issue, nil)),
		d.Set("level", utils.PathSearch("data.level", issue, nil)),
		d.Set("root_cause_cloud_service", utils.PathSearch("data.root_cause_cloud_service", issue, nil)),
		d.Set("root_cause_type", utils.PathSearch("data.root_cause_type", issue, nil)),
		d.Set("current_cloud_service_id", utils.PathSearch("data.current_cloud_service_id", issue, nil)),
		d.Set("issue_contact_person", utils.PathSearch("data.issue_contact_person", issue, nil)),
		d.Set("issue_version", utils.PathSearch("data.issue_version", issue, nil)),
		d.Set("source", utils.PathSearch("data.source", issue, nil)),
		d.Set("commit_upload_attachment", utils.PathSearch("data.commit_upload_attachment", issue, nil)),
		d.Set("source_id", utils.PathSearch("data.source_id", issue, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("data.enterprise_project", issue, nil)),
		d.Set("virtual_schedule_type", utils.PathSearch("data.virtual_schedule_type", issue, nil)),
		d.Set("title", utils.PathSearch("data.title", issue, nil)),
		d.Set("regions", strings.Split(utils.PathSearch("data.regions", issue, "").(string), ",")),
		d.Set("description", utils.PathSearch("data.description", issue, nil)),
		d.Set("root_cause_comment", utils.PathSearch("data.root_cause_comment", issue, nil)),
		d.Set("solution", utils.PathSearch("data.solution", issue, nil)),
		d.Set("regions_search", utils.PathSearch("data.regions_serch", issue, nil)),
		d.Set("level_approve_config", utils.PathSearch("data.level_approve_config", issue, nil)),
		d.Set("suspension_approve_config", utils.PathSearch("data.suspension_approve_config", issue, nil)),
		d.Set("handle_time", utils.PathSearch("data.handle_time", issue, nil)),
		d.Set("found_time", utils.PathSearch("data.fount_time", issue, nil)),
		d.Set("is_common_issue", utils.PathSearch("data.is_common_issue", issue, nil)),
		d.Set("is_need_change", utils.PathSearch("data.is_need_change", issue, nil)),
		d.Set("is_enable_suspension", utils.PathSearch("data.is_enable_suspension", issue, nil)),
		d.Set("is_start_process_async", utils.PathSearch("data.is_start_process_async", issue, nil)),
		d.Set("is_update_null", utils.PathSearch("data.is_update_null", issue, nil)),
		d.Set("creator", utils.PathSearch("data.creator", issue, nil)),
		d.Set("operator", utils.PathSearch("data.operator", issue, nil)),
		d.Set("is_return_full_info", utils.PathSearch("data.is_return_full_info", issue, nil)),
		d.Set("is_start_process", utils.PathSearch("data.is_start_process", issue, nil)),
		d.Set("ticket_id", utils.PathSearch("data.ticket_id", issue, nil)),
		d.Set("real_ticket_id", utils.PathSearch("data.real_ticket_id", issue, nil)),
		d.Set("assignee", utils.PathSearch("data.assignee", issue, nil)),
		d.Set("participator", utils.PathSearch("data.participator", issue, nil)),
		d.Set("work_flow_status", utils.PathSearch("data.work_flow_status", issue, nil)),
		d.Set("engine_error_msg", utils.PathSearch("data.engine_error_msg", issue, nil)),
		d.Set("baseline_status", utils.PathSearch("data.baseline_status", issue, nil)),
		d.Set("ticket_type", utils.PathSearch("data.ticket_type", issue, nil)),
		d.Set("phase", utils.PathSearch("data.phase", issue, nil)),
		d.Set("sub_tickets", flattenCoIssueSubTickets(
			utils.PathSearch("data.sub_tickets", issue, nil))),
		d.Set("enum_data_list", flattenCoIssueEnumDataList(
			utils.PathSearch("data.enum_data_list", issue, nil))),
		d.Set("meta_data_version", utils.PathSearch("data.meta_data_version", issue, nil)),
		d.Set("update_time", utils.PathSearch("data.update_time", issue, nil)),
		d.Set("create_time", utils.PathSearch("data.create_time", issue, nil)),
		d.Set("is_deleted", utils.PathSearch("data.is_deleted", issue, nil)),
		d.Set("ticket_type_id", utils.PathSearch("data.ticket_type_id", issue, nil)),
		d.Set("form_info", utils.JsonToString(utils.PathSearch("data.properties._form_info", issue, nil))),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error getting COC issue ticket fields: %s", err)
	}

	return nil
}

func GetIssueTicketDetail(client *golangsdk.ServiceClient, ticketID string) (interface{}, error) {
	httpUrl := "v1/issues_mgmt/tickets/{ticket_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{ticket_id}", ticketID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenCoIssueSubTickets(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			m := map[string]interface{}{
				"is_deleted":       utils.PathSearch("is_deleted", params, nil),
				"id":               utils.PathSearch("id", params, nil),
				"main_ticket_id":   utils.PathSearch("main_ticket_id", params, nil),
				"parent_ticket_id": utils.PathSearch("parent_ticket_id", params, nil),
				"ticket_id":        utils.PathSearch("ticket_id", params, nil),
				"real_ticket_id":   utils.PathSearch("real_ticket_id", params, nil),
				"ticket_path":      utils.PathSearch("ticket_path", params, nil),
				"target_value":     utils.PathSearch("target_value", params, nil),
				"target_type":      utils.PathSearch("target_type", params, nil),
				"create_time":      utils.PathSearch("create_time", params, nil),
				"update_time":      utils.PathSearch("update_time", params, nil),
				"creator":          utils.PathSearch("creator", params, nil),
				"operator":         utils.PathSearch("operator", params, nil),
			}
			rst = append(rst, m)
		}
		return rst
	}
	return nil
}

func flattenCoIssueEnumDataList(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			m := map[string]interface{}{
				"is_deleted":     utils.PathSearch("is_deleted", params, nil),
				"match_type":     utils.PathSearch("match_type", params, nil),
				"ticket_id":      utils.PathSearch("ticket_id", params, nil),
				"real_ticket_id": utils.PathSearch("real_ticket_id", params, nil),
				"name_zh":        utils.PathSearch("name_zh", params, nil),
				"name_en":        utils.PathSearch("name_en", params, nil),
				"biz_id":         utils.PathSearch("biz_id", params, nil),
				"prop_id":        utils.PathSearch("prop_id", params, nil),
				"model_id":       utils.PathSearch("model_id", params, nil),
			}
			rst = append(rst, m)
		}
		return rst
	}
	return nil
}

func resourceIssueUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceIssueDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting issue resource is not supported. The issue resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
