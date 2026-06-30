package vpn

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

// @API VPN GET /V1.0/{project_id}/metrics
func DataSourceVpnMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpnMetricsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dim": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metrics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     vpnMetricsSchema(),
			},
		},
	}
}

func vpnMetricsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dimensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     vpnMetricsDimensionsSchema(),
			},
		},
	}
}

func vpnMetricsDimensionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceVpnMetricsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	httpUrl := "V1.0/{project_id}/metrics"

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var allMetrics []interface{}
	start := ""
	for {
		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath += buildGetVpnMetricsQueryParams(d, start)

		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving VPN metrics: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		metrics := flattenGetVpnMetricsBody(listRespBody)
		if len(metrics) == 0 {
			break
		}
		allMetrics = append(allMetrics, metrics...)

		marker := utils.PathSearch("meta_data.marker", listRespBody, "").(string)
		if marker == "" {
			break
		}
		start = marker
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("metrics", allMetrics),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetVpnMetricsQueryParams(d *schema.ResourceData, start string) string {
	res := "?limit=1000"
	if v, ok := d.GetOk("namespace"); ok {
		res = fmt.Sprintf("%s&namespace=%v", res, v)
	}
	if v, ok := d.GetOk("metric_name"); ok {
		res = fmt.Sprintf("%s&metric_name=%v", res, v)
	}
	if v, ok := d.GetOk("dim"); ok {
		for i, dim := range v.([]interface{}) {
			res = fmt.Sprintf("%s&dim.%d=%v", res, i, dim)
		}
	}
	if start != "" {
		res = fmt.Sprintf("%s&start=%v", res, start)
	}
	if v, ok := d.GetOk("order"); ok {
		res = fmt.Sprintf("%s&order=%v", res, v)
	}

	return res
}

func flattenGetVpnMetricsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("metrics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"namespace":   utils.PathSearch("namespace", v, nil),
			"metric_name": utils.PathSearch("metric_name", v, nil),
			"unit":        utils.PathSearch("unit", v, nil),
			"dimensions":  flattenGetVpnMetricsDimensionsBody(v),
		})
	}
	return res
}

func flattenGetVpnMetricsDimensionsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dimensions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return res
}
