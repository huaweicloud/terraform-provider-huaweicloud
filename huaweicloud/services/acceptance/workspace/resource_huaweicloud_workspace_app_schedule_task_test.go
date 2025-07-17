package workspace

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getResourceAppScheduleTaskFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}

	return workspace.GetAppScheduleTaskById(client, state.Primary.ID)
}

func TestAccResourceAppScheduleTask_basic(t *testing.T) {
	var (
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		withDayToMonth        = "huaweicloud_workspace_app_schedule_task.with_day_to_month"
		withWeekToFixedTime   = "huaweicloud_workspace_app_schedule_task.with_week_to_fixed_time"
		scheduleTask          interface{}
		rcWithDayToMonth      = acceptance.InitResourceCheck(withDayToMonth, &scheduleTask, getResourceAppScheduleTaskFunc)
		rcWithWeekToFixedTime = acceptance.InitResourceCheck(withWeekToFixedTime, &scheduleTask, getResourceAppScheduleTaskFunc)

		scheduledTime = time.Now().AddDate(0, 1, 0)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithDayToMonth.CheckResourceDestroy(),
			rcWithWeekToFixedTime.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppScheduleTask_basic_step1(name, scheduledTime),
				Check: resource.ComposeTestCheckFunc(
					rcWithDayToMonth.CheckResourceExists(),
					resource.TestCheckResourceAttr(withDayToMonth, "task_name", name),
					resource.TestCheckResourceAttr(withDayToMonth, "task_type", "STOP_SERVER"),
					resource.TestCheckResourceAttr(withDayToMonth, "scheduled_type", "DAY"),
					resource.TestCheckResourceAttr(withDayToMonth, "scheduled_time", "00:00:00"),
					resource.TestCheckResourceAttr(withDayToMonth, "day_interval", "2"),
					resource.TestCheckResourceAttr(withDayToMonth, "time_zone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(withDayToMonth, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withDayToMonth, "schedule_task_policy.#", "0"),
					resource.TestCheckResourceAttr(withDayToMonth, "target_infos.#", "2"),
					resource.TestCheckResourceAttr(withDayToMonth, "is_enable", "true"),
					rcWithWeekToFixedTime.CheckResourceExists(),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "task_type", "RESTART_SERVER"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "scheduled_type", "WEEK"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "scheduled_time", "02:00:00"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "week_list", "1,2,3,4,5"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "time_zone", "Asia/Shanghai"),
					resource.TestCheckResourceAttrSet(withWeekToFixedTime, "expire_time"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "schedule_task_policy.0.enforcement_enable", "false"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "target_infos.#", "1"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "target_infos.0.target_type", "SERVER_GROUP"),
					resource.TestCheckResourceAttrPair(withWeekToFixedTime, "target_infos.0.target_id",
						"huaweicloud_workspace_app_server_group.test.0", "id"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "is_enable", "false"),
				),
			},
			{
				Config: testAccResourceAppScheduleTask_basic_step2(name, updateName, scheduledTime),
				Check: resource.ComposeTestCheckFunc(
					rcWithDayToMonth.CheckResourceExists(),
					resource.TestCheckResourceAttr(withDayToMonth, "task_name", updateName),
					resource.TestCheckResourceAttr(withDayToMonth, "task_type", "START_SERVER"),
					resource.TestCheckResourceAttr(withDayToMonth, "scheduled_type", "MONTH"),
					resource.TestCheckResourceAttr(withDayToMonth, "scheduled_time", "01:00:00"),
					resource.TestCheckResourceAttrSet(withDayToMonth, "month_list"),
					resource.TestCheckResourceAttr(withDayToMonth, "week_list", ""),
					resource.TestCheckResourceAttr(withDayToMonth, "date_list", "L"),
					resource.TestCheckResourceAttr(withDayToMonth, "scheduled_date", ""),
					resource.TestCheckResourceAttr(withDayToMonth, "time_zone", "Asia/Shanghai"),
					resource.TestCheckResourceAttr(withDayToMonth, "description", ""),
					resource.TestCheckResourceAttr(withDayToMonth, "target_infos.#", "1"),
					resource.TestCheckResourceAttrPair(withDayToMonth, "target_infos.0.target_id",
						"huaweicloud_workspace_app_server_group.test.1", "id"),
					resource.TestCheckResourceAttr(withDayToMonth, "is_enable", "false"),
					rcWithWeekToFixedTime.CheckResourceExists(),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "task_type", "REINSTALL_OS"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "scheduled_type", "FIXED_TIME"),
					resource.TestCheckResourceAttrSet(withWeekToFixedTime, "scheduled_date"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "week_list", ""),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "month_list", ""),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "date_list", ""),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "expire_time", ""),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "description", "Updated by terraform script"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "schedule_task_policy.0.enforcement_enable", "true"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "target_infos.#", "1"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "target_infos.0.target_type", "SERVER_GROUP"),
					resource.TestCheckResourceAttrPair(withWeekToFixedTime, "target_infos.0.target_id",
						"huaweicloud_workspace_app_server_group.test.1", "id"),
					resource.TestCheckResourceAttr(withWeekToFixedTime, "is_enable", "true"),
				),
			},
			{
				ResourceName:      withDayToMonth,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				ResourceName:      withWeekToFixedTime,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceAppScheduleTask_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  count = 2

  name             = "%[1]s${count.index}"
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
`, name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testAccResourceAppScheduleTask_basic_step1(name string, scheduledTime time.Time) string {
	expireTime := scheduledTime.Format("2006-01-02T15:04:05Z")
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_schedule_task" "with_day_to_month" {
  task_name      = "%[2]s"
  task_type      = "STOP_SERVER"
  scheduled_type = "DAY"
  scheduled_time = "00:00:00"
  day_interval   = 2
  description    = "Created by terraform script"

  dynamic "target_infos" {
    for_each = huaweicloud_workspace_app_server_group.test[*].id

    content {
      target_type = "SERVER_GROUP"
      target_id   = target_infos.value
    }
  }
}

resource "huaweicloud_workspace_app_schedule_task" "with_week_to_fixed_time" {
  task_name      = "%[2]s_week"
  task_type      = "RESTART_SERVER"
  scheduled_type = "WEEK"
  scheduled_time = "02:00:00"
  week_list      = "1,2,3,4,5"
  time_zone      = "Asia/Shanghai"
  description    = "Created by terraform script"
  expire_time    = "%[3]s"
  is_enable      = false

  schedule_task_policy {
    enforcement_enable = false
  }

  target_infos {
    target_type = "SERVER_GROUP"
    target_id   = huaweicloud_workspace_app_server_group.test[0].id
  }
}
`, testAccResourceAppScheduleTask_base(name), name, expireTime)
}

