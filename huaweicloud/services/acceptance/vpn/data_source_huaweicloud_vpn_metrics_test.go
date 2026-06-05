package vpn

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpnMetrics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_vpn_metrics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceVpnMetrics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.metric_name"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.unit"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.dimensions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.dimensions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "metrics.0.dimensions.0.value"),
					resource.TestCheckOutput("namespace_filter_is_useful", "true"),
					resource.TestCheckOutput("metric_name_filter_is_useful", "true"),
					resource.TestCheckOutput("dim_filter_is_useful", "true"),
					resource.TestCheckOutput("order_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceVpnMetrics_basic() string {
	return `
data "huaweicloud_vpn_metrics" "test" {}

data "huaweicloud_vpn_metrics" "namespace_filter" {
  namespace = data.huaweicloud_vpn_metrics.test.metrics[0].namespace
}
locals {
  namespace = data.huaweicloud_vpn_metrics.test.metrics[0].namespace
}
output "namespace_filter_is_useful" {
  value = length(data.huaweicloud_vpn_metrics.namespace_filter.metrics) > 0 && alltrue(
  [for v in data.huaweicloud_vpn_metrics.namespace_filter.metrics[*].namespace : v == local.namespace]
  )
}

data "huaweicloud_vpn_metrics" "metric_name_filter" {
  metric_name = data.huaweicloud_vpn_metrics.test.metrics[0].metric_name
}
locals {
  metric_name = data.huaweicloud_vpn_metrics.test.metrics[0].metric_name
}
output "metric_name_filter_is_useful" {
  value = length(data.huaweicloud_vpn_metrics.metric_name_filter.metrics) > 0 && alltrue(
  [for v in data.huaweicloud_vpn_metrics.metric_name_filter.metrics[*].metric_name : v == local.metric_name]
  )
}

data "huaweicloud_vpn_metrics" "dim_filter" {
  dim = [format("%s,%s", local.dim_name,local.dim_value)]
}
locals {
  dim_name  = data.huaweicloud_vpn_metrics.test.metrics[0].dimensions[0].name
  dim_value = data.huaweicloud_vpn_metrics.test.metrics[0].dimensions[0].value
}
output "dim_filter_is_useful" {
  value = length(data.huaweicloud_vpn_metrics.dim_filter.metrics) > 0 && alltrue(
  [for v in data.huaweicloud_vpn_metrics.dim_filter.metrics[*].dimensions : 
    alltrue([for vv in v[*].name : vv == local.dim_name]) && alltrue([for vv in v[*].value : vv == local.dim_value])]
  )
}

data "huaweicloud_vpn_metrics" "order_filter" {
  order = "desc"
}
output "order_filter_is_useful" {
  value = length(data.huaweicloud_vpn_metrics.order_filter.metrics) > 0
}
`
}
