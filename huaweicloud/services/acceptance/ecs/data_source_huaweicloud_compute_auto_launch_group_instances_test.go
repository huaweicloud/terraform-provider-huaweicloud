package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCmsComputeAutoLaunchGroupInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_auto_launch_group_instances.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSLaunchGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCmsComputeAutoLaunchGroupInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instances.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.flavor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.availability_zone_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "instances.0.sell_mode"),
				),
			},
		},
	})
}

func testDataSourceCmsComputeAutoLaunchGroupInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_compute_auto_launch_group_instances" "test" {
  auto_launch_group_id = "%s"
}
`, acceptance.HW_ECS_LAUNCH_GROUP_ID)
}