func testAccResourceAppScheduleTask_basic_step2(name, updateName string, scheduledTime time.Time) string {
	scheduledMonth := fmt.Sprintf("%d", scheduledTime.Month())
	scheduledDate := scheduledTime.Format("2006-01-02")
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_schedule_task" "with_day_to_month" {
  task_name      = "%[2]s"
  task_type      = "START_SERVER"
  scheduled_type = "MONTH"
  scheduled_time = "01:00:00"
  month_list     = "%[3]s"
  date_list      = "L"
  is_enable      = false

  schedule_task_policy {
    enforcement_enable = true
  }

  target_infos {
    target_type = "SERVER_GROUP"
    target_id   = huaweicloud_workspace_app_server_group.test[1].id
  }
}

resource "huaweicloud_workspace_app_schedule_task" "with_week_to_fixed_time" {
  task_name      = "%[2]s_fixed"
  task_type      = "REINSTALL_OS"
  scheduled_type = "FIXED_TIME"
  scheduled_time = "03:00:00"
  scheduled_date = "%[4]s"
  time_zone      = "Asia/Shanghai"
  description    = "Updated by terraform script"

  schedule_task_policy {
    enforcement_enable = true
  }

  target_infos {
    target_type = "SERVER_GROUP"
    target_id   = huaweicloud_workspace_app_server_group.test[1].id
  }
}
`, testAccResourceAppScheduleTask_base(name), updateName, scheduledMonth, scheduledDate)
}
