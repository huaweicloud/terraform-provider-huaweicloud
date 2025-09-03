package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcBandwidthPackageLines_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_bandwidth_package_lines.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcBandwidthPackageLines_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.#"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.local_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.remote_region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.support_levels.#"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.0.spec_code"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.0.name_cn"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.0.max_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.0.min_bandwidth"),
					resource.TestCheckResourceAttrSet(dataSource, "bandwidth_package_lines.0.spec_codes.0.support_billing_modes.#"),
				),
			},
		},
	})
}

func testDataSourceCcBandwidthPackageLines_basic() string {
	return `
data "huaweicloud_cc_bandwidth_package_lines" "test" {}
`
}
