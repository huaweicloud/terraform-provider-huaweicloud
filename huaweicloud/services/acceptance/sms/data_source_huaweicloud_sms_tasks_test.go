package sms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSmsTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_sms_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceSmsTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.create_date"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.priority"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.speed_limit"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.migrate_speed"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.compress_rate"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.start_target_server"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.error_json"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.total_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.migration_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.sub_tasks.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.sub_tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.sub_tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.sub_tasks.0.progress"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.sub_tasks.0.start_date"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.sub_tasks.0.end_date"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.0.ip"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.0.os_version"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.0.oem_system"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.source_server.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.target_server.#"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.target_server.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.target_server.0.vm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.target_server.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.target_server.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.log_collect_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.syncing"),
					resource.TestCheckOutput("state_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("id_filter_is_useful", "true"),
					resource.TestCheckOutput("source_server_id_filter_is_useful", "true"),
					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceSmsTasks_basic() string {
	return `
data "huaweicloud_sms_tasks" "test" {}

locals {
  state = [for v in data.huaweicloud_sms_tasks.test.tasks[*].state : v if v != ""][0]
}

data "huaweicloud_sms_tasks" "state_filter" {
  state = local.state
}

output "state_filter_is_useful" {
  value = length(data.huaweicloud_sms_tasks.state_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_sms_tasks.state_filter.tasks[*].state : v == local.state]
  )
}

locals {
  name = [for v in data.huaweicloud_sms_tasks.test.tasks[*].name : v if v != ""][0]
}

data "huaweicloud_sms_tasks" "name_filter" {
  name = local.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_sms_tasks.name_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_sms_tasks.name_filter.tasks[*].name : v == local.name]
  )
}

locals {
  task_id = [for v in data.huaweicloud_sms_tasks.test.tasks[*].id : v if v != ""][0]
}

data "huaweicloud_sms_tasks" "id_filter" {
  task_id = local.task_id
}

output "id_filter_is_useful" {
  value = length(data.huaweicloud_sms_tasks.id_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_sms_tasks.id_filter.tasks[*].id : v == local.task_id]
  )
}

locals {
  source_server_id = [for v in data.huaweicloud_sms_tasks.test.tasks[*].source_server[0].id : v if v != ""][0]
}

data "huaweicloud_sms_tasks" "source_server_id_filter" {
  source_server_id = local.source_server_id
}

output "source_server_id_filter_is_useful" {
  value = length(data.huaweicloud_sms_tasks.source_server_id_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_sms_tasks.source_server_id_filter.tasks[*].source_server[0].id : v == local.source_server_id]
  )
}

locals {
  enterprise_project_id = [for v in data.huaweicloud_sms_tasks.test.tasks[*].enterprise_project_id : v if v != ""][0]
}

data "huaweicloud_sms_tasks" "enterprise_project_id_filter" {
  enterprise_project_id = local.enterprise_project_id
}

output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_sms_tasks.enterprise_project_id_filter.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_sms_tasks.enterprise_project_id_filter.tasks[*].enterprise_project_id :
      v == local.enterprise_project_id]
  )
}
`
}
