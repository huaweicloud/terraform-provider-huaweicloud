package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussFlavors_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_flavors.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceGaussdbOpengaussFlavors_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.group_type"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.vcpus"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.ram"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.availability_zone.#"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.availability_zone.0"),
					resource.TestCheckResourceAttrSet(dataSource, "flavors.0.az_status.%"),
					resource.TestCheckOutput("version_filter_is_useful", "true"),
					resource.TestCheckOutput("spec_code_filter_is_useful", "true"),
					resource.TestCheckOutput("ha_mode_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceGaussdbOpengaussFlavors_basic() string {
	return `
data "huaweicloud_gaussdb_opengauss_flavors" "test" {}

locals {
  version = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].version
}
data "huaweicloud_gaussdb_opengauss_flavors" "version_filter" {
  version = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].version
}

output "version_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_flavors.version_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_flavors.version_filter.flavors[*].version : v == local.version]
  )
}

locals {
  spec_code = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
}
data "huaweicloud_gaussdb_opengauss_flavors" "spec_code_filter" {
  spec_code = data.huaweicloud_gaussdb_opengauss_flavors.test.flavors[0].spec_code
}

output "spec_code_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_flavors.spec_code_filter.flavors) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_flavors.spec_code_filter.flavors[*].spec_code : v == local.spec_code]
  )
}

locals {
  ha_mode = "centralization_standard"
}
data "huaweicloud_gaussdb_opengauss_flavors" "ha_mode_filter" {
  ha_mode = "centralization_standard"
}

output "ha_mode_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_flavors.ha_mode_filter.flavors) > 0
}
`
}
