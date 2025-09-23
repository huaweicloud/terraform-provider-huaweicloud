package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test case, you need to prepare an AS group ID with hook instances.
func TestAccDataSourceHookInstances_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_as_hook_instances.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byInstanceId   = "data.huaweicloud_as_hook_instances.filter_by_instance_id"
		dcByInstanceId = acceptance.InitDataSourceCheck(byInstanceId)

		byLifecycleHookName   = "data.huaweicloud_as_hook_instances.filter_by_name"
		dcByLifecycleHookName = acceptance.InitDataSourceCheck(byLifecycleHookName)

		byLifecycleHookStatus   = "data.huaweicloud_as_hook_instances.filter_by_status"
		dcByLifecycleHookStatus = acceptance.InitDataSourceCheck(byLifecycleHookStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AS group in the hooked state in advance and configure the AS group ID into the
			// environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHookInstances_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.0.lifecycle_hook_name"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.0.lifecycle_hook_status"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.0.lifecycle_action_key"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.0.default_result"),
					resource.TestCheckResourceAttrSet(dataSource, "instance_hanging_info.0.timeout"),

					dcByInstanceId.CheckResourceExists(),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),

					dcByLifecycleHookName.CheckResourceExists(),
					resource.TestCheckOutput("lifecycle_hook_name_filter_is_useful", "true"),

					dcByLifecycleHookStatus.CheckResourceExists(),
					resource.TestCheckOutput("lifecycle_hook_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceHookInstances_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_hook_instances" "test" {
  scaling_group_id = "%[1]s"
}

locals {
  instance_id = data.huaweicloud_as_hook_instances.test.instance_hanging_info[0].instance_id
}

data "huaweicloud_as_hook_instances" "filter_by_instance_id" {
  scaling_group_id = "%[1]s"
  instance_id      = local.instance_id
}

locals {
  instance_id_filter_result = [
    for v in data.huaweicloud_as_hook_instances.filter_by_instance_id.instance_hanging_info[*].instance_id : 
    v == local.instance_id
  ]
}

output "instance_id_filter_is_useful" {
  value = alltrue(local.instance_id_filter_result) && length(local.instance_id_filter_result) > 0
}

locals {
  lifecycle_hook_name = data.huaweicloud_as_hook_instances.test.instance_hanging_info[0].lifecycle_hook_name
}

data "huaweicloud_as_hook_instances" "filter_by_name" {
  scaling_group_id    = "%[1]s"
  lifecycle_hook_name = local.lifecycle_hook_name
}

locals {
  lifecycle_hook_name_filter_result = [
    for v in data.huaweicloud_as_hook_instances.filter_by_name.instance_hanging_info[*].lifecycle_hook_name : 
    v == local.lifecycle_hook_name
  ]
}

output "lifecycle_hook_name_filter_is_useful" {
  value = alltrue(local.lifecycle_hook_name_filter_result) && length(local.lifecycle_hook_name_filter_result) > 0
}

locals {
  lifecycle_hook_status = data.huaweicloud_as_hook_instances.test.instance_hanging_info[0].lifecycle_hook_status
}

data "huaweicloud_as_hook_instances" "filter_by_status" {
  scaling_group_id      = "%[1]s"
  lifecycle_hook_status = local.lifecycle_hook_status
}

locals {
  lifecycle_hook_status_filter_result = [
    for v in data.huaweicloud_as_hook_instances.filter_by_status.instance_hanging_info[*].lifecycle_hook_status : 
    v == local.lifecycle_hook_status
  ]
}

output "lifecycle_hook_status_filter_is_useful" {
  value = alltrue(local.lifecycle_hook_status_filter_result) && length(local.lifecycle_hook_status_filter_result) > 0
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
