package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetApps_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_apps.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssetApps_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.app_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.recent_scan_time"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_app_name_filter_useful", "true"),
					resource.TestCheckOutput("is_host_ip_filter_useful", "true"),
					resource.TestCheckOutput("is_version_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testDataSourceAssetApps_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_asset_apps" "test" {}

# Filter using host ID.
data "huaweicloud_hss_asset_apps" "host_id_filter" {
  host_id = "%[1]s"
}

output "is_host_id_filter_useful" {
  value = length(data.huaweicloud_hss_asset_apps.host_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_apps.host_id_filter.data_list[*].host_id : v == "%[1]s"]
  )
}

# Filter using host name.
locals {
  host_name = data.huaweicloud_hss_asset_apps.test.data_list[0].host_name
}

data "huaweicloud_hss_asset_apps" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_apps.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_apps.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

# Filter using app name.
locals {
  app_name = data.huaweicloud_hss_asset_apps.test.data_list[0].app_name
}

data "huaweicloud_hss_asset_apps" "app_name_filter" {
  app_name = local.app_name
}

output "is_app_name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_apps.app_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_apps.app_name_filter.data_list[*].app_name : v == local.app_name]
  )
}

# Filter using host IP.
locals {
  host_ip = data.huaweicloud_hss_asset_apps.test.data_list[0].host_ip
}

data "huaweicloud_hss_asset_apps" "host_ip_filter" {
  host_ip = local.host_ip
}

output "is_host_ip_filter_useful" {
  value = length(data.huaweicloud_hss_asset_apps.host_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_apps.host_ip_filter.data_list[*].host_ip : v == local.host_ip]
  )
}

# Filter using version.
locals {
  version = data.huaweicloud_hss_asset_apps.test.data_list[0].version
}

data "huaweicloud_hss_asset_apps" "version_filter" {
  version = local.version
}

output "is_version_filter_useful" {
  value = length(data.huaweicloud_hss_asset_apps.version_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_apps.version_filter.data_list[*].version : v == local.version]
  )
}

# Filter using non existent app name.
data "huaweicloud_hss_asset_apps" "not_found" {
  app_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_asset_apps.not_found.data_list) == 0
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
