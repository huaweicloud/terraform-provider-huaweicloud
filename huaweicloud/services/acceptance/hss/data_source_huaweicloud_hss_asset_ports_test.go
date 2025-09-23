package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetPorts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_ports.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetPorts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.pid"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),

					resource.TestCheckOutput("is_port_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAssetPorts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_asset_ports" "test" {
  host_id = "%[1]s"
}

locals {
  port = data.huaweicloud_hss_asset_ports.test.data_list[0].port
}

data "huaweicloud_hss_asset_ports" "port_filter" {
  host_id = "%[1]s"
  port    = local.port
}

output "is_port_filter_useful" {
  value = length(data.huaweicloud_hss_asset_ports.port_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_ports.port_filter.data_list[*].port : v == local.port]
  )
}

locals {
  type = data.huaweicloud_hss_asset_ports.test.data_list[0].type
}

data "huaweicloud_hss_asset_ports" "type_filter" {	
  host_id = "%[1]s"
  type    = local.type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_hss_asset_ports.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_ports.type_filter.data_list[*].type : v == local.type]
  )
}

locals {
  status = data.huaweicloud_hss_asset_ports.test.data_list[0].status	
}

data "huaweicloud_hss_asset_ports" "status_filter" {
  host_id = "%[1]s"
  status  = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_hss_asset_ports.status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_ports.status_filter.data_list[*].status : v == local.status]
  )
}

data "huaweicloud_hss_asset_ports" "category_filter" {
  host_id  = "%[1]s"
  category = "host"
}

output "is_category_filter_useful" {
  value = length(data.huaweicloud_hss_asset_ports.category_filter.data_list) > 0
}

data "huaweicloud_hss_asset_ports" "eps_filter" {
  host_id               = "%[1]s"
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_asset_ports.eps_filter.data_list) > 0
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
