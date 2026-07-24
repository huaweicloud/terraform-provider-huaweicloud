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

// @API DSC GET /v1/{project_id}/devices/alarms
func DataSourceDscAlarmInfos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAlarmInfosRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alarm_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscAlarmInfoSchema(),
			},
		},
	}
}

func dscAlarmInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"module": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildDscAlarmInfosQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func dataSourceDscAlarmInfosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/devices/alarms"
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
		currentPath := requestPath + buildDscAlarmInfosQueryParams(limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC alarm infos: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		alarmInfos := utils.PathSearch("alarm_infos", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, alarmInfos...)

		if len(alarmInfos) < limit {
			break
		}

		offset += len(alarmInfos)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("alarm_infos", flattenDscAlarmInfos(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscAlarmInfos(alarmInfos []interface{}) []interface{} {
	if len(alarmInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(alarmInfos))
	for _, v := range alarmInfos {
		rst = append(rst, map[string]interface{}{
			"count":       utils.PathSearch("count", v, nil),
			"create_time": utils.PathSearch("create_time", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"device_ip":   utils.PathSearch("device_ip", v, nil),
			"id":          utils.PathSearch("id", v, nil),
			"module":      utils.PathSearch("module", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"severity":    utils.PathSearch("severity", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"type":        utils.PathSearch("type", v, nil),
		})
	}

	return rst
}
