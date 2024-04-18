package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcGlobalConnectionBandwidthLineLevels_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_global_connection_bandwidth_line_levels.test"

	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcGlobalConnectionBandwidthLineLevels_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.local_area"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.remote_area"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.levels.#"),
				),
			},
			{
				Config: testAccDataSourceCcGlobalConnectionBandwidthLineLevels_localArea(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "line_levels.0.local_area", "cn-south-guangzhou"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.remote_area"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.levels.#"),
				),
			},
			{
				Config: testAccDataSourceCcGlobalConnectionBandwidthLineLevels_levels(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "line_levels.0.levels.0", "Ag"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.local_area"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.remote_area"),
					resource.TestCheckResourceAttrSet(dataSource, "line_levels.0.levels.#"),
				),
			},
		},
	})
}

func testDataSourceCcGlobalConnectionBandwidthLineLevels_basic() string {
	return `data "huaweicloud_cc_global_connection_bandwidth_line_levels" "test" {}`
}

func testAccDataSourceCcGlobalConnectionBandwidthLineLevels_localArea() string {
	return `
data "huaweicloud_cc_global_connection_bandwidth_line_levels" "test" {
  local_area = "cn-south-guangzhou"
}
`
}

func testAccDataSourceCcGlobalConnectionBandwidthLineLevels_levels() string {
	return `
data "huaweicloud_cc_global_connection_bandwidth_line_levels" "test" {
  levels = "Ag"
}
`
}
