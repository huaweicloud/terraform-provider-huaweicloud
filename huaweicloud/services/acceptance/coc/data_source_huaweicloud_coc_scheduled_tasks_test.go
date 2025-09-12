package coc

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocScheduledTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_scheduled_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source:            "hashicorp/time",
				VersionConstraint: "0.12.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocScheduledTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.scheduled_type"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.task_type"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.associated_task_type"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.risk_level"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.created_by"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.update_by"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.created_user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.approve_status"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.execution_count"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.created_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.modified_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.region_id"),
					resource.TestCheckResourceAttrSet(dataSource, "scheduled_tasks.0.associated_task_name"),
					resource.TestCheckOutput("task_id_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
					resource.TestCheckOutput("task_name_filter_is_useful", "true"),
					resource.TestCheckOutput("scheduled_type_filter_is_useful", "true"),
					resource.TestCheckOutput("task_type_filter_is_useful", "true"),
					resource.TestCheckOutput("associated_task_type_filter_is_useful", "true"),
					resource.TestCheckOutput("risk_level_filter_is_useful", "true"),
					resource.TestCheckOutput("created_by_filter_is_useful", "true"),
					resource.TestCheckOutput("approve_status_filter_is_useful", "true"),
					resource.TestCheckOutput("region_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocScheduledTasks_basic(name string) string {
	currentTime := time.Now()
	tenMinutesAgo := currentTime.Add(-10*time.Minute).Unix() * 1e3
	tenMinutesLater := currentTime.Add(10*time.Minute).Unix() * 1e3
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_coc_scheduled_tasks" "test" {
  task_id = huaweicloud_coc_scheduled_task.test.id
}

output "task_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.test.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.test.scheduled_tasks[*].id : v == huaweicloud_coc_scheduled_task.test.id]
  )
}

data "huaweicloud_coc_scheduled_tasks" "enterprise_project_id_filter" {
  depends_on = [huaweicloud_coc_scheduled_task.test]

  enterprise_project_id = "0"
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.enterprise_project_id_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.enterprise_project_id_filter.scheduled_tasks[*].enterprise_project_id :
      v == "0"]
  )
}

data "huaweicloud_coc_scheduled_tasks" "task_name_filter" {
  task_name = huaweicloud_coc_scheduled_task.test.name
}

output "task_name_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.task_name_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.task_name_filter.scheduled_tasks[*].name :
      v == huaweicloud_coc_scheduled_task.test.name]
  )
}

data "huaweicloud_coc_scheduled_tasks" "scheduled_type_filter" {
  scheduled_type = huaweicloud_coc_scheduled_task.test.trigger_time[0].policy
}

output "scheduled_type_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.scheduled_type_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.scheduled_type_filter.scheduled_tasks[*].scheduled_type :
      v == huaweicloud_coc_scheduled_task.test.trigger_time[0].policy]
  )
}

data "huaweicloud_coc_scheduled_tasks" "task_type_filter" {
  task_type = huaweicloud_coc_scheduled_task.test.task_type
}

output "task_type_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.task_type_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.task_type_filter.scheduled_tasks[*].task_type :
      v == huaweicloud_coc_scheduled_task.test.task_type]
  )
}

data "huaweicloud_coc_scheduled_tasks" "associated_task_type_filter" {
  associated_task_type = huaweicloud_coc_scheduled_task.test.associated_task_type
}

output "associated_task_type_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.associated_task_type_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.associated_task_type_filter.scheduled_tasks[*].associated_task_type :
      v == huaweicloud_coc_scheduled_task.test.associated_task_type]
  )
}

data "huaweicloud_coc_scheduled_tasks" "risk_level_filter" {
  risk_level = huaweicloud_coc_scheduled_task.test.risk_level
}

output "risk_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.risk_level_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.risk_level_filter.scheduled_tasks[*].risk_level :
      v == huaweicloud_coc_scheduled_task.test.risk_level]
  )
}

locals {
  created_user_name = [for v in data.huaweicloud_coc_scheduled_tasks.test.scheduled_tasks[*].created_user_name :
    v if v != ""][0]
}

data "huaweicloud_coc_scheduled_tasks" "created_by_filter" {
  created_by = local.created_user_name
}

output "created_by_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.created_by_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.created_by_filter.scheduled_tasks[*].created_user_name :
      v == local.created_user_name]
  )
}

data "huaweicloud_coc_scheduled_tasks" "approve_status_filter" {
  approve_status = huaweicloud_coc_scheduled_task.test.approve_status
}

output "approve_status_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.approve_status_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.approve_status_filter.scheduled_tasks[*].approve_status :
      v == huaweicloud_coc_scheduled_task.test.approve_status]
  )
}

locals {
  region_id = [for v in data.huaweicloud_coc_scheduled_tasks.test.scheduled_tasks[*].region_id : v if v != ""][0]
}

data "huaweicloud_coc_scheduled_tasks" "region_id_filter" {
  region_id = local.region_id
}

output "region_id_filter_is_useful" {
  value = length(data.huaweicloud_coc_scheduled_tasks.region_id_filter.scheduled_tasks) > 0 && alltrue(
    [for v in data.huaweicloud_coc_scheduled_tasks.region_id_filter.scheduled_tasks[*].region_id : v == local.region_id]
  )
}
`, testScheduledTask_basic(name), tenMinutesAgo, tenMinutesLater)
}
