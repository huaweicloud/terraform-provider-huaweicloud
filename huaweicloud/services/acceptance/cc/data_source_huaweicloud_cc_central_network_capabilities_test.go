package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCentralNetworkCapabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_central_network_capabilities.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCentralNetworkCapabilities_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.capability"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.specifications"),
				),
			},
			{
				Config: testDataSourceCcCentralNetworkCapabilities_useFilter(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "capabilities.0.capability", "central-network.is-support"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.specifications"),
				),
			},
		},
	})
}

func testDataSourceCcCentralNetworkCapabilities_basic() string {
	return `data "huaweicloud_cc_central_network_capabilities" "test" {}`
}

func testDataSourceCcCentralNetworkCapabilities_useFilter() string {
	return `
data "huaweicloud_cc_central_network_capabilities" "test" {
  capability = "central-network.is-support"
}
`
}
