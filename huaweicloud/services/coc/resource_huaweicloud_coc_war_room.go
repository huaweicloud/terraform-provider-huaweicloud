package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var warRoomNonUpdatableParams = []string{"war_room_name", "description", "region_code_list", "application_id_list",
	"incident_number", "schedule_group", "schedule_group.role_id", "schedule_group.scene_id", "participant",
	"war_room_admin", "application_names", "region_names", "enterprise_project_id", "notification_type"}

// @API COC POST /v1/external/warrooms
// @API COC POST /v1/external/warrooms/list
func ResourceWarRoom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWarRoomCreate,
		ReadContext:   resourceWarRoomRead,
		UpdateContext: resourceWarRoomUpdate,
		DeleteContext: resourceWarRoomDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(warRoomNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"war_room_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"incident_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"schedule_group": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"role_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"scene_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"war_room_admin": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_code_list": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"participant": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"application_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notification_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"war_room_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"recover_member": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"recover_leader": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"incident": buildWarRoomsReqIncident(),
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"change_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"occur_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"recover_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fault_cause": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"first_report_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"recovery_notification_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fault_impact": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"circular_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"war_room_status": buildWarRoomsReqWarRoomStatus(),
			"processing_duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"restoration_duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildWarRoomsReqIncident() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"incident_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"is_change_event": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"failure_level": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"incident_url": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func buildWarRoomsReqWarRoomStatus() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
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
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func buildWarRoomCreateOpts(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"war_room_name":         d.Get("war_room_name"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"region_code_list":      utils.ValueIgnoreEmpty(d.Get("region_code_list").(*schema.Set).List()),
		"application_id_list":   d.Get("application_id_list").(*schema.Set).List(),
		"incident_number":       d.Get("incident_number"),
		"schedule_group":        buildWarRoomScheduleGroupParamsBody(d.Get("schedule_group").(*schema.Set).List()),
		"participant":           utils.ValueIgnoreEmpty(d.Get("participant").(*schema.Set).List()),
		"war_room_admin":        d.Get("war_room_admin"),
		"application_names":     utils.ValueIgnoreEmpty(d.Get("application_names").(*schema.Set).List()),
		"region_names":          utils.ValueIgnoreEmpty(d.Get("region_names").(*schema.Set).List()),
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"notification_type":     utils.ValueIgnoreEmpty(d.Get("notification_type")),
	}

	return bodyParams
}

func buildWarRoomScheduleGroupParamsBody(rawParams []interface{}) []interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	rst := make([]interface{}, len(rawParams))
	for i, v := range rawParams {
		rst[i] = map[string]interface{}{
			"role_id":  utils.PathSearch("role_id", v, nil),
			"scene_id": utils.PathSearch("scene_id", v, nil),
		}
	}

	return rst
}

func resourceWarRoomCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	httpUrl := "v1/external/warrooms"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildWarRoomCreateOpts(d)),
	}

	createWarRoomResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC war room: %s", err)
	}

	createWarRoomRespBody, err := utils.FlattenResponse(createWarRoomResp)
	if err != nil {
		return diag.FromErr(err)
	}

	warRoomNum := utils.PathSearch("data", createWarRoomRespBody, "").(string)
	if warRoomNum == "" {
		return diag.Errorf("unable to find the COC war room ID from the API response")
	}
	d.SetId(warRoomNum)

	return resourceWarRoomRead(ctx, d, meta)
}

func resourceWarRoomRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	readWarRoomRespBody, err := GetWarRoom(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving war room (%s)", d.Id()))
	}

	mErr = multierror.Append(
		mErr,
		d.Set("war_room_name", utils.PathSearch("title", readWarRoomRespBody, nil)),
		d.Set("application_id_list", utils.PathSearch("impacted_application[*].id", readWarRoomRespBody, nil)),
		d.Set("incident_number", utils.PathSearch("incident.incident_id", readWarRoomRespBody, nil)),
		d.Set("war_room_admin", utils.PathSearch("admin", readWarRoomRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", readWarRoomRespBody, nil)),
		d.Set("description", utils.PathSearch("description", readWarRoomRespBody, nil)),
		d.Set("region_code_list", utils.PathSearch("regions[*].code", readWarRoomRespBody, nil)),
		d.Set("application_names", utils.PathSearch("impacted_application[*].name", readWarRoomRespBody, nil)),
		d.Set("region_names", utils.PathSearch("regions[*].name", readWarRoomRespBody, nil)),
		d.Set("war_room_id", utils.PathSearch("id", readWarRoomRespBody, nil)),
		d.Set("recover_member", utils.PathSearch("recover_member", readWarRoomRespBody, nil)),
		d.Set("recover_leader", utils.PathSearch("recover_leader", readWarRoomRespBody, nil)),
		d.Set("incident", flattenCocListWarRoomsIncident(
			utils.PathSearch("incident", readWarRoomRespBody, nil))),
		d.Set("source", utils.PathSearch("source", readWarRoomRespBody, nil)),
		d.Set("change_num", utils.PathSearch("change_num", readWarRoomRespBody, nil)),
		d.Set("occur_time", utils.PathSearch("occur_time", readWarRoomRespBody, nil)),
		d.Set("recover_time", utils.PathSearch("recover_time", readWarRoomRespBody, nil)),
		d.Set("fault_cause", utils.PathSearch("fault_cause", readWarRoomRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", readWarRoomRespBody, nil)),
		d.Set("first_report_time", utils.PathSearch("first_report_time", readWarRoomRespBody, nil)),
		d.Set("recovery_notification_time", utils.PathSearch("recovery_notification_time", readWarRoomRespBody, nil)),
		d.Set("fault_impact", utils.PathSearch("fault_impact", readWarRoomRespBody, nil)),
		d.Set("circular_level", utils.PathSearch("circular_level", readWarRoomRespBody, nil)),
		d.Set("war_room_status", flattenCocListWarRoomsWarRoomStatus(
			utils.PathSearch("war_room_status", readWarRoomRespBody, nil))),
		d.Set("processing_duration", utils.PathSearch("processing_duration", readWarRoomRespBody, nil)),
		d.Set("restoration_duration", utils.PathSearch("restoration_duration", readWarRoomRespBody, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting COC war room fields: %s", err)
	}

	return nil
}

func flattenCocListWarRoomsIncident(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"id":              utils.PathSearch("id", param, nil),
			"incident_id":     utils.PathSearch("incident_id", param, nil),
			"is_change_event": utils.PathSearch("is_change_event", param, nil),
			"failure_level":   utils.PathSearch("failure_level", param, nil),
			"incident_url":    utils.PathSearch("incident_url", param, nil),
		},
	}

	return rst
}

func flattenCocListWarRoomsWarRoomStatus(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"id":      utils.PathSearch("id", param, nil),
			"name_zh": utils.PathSearch("name_zh", param, nil),
			"name_en": utils.PathSearch("name_en", param, nil),
			"type":    utils.PathSearch("type", param, nil),
		},
	}

	return rst
}

func GetWarRoom(client *golangsdk.ServiceClient, warRoomNum string) (interface{}, error) {
	httpUrl := "v1/external/warrooms/list"
	readPath := client.Endpoint + httpUrl

	readOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	readOpt.JSONBody = map[string]interface{}{
		// have to send limit and offset, otherwise the list returns empty
		"limit":        1,
		"offset":       0,
		"war_room_num": warRoomNum,
	}

	readWarRoomResp, err := client.Request("POST", readPath, &readOpt)
	if err != nil {
		return nil, err
	}
	readWarRoomRespBody, err := utils.FlattenResponse(readWarRoomResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening war room: %s", err)
	}

	warRoomRespBody := utils.PathSearch(fmt.Sprintf("data.list[?war_room_num=='%s']|[0]", warRoomNum), readWarRoomRespBody, nil)
	if warRoomRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return warRoomRespBody, nil
}

func resourceWarRoomUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWarRoomDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting war room resource is not supported. The war room resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
