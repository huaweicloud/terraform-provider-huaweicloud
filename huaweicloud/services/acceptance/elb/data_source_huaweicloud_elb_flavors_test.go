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
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.shared"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.flavor_sold_out"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.public_border_group"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.max_connections"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.cps"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.qps"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.lcu"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.https_cps"),
					resource.TestCheckResourceAttrSet(dataSource, "ids.#"),

					resource.TestCheckOutput("flavor_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("shared_filter_is_useful", "true"),
					resource.TestCheckOutput("public_border_group_filter_is_useful", "true"),
					resource.TestCheckOutput("category_filter_is_useful", "true"),
					resource.TestCheckOutput("flavor_sold_out_filter_is_useful", "true"),
					resource.TestCheckOutput("list_all_filter_is_useful", "true"),
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
  flavor_id = data.huaweicloud_elb_flavors.test.flavors[0].id
}
data "huaweicloud_elb_flavors" "flavor_id_filter" {
  flavor_id = data.huaweicloud_elb_flavors.test.flavors[0].id
}
output "flavor_id_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.flavor_id_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.flavor_id_filter.flavors[*].id : v == local.flavor_id]
  )
}

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
  name = data.huaweicloud_elb_flavors.test.flavors[0].name
}
data "huaweicloud_elb_flavors" "name_filter" {
  name = data.huaweicloud_elb_flavors.test.flavors[0].name
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.name_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.name_filter.flavors[*].name : v == local.name]
  )
}

locals {
  shared = true
}
data "huaweicloud_elb_flavors" "shared_filter" {
  shared = "true"
}
output "shared_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.shared_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.shared_filter.flavors[*].shared : v == local.shared]
  )
}

locals {
  public_border_group = data.huaweicloud_elb_flavors.test.flavors[0].public_border_group
}
data "huaweicloud_elb_flavors" "public_border_group_filter" {
  public_border_group = data.huaweicloud_elb_flavors.test.flavors[0].public_border_group
}
output "public_border_group_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.public_border_group_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.public_border_group_filter.flavors[*].public_border_group : v == local.public_border_group]
  )
}

locals {
  category = data.huaweicloud_elb_flavors.test.flavors[0].category
}
data "huaweicloud_elb_flavors" "category_filter" {
  category = data.huaweicloud_elb_flavors.test.flavors[0].category
}
output "category_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.category_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.category_filter.flavors[*].category : v == local.category]
  )
}

locals {
  flavor_sold_out = data.huaweicloud_elb_flavors.test.flavors[0].flavor_sold_out
}
data "huaweicloud_elb_flavors" "flavor_sold_out_filter" {
  flavor_sold_out = data.huaweicloud_elb_flavors.test.flavors[0].flavor_sold_out
}
output "flavor_sold_out_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.flavor_sold_out_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.flavor_sold_out_filter.flavors[*].flavor_sold_out : v == local.flavor_sold_out]
  )
}

locals {
  list_all = "true"
}
data "huaweicloud_elb_flavors" "list_all_filter" {
  list_all = "true"
}
output "list_all_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.list_all_filter.flavors) > 0 
}

locals {
  max_connections = data.huaweicloud_elb_flavors.test.flavors[0].max_connections
}
data "huaweicloud_elb_flavors" "max_connections_filter" {
  max_connections = data.huaweicloud_elb_flavors.test.flavors[0].max_connections
}
output "max_connections_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.max_connections_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.max_connections_filter.flavors[*].max_connections : v == local.max_connections]
  )
}

locals {
  cps = data.huaweicloud_elb_flavors.test.flavors[0].cps
}
data "huaweicloud_elb_flavors" "cps_filter" {
  cps = data.huaweicloud_elb_flavors.test.flavors[0].cps
}
output "cps_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.cps_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.cps_filter.flavors[*].cps : v == local.cps]
  )
}

locals {
  qps = data.huaweicloud_elb_flavors.test.flavors[0].qps
}
data "huaweicloud_elb_flavors" "qps_filter" {
  qps = data.huaweicloud_elb_flavors.test.flavors[0].qps
}
output "qps_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.qps_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.qps_filter.flavors[*].qps : v == local.qps]
  )
}

locals {
  bandwidth = data.huaweicloud_elb_flavors.test.flavors[0].bandwidth
}
data "huaweicloud_elb_flavors" "bandwidth_filter" {
  bandwidth = data.huaweicloud_elb_flavors.test.flavors[0].bandwidth
}
output "bandwidth_filter_is_useful" {
  value = length(data.huaweicloud_elb_flavors.bandwidth_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_elb_flavors.bandwidth_filter.flavors[*].bandwidth : v == local.bandwidth]
  )
}
`
}
