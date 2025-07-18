package workspace

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppScheduleTasks_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dataSourceName = "data.huaweicloud_workspace_app_schedule_tasks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byTaskName   = "data.huaweicloud_workspace_app_schedule_tasks.filter_by_task_name"
		dcByTaskName = acceptance.InitDataSourceCheck(byTaskName)

		byWeekTaskName   = "data.huaweicloud_workspace_app_schedule_tasks.filter_by_week_task_name"
		dcByWeekTaskName = acceptance.InitDataSourceCheck(byWeekTaskName)

		byMonthTaskName   = "data.huaweicloud_workspace_app_schedule_tasks.filter_by_month_task_name"
		dcByMonthTaskName = acceptance.InitDataSourceCheck(byMonthTaskName)

		byFixedTimeTaskName   = "data.huaweicloud_workspace_app_schedule_tasks.filter_by_fixed_time_task_name"
		dcByFixedTimeTaskName = acceptance.InitDataSourceCheck(byFixedTimeTaskName)

		byTaskType   = "data.huaweicloud_workspace_app_schedule_tasks.filter_by_task_type"
		dcByTaskType = acceptance.InitDataSourceCheck(byTaskType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppScheduleTasks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSourceName, "tasks.#", regexp.MustCompile(`^[0-9]+$`)),
					dcByTaskName.CheckResourceExists(),
					resource.TestCheckOutput("is_task_name_filter_useful", "true"),
					resource.TestCheckResourceAttr(byTaskName, "tasks.0.scheduled_type", "DAY"),
					resource.TestCheckResourceAttr(byTaskName, "tasks.0.is_enable", "true"),
					resource.TestCheckResourceAttr(byTaskName, "tasks.0.time_zone", "Asia/Shanghai"),
					resource.TestCheckResourceAttrPair(byTaskName, "tasks.0.scheduled_time",
						"huaweicloud_workspace_app_schedule_task.with_day", "scheduled_time"),
					resource.TestCheckResourceAttrPair(byTaskName, "tasks.0.day_interval",
						"huaweicloud_workspace_app_schedule_task.with_day", "day_interval"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.expire_time"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.description"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.next_execution_time"),
					resource.TestCheckResourceAttrSet(byTaskName, "tasks.0.task_cron"),
					resource.TestMatchResourceAttr(byTaskName, "tasks.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byTaskName, "tasks.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByWeekTaskName.CheckResourceExists(),
					resource.TestCheckResourceAttr(byWeekTaskName, "tasks.0.scheduled_type", "WEEK"),
					resource.TestCheckResourceAttrSet(byWeekTaskName, "tasks.0.week_list"),
					dcByMonthTaskName.CheckResourceExists(),
					resource.TestCheckResourceAttr(byMonthTaskName, "tasks.0.scheduled_type", "MONTH"),
					resource.TestCheckResourceAttrSet(byMonthTaskName, "tasks.0.month_list"),
					resource.TestCheckResourceAttrSet(byMonthTaskName, "tasks.0.date_list"),
					dcByFixedTimeTaskName.CheckResourceExists(),
					resource.TestCheckResourceAttr(byFixedTimeTaskName, "tasks.0.scheduled_type", "FIXED_TIME"),
					resource.TestCheckResourceAttrSet(byFixedTimeTaskName, "tasks.0.scheduled_date"),
					dcByTaskType.CheckResourceExists(),
					resource.TestCheckOutput("is_task_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAppScheduleTasks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_workspace_app_schedule_tasks" "test" {
  depends_on = [huaweicloud_workspace_app_schedule_task.with_day]
}

locals {
  task_name = huaweicloud_workspace_app_schedule_task.with_day.task_name
  task_type = huaweicloud_workspace_app_schedule_task.with_day.task_type
}

# Exact match by task name.
data "huaweicloud_workspace_app_schedule_tasks" "filter_by_task_name" {
  task_name = local.task_name

  depends_on = [huaweicloud_workspace_app_schedule_task.with_day]
}

data "huaweicloud_workspace_app_schedule_tasks" "filter_by_week_task_name" {
  task_name = huaweicloud_workspace_app_schedule_task.with_week.task_name

  depends_on = [huaweicloud_workspace_app_schedule_task.with_week]
}

data "huaweicloud_workspace_app_schedule_tasks" "filter_by_month_task_name" {
  task_name = huaweicloud_workspace_app_schedule_task.with_month.task_name

  depends_on = [huaweicloud_workspace_app_schedule_task.with_month]
}

data "huaweicloud_workspace_app_schedule_tasks" "filter_by_fixed_time_task_name" {
  task_name = huaweicloud_workspace_app_schedule_task.with_fixed_time.task_name

  depends_on = [huaweicloud_workspace_app_schedule_task.with_fixed_time]
}

locals {
  task_name_filter_result = [for v in data.huaweicloud_workspace_app_schedule_tasks.filter_by_task_name.tasks : v.task_name == local.task_name]
}

output "is_task_name_filter_useful" {
  value = length(local.task_type_filter_result) > 0 && alltrue(local.task_type_filter_result)
}

# Filter by task type.
data "huaweicloud_workspace_app_schedule_tasks" "filter_by_task_type" {
  task_type = local.task_type

  depends_on = [huaweicloud_workspace_app_schedule_task.with_day]
}

locals {
  task_type_filter_result = [for v in data.huaweicloud_workspace_app_schedule_tasks.filter_by_task_type.tasks : v.task_type == local.task_type]
}

output "is_task_type_filter_useful" {
  value = length(local.task_type_filter_result) > 0 && alltrue(local.task_type_filter_result)
}
`, testAccResourceAppScheduleTasks_base(name))
}

func testAccResourceAppScheduleTasks_base(name string) string {
	scheduledTime := time.Now().AddDate(0, 1, 0)
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = try(data.huaweicloud_workspace_service.test.network_ids[0], "")
  system_disk_type = "SATA"
  system_disk_size = 80
  is_vdi           = true
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
}

resource "huaweicloud_workspace_app_schedule_task" "with_day" {
  task_name      = "%[1]s-day"
  task_type      = "RESTART_SERVER"
  scheduled_time = "09:00:00"
  scheduled_type = "DAY"
  day_interval   = 2
  expire_time    = "%[5]s"
  description    = "Created by terraform script"

  target_infos {
    target_type = "SERVER_GROUP"
    target_id   = huaweicloud_workspace_app_server_group.test.id
  }
}

resource "huaweicloud_workspace_app_schedule_task" "with_week" {
  task_name      = "%[1]s-week"
  task_type      = "RESTART_SERVER"
  scheduled_time = "09:00:00"
  scheduled_type = "WEEK"
  week_list      = "7,1"
  description    = "Created by terraform script"

  target_infos {
    target_type = "SERVER_GROUP"
    target_id   = huaweicloud_workspace_app_server_group.test.id
  }
}

resource "huaweicloud_workspace_app_schedule_task" "with_month" {
  task_name      = "%[1]s-month"
  task_type      = "RESTART_SERVER"
  scheduled_time = "09:00:00"
  scheduled_type = "MONTH"
  month_list     = "%[6]s"
  date_list      = "24,28"
  description    = "Created by terraform script"
  is_enable      = false

  target_infos {
    target_type = "SERVER_GROUP"
    target_id   = huaweicloud_workspace_app_server_group.test.id
  }
}

resource "huaweicloud_workspace_app_schedule_task" "with_fixed_time" {
  task_name      = "%[1]s-fixed-time"
  task_type      = "RESTART_SERVER"
  scheduled_time = "09:00:00"
  scheduled_type = "FIXED_TIME"
  scheduled_date = "%[7]s"
  description    = "Created by terraform script"

  target_infos {
    target_type = "SERVER_GROUP"
    target_id   = huaweicloud_workspace_app_server_group.test.id
  }
}
`, name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		scheduledTime.Format("2006-01-02T15:04:05Z"),
		fmt.Sprintf("%d", scheduledTime.Month()),
		scheduledTime.Format("2006-01-02"))
}
