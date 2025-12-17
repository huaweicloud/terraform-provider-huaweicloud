package bms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAvailableResources_basic(t *testing.T) {
	rName := "data.huaweicloud_bms_available_resources.test"
	dc := acceptance.InitDataSourceCheck(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAvailableResources_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "available_resource.#"),
					resource.TestCheckResourceAttrSet(rName, "available_resource.0.availability_zone"),
					resource.TestCheckResourceAttrSet(rName, "available_resource.0.flavors.#"),
					resource.TestCheckResourceAttrSet(rName, "available_resource.0.flavors.0.flavor_id"),
					resource.TestCheckResourceAttrSet(rName, "available_resource.0.flavors.0.status"),
					resource.TestCheckOutput("availability_zone_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceAvailableResources_basic() string {
	return `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_bms_available_resources" "test" {
  availability_zone = [data.huaweicloud_availability_zones.test.names[0]]
}

locals{
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}
output "availability_zone_filter_is_useful" {
  value = length(data.huaweicloud_bms_available_resources.test.available_resource) > 0 && alltrue(
    [for v in data.huaweicloud_bms_available_resources.test.available_resource[*].availability_zone : v == local.availability_zone]
  )
}
`
}
