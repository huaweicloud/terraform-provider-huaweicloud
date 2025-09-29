package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingTwoFactorLoginHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_two_factor_login_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, you need to prepare a host with host protection enabled.
			// And setting two-factor authentication on console for the host.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceSettingTwoFactorHosts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.auth_switch"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.auth_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.topic_display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.outside_host"),

					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_display_name_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSettingTwoFactorHosts_basic = `
data "huaweicloud_hss_setting_two_factor_login_hosts" "test" {}

locals {
  host_name = data.huaweicloud_hss_setting_two_factor_login_hosts.test.data_list[0].host_name
}

data "huaweicloud_hss_setting_two_factor_login_hosts" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_setting_two_factor_login_hosts.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_two_factor_login_hosts.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

locals {
  display_name = data.huaweicloud_hss_setting_two_factor_login_hosts.test.data_list[0].topic_display_name
}

data "huaweicloud_hss_setting_two_factor_login_hosts" "display_name_filter" {
  display_name = local.display_name
}

output "is_display_name_filter_useful" {
  value = length(data.huaweicloud_hss_setting_two_factor_login_hosts.display_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_two_factor_login_hosts.display_name_filter.data_list[*].topic_display_name : v == local.display_name]
  )
}

data "huaweicloud_hss_setting_two_factor_login_hosts" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_setting_two_factor_login_hosts.eps_filter.data_list) > 0
}
`
