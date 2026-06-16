package gaussdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/alarm-history-record
func DataSourceAlarms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"history_records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbAlarmsSchema(),
			},
		},
	}
}

func gaussDbAlarmsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alarm_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alarm_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildAlarmsParams(d *schema.ResourceData) string {
	startTime := strings.ReplaceAll(d.Get("start_time").(string), "+", "%2B")
	queryParams := fmt.Sprintf("&start_time=%s", startTime)

	if v, ok := d.GetOk("level"); ok {
		queryParams = fmt.Sprintf("%s&level=%v", queryParams, v)
	}

	return queryParams
}

func listAlarms(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {

	var (
		httpUrl = "v3/{project_id}/alarm-history-record?limit={limit}"
		result  = make([]interface{}, 0)
		limit   = 100
		offset  = 0
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildAlarmsParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		alarms := utils.PathSearch("history_records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, alarms...)
		if len(alarms) < limit {
			break
		}

		offset += len(alarms)
	}

	return result, nil
}

func dataSourceAlarmsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	alarms, err := listAlarms(client, d)
	if err != nil {
		return diag.Errorf("error querying GaussDB alarm history records: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("history_records", flattenAlarms(alarms)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAlarms(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("history_records", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(
			res, map[string]interface{}{
				"alarm_id":      utils.PathSearch("alarm_id", v, nil),
				"name":          utils.PathSearch("name", v, nil),
				"status":        utils.PathSearch("status", v, nil),
				"alarm_type":    utils.PathSearch("alarm_type", v, nil),
				"level":         utils.PathSearch("level", v, nil),
				"instance_id":   utils.PathSearch("instance_id", v, nil),
				"instance_name": utils.PathSearch("instance_name", v, nil),
				"begin_time":    utils.PathSearch("begin_time", v, nil),
				"update_time":   utils.PathSearch("update_time", v, nil),
			},
		)
	}
	return res
}
