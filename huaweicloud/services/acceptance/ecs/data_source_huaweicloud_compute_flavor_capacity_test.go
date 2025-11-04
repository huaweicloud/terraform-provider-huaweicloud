package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeFlavorCapacity_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_flavor_capacity.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeFlavorCapacity_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "resources.#"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.availability_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "resources.0.prefer"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeFlavorCapacity_basic() string {
	return `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_compute_flavor_capacity" "test" {
  flavor_id = data.huaweicloud_compute_flavors.test.flavors[0].id
}
`
}
