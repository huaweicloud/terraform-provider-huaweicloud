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

// @API DSC GET /v2/{project_id}/sec-ops/alarms
func DataSourceDscAlarms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAlarmsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alarm_level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the alarm level.",
			},
			"alarm_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the alarm name.",
			},
			"alarm_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the alarm status.",
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
				Description: "Specifies the responsible person for the alarm.",
			},
			"source_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the alarm source name.",
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the alarm source type.",
			},
			"verification_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the verification status.",
			},
			"alarms": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The alarm information list.",
				Elem:        dscAlarmsRecordsSchema(),
			},
		},
	}
}

func dscAlarmsRecordsSchema() *schema.Resource {
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
				Elem:        dscAlarmsResponsiblePersonSchema(),
			},
			"source_instance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The alarm source instance information.",
				Elem:        dscAlarmsSourceInstanceSchema(),
			},
		},
	}
}

func dscAlarmsResponsiblePersonSchema() *schema.Resource {
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

func dscAlarmsSourceInstanceSchema() *schema.Resource {
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

func buildDscAlarmsQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := ""

	if v, ok := d.GetOk("alarm_level"); ok {
		queryParams = fmt.Sprintf("%s&alarm_level=%v", queryParams, v)
	}
	if v, ok := d.GetOk("alarm_name"); ok {
		queryParams = fmt.Sprintf("%s&alarm_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("alarm_status"); ok {
		queryParams = fmt.Sprintf("%s&alarm_status=%v", queryParams, v)
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

	if offset > 0 {
		queryParams = fmt.Sprintf("%s&offset=%v", queryParams, offset)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceDscAlarmsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/alarms"
		offset  = 0
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
	}

	for {
		currentPath := requestPath + buildDscAlarmsQueryParams(d, offset)
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC alarms: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		alarmInfoList := utils.PathSearch("alarm_info_list", requestRespBody, make([]interface{}, 0)).([]interface{})
		if len(alarmInfoList) == 0 {
			break
		}

		result = append(result, alarmInfoList...)
		offset += len(alarmInfoList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("alarms", flattenDscAlarmsRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscAlarmsRecords(alarmInfoList []interface{}) []interface{} {
	if len(alarmInfoList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(alarmInfoList))
	for _, v := range alarmInfoList {
		rst = append(rst, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"alarm_name":          utils.PathSearch("alarm_name", v, nil),
			"alarm_level":         utils.PathSearch("alarm_level", v, nil),
			"alarm_status":        utils.PathSearch("alarm_status", v, nil),
			"alarm_type":          utils.PathSearch("alarm_type", v, nil),
			"close_reason":        utils.PathSearch("close_reason", v, nil),
			"create_time":         utils.PathSearch("create_time", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"disposal_suggestion": utils.PathSearch("disposal_suggestion", v, nil),
			"domain_id":           utils.PathSearch("domain_id", v, nil),
			"occur_time":          utils.PathSearch("occur_time", v, nil),
			"project_id":          utils.PathSearch("project_id", v, nil),
			"source_module":       utils.PathSearch("source_module", v, nil),
			"source_sub_type":     utils.PathSearch("source_sub_type", v, nil),
			"source_type":         utils.PathSearch("source_type", v, nil),
			"verification_status": utils.PathSearch("verification_status", v, nil),
			"affected_asset":      utils.PathSearch("affected_asset", v, nil),
			"responsible_person":  flattenDscAlarmsResponsiblePerson(v),
			"source_instance":     flattenDscAlarmsSourceInstance(v),
		})
	}

	return rst
}

func flattenDscAlarmsResponsiblePerson(respBody interface{}) []interface{} {
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

func flattenDscAlarmsSourceInstance(respBody interface{}) []interface{} {
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
