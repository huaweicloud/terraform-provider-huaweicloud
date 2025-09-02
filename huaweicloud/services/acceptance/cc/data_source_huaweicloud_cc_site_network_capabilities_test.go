package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcSiteNetworkCapabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_site_network_capabilities.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcSiteNetworkCapabilities_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.specification"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.is_support_enterprise_project"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.is_support_tag"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.is_support_intra_region"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.is_support"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_locations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_freeze_regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_topologies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_dscp_regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.charge_mode.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.size_range.#"),
					resource.TestCheckOutput("specification_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcSiteNetworkCapabilities_basic() string {
	return `
data "huaweicloud_cc_site_network_capabilities" "test" {}

data "huaweicloud_cc_site_network_capabilities" "specification_filter" {
  specification = [data.huaweicloud_cc_site_network_capabilities.test.capabilities[0].specification]
}
locals{
  specification = data.huaweicloud_cc_site_network_capabilities.test.capabilities[0].specification
}
output "specification_filter_is_useful" {
  value = length(data.huaweicloud_cc_site_network_capabilities.specification_filter.capabilities) > 0 && alltrue(
  [for v in data.huaweicloud_cc_site_network_capabilities.specification_filter.capabilities[*].specification :
  v == local.specification]
  )
}
`
}
