package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsInterconnectedServices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rms_interconnected_services.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRmsInterconnectedServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.provider"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.category_display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.global"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.console_endpoint_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.console_detail_url"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.console_list_url"),
					resource.TestCheckResourceAttrSet(dataSource, "resource_providers.0.resource_types.0.track"),
				),
			},
		},
	})
}

func testDataSourceRmsInterconnectedServices_basic() string {
	return `
data "huaweicloud_rms_interconnected_services" "test" {}
`
}
