package cbh

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceAvailabilityZones_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cbh_availability_zones.test"
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAvailabilityZones_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "availability_zones.#"),
					resource.TestCheckResourceAttrSet(rName, "availability_zones.0.name"),
					resource.TestCheckResourceAttrSet(rName, "availability_zones.0.region_id"),
					resource.TestCheckResourceAttrSet(rName, "availability_zones.0.display_name"),
					resource.TestCheckResourceAttrSet(rName, "availability_zones.0.type"),
					resource.TestCheckResourceAttrSet(rName, "availability_zones.0.status"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("display_name_filter_is_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

func testAccDatasourceAvailabilityZones_basic() string {
	return `
data "huaweicloud_cbh_availability_zones" "test" {}

locals {
  name = data.huaweicloud_cbh_availability_zones.test.availability_zones[0].name
}

# Filter using name.
data "huaweicloud_cbh_availability_zones" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_cbh_availability_zones.name_filter.availability_zones) > 0 && alltrue(
    [for v in data.huaweicloud_cbh_availability_zones.name_filter.availability_zones[*].name : v == local.name]
  )
}

locals {
  display_name = data.huaweicloud_cbh_availability_zones.test.availability_zones[0].display_name
}

# Filter using display_name.
data "huaweicloud_cbh_availability_zones" "display_name_filter" {
  display_name = local.display_name
}

output "display_name_filter_is_useful" {
  value = length(data.huaweicloud_cbh_availability_zones.display_name_filter.availability_zones) > 0 && alltrue(
    [for v in data.huaweicloud_cbh_availability_zones.display_name_filter.availability_zones[*].display_name : v == local.display_name]
  )
}

# Filter using non existent name.
data "huaweicloud_cbh_availability_zones" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_cbh_availability_zones.not_found.availability_zones) == 0
}
`
}
