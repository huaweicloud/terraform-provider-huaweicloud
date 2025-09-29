package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingLoginWhiteIps_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_login_white_ips.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, you need to prepare a host with host protection enabled.
			// And setting SSH login ip whitelist on console for the host.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceSettingLoginWhiteIps_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.white_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.total_num"),

					resource.TestCheckOutput("is_white_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSettingLoginWhiteIps_basic = `
data "huaweicloud_hss_setting_login_white_ips" "test" {}

locals {
  white_ip = data.huaweicloud_hss_setting_login_white_ips.test.data_list[0].white_ip
}

data "huaweicloud_hss_setting_login_white_ips" "white_ip_filter" {
  white_ip = local.white_ip
}

output "is_white_ip_filter_useful" {
  value = length(data.huaweicloud_hss_setting_login_white_ips.white_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_login_white_ips.white_ip_filter.data_list[*].white_ip : v == local.white_ip]
  )
}

data "huaweicloud_hss_setting_login_white_ips" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_setting_login_white_ips.eps_filter.data_list) > 0
}
`
