package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussStorageTypes_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_storage_types.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussStorageTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.#"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.0.az_status.%"),
					resource.TestCheckResourceAttrSet(dataSource, "storage_type.0.support_compute_group_type.#"),
					resource.TestCheckOutput("ha_mode_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussStorageTypes_basic() string {
	return `
data "huaweicloud_gaussdb_opengauss_storage_types" "test" {
  version = "8.201"
  ha_mode = ""
}

locals {
  ha_mode = "centralization_standard"
}
data "huaweicloud_gaussdb_opengauss_storage_types" "ha_mode_filter" {
  version = "8.201"
  ha_mode = "centralization_standard"
}

output "ha_mode_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_storage_types.ha_mode_filter.storage_type) > 0
}
`
}
