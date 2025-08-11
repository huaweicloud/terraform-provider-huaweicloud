package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppScheduleTaskExecuteDetails_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_workspace_app_schedule_task_execute_details.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppScheduleTaskExecuteDetails_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "execute_details.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "execute_details.0.id"),
					resource.TestCheckResourceAttrPair(dataSource, "execute_details.0.execute_id",
						"data.huaweicloud_workspace_app_schedule_task_executions.test", "executions.0.id"),
					resource.TestCheckResourceAttrPair(dataSource, "execute_details.0.server_id",
						"huaweicloud_workspace_app_server.test", "id"),
					resource.TestCheckResourceAttrPair(dataSource, "execute_details.0.server_name",
						"huaweicloud_workspace_app_server.test", "name"),
					resource.TestCheckResourceAttrPair(dataSource, "execute_details.0.task_type",
						"huaweicloud_workspace_app_schedule_task.test", "task_type"),
					resource.TestCheckResourceAttrPair(dataSource, "execute_details.0.time_zone",
						"huaweicloud_workspace_app_schedule_task.test", "time_zone"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_details.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "execute_details.0.begin_time"),
					// The task may not be completed yet, `end_time` is not checked.
				),
			},
		},
	})
}

func testAccDataAppScheduleTaskExecuteDetails_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = try(data.huaweicloud_workspace_service.test.network_ids[0], null)
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = true
  app_type         = "COMMON_APP"
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
}

resource "huaweicloud_workspace_app_server" "test" {
  name                = "%[1]s" 
  server_group_id     = huaweicloud_workspace_app_server_group.test.id
  type                = "createApps"
  flavor_id           = huaweicloud_workspace_app_server_group.test.flavor_id
  vpc_id              = huaweicloud_workspace_app_server_group.test.vpc_id
  subnet_id           = huaweicloud_workspace_app_server_group.test.subnet_id
  update_access_agent = false

  root_volume {
    type = huaweicloud_workspace_app_server_group.test.system_disk_type
    size = huaweicloud_workspace_app_server_group.test.system_disk_size
  }
}

resource "huaweicloud_workspace_app_schedule_task" "test" {
  task_name      = "%[1]s"
  task_type      = "STOP_SERVER"
  scheduled_time = formatdate("HH:mm:ss", timeadd(timestamp(), "1m"))
  scheduled_type = "DAY"
  day_interval   = 1
  time_zone      = "Coordinated Universal Time"

  target_infos {
    target_type = "SERVER"
    target_id   = huaweicloud_workspace_app_server.test.id
  }

  # The 'scheduled_time' field is generated using a timestamp, so its changes should be ignored.
  lifecycle {
    ignore_changes = [scheduled_time]
  }
}

# Wait for the scheduled task to be executed.
resource "null_resource" "test" {
  provisioner "local-exec" {
    command = "sleep 90"
  }

  depends_on = [huaweicloud_workspace_app_schedule_task.test]
}

data "huaweicloud_workspace_app_schedule_task_executions" "test" {
  depends_on = [null_resource.test]

  task_id = huaweicloud_workspace_app_schedule_task.test.id
}
`, name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testAccDataAppScheduleTaskExecuteDetails_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_schedule_task_execute_details" "test" {
  depends_on = [null_resource.test]

  execute_history_id = try(data.huaweicloud_workspace_app_schedule_task_executions.test.executions[0].id, "")
}
`, testAccDataAppScheduleTaskExecuteDetails_base(name))
}
