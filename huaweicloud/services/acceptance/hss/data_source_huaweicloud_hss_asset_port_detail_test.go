package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetPortdetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_port_detail.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need a host with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetPortdetail_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.pid"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.port"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),

					resource.TestCheckOutput("host_name_filter_useful", "true"),
					resource.TestCheckOutput("host_ip_filter_useful", "true"),
					resource.TestCheckOutput("type_filter_useful", "true"),
					resource.TestCheckOutput("category_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAssetPortdetail_basic = `
data "huaweicloud_hss_asset_port_detail" "test" {
  port = 22
}

locals {
  host_name = data.huaweicloud_hss_asset_port_detail.test.data_list[0].host_name
  host_ip   = data.huaweicloud_hss_asset_port_detail.test.data_list[0].host_ip
  type      = data.huaweicloud_hss_asset_port_detail.test.data_list[0].type
}

data "huaweicloud_hss_asset_port_detail" "host_name_filter" {
  port      = 22
  host_name = local.host_name
}

output "host_name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_detail.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_port_detail.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

data "huaweicloud_hss_asset_port_detail" "host_ip_filter" {
  port    = 22
  host_ip = local.host_ip
}

output "host_ip_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_detail.host_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_port_detail.host_ip_filter.data_list[*].host_ip : v == local.host_ip]
  )
}

data "huaweicloud_hss_asset_port_detail" "type_filter" {	
  port = 22
  type = local.type
}

output "type_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_detail.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_port_detail.type_filter.data_list[*].type : v == local.type]
  )
}

data "huaweicloud_hss_asset_port_detail" "category_filter" {
  port     = 22
  category = "host"
}

output "category_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_detail.category_filter.data_list) > 0
}

data "huaweicloud_hss_asset_port_detail" "eps_filter" {
  port                  = 22
  enterprise_project_id = "0"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_asset_port_detail.eps_filter.data_list) > 0
}
`
