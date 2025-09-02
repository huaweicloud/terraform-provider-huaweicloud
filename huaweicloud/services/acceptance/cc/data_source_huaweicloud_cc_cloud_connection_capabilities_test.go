package cc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcCloudConnectionCapabilities_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_cloud_connection_capabilities.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcCloudConnectionCapabilities_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "capabilities.0.support_regions.#"),
					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcCloudConnectionCapabilities_basic() string {
	return `
data "huaweicloud_cc_cloud_connection_capabilities" "test" {}

data "huaweicloud_cc_cloud_connection_capabilities" "resource_type_filter" {
  resource_type = data.huaweicloud_cc_cloud_connection_capabilities.test.capabilities[0].resource_type
}

locals {
  resource_type = data.huaweicloud_cc_cloud_connection_capabilities.test.capabilities[0].resource_type
}

output "resource_type_filter_is_useful" {
  value = length(data.huaweicloud_cc_cloud_connection_capabilities.resource_type_filter.capabilities) > 0 && alltrue(
    [for v in data.huaweicloud_cc_cloud_connection_capabilities.resource_type_filter.capabilities[*].resource_type :
    v == local.resource_type]
  )
}
`
}
