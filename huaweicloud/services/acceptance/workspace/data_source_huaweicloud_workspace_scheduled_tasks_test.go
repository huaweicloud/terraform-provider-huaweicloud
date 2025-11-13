package workspace

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataScheduledTasks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_workspace_scheduled_tasks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byTaskName   = "data.huaweicloud_workspace_scheduled_tasks.filter_by_task_name"
		dcByTaskName = acceptance.InitDataSourceCheck(byTaskName)

		byTaskType   = "data.huaweicloud_workspace_scheduled_tasks.filter_by_task_type"
		dcByTaskType = acceptance.InitDataSourceCheck(byTaskType)

		byScheduledType   = "data.huaweicloud_workspace_scheduled_tasks.filter_by_scheduled_type"
		dcByScheduledType = acceptance.InitDataSourceCheck(byScheduledType)

		byLastStatus   = "data.huaweicloud_workspace_scheduled_tasks.filter_by_last_status"
		dcByLastStatus = acceptance.InitDataSourceCheck(byLastStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataScheduledTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "tasks.#", regexp.MustCompile(`^[0-9]+$`)),
					dcByTaskName.CheckResourceExists(),
					resource.TestCheckOutput("is_task_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.type"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.scheduled_type"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.enable"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.time_zone"),
					dcByTaskType.CheckResourceExists(),
					resource.TestCheckOutput("is_task_type_filter_useful", "true"),
					dcByScheduledType.CheckResourceExists(),
					resource.TestCheckOutput("is_scheduled_type_filter_useful", "true"),
					dcByLastStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_last_status_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataScheduledTasks_basic = `
# Query all scheduled tasks without any filter parameters.
data "huaweicloud_workspace_scheduled_tasks" "test" {}

locals {
  task_name      = try(data.huaweicloud_workspace_scheduled_tasks.test.tasks[0].name, "NOT_FOUND")
  task_type      = try(data.huaweicloud_workspace_scheduled_tasks.test.tasks[0].type, "NOT_FOUND")
  scheduled_type = try(data.huaweicloud_workspace_scheduled_tasks.test.tasks[0].scheduled_type, "NOT_FOUND")
  last_status    = try([for v in data.huaweicloud_workspace_scheduled_tasks.test.tasks : v.last_status if v.last_status != ""][0], "NOT_FOUND")
}

# Filter by task name.
data "huaweicloud_workspace_scheduled_tasks" "filter_by_task_name" {
  task_name = local.task_name
}

locals {
  task_name_filter_result = [
    for v in data.huaweicloud_workspace_scheduled_tasks.filter_by_task_name.tasks : strcontains(v.name, local.task_name)
  ]
}

output "is_task_name_filter_useful" {
  value = length(local.task_name_filter_result) > 0 && alltrue(local.task_name_filter_result)
}

# Filter by task type.
data "huaweicloud_workspace_scheduled_tasks" "filter_by_task_type" {
  task_type = local.task_type
}

locals {
  task_type_filter_result = [
    for v in data.huaweicloud_workspace_scheduled_tasks.filter_by_task_type.tasks : v.type == local.task_type
  ]
}

output "is_task_type_filter_useful" {
  value = length(local.task_type_filter_result) > 0 && alltrue(local.task_type_filter_result)
}

# Filter by scheduled type.
data "huaweicloud_workspace_scheduled_tasks" "filter_by_scheduled_type" {
  scheduled_type = local.scheduled_type
}

locals {
  scheduled_type_filter_result = [
    for v in data.huaweicloud_workspace_scheduled_tasks.filter_by_scheduled_type.tasks : v.scheduled_type == local.scheduled_type
  ]
}

output "is_scheduled_type_filter_useful" {
  value = length(local.scheduled_type_filter_result) > 0 && alltrue(local.scheduled_type_filter_result)
}

# Filter by last status.
data "huaweicloud_workspace_scheduled_tasks" "filter_by_last_status" {
  last_status = local.last_status
}

locals {
  last_status_filter_result = [
    for v in data.huaweicloud_workspace_scheduled_tasks.filter_by_last_status.tasks : v.last_status == local.last_status
  ]
}

output "is_last_status_filter_useful" {
  value = length(local.last_status_filter_result) > 0 && alltrue(local.last_status_filter_result)
}
`
