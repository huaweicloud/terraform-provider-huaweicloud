package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAvailableUpgradeSmallVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_available_upgrade_small_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAvailableUpgradeSmallVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_stores.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_stores.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_stores.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_stores.0.favored"),
				),
			},
		},
	})
}

func testDataSourceAvailableUpgradeSmallVersions_basic() string {
	return `
data "huaweicloud_rds_available_upgrade_small_versions" "test" {
  database_name = "postgresql"
  version       = "16"
}
`
}
