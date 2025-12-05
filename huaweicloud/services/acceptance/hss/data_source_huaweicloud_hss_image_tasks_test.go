package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceImageTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_image_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{

			{
				Config: testAccDataSourceImageTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.task_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.task_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.failed_images.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.failed_images.0.registry_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.failed_images.0.registry_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.failed_images.0.registry_type"),

					resource.TestCheckOutput("task_type_filter_useful", "true"),
					resource.TestCheckOutput("task_id_filter_useful", "true"),
					resource.TestCheckOutput("time_filter_useful", "true"),
					resource.TestCheckOutput("task_status_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceImageTasks_basic = `
data "huaweicloud_hss_image_tasks" "test" {
  type = "image_sync"
}

locals {
  task_type = data.huaweicloud_hss_image_tasks.test.data_list[0].task_type
  task_id   = data.huaweicloud_hss_image_tasks.test.data_list[0].task_id
}

data "huaweicloud_hss_image_tasks" "task_type_filter" {
  type      = "image_sync"
  task_type = local.task_type
}

output "task_type_filter_useful" {
  value = length(data.huaweicloud_hss_image_tasks.task_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_tasks.task_type_filter.data_list[*].task_type : v == local.task_type]
  )
}

data "huaweicloud_hss_image_tasks" "task_id_filter" {
  type    = "image_sync"
  task_id = local.task_id
}

output "task_id_filter_useful" {
  value = length(data.huaweicloud_hss_image_tasks.task_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_image_tasks.task_id_filter.data_list[*].task_id : v == local.task_id]
  )
}

locals {
  create_time = data.huaweicloud_hss_image_tasks.test.data_list[0].begin_time
  end_time    = data.huaweicloud_hss_image_tasks.test.data_list[0].end_time
}

data "huaweicloud_hss_image_tasks" "time_filter" {
  type        = "image_sync"
  create_time = local.create_time
  end_time    = local.end_time
}

output "time_filter_useful" {
  value = length(data.huaweicloud_hss_image_tasks.time_filter.data_list) > 0
}

data "huaweicloud_hss_image_tasks" "task_status_filter" {
  type        = "image_sync"
  task_status = "finished"
}

output "task_status_filter_useful" {
  value = length(data.huaweicloud_hss_image_tasks.task_status_filter.data_list) > 0
}

data "huaweicloud_hss_image_tasks" "eps_filter" {
  type                  = "image_sync"
  enterprise_project_id = "0"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_image_tasks.eps_filter.data_list) > 0
}
`
