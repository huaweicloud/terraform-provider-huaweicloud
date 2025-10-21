package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCmsComputeAutoLaunchGroups_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_auto_launch_groups.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSLaunchTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCmsComputeAutoLaunchGroups_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_groups.#"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_groups.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_groups.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_groups.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_groups.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_groups.0.task_state"),
					resource.TestCheckResourceAttrSet(dataSource, "auto_launch_groups.0.valid_since"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("valid_since_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCmsComputeAutoLaunchGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_compute_auto_launch_groups" "test" {
  depends_on = [huaweicloud_compute_auto_launch_group.test]
}

locals {
  name = "%[2]s"
}
data "huaweicloud_compute_auto_launch_groups" "name_filter" {
  depends_on = [huaweicloud_compute_auto_launch_group.test]

  name = "%[2]s"
}
output "name_filter_is_useful" {
  value = length(data.huaweicloud_compute_auto_launch_groups.name_filter.auto_launch_groups) > 0 && alltrue(
  [for v in data.huaweicloud_compute_auto_launch_groups.name_filter.auto_launch_groups[*].name : v == local.name]
  )
}

data "huaweicloud_compute_auto_launch_groups" "valid_since_filter" {
  depends_on = [huaweicloud_compute_auto_launch_group.test]

  valid_since = data.huaweicloud_compute_auto_launch_groups.test.auto_launch_groups[0].valid_since
}
output "valid_since_filter_is_useful" {
  value = length(data.huaweicloud_compute_auto_launch_groups.valid_since_filter.auto_launch_groups) > 0
}
`, testAccAutoLaunchGroup_basic(name), name)
}
