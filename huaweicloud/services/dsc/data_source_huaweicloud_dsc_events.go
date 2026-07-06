package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v2/{project_id}/sec-ops/events
func DataSourceDscEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"event_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the event risk level.",
			},
			"event_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the event name.",
			},
			"event_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the event handling status.",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the end time of the query (timestamp).",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the start time of the query (timestamp).",
			},
			"responsible_person": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the responsible person for the event.",
			},
			"source_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the event source name.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the event source type.",
			},
			"verification_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the verification status.",
			},
			"events": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The event information list.",
				Elem:        dscEventsRecordsSchema(),
			},
		},
	}
}

func dscEventsRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The event ID.",
			},
			"event_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event name.",
			},
			"event_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event risk level.",
			},
			"event_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event handling status.",
			},
			"event_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event type.",
			},
			"close_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The reason for closing the event.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The event creation time.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event description.",
			},
			"disposal_suggestion": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disposal suggestion.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain ID.",
			},
			"occur_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The event occurrence time.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project ID.",
			},
			"scheduled_close_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The scheduled close time.",
			},
			"source_module": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event source module.",
			},
			"verification_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The verification status.",
			},
			"affected_asset": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The affected assets.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"responsible_person": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The responsible person information.",
				Elem:        dscEventsResponsiblePersonSchema(),
			},
			"source_instance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The event source instance information.",
				Elem:        dscEventsSourceInstanceSchema(),
			},
			"related_alarm_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The related alarm information list.",
				Elem:        dscEventsRelatedAlarmSchema(),
			},
		},
	}
}

func dscEventsResponsiblePersonSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name.",
			},
		},
	}
}

func dscEventsSourceInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance name.",
			},
		},
	}
}

func dscEventsRelatedAlarmSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The alarm ID.",
			},
			"alarm_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm name.",
			},
			"alarm_level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm level.",
			},
			"alarm_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm status.",
			},
			"alarm_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm type.",
			},
			"close_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The reason for closing the alarm.",
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The alarm creation time.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm description.",
			},
			"disposal_suggestion": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The disposal suggestion.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The domain ID.",
			},
			"occur_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The alarm occurrence time.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The project ID.",
			},
			"source_module": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm source module.",
			},
			"source_sub_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm source sub type.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm source type.",
			},
			"verification_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The verification status.",
			},
			"affected_asset": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The affected assets.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"responsible_person": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The responsible person information.",
				Elem:        dscEventsRelatedAlarmResponsiblePersonSchema(),
			},
			"source_instance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The alarm source instance information.",
				Elem:        dscEventsRelatedAlarmSourceInstanceSchema(),
			},
		},
	}
}

func dscEventsRelatedAlarmResponsiblePersonSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user ID.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name.",
			},
		},
	}
}

func dscEventsRelatedAlarmSourceInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance ID.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance name.",
			},
		},
	}
}

func buildDscEventsQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	if v, ok := d.GetOk("event_level"); ok {
		queryParams = fmt.Sprintf("%s&event_level=%v", queryParams, v)
	}
	if v, ok := d.GetOk("event_name"); ok {
		queryParams = fmt.Sprintf("%s&event_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("event_status"); ok {
		queryParams = fmt.Sprintf("%s&event_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		queryParams = fmt.Sprintf("%s&start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("responsible_person"); ok {
		queryParams = fmt.Sprintf("%s&responsible_person=%v", queryParams, v)
	}
	if v, ok := d.GetOk("source_name"); ok {
		queryParams = fmt.Sprintf("%s&source_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("source_type"); ok {
		queryParams = fmt.Sprintf("%s&source_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("verification_status"); ok {
		queryParams = fmt.Sprintf("%s&verification_status=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDscEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/events"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDscEventsQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC events: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		eventInfoList := utils.PathSearch("event_info_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, eventInfoList...)

		if len(eventInfoList) < limit {
			break
		}

		offset += len(eventInfoList)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("events", flattenDscEventsRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscEventsRecords(eventInfoList []interface{}) []interface{} {
	if len(eventInfoList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(eventInfoList))
	for _, v := range eventInfoList {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"event_name":           utils.PathSearch("event_name", v, nil),
			"event_level":          utils.PathSearch("event_level", v, nil),
			"event_status":         utils.PathSearch("event_status", v, nil),
			"event_type":           utils.PathSearch("event_type", v, nil),
			"close_reason":         utils.PathSearch("close_reason", v, nil),
			"create_time":          utils.PathSearch("create_time", v, nil),
			"description":          utils.PathSearch("description", v, nil),
			"disposal_suggestion":  utils.PathSearch("disposal_suggestion", v, nil),
			"domain_id":            utils.PathSearch("domain_id", v, nil),
			"occur_time":           utils.PathSearch("occur_time", v, nil),
			"project_id":           utils.PathSearch("project_id", v, nil),
			"scheduled_close_time": utils.PathSearch("scheduled_close_time", v, nil),
			"source_module":        utils.PathSearch("source_module", v, nil),
			"verification_status":  utils.PathSearch("verification_status", v, nil),
			"affected_asset":       utils.PathSearch("affected_asset", v, nil),
			"responsible_person":   flattenDscEventsResponsiblePerson(v),
			"source_instance":      flattenDscEventsSourceInstance(v),
			"related_alarm_list":   flattenDscEventsRelatedAlarmList(v),
		})
	}

	return rst
}

func flattenDscEventsResponsiblePerson(respBody interface{}) []interface{} {
	responsiblePersonResp := utils.PathSearch("responsible_person", respBody, nil)
	if responsiblePersonResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", responsiblePersonResp, nil),
			"user_name": utils.PathSearch("user_name", responsiblePersonResp, nil),
		},
	}
}

func flattenDscEventsSourceInstance(respBody interface{}) []interface{} {
	sourceInstanceResp := utils.PathSearch("source_instance", respBody, nil)
	if sourceInstanceResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", sourceInstanceResp, nil),
			"instance_name": utils.PathSearch("instance_name", sourceInstanceResp, nil),
		},
	}
}

func flattenDscEventsRelatedAlarmList(respBody interface{}) []interface{} {
	alarmList := utils.PathSearch("related_alarm_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(alarmList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(alarmList))
	for _, alarm := range alarmList {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", alarm, nil),
			"alarm_name":          utils.PathSearch("alarm_name", alarm, nil),
			"alarm_level":         utils.PathSearch("alarm_level", alarm, nil),
			"alarm_status":        utils.PathSearch("alarm_status", alarm, nil),
			"alarm_type":          utils.PathSearch("alarm_type", alarm, nil),
			"close_reason":        utils.PathSearch("close_reason", alarm, nil),
			"create_time":         utils.PathSearch("create_time", alarm, nil),
			"description":         utils.PathSearch("description", alarm, nil),
			"disposal_suggestion": utils.PathSearch("disposal_suggestion", alarm, nil),
			"domain_id":           utils.PathSearch("domain_id", alarm, nil),
			"occur_time":          utils.PathSearch("occur_time", alarm, nil),
			"project_id":          utils.PathSearch("project_id", alarm, nil),
			"source_module":       utils.PathSearch("source_module", alarm, nil),
			"source_sub_type":     utils.PathSearch("source_sub_type", alarm, nil),
			"source_type":         utils.PathSearch("source_type", alarm, nil),
			"verification_status": utils.PathSearch("verification_status", alarm, nil),
			"affected_asset":      utils.PathSearch("affected_asset", alarm, nil),
			"responsible_person":  flattenDscEventsRelatedAlarmResponsiblePerson(alarm),
			"source_instance":     flattenDscEventsRelatedAlarmSourceInstance(alarm),
		})
	}

	return rst
}

func flattenDscEventsRelatedAlarmResponsiblePerson(alarm interface{}) []interface{} {
	responsiblePersonResp := utils.PathSearch("responsible_person", alarm, nil)
	if responsiblePersonResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", responsiblePersonResp, nil),
			"user_name": utils.PathSearch("user_name", responsiblePersonResp, nil),
		},
	}
}

func flattenDscEventsRelatedAlarmSourceInstance(alarm interface{}) []interface{} {
	sourceInstanceResp := utils.PathSearch("source_instance", alarm, nil)
	if sourceInstanceResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", sourceInstanceResp, nil),
			"instance_name": utils.PathSearch("instance_name", sourceInstanceResp, nil),
		},
	}
}
