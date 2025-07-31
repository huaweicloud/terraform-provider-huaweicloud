package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAntivirusVirusScanTasks_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_antivirus_virus_scan_tasks.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires preparation of virus scanning task data.
			acceptance.TestAccPreCheckHSSAntivirusEnabled(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAntivirusVirusScanTasks_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.task_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.scan_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.start_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.action"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.task_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.success_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.fail_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cancel_host_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.run_duration"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.scan_progress"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.virus_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.scan_file_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.host_task_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_info_list.0.whether_using_quota"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rescan"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.whether_paid_task"),

					resource.TestCheckOutput("is_task_name_filter_useful", "true"),
					resource.TestCheckOutput("is_last_days_filter_useful", "true"),
					resource.TestCheckOutput("is_task_status_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAntivirusVirusScanTasks_basic string = `
data "huaweicloud_hss_antivirus_virus_scan_tasks" "test" {
  whether_paid_task = false
}

# Filter using task_name.
locals {
  task_name = data.huaweicloud_hss_antivirus_virus_scan_tasks.test.data_list[0].task_name
}

data "huaweicloud_hss_antivirus_virus_scan_tasks" "task_name_filter" {
  whether_paid_task = false
  task_name         = local.task_name
}

# The task_name is fuzzy match.
output "is_task_name_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_virus_scan_tasks.task_name_filter.data_list) > 0
}

# Filter using last_days.
data "huaweicloud_hss_antivirus_virus_scan_tasks" "last_days_filter" {
  whether_paid_task = false
  last_days         = 30
}

output "is_last_days_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_virus_scan_tasks.last_days_filter.data_list) > 0
}

# Filter using task_status.
locals {
  task_status = data.huaweicloud_hss_antivirus_virus_scan_tasks.test.data_list[0].task_status
}

data "huaweicloud_hss_antivirus_virus_scan_tasks" "task_status_filter" {
  whether_paid_task = false
  task_status       = local.task_status
}

output "is_task_status_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_virus_scan_tasks.task_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_antivirus_virus_scan_tasks.task_status_filter.data_list[*].task_status : v == local.task_status]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_antivirus_virus_scan_tasks" "enterprise_project_id_filter" {
  whether_paid_task     = false
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_antivirus_virus_scan_tasks.enterprise_project_id_filter.data_list) > 0
}

# Filter using non existent task_name.
data "huaweicloud_hss_antivirus_virus_scan_tasks" "not_found" {
  whether_paid_task = false
  task_name         = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_antivirus_virus_scan_tasks.not_found.data_list) == 0
}
`
