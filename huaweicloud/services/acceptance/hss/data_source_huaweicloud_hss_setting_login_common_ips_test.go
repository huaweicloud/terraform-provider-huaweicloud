package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingLoginCommonIps_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_login_common_ips.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, you need to prepare a host with host protection enabled.
			// And setting login common ips on console for the host.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceSettingLoginCommonIps_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.ip_addr"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.total_num"),

					resource.TestCheckOutput("is_ip_addr_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSettingLoginCommonIps_basic = `
data "huaweicloud_hss_setting_login_common_ips" "test" {}

locals {
  ip_addr = data.huaweicloud_hss_setting_login_common_ips.test.data_list[0].ip_addr
}

data "huaweicloud_hss_setting_login_common_ips" "ip_addr_filter" {
  ip_addr  = local.ip_addr
}

output "is_ip_addr_filter_useful" {
  value = length(data.huaweicloud_hss_setting_login_common_ips.ip_addr_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_login_common_ips.ip_addr_filter.data_list[*].ip_addr : v == local.ip_addr]
  )
}

data "huaweicloud_hss_setting_login_common_ips" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_setting_login_common_ips.eps_filter.data_list) > 0
}
`
