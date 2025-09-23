package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/event/white-list/alarm
func DataSourceEventAlarmWhiteLists() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventAlarmWhiteListsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hash": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_type": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"remain_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"limit_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"event_type_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hash": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"white_field": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"field_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"judge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildEventAlarmWhiteListsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("hash"); ok {
		queryParams = fmt.Sprintf("%s&hash=%v", queryParams, v)
	}
	if v, ok := d.GetOk("event_type"); ok {
		queryParams = fmt.Sprintf("%s&event_type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceEventAlarmWhiteListsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		product       = "hss"
		epsId         = cfg.GetEnterpriseProjectID(d)
		result        = make([]interface{}, 0)
		eventTypeList = make([]interface{}, 0)
		offset        = 0
		totalNum      float64
		remainNum     float64
		limitNum      float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/event/white-list/alarm"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEventAlarmWhiteListsQueryParams(d, epsId)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS event alarm white lists: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		remainNum = utils.PathSearch("remain_num", respBody, float64(0)).(float64)
		limitNum = utils.PathSearch("limit_num", respBody, float64(0)).(float64)

		// Get `event_type_list` only once.
		if offset == 0 {
			eventTypeList = utils.PathSearch("event_type_list", respBody, make([]interface{}, 0)).([]interface{})
		}

		dataResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	eventTypeIntList := make([]int, 0, len(eventTypeList))
	for _, v := range eventTypeList {
		if fv, ok := v.(float64); ok {
			eventTypeIntList = append(eventTypeIntList, int(fv))
		}
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("remain_num", remainNum),
		d.Set("limit_num", limitNum),
		d.Set("event_type_list", eventTypeIntList),
		d.Set("data_list", flattenEventAlarmWhiteListsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEventAlarmWhiteListsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"enterprise_project_name": utils.PathSearch("enterprise_project_name", v, nil),
			"hash":                    utils.PathSearch("hash", v, nil),
			"description":             utils.PathSearch("description", v, nil),
			"event_type":              utils.PathSearch("event_type", v, nil),
			"white_field":             utils.PathSearch("white_field", v, nil),
			"field_value":             utils.PathSearch("field_value", v, nil),
			"judge_type":              utils.PathSearch("judge_type", v, nil),
			"update_time":             utils.PathSearch("update_time", v, nil),
		})
	}

	return rst
}
