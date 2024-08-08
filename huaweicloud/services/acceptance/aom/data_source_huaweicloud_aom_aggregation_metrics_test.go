package aom

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAggregationMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_aggregation_metrics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAggregationMetrics_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "service_metrics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "service_metrics.0.service"),
					resource.TestCheckResourceAttrSet(dataSource, "service_metrics.0.metrics.#"),
				),
			},
		},
	})
}

const testDataSourceAggregationMetrics_basic string = `data "huaweicloud_aom_aggregation_metrics" "test" {}`
