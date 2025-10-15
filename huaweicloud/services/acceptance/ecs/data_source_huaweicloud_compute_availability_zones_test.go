package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeAvailabilityZones_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_availability_zones.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeAvailabilityZones_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.#"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.availability_zone_id"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.category"),
					resource.TestCheckResourceAttrSet(dataSource, "availability_zones.0.az_group_ids.#"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeAvailabilityZones_basic() string {
	return `
data "huaweicloud_compute_availability_zones" "test" {}
`
}
