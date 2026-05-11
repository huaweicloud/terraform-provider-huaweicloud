package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsEvents_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_events.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsEvents_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_count"),
					resource.TestCheckResourceAttrSet(dataSource, "inquiring_count"),
					resource.TestCheckResourceAttrSet(dataSource, "schedule_count"),
					resource.TestCheckResourceAttrSet(dataSource, "executing_count"),
					resource.TestCheckResourceAttrSet(dataSource, "failed_count"),
					resource.TestCheckResourceAttrSet(dataSource, "events.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.db_type"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.created_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.impact"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.reason"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.latest_execution_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execution_time_window.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execution_time_window.0.planned_execution_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execution_time_window.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execution_time_window.0.end_time"),
					resource.TestCheckOutput("event_id_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("level_filter", "true"),
					resource.TestCheckOutput("sort_filter", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsEvents_basic() string {
	return `
data "huaweicloud_rds_events" "test" {}

data "huaweicloud_rds_events" "event_id_filter" {
  event_id = data.huaweicloud_rds_events.test.events[0].id
}
locals {
  event_id = data.huaweicloud_rds_events.test.events[0].id
}
output "event_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_events.event_id_filter.events) > 0 && alltrue(
  [for v in data.huaweicloud_rds_events.event_id_filter.events[*].id : v == local.event_id]
  )
}

data "huaweicloud_rds_events" "instance_id_filter" {
  instance_id = data.huaweicloud_rds_events.test.events[0].id
}
locals {
  instance_id = data.huaweicloud_rds_events.test.events[0].instance_id
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_events.instance_id_filter.events) > 0 && alltrue(
  [for v in data.huaweicloud_rds_events.instance_id_filter.events[*].instance_id : v == local.instance_id]
  )
}

data "huaweicloud_rds_events" "status_filter" {
  status = data.huaweicloud_rds_events.test.events[0].status
}
locals {
  status = data.huaweicloud_rds_events.test.events[0].status
}
output "status_filter_is_useful" {
  value = length(data.huaweicloud_rds_events.status_filter.events) > 0 && alltrue(
  [for v in data.huaweicloud_rds_events.status_filter.events[*].status : v == local.status]
  )
}

data "huaweicloud_rds_events" "type_filter" {
  type = data.huaweicloud_rds_events.test.events[0].type
}
locals {
  type = data.huaweicloud_rds_events.test.events[0].type
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_rds_events.type_filter.events) > 0 && alltrue(
  [for v in data.huaweicloud_rds_events.type_filter.events[*].type : v == local.type]
  )
}

data "huaweicloud_rds_events" "level_filter" {
  level = data.huaweicloud_rds_events.test.events[0].level
}
locals {
  level = data.huaweicloud_rds_events.test.events[0].level
}
output "level_filter_is_useful" {
  value = length(data.huaweicloud_rds_events.level_filter.events) > 0 && alltrue(
  [for v in data.huaweicloud_rds_events.level_filter.events[*].level : v == local.level]
  )
}

data "huaweicloud_rds_events" "sort_filter" {
  sort_field = "latest_execution_time"
  order      = "ASC"
}
output "sort_filter_is_useful" {
  value = length(data.huaweicloud_rds_events.sort_filter.events) > 0
}
`
}
