package aom

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

// @API AOM GET /v1/{project_id}/aom/aggr-metrics
func DataSourceAggregationMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAggregationMetricsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_metrics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"metrics": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAggregationMetricsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	listHttpUrl := "v1/{project_id}/aom/aggr-metrics"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving aggregation metrics: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.Errorf("error flattening aggregation metrics: %s", err)
	}
	serviceMetrics := utils.PathSearch("service_metrics", listRespBody, make([]interface{}, 0)).([]interface{})

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("service_metrics", flattenMetrics(serviceMetrics)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMetrics(metrics []interface{}) []interface{} {
	if len(metrics) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		result = append(result, map[string]interface{}{
			"service": utils.PathSearch("service", metric, nil),
			"metrics": utils.PathSearch("metrics", metric, nil),
		})
	}
	return result
}
