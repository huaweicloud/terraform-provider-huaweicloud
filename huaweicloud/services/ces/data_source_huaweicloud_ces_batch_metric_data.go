package ces

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CES POST /v2/{project_id}/batch-query-metric-data
func DataSourceCesBatchMetricData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesBatchMetricDataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_dimension": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_points": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     batchMetricDataDataPointSchema(),
			},
		},
	}
}

func batchMetricDataDataPointSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"dimensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     batchMetricDataDataPointsDimensionSchema(),
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func batchMetricDataDataPointsDimensionSchema() *schema.Resource {
	sc := schema.Resource{
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
	return &sc
}

func dataSourceCesBatchMetricDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v2/{project_id}/batch-query-metric-data"
		product = "ces"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	offset := 0
	res := make([]interface{}, 0)
	for {
		listOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		bodyParams, err := buildListBatchMetricDataQueryParams(d, offset)
		if err != nil {
			return diag.FromErr(err)
		}
		listOpt.JSONBody = bodyParams

		listResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving CES batch metric data: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		dataPoints := flattenCesBatchMetricDataBody(listRespBody)
		if len(dataPoints) < 1 {
			break
		}
		res = append(res, dataPoints...)
		offset += 1000
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_points", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListBatchMetricDataQueryParams(d *schema.ResourceData, offset int) (interface{}, error) {
	bodyParams := map[string]interface{}{
		"namespace":        d.Get("namespace").(string),
		"metric_name":      d.Get("metric_name").(string),
		"metric_dimension": d.Get("metric_dimension").(string),
		"limit":            1000,
		"offset":           offset,
	}
	if v, ok := d.GetOk("from"); ok {
		fromTime, err := utils.FormatUTCTimeStamp(v.(string))
		if err != nil {
			return nil, err
		}
		bodyParams["from"] = fromTime * 1000
	}
	if v, ok := d.GetOk("to"); ok {
		fromTime, err := utils.FormatUTCTimeStamp(v.(string))
		if err != nil {
			return nil, err
		}
		bodyParams["to"] = fromTime * 1000
	}
	return bodyParams, nil
}

func flattenCesBatchMetricDataBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("data_points", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"dimensions": flattenCesBatchMetricDataDimensionBody(v),
			"timestamp":  utils.PathSearch("timestamp", v, nil),
			"value":      utils.PathSearch("value", v, nil),
			"unit":       utils.PathSearch("unit", v, nil),
		})
	}
	return rst
}

func flattenCesBatchMetricDataDimensionBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dimensions", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
