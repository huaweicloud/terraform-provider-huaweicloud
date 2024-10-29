package ga

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaAvailabilityZones_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_ga_availability_zones.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byArea   = "data.huaweicloud_ga_availability_zones.filter_by_area"
		dcByArea = acceptance.InitDataSourceCheck(byArea)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAvailabilityZones_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "regions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "regions.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "regions.0.area"),
					resource.TestCheckResourceAttrSet(dataSourceName, "regions.0.endpoint_types.#"),

					dcByArea.CheckResourceExists(),
					resource.TestCheckOutput("area_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceAvailabilityZones_basic = `
data "huaweicloud_ga_availability_zones" "test" {}

locals {
  area = data.huaweicloud_ga_availability_zones.test.regions[0].area
}

data "huaweicloud_ga_availability_zones" "filter_by_area" {
  area = local.area
}

locals {
  area_filter_result = [
    for v in data.huaweicloud_ga_availability_zones.filter_by_area.regions[*].area : v == local.area
  ]
}

output "area_filter_is_useful" {
  value = alltrue(local.area_filter_result) && length(local.area_filter_result) > 0
}
`
