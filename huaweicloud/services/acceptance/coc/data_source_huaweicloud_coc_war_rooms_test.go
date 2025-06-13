package coc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCocWarRooms_basic(t *testing.T) {
	dataSource := "data.huaweicloud_coc_war_rooms.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCocWarRooms_basic(),
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
					resource.TestCheckResourceAttrSet(dataSource, "list.0.description"),
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

func testDataSourceDataSourceCocWarRooms_basic() string {
	return `
data "huaweicloud_coc_war_rooms" "test" {}

locals {
  incident_num = [for v in data.huaweicloud_coc_war_rooms.test.list[*].incident[0].incident_id : v if v != ""][0]
}

data "huaweicloud_coc_war_rooms" "incident_num_filter" {
  incident_num = local.incident_num
}

output "incident_num_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.incident_num_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.incident_num_filter.list[*].incident[0].incident_id :
      v == local.incident_num]
  )
}

locals {
  title = [for v in data.huaweicloud_coc_war_rooms.test.list[*].title : v if v != ""][0]
}

data "huaweicloud_coc_war_rooms" "title_filter" {
  title = local.title
}

output "title_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.title_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.title_filter.list[*].title : v == local.title]
  )
}

locals {
  region_code_list = [[for v in data.huaweicloud_coc_war_rooms.test.list[*].regions[*].code : v if v != ""][0][0]]
}

data "huaweicloud_coc_war_rooms" "region_code_list_filter" {
  region_code_list = local.region_code_list
}

output "region_code_list_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.region_code_list_filter.list) > 0 && alltrue(
    [for v in flatten([for v in data.huaweicloud_coc_war_rooms.region_code_list_filter.list[*].regions[*].code :
      v if v != ""]) : contains(local.region_code_list, v) ]
  )
}

data "huaweicloud_coc_war_rooms" "incident_levels_filter" {
  incident_levels = ["level_50"]
}

output "incident_levels_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.incident_levels_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.incident_levels_filter.list[*].circular_level : v == "level_50"]
  )
}

locals {
  impacted_application_ids = [[for v in data.huaweicloud_coc_war_rooms.test.list[*].impacted_application[0].id :
    v if v != ""][0]]
}

data "huaweicloud_coc_war_rooms" "impacted_application_ids_filter" {
  impacted_application_ids = local.impacted_application_ids
}

output "impacted_application_ids_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.impacted_application_ids_filter.list) > 0 && alltrue(
    [for v in flatten(
      [for v in data.huaweicloud_coc_war_rooms.impacted_application_ids_filter.list[*].impacted_application[0].id :
        v if v != ""]) : contains(local.impacted_application_ids, v) ]
  )
}

locals {
  admins = [[for v in data.huaweicloud_coc_war_rooms.test.list[*].admin : v if v != ""][0]]
}

data "huaweicloud_coc_war_rooms" "admin_filter" {
  admin = local.admins
}

output "admin_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.admin_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.admin_filter.list[*].admin : contains(local.admins, v)]
  )
}

data "huaweicloud_coc_war_rooms" "status_filter" {
  status = ["1"]
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
}

output "time_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.time_filter.list) > 0
}

data "huaweicloud_coc_war_rooms" "notification_level_filter" {
  notification_level = ["level_50"]
}

output "notification_level_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.notification_level_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.notification_level_filter.list[*].circular_level : v == "level_50"]
  )
}

data "huaweicloud_coc_war_rooms" "enterprise_project_ids_filter" {
  enterprise_project_ids = ["0"]
}

output "enterprise_project_ids_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.enterprise_project_ids_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.enterprise_project_ids_filter.list[*].enterprise_project_id : v == "0"]
  )
}

locals {
  war_room_num = [for v in data.huaweicloud_coc_war_rooms.test.list[*].war_room_num : v if v != ""][0]
}

data "huaweicloud_coc_war_rooms" "war_room_num_filter" {
  war_room_num = local.war_room_num
}

output "war_room_num_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.war_room_num_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.war_room_num_filter.list[*].war_room_num : v == local.war_room_num]
  )
}

locals {
  current_users = [[for v in data.huaweicloud_coc_war_rooms.test.list[*].admin : v if v != ""][0]]
}

data "huaweicloud_coc_war_rooms" "current_users_filter" {
  current_users = local.current_users
}

output "current_users_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.current_users_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.current_users_filter.list[*].admin : contains(local.current_users, v)]
  )
}

locals {
  war_room_nums = [[for v in data.huaweicloud_coc_war_rooms.test.list[*].war_room_num : v if v != ""][0]]
}

data "huaweicloud_coc_war_rooms" "war_room_nums_filter" {
  war_room_nums = local.war_room_nums
}

output "war_room_nums_filter_is_useful" {
  value = length(data.huaweicloud_coc_war_rooms.war_room_nums_filter.list) > 0 && alltrue(
    [for v in data.huaweicloud_coc_war_rooms.war_room_nums_filter.list[*].war_room_num :
      contains(local.war_room_nums, v)]
  )
}
`
}
