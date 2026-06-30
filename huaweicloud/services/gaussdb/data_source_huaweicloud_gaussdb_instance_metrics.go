package gaussdb

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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/metric-data
func DataSourceInstanceMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceMetricsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_id": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"component_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"metrics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceMetricsMetricsSchema(),
			},
		},
	}
}

func instanceMetricsMetricsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datapoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceMetricsDatapointsSchema(),
			},
			"timestamps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func instanceMetricsDatapointsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"datapoint_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datapoint_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceInstanceMetricsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/metric-data"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildGetInstanceMetricsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance metrics: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("metrics", flattenGetInstanceMetricsBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetInstanceMetricsQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&start_time=%v", res, d.Get("start_time"))
	res = fmt.Sprintf("%s&end_time=%v", res, d.Get("end_time"))

	if v, ok := d.GetOk("metric"); ok {
		for _, val := range v.([]interface{}) {
			res = fmt.Sprintf("%s&metric=%v", res, val)
		}
	}

	if v, ok := d.GetOk("node_id"); ok {
		for _, val := range v.([]interface{}) {
			res = fmt.Sprintf("%s&node_id=%v", res, val)
		}
	}

	if v, ok := d.GetOk("component_id"); ok {
		for _, val := range v.([]interface{}) {
			res = fmt.Sprintf("%s&component_id=%v", res, val)
		}
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetInstanceMetricsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("metrics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"metric":     utils.PathSearch("metric", v, nil),
			"type":       utils.PathSearch("type", v, nil),
			"unit":       utils.PathSearch("unit", v, nil),
			"datapoints": flattenGetInstanceMetricsDatapoints(v),
			"timestamps": utils.PathSearch("timestamps", v, nil),
		})
	}
	return rst
}

func flattenGetInstanceMetricsDatapoints(resp interface{}) []interface{} {
	curJson := utils.PathSearch("datapoints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"datapoint_name":   utils.PathSearch("datapoint_name", v, nil),
			"datapoint_values": utils.PathSearch("datapoint_values", v, nil),
		})
	}
	return rst
}
