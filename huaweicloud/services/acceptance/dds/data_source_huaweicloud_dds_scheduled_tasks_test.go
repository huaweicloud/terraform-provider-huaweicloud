package dds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDdsScheduledTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_scheduled_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSScheduledTasksEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDdsScheduledTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.job_name"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.instance_name"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.instance_status"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "schedules.0.job_status"),

					resource.TestCheckOutput("instance_id_filter_is_useful", "true"),
					resource.TestCheckOutput("job_status_filter_is_useful", "true"),
					resource.TestCheckOutput("job_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceDdsScheduledTasks_basic string = `
data "huaweicloud_dds_scheduled_tasks" "test" {}

locals {
  instance_id = data.huaweicloud_dds_scheduled_tasks.test.schedules.0.instance_id
  job_status  = data.huaweicloud_dds_scheduled_tasks.test.schedules.0.job_status
  job_name    = data.huaweicloud_dds_scheduled_tasks.test.schedules.0.job_name
}

data "huaweicloud_dds_scheduled_tasks" "filter" {
  instance_id = local.instance_id
  job_status  = local.job_status
  job_name    = local.job_name
}

output "instance_id_filter_is_useful" {
  value = length(data.huaweicloud_dds_scheduled_tasks.filter.schedules) > 0 && alltrue(
    [for v in data.huaweicloud_dds_scheduled_tasks.filter.schedules[*].instance_id : v == local.instance_id]
  )
}

output "job_status_filter_is_useful" {
  value = length(data.huaweicloud_dds_scheduled_tasks.filter.schedules) > 0 && alltrue(
    [for v in data.huaweicloud_dds_scheduled_tasks.filter.schedules[*].job_status : v == local.job_status]
  )
}

output "job_name_filter_is_useful" {
  value = length(data.huaweicloud_dds_scheduled_tasks.filter.schedules) > 0 && alltrue(
    [for v in data.huaweicloud_dds_scheduled_tasks.filter.schedules[*].job_name : v == local.job_name]
  )
}
`
