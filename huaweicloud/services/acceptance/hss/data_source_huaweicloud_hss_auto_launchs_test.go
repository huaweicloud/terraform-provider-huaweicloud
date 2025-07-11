package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAutoLaunchs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_auto_launchs.test"
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
				Config: testAccDataSourceAutoLaunchs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.run_user"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.recent_scan_time"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAutoLaunchs_basic = `
data "huaweicloud_hss_auto_launchs" "test" {}

locals {
  host_id = data.huaweicloud_hss_auto_launchs.test.data_list[0].host_id
}

data "huaweicloud_hss_auto_launchs" "host_id_filter" {
  host_id = local.host_id
}

output "is_host_id_filter_useful" {
  value = length(data.huaweicloud_hss_auto_launchs.host_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_auto_launchs.host_id_filter.data_list[*].host_id : v == local.host_id]
  )
}

locals {
  host_name = data.huaweicloud_hss_auto_launchs.test.data_list[0].host_name
}

data "huaweicloud_hss_auto_launchs" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_auto_launchs.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_auto_launchs.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

locals {
  name = data.huaweicloud_hss_auto_launchs.test.data_list[0].name
}

data "huaweicloud_hss_auto_launchs" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_hss_auto_launchs.name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_auto_launchs.name_filter.data_list[*].name : v == local.name]
  )
}

locals {
  type = data.huaweicloud_hss_auto_launchs.test.data_list[0].type
}

data "huaweicloud_hss_auto_launchs" "type_filter" {
  type = local.type
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_hss_auto_launchs.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_auto_launchs.type_filter.data_list[*].type : v == local.type]
  )
}

data "huaweicloud_hss_auto_launchs" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_auto_launchs.eps_filter.data_list) > 0
}
`
