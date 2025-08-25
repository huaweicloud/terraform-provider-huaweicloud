package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocWarRooms_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_war_rooms.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUserId(t)
			acceptance.TestAccPreCheckCocRoleID(t)
			acceptance.TestAccPreCheckCocSceneID(t)
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
				Config: testDataSourceDataSourceCocWarRooms_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.admin"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.incident.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.incident.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.incident.0.incident_id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.incident.0.is_change_event"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.incident.0.failure_level"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.regions.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.occur_time"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.circular_level"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.war_room_status.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.war_room_status.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.war_room_status.0.name_zh"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.war_room_status.0.name_en"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.war_room_status.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.impacted_application.#"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.impacted_application.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.war_room_num"),
					resource.TestCheckResourceAttrSet(dataSource, "list.0.enterprise_project_id"),
					resource.TestCheckOutput("admin_filter_is_useful", "true"),
					resource.TestCheckOutput("current_users_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_ids_filter_is_useful", "true"),
					resource.TestCheckOutput("impacted_application_ids_filter_is_useful", "true"),
					resource.TestCheckOutput("incident_levels_filter_is_useful", "true"),
					resource.TestCheckOutput("incident_num_filter_is_useful", "true"),
					resource.TestCheckOutput("notification_level_filter_is_useful", "true"),
					resource.TestCheckOutput("region_code_list_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("time_filter_is_useful", "true"),
					resource.TestCheckOutput("title_filter_is_useful", "true"),
					resource.TestCheckOutput("war_room_num_filter_is_useful", "true"),
					resource.TestCheckOutput("war_room_nums_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCocWarRooms_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_coc_war_room" "test" {
  war_room_name         = "%[2]s"
  application_id_list   = [huaweicloud_coc_application.test.id]
  incident_number       = huaweicloud_coc_incident.test.id
  war_room_admin        = "%[3]s"
  enterprise_project_id = "0"
  region_code_list      = ["%[4]s"]

  schedule_group {
    role_id  = "%[5]s"
    scene_id = "%[6]s"
  }

  depends_on = [huaweicloud_coc_incident_handle.test]
}

data "huaweicloud_coc_war_rooms" "test" {
  depends_on = [huaweicloud_coc_incident.test]
}

data "huaweicloud_coc_war_rooms" "incident_num_filter" {
  incident_num = huaweicloud_coc_incident.test.id
  depends_on   = [huaweicloud_coc_war_room.test]
}

output "incident_num_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.incident_num_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.incident_num_filter.list[*].incident[0].incident_id :
      v == huaweicloud_coc_incident.test.id]
  )
}

data "huaweicloud_coc_war_rooms" "title_filter" {
  title = huaweicloud_coc_war_room.test.war_room_name
}

output "title_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.title_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.title_filter.list[*].title :
      v == huaweicloud_coc_war_room.test.war_room_name]
  )
}

data "huaweicloud_coc_war_rooms" "region_code_list_filter" {
  region_code_list = huaweicloud_coc_war_room.test.region_code_list
}

output "region_code_list_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.region_code_list_filter.list) > 0 && alltrue(
    [for v in flatten([for v in data.huaweicloud_coc_war_rooms.region_code_list_filter.list[*].regions[*].code :
      v if v != ""]) : contains(huaweicloud_coc_war_room.test.region_code_list, v) ]
  )
}

data "huaweicloud_coc_war_rooms" "incident_levels_filter" {
  incident_levels = ["level_50"]
  depends_on      = [huaweicloud_coc_war_room.test]
}

output "incident_levels_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.incident_levels_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.incident_levels_filter.list[*].circular_level : v == "level_50"]
  )
}

data "huaweicloud_coc_war_rooms" "impacted_application_ids_filter" {
  impacted_application_ids = huaweicloud_coc_war_room.test.application_id_list
}

output "impacted_application_ids_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.impacted_application_ids_filter.list) > 0 && alltrue(
    [for v in [for v in data.huaweicloud_coc_war_rooms.impacted_application_ids_filter.list[*].impacted_application[*].id :
        v if v != []] : contains(distinct(v), flatten(huaweicloud_coc_war_room.test.application_id_list)[0]) ]
  )
}

data "huaweicloud_coc_war_rooms" "admin_filter" {
  admin = [huaweicloud_coc_war_room.test.war_room_admin]
}

output "admin_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.admin_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.admin_filter.list[*].admin :
      contains([huaweicloud_coc_war_room.test.war_room_admin], v)]
  )
}

data "huaweicloud_coc_war_rooms" "status_filter" {
  status     = ["1"]
  depends_on = [huaweicloud_coc_war_room.test]
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.status_filter.list) > 0 && alltrue(
    [for v in flatten([for v in data.huaweicloud_coc_war_rooms.status_filter.list[*].war_room_status[0].id :
      v if v != ""]) : v == "1" ]
  )
}

data "huaweicloud_coc_war_rooms" "time_filter" {
  triggered_start_time = 1
  triggered_end_time   = 9223372036854775807
  occur_start_time     = 1
  occur_end_time       = 9223372036854775807
  recover_start_time   = 1
  recover_end_time     = 9223372036854775807
  depends_on           = [huaweicloud_coc_war_room.test]
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.time_filter.list) > 0
}

data "huaweicloud_coc_war_rooms" "notification_level_filter" {
  notification_level = ["level_50"]
  depends_on         = [huaweicloud_coc_war_room.test]
}

output "notification_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.notification_level_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.notification_level_filter.list[*].circular_level : v == "level_50"]
  )
}

data "huaweicloud_coc_war_rooms" "enterprise_project_ids_filter" {
  enterprise_project_ids = ["0"]
  depends_on             = [huaweicloud_coc_war_room.test]
}

output "enterprise_project_ids_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.enterprise_project_ids_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.enterprise_project_ids_filter.list[*].enterprise_project_id : v == "0"]
  )
}

data "huaweicloud_coc_war_rooms" "war_room_num_filter" {
  war_room_num = huaweicloud_coc_war_room.test.id
}

output "war_room_num_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.war_room_num_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.war_room_num_filter.list[*].war_room_num :
      v == huaweicloud_coc_war_room.test.id]
  )
}

data "huaweicloud_coc_war_rooms" "current_users_filter" {
  current_users = [huaweicloud_coc_war_room.test.war_room_admin]
}

output "current_users_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.current_users_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.current_users_filter.list[*].admin :
      contains([huaweicloud_coc_war_room.test.war_room_admin], v)]
  )
}

data "huaweicloud_coc_war_rooms" "war_room_nums_filter" {
  war_room_nums = [huaweicloud_coc_war_room.test.id]
}

output "war_room_nums_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.war_room_nums_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.war_room_nums_filter.list[*].war_room_num :
      contains([huaweicloud_coc_war_room.test.id], v)]
  )
}
`, testIncidentHandle_basic(name), name, acceptance.HW_USER_ID, acceptance.HW_REGION_NAME, acceptance.HW_COC_ROLE_ID,
		acceptance.HW_COC_SCENE_ID)
}
