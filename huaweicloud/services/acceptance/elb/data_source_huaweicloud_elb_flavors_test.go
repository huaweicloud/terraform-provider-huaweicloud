package elb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccElbFlavorsDataSource_basic(t *testing.T) {
	dataSource := "data.huaweicloud_elb_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccElbFlavorsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.max_connections"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.cps"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.qps"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "ids.#"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("max_connections_filter_is_useful", "true"),
					resource.TestCheckOutput("cps_filter_is_useful", "true"),
					resource.TestCheckOutput("qps_filter_is_useful", "true"),
					resource.TestCheckOutput("bandwidth_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccElbFlavorsDataSource_basic() string {
	return `
data "huaweicloud_elb_flavors" "test" {}

locals {
  type = "L4"
}
data "huaweicloud_elb_flavors" "type_filter" {
  type = "L4"
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.type_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.type_filter.flavors[*].type : v == local.type]
  )
}

locals {
  name = "L4_flavor.elb.s2.medium"
}
data "huaweicloud_elb_flavors" "name_filter" {
  name = "L4_flavor.elb.s2.medium"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.name_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.name_filter.flavors[*].name : v == local.name]
  )
}

locals {
  max_connections = 4000000
}
data "huaweicloud_elb_flavors" "max_connections_filter" {
  max_connections = 4000000
}
output "max_connections_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.max_connections_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.max_connections_filter.flavors[*].max_connections : v == local.max_connections]
  )
}

locals {
  cps = 10000
}
data "huaweicloud_elb_flavors" "cps_filter" {
  cps = 10000
}
output "cps_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.cps_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.cps_filter.flavors[*].cps : v == local.cps]
  )
}

locals {
  qps = 80000
}
data "huaweicloud_elb_flavors" "qps_filter" {
  qps = 80000
}
output "qps_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.qps_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.qps_filter.flavors[*].qps : v == local.qps]
  )
}

locals {
  bandwidth = 1000
}
data "huaweicloud_elb_flavors" "bandwidth_filter" {
  bandwidth = 1000
}
output "bandwidth_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.bandwidth_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.bandwidth_filter.flavors[*].bandwidth : v == local.bandwidth]
  )
}
`
}
