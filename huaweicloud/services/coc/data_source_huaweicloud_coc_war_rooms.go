package coc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCocWarRooms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocWarRoomsRead,

		Schema: map[string]*schema.Schema{
			"incident_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_code_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"incident_levels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"impacted_application_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"admin": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"status": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"triggered_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"triggered_end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"occur_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"occur_end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"recover_start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"recover_end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"notification_level": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"war_room_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"statistic_flag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"current_users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"war_room_nums": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"admin": {
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
						"regions": buildWarRoomsReqRegions(),
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"circular_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"war_room_status":      buildWarRoomsReqWarRoomStatus(),
						"impacted_application": buildWarRoomsReqImpactedApplication(),
						"processing_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"restoration_duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"war_room_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"running_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"closed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildWarRoomsReqRegions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"code": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func buildWarRoomsReqImpactedApplication() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// @API COC POST /v1/external/warrooms/list
func dataSourceCocWarRoomsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	listHttpUrl := "v1/external/warrooms/list"
	listPath := client.Endpoint + listHttpUrl

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	currentOffset := 0
	list := make([]map[string]interface{}, 0)
	for {
		listOpt.JSONBody = utils.RemoveNil(buildListWarRoomsBodyParams(d, currentOffset))
		listResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving COC war rooms: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		if currentOffset == 0 {
			mErr = multierror.Append(
				d.Set("running_num", utils.PathSearch("data.running_num", listRespBody, nil)),
				d.Set("closed_num", utils.PathSearch("data.closed_num", listRespBody, nil)),
				d.Set("total_num", utils.PathSearch("data.total_num", listRespBody, nil)),
			)
		}
		listWarRooms := flattenCocListWarRooms(listRespBody)
		if len(listWarRooms) < 1 {
			break
		}
		list = append(list, listWarRooms...)
		currentOffset += len(listWarRooms)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("list", list),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListWarRoomsBodyParams(d *schema.ResourceData, currentOffset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"limit":                    100,
		"offset":                   currentOffset,
		"admin":                    utils.ValueIgnoreEmpty(d.Get("admin")),
		"current_users":            utils.ValueIgnoreEmpty(d.Get("current_users")),
		"enterprise_project_ids":   utils.ValueIgnoreEmpty(d.Get("enterprise_project_ids")),
		"impacted_application_ids": utils.ValueIgnoreEmpty(d.Get("impacted_application_ids")),
		"incident_levels":          utils.ValueIgnoreEmpty(d.Get("incident_levels")),
		"incident_num":             utils.ValueIgnoreEmpty(d.Get("incident_num")),
		"notification_level":       utils.ValueIgnoreEmpty(d.Get("notification_level")),
		"occur_end_time":           utils.ValueIgnoreEmpty(d.Get("occur_end_time")),
		"occur_start_time":         utils.ValueIgnoreEmpty(d.Get("occur_start_time")),
		"recover_end_time":         utils.ValueIgnoreEmpty(d.Get("recover_end_time")),
		"recover_start_time":       utils.ValueIgnoreEmpty(d.Get("recover_start_time")),
		"region_code_list":         utils.ValueIgnoreEmpty(d.Get("region_code_list")),
		"statistic_flag":           utils.ValueIgnoreEmpty(d.Get("statistic_flag")),
		"status":                   utils.ValueIgnoreEmpty(d.Get("status")),
		"title":                    utils.ValueIgnoreEmpty(d.Get("title")),
		"triggered_end_time":       utils.ValueIgnoreEmpty(d.Get("triggered_end_time")),
		"triggered_start_time":     utils.ValueIgnoreEmpty(d.Get("triggered_start_time")),
		"war_room_num":             utils.ValueIgnoreEmpty(d.Get("war_room_num")),
		"war_room_nums":            utils.ValueIgnoreEmpty(d.Get("war_room_nums")),
	}
	return bodyParams
}

func flattenCocListWarRooms(resp interface{}) []map[string]interface{} {
	warRoomsJson := utils.PathSearch("data.list", resp, make([]interface{}, 0))
	warRoomsArray := warRoomsJson.([]interface{})
	if len(warRoomsArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(warRoomsArray))
	for _, warRoom := range warRoomsArray {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", warRoom, nil),
			"title":          utils.PathSearch("title", warRoom, nil),
			"admin":          utils.PathSearch("admin", warRoom, nil),
			"recover_member": utils.PathSearch("recover_member", warRoom, nil),
			"recover_leader": utils.PathSearch("recover_leader", warRoom, nil),
			"incident": flattenCocListWarRoomsIncident(
				utils.PathSearch("incident", warRoom, nil)),
			"source": utils.PathSearch("source", warRoom, nil),
			"regions": flattenCocListWarRommsRegions(
				utils.PathSearch("regions", warRoom, make([]interface{}, 0)).([]interface{})),
			"change_num":                 utils.PathSearch("change_num", warRoom, nil),
			"occur_time":                 utils.PathSearch("occur_time", warRoom, nil),
			"recover_time":               utils.PathSearch("recover_time", warRoom, nil),
			"fault_cause":                utils.PathSearch("fault_cause", warRoom, nil),
			"create_time":                utils.PathSearch("create_time", warRoom, nil),
			"first_report_time":          utils.PathSearch("first_report_time", warRoom, nil),
			"recovery_notification_time": utils.PathSearch("recovery_notification_time", warRoom, nil),
			"fault_impact":               utils.PathSearch("fault_impact", warRoom, nil),
			"description":                utils.PathSearch("description", warRoom, nil),
			"circular_level":             utils.PathSearch("circular_level", warRoom, nil),
			"war_room_status": flattenCocListWarRoomsWarRoomStatus(
				utils.PathSearch("war_room_status", warRoom, nil)),
			"impacted_application": flattenCocListWarRommsImpactedApplication(utils.PathSearch(
				"impacted_application", warRoom, make([]interface{}, 0)).([]interface{})),
			"processing_duration":   utils.PathSearch("processing_duration", warRoom, nil),
			"restoration_duration":  utils.PathSearch("restoration_duration", warRoom, nil),
			"war_room_num":          utils.PathSearch("war_room_num", warRoom, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", warRoom, nil),
		})
	}

	return result
}

func flattenCocListWarRommsRegions(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"code": utils.PathSearch("code", params, nil),
			"name": utils.PathSearch("name", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenCocListWarRommsImpactedApplication(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"id":   utils.PathSearch("id", params, nil),
			"name": utils.PathSearch("name", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}
