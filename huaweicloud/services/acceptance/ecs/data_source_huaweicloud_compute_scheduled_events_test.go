package ecs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEcsComputeScheduledEvents_basic(t *testing.T) {
	dataSource := "data.huaweicloud_compute_scheduled_events.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckECSID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEcsComputeScheduledEvents_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "events.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.publish_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.finish_time"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.not_before"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.not_after"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.not_before_deadline"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_options.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_options.0.device"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_options.0.wwn"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_options.0.serial_number"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_options.0.resize_target_flavor_id"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_options.0.migrate_policy"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.execute_options.0.executor"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.source.#"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.source.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "events.0.source.0.host_scheduled_event_id"),
					resource.TestCheckOutput("event_id_filter_is_useful", "true"),
					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("state_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEcsComputeScheduledEvents_basic() string {
	return `
data "huaweicloud_compute_scheduled_events" "test" {}

locals {
  event_id = data.huaweicloud_compute_scheduled_events.test.events[0].id
} 
data "huaweicloud_compute_scheduled_events" "event_id_filter" {
  event_id = data.huaweicloud_compute_scheduled_events.test.events[0].id
}
output "event_id_filter_is_useful" {
  value = length(data.huaweicloud_compute_scheduled_events.event_id_filter.events) > 0 alltrue(
  [for v in data.huaweicloud_compute_scheduled_events.event_id_filter.events[*].id : v == local.event_id]
  )
}

locals {
  instance_id = data.huaweicloud_compute_scheduled_events.test.events[0].instance_id
} 
data "huaweicloud_compute_scheduled_events" "instance_id_filter" {
  instance_id = [data.huaweicloud_compute_scheduled_events.test.events[0].instance_id]
}
output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_compute_scheduled_events.instance_id_filter.events) > 0 alltrue(
  [for v in data.huaweicloud_compute_scheduled_events.instance_id_filter.events[*].instance_id : v == local.instance_id]
  )
}

locals {
  type = data.huaweicloud_compute_scheduled_events.test.events[0].type
} 
data "huaweicloud_compute_scheduled_events" "type_filter" {
  type = [data.huaweicloud_compute_scheduled_events.test.events[0].type]
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_compute_scheduled_events.type_filter.events) > 0 alltrue(
  [for v in data.huaweicloud_compute_scheduled_events.type_filter.events[*].instance_id : v == local.instance_id]
  )
}

locals {
  state = data.huaweicloud_compute_scheduled_events.test.events[0].state
} 
data "huaweicloud_compute_scheduled_events" "state_filter" {
  state = [data.huaweicloud_compute_scheduled_events.test.events[0].state]
}
output "state_filter_is_useful" {
  value = length(data.huaweicloud_compute_scheduled_events.state_filter.events) > 0 alltrue(
  [for v in data.huaweicloud_compute_scheduled_events.state_filter.events[*].state : v == local.state]
  )
}
`
}
