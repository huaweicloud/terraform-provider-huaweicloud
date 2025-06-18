package dew

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCsmsTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_csms_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the csms task before running the test
			acceptance.TestAccPrecheckCsmsTask(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCsmsTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.secret_name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_status"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.operate_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.rotation_func_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_error_code"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_error_msg"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.attempt_nums"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.task_time"),

					resource.TestCheckOutput("is_secret_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_task_id_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceCsmsTasks_basic = `
data "huaweicloud_csms_tasks" "test" {}

locals {
  secret_name = data.huaweicloud_csms_tasks.test.tasks.0.secret_name
  status      = data.huaweicloud_csms_tasks.test.tasks.0.task_status
  task_id     = data.huaweicloud_csms_tasks.test.tasks.0.task_id
}

# Filter by secret_name
data "huaweicloud_csms_tasks" "secret_name_filter" {
  secret_name = local.secret_name
}

locals {
  secret_name_filter_result = [
    for v in data.huaweicloud_csms_tasks.secret_name_filter.tasks[*].secret_name : v == local.secret_name
  ]
}

output "is_secret_name_filter_useful" {
  value = length(local.secret_name_filter_result) > 0 && alltrue(local.secret_name_filter_result)
}

# Filter by status
data "huaweicloud_csms_tasks" "status_filter" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_csms_tasks.status_filter.tasks[*].task_status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

# Filter by task_id
data "huaweicloud_csms_tasks" "task_id_filter" {
  task_id = local.task_id
}

locals {
  task_id_filter_result = [
    for v in data.huaweicloud_csms_tasks.task_id_filter.tasks[*].task_id : v == local.task_id
  ]
}

output "is_task_id_filter_useful" {
  value = length(local.task_id_filter_result) > 0 && alltrue(local.task_id_filter_result)
}
`
