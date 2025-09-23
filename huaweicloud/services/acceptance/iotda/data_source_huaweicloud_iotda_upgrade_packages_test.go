package iotda

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceUpgradePackages_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_iotda_upgrade_packages.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		name           = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHWIOTDAAccessAddress(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceUpgradePackages_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.space_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.product_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.support_source_versions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.custom_info"),
					resource.TestCheckResourceAttrSet(dataSourceName, "packages.0.created_at"),

					resource.TestCheckOutput("space_id_filter_is_useful", "true"),
					resource.TestCheckOutput("product_id_filter_is_useful", "true"),
					resource.TestCheckOutput("version_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceUpgradePackages_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_iotda_upgrade_packages" "test" {
  type = huaweicloud_iotda_upgrade_package.test.type
}

locals {
  space_id = data.huaweicloud_iotda_upgrade_packages.test.packages[0].space_id
}

data "huaweicloud_iotda_upgrade_packages" "space_id_filter" {
  type     = huaweicloud_iotda_upgrade_package.test.type
  space_id = local.space_id
}

output "space_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_upgrade_packages.space_id_filter.packages) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_upgrade_packages.space_id_filter.packages[*].space_id : v == local.space_id]
  )
}

locals {
  product_id = data.huaweicloud_iotda_upgrade_packages.test.packages[0].product_id
}

data "huaweicloud_iotda_upgrade_packages" "product_id_filter" {
  type       = huaweicloud_iotda_upgrade_package.test.type
  product_id = local.product_id
}

output "product_id_filter_is_useful" {
  value = length(data.huaweicloud_iotda_upgrade_packages.product_id_filter.packages) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_upgrade_packages.product_id_filter.packages[*].product_id : v == local.product_id]
  )
}

locals {
  version = data.huaweicloud_iotda_upgrade_packages.test.packages[0].version
}

data "huaweicloud_iotda_upgrade_packages" "version_filter" {
  type    = huaweicloud_iotda_upgrade_package.test.type
  version = local.version
}

output "version_filter_is_useful" {
  value = length(data.huaweicloud_iotda_upgrade_packages.version_filter.packages) > 0 && alltrue(
    [for v in data.huaweicloud_iotda_upgrade_packages.version_filter.packages[*].version : v == local.version]
  )
}
`, testUpgradePackage_basic(name))
}
