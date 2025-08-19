package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRmsCloudServices_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_rms_cloud_services.test"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRmsCloudServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.provider"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.category_display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.global"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.regions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.console_endpoint_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.console_list_url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.console_detail_url"),
					resource.TestCheckResourceAttrSet(dataSourceName, "resource_providers.0.resource_types.0.track"),
					resource.TestCheckOutput("track_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRmsCloudServices_basic() string {
	return `
data "huaweicloud_rms_cloud_services" "test" {}

data "huaweicloud_rms_cloud_services" "track_filter" {
  track = "tracked"
}

output "track_filter_is_useful" {
  value = length(data.huaweicloud_rms_cloud_services.track_filter.resource_providers) > 0 && alltrue(
    [for v in data.huaweicloud_rms_cloud_services.track_filter.resource_providers : alltrue(
    [for vv in v.resource_types : vv.track == "tracked"]
    )]
  )
}
`
}
