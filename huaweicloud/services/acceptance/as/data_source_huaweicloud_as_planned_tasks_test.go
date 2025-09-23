package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePlannedTasks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_planned_tasks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byTaskId   = "data.huaweicloud_as_planned_tasks.filter_by_task_id"
		dcByTaskId = acceptance.InitDataSourceCheck(byTaskId)

		byName   = "data.huaweicloud_as_planned_tasks.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AS group containing the planned tasks in advance and configure the AS group ID into
			// the environment variable.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePlannedTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "scheduled_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scheduled_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scheduled_tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scheduled_tasks.0.scaling_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scheduled_tasks.0.scheduled_policy.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scheduled_tasks.0.instance_number.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "scheduled_tasks.0.created_at"),

					dcByTaskId.CheckResourceExists(),
					resource.TestCheckOutput("task_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePlannedTasks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_planned_tasks" "test" {
  scaling_group_id = "%[1]s"
}

# Filter by task_id
locals {
  task_id = data.huaweicloud_as_planned_tasks.test.scheduled_tasks[0].id
}

data "huaweicloud_as_planned_tasks" "filter_by_task_id" {
  scaling_group_id = "%[1]s"
  task_id          = local.task_id
}

locals {
  task_id_filter_result = [
    for v in data.huaweicloud_as_planned_tasks.filter_by_task_id.scheduled_tasks[*].id : v == local.task_id
  ]
}

output "task_id_filter_is_useful" {
  value = alltrue(local.task_id_filter_result) && length(local.task_id_filter_result) > 0
}

# Filter by name
locals {
  name = data.huaweicloud_as_planned_tasks.test.scheduled_tasks[0].name
}

data "huaweicloud_as_planned_tasks" "filter_by_name" {
  scaling_group_id = "%[1]s"
  name             = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_as_planned_tasks.filter_by_name.scheduled_tasks[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
