package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageClasses_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_classes.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageClasses_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.#"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.name_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.display_priority"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_levels.0.updated_at"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackageClasses_basic() string {
	return `
data "huaweicloud_cc_bandwidth_package_classes" "test" {}
`
}
