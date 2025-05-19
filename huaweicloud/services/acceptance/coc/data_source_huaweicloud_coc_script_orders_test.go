package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocScriptOrders_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_script_orders.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocScriptOrderID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocScriptOrders_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.order_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.order_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.execute_uuid"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.gmt_created"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.gmt_finished"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.execute_costs"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.properties.0.region_ids"),
					resource.TestCheckOutput("start_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
					resource.TestCheckOutput("creator_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocScriptOrders_basic() string {
	return `
data "huaweicloud_coc_script_orders" "test" {}

locals {
  start_time = [for v in data.huaweicloud_coc_script_orders.test.data[*].gmt_created : v if v != ""][0]
}

data "huaweicloud_coc_script_orders" "start_time_filter" {
  start_time = local.start_time
}

output "start_time_filter_is_useful" {
  value = length(data.huaweicloud_coc_script_orders.start_time_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_script_orders.start_time_filter.data[*].gmt_created : v >= local.start_time]
  )
}

locals {
  end_time = [for v in data.huaweicloud_coc_script_orders.test.data[*].gmt_finished : v if v != ""][0]
}

data "huaweicloud_coc_script_orders" "end_time_filter" {
  end_time = local.end_time
}

output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_coc_script_orders.end_time_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_script_orders.end_time_filter.data[*].gmt_finished : v <= local.end_time]
  )
}

locals {
  creator = [for v in data.huaweicloud_coc_script_orders.test.data[*].creator : v if v != ""][0]
}

data "huaweicloud_coc_script_orders" "creator_filter" {
  creator = local.creator
}

output "creator_filter_is_useful" {
  value = length(data.huaweicloud_coc_script_orders.creator_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_script_orders.creator_filter.data[*].creator : v == local.creator]
  )
}

locals {
  status = [for v in data.huaweicloud_coc_script_orders.test.data[*].status : v if v != ""][0]
}

data "huaweicloud_coc_script_orders" "status_filter" {
  status = local.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_coc_script_orders.status_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_coc_script_orders.status_filter.data[*].status : v == local.status]
  )
}
`
}
