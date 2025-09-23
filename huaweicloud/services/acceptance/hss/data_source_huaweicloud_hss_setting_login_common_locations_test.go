package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingLoginCommonLocations_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_login_common_locations.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, you need to prepare a host with host protection enabled.
			// And setting common login locations on console for the host.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceSettingLoginCommonLocations_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.area_code"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.total_num"),

					resource.TestCheckOutput("is_area_code_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSettingLoginCommonLocations_basic = `
data "huaweicloud_hss_setting_login_common_locations" "test" {}

locals {
  area_code = data.huaweicloud_hss_setting_login_common_locations.test.data_list[0].area_code
}

data "huaweicloud_hss_setting_login_common_locations" "area_code_filter" {
  area_code = local.area_code
}

output "is_area_code_filter_useful" {
  value = length(data.huaweicloud_hss_setting_login_common_locations.area_code_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_login_common_locations.area_code_filter.data_list[*].area_code : v == local.area_code]
  )
}

data "huaweicloud_hss_setting_login_common_locations" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_setting_login_common_locations.eps_filter.data_list) > 0
}
`
