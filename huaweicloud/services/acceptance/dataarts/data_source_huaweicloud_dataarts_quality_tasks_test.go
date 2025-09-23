package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before testing, make sure that at least one quality task is prepared in the environment and at least one quality
// task has been run.
func TestAccDataSourceQualityTasks_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_quality_tasks.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dataarts_quality_tasks.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_dataarts_quality_tasks.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byCategoryId   = "data.huaweicloud_dataarts_quality_tasks.filter_by_category_id"
		dcByCategoryId = acceptance.InitDataSourceCheck(byCategoryId)

		byScheduleStatus   = "data.huaweicloud_dataarts_quality_tasks.filter_by_schedule_status"
		dcByScheduleStatus = acceptance.InitDataSourceCheck(byScheduleStatus)

		byStartTime   = "data.huaweicloud_dataarts_quality_tasks.filter_by_start_time"
		dcByStartTime = acceptance.InitDataSourceCheck(byStartTime)

		byCreator   = "data.huaweicloud_dataarts_quality_tasks.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceQualityTasks_basic_step1(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "tasks.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
			{
				Config: testAccDataSourceQualityTasks_basic_step2(),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
					dcByCategoryId.CheckResourceExists(),
					resource.TestCheckOutput("is_category_id_filter_useful", "true"),
					dcByScheduleStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_schedule_status_filter_useful", "true"),
					dcByStartTime.CheckResourceExists(),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceQualityTasks_basic_step1() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_quality_tasks" "test" {
  workspace_id = "%[1]s"
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDataSourceQualityTasks_basic_step2() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_quality_tasks" "test" {
  workspace_id = "%[1]s"
}

# Filter by name
locals {
  task_name = data.huaweicloud_dataarts_quality_tasks.test.tasks[0].name
}

data "huaweicloud_dataarts_quality_tasks" "filter_by_name" {
  workspace_id = "%[1]s"

  name = local.task_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_quality_tasks.filter_by_name.tasks[*].name : v == local.task_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by not found name
data "huaweicloud_dataarts_quality_tasks" "filter_by_not_found_name" {
  workspace_id = "%[1]s"

  name = "not_found"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_dataarts_quality_tasks.filter_by_not_found_name.tasks) < 1
}

# Filter by category ID
locals {
  category_id = data.huaweicloud_dataarts_quality_tasks.test.tasks[0].category_id
}

data "huaweicloud_dataarts_quality_tasks" "filter_by_category_id" {
  workspace_id = "%[1]s"

  category_id = local.category_id
}

locals {
  category_id_filter_result = [
    for v in data.huaweicloud_dataarts_quality_tasks.filter_by_name.tasks[*].category_id : v == local.category_id
  ]
}

output "is_category_id_filter_useful" {
  value = length(local.category_id_filter_result) > 0 && alltrue(local.category_id_filter_result)
}

# Filter by schedule status
locals {
  schedule_status = data.huaweicloud_dataarts_quality_tasks.test.tasks[0].schedule_status
}

data "huaweicloud_dataarts_quality_tasks" "filter_by_schedule_status" {
  workspace_id = "%[1]s"

  schedule_status = local.schedule_status
}

locals {
  schedule_status_filter_result = [
    for v in data.huaweicloud_dataarts_quality_tasks.filter_by_schedule_status.tasks[*].schedule_status : v == local.schedule_status
  ]
}

output "is_schedule_status_filter_useful" {
  value = length(local.schedule_status_filter_result) > 0 && alltrue(local.schedule_status_filter_result)
}

# Filter by start time
locals {
  start_time = timeadd(compact(data.huaweicloud_dataarts_quality_tasks.test.tasks[*].last_run_time)[0], "-1h")
}

output "start_time" {
  value = local.start_time
}

data "huaweicloud_dataarts_quality_tasks" "filter_by_start_time" {
  workspace_id = "%[1]s"

  start_time = local.start_time
}

locals {
  start_time_filter_result = [
    # If the result of the expression timecmp(A, B) is greater than 0, indicating that time A is after time B.
    for v in data.huaweicloud_dataarts_quality_tasks.filter_by_start_time.tasks[*].last_run_time : timecmp(v, local.start_time) > 0
  ]
}

output "is_start_time_filter_useful" {
  value = length(local.start_time_filter_result) > 0 && alltrue(local.start_time_filter_result)
}

# Filter by creator
locals {
  task_creator = data.huaweicloud_dataarts_quality_tasks.test.tasks[0].creator
}

data "huaweicloud_dataarts_quality_tasks" "filter_by_creator" {
  workspace_id = "%[1]s"

  creator = local.task_creator
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_quality_tasks.filter_by_creator.tasks[*].creator : v == local.task_creator
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
