package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_statistics.test"
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
				Config: testDataSourceAssetStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "account_num"),
					resource.TestCheckResourceAttrSet(dataSource, "port_num"),
					resource.TestCheckResourceAttrSet(dataSource, "process_num"),
					resource.TestCheckResourceAttrSet(dataSource, "app_num"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_num"),
					resource.TestCheckResourceAttrSet(dataSource, "web_framework_num"),
					resource.TestCheckResourceAttrSet(dataSource, "web_site_num"),
					resource.TestCheckResourceAttrSet(dataSource, "jar_package_num"),
					resource.TestCheckResourceAttrSet(dataSource, "kernel_module_num"),
					resource.TestCheckResourceAttrSet(dataSource, "web_service_num"),
					resource.TestCheckResourceAttrSet(dataSource, "web_app_num"),
					resource.TestCheckResourceAttrSet(dataSource, "database_num"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAssetStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_asset_statistics" "test" {}

# Filter using host_id.
data "huaweicloud_hss_asset_statistics" "host_id_filter" {
  host_id = "%s"
}

output "is_host_id_filter_useful" {
  value = data.huaweicloud_hss_asset_statistics.host_id_filter.account_num > 0
}

# Filter using category.
data "huaweicloud_hss_asset_statistics" "category_filter" {
  category = "host"
}

output "is_category_filter_useful" {
  value = data.huaweicloud_hss_asset_statistics.category_filter.account_num > 0
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
