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
		dc    = acceptance.InitDataSourceCheck(rName)

		byLifeCycleState   = "data.huaweicloud_as_instances.life_cycle_state_filter"
		dcByLifeCycleState = acceptance.InitDataSourceCheck(byLifeCycleState)

		byHealthStatus   = "data.huaweicloud_as_instances.health_status_filter"
		dcByHealthStatus = acceptance.InitDataSourceCheck(byHealthStatus)

		byProtectFromScalingDown   = "data.huaweicloud_as_instances.protect_from_scaling_down_filter"
		dcByProtectFromScalingDown = acceptance.InitDataSourceCheck(byProtectFromScalingDown)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AS group containing the instance in advance and configure the AS group ID into the
			// environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceASInstances_basic(),
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

					dcByLifeCycleState.CheckResourceExists(),
					resource.TestCheckOutput("life_cycle_state_filter_is_useful", "true"),

					dcByHealthStatus.CheckResourceExists(),
					resource.TestCheckOutput("health_status_filter_is_useful", "true"),

					dcByProtectFromScalingDown.CheckResourceExists(),
					resource.TestCheckOutput("protect_from_scaling_down_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceASInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_instances" "test" {
  scaling_group_id = "%[1]s"
}

# Test with life_cycle_state
locals {
  life_cycle_state = data.huaweicloud_as_instances.test.instances[0].life_cycle_state
}

data "huaweicloud_as_instances" "life_cycle_state_filter" {
  scaling_group_id = "%[1]s"
  life_cycle_state = local.life_cycle_state
}

locals {
  life_cycle_state_filter_result = [
    for v in data.huaweicloud_as_instances.life_cycle_state_filter.instances[*].life_cycle_state : v == local.life_cycle_state
  ]
}

output "life_cycle_state_filter_is_useful" {
  value = length(local.life_cycle_state_filter_result) > 0 && alltrue(local.life_cycle_state_filter_result)  
}

# Test with health_status
locals {
  health_status = data.huaweicloud_as_instances.test.instances[0].health_status
}

data "huaweicloud_as_instances" "health_status_filter" {
  scaling_group_id = "%[1]s"
  health_status    = local.health_status
}

locals {
  health_status_filter_result = [
    for v in data.huaweicloud_as_instances.health_status_filter.instances[*].health_status : v == local.health_status
  ]
}

output "health_status_filter_is_useful" {
  value = length(local.health_status_filter_result) > 0 && alltrue(local.health_status_filter_result)  
}

# Test with protect_from_scaling_down
locals {
  protect_from_scaling_down = data.huaweicloud_as_instances.test.instances[0].protect_from_scaling_down
}

data "huaweicloud_as_instances" "protect_from_scaling_down_filter" {
  scaling_group_id          = "%[1]s"
  protect_from_scaling_down = local.protect_from_scaling_down
}

locals {
  protect_from_scaling_down_filter_result = [
    for v in data.huaweicloud_as_instances.protect_from_scaling_down_filter.instances[*].protect_from_scaling_down :
    v == local.protect_from_scaling_down
  ]
}

output "protect_from_scaling_down_filter_is_useful" {
  value = length(local.protect_from_scaling_down_filter_result) > 0 && alltrue(local.protect_from_scaling_down_filter_result)  
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
