package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceASInstances_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_as_instances.test"
		name  = acceptance.RandomAccResourceName()
		dc    = acceptance.InitDataSourceCheck(rName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceASInstances_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.instance_name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.scaling_group_name"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.life_cycle_state"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.health_status"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "instances.0.protect_from_scaling_down"),

					resource.TestCheckOutput("life_cycle_state_filter_is_useful", "true"),

					resource.TestCheckOutput("health_status_filter_is_useful", "true"),

					resource.TestCheckOutput("protect_from_scaling_down_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceASInstances_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_as_instances" "test" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
}

# Test with life_cycle_state
locals {
  life_cycle_state = data.huaweicloud_as_instances.test.instances[0].life_cycle_state
}

data "huaweicloud_as_instances" "life_cycle_state_filter" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  life_cycle_state = local.life_cycle_state
}

output "life_cycle_state_filter_is_useful" {
  value = length(data.huaweicloud_as_instances.life_cycle_state_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_as_instances.life_cycle_state_filter.instances[*].life_cycle_state : v == local.life_cycle_state]
  )  
}

# Test with health_status
locals {
  health_status = data.huaweicloud_as_instances.test.instances[0].health_status
}

data "huaweicloud_as_instances" "health_status_filter" {
  scaling_group_id = huaweicloud_as_group.acc_as_group.id
  health_status    = local.health_status
}

output "health_status_filter_is_useful" {
  value = length(data.huaweicloud_as_instances.health_status_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_as_instances.health_status_filter.instances[*].health_status : v == local.health_status]
  )  
}

# Test with protect_from_scaling_down
locals {
  protect_from_scaling_down = data.huaweicloud_as_instances.test.instances[0].protect_from_scaling_down
}

data "huaweicloud_as_instances" "protect_from_scaling_down_filter" {
  scaling_group_id          = huaweicloud_as_group.acc_as_group.id
  protect_from_scaling_down = local.protect_from_scaling_down
}

output "protect_from_scaling_down_filter_is_useful" {
  value = length(data.huaweicloud_as_instances.protect_from_scaling_down_filter.instances) > 0 && alltrue(
    [for v in data.huaweicloud_as_instances.protect_from_scaling_down_filter.instances[*].protect_from_scaling_down :
    v == local.protect_from_scaling_down]
  )  
}
`, testASGroup_forceDelete(name))
}
