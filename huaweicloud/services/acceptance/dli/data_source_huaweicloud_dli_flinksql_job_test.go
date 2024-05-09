package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliFlinkSQLJobs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dli_flinksql_jobs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	rName := acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliFlinkSQLJobs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.queue_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.cu_num"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.parallel_num"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.manager_cu_num"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.tm_cu_num"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.tm_slot_num"),

					resource.TestCheckOutput("job_id_filter_is_useful", "true"),
					resource.TestCheckOutput("queue_name_filter_is_useful", "true"),
					resource.TestCheckOutput("cu_num_filter_is_useful", "true"),
					resource.TestCheckOutput("parallel_num_filter_is_useful", "true"),
					resource.TestCheckOutput("manager_cu_num_filter_is_useful", "true"),
					resource.TestCheckOutput("tm_cu_num_filter_is_useful", "true"),
					resource.TestCheckOutput("tm_slot_num_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDliFlinkSQLJobs_basic(name string) string {
	return fmt.Sprintf(`
%s
data "huaweicloud_dli_flinksql_jobs" "test" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]
}
data "huaweicloud_dli_flinksql_jobs" "job_id_filter" {
  job_id = local.job_id
}
  
locals {
  job_id = data.huaweicloud_dli_flinksql_jobs.test.jobs[0].id
}
  
output "job_id_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.job_id_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.job_id_filter.jobs[*].id : v == local.job_id]
  )
}
data "huaweicloud_dli_flinksql_jobs" "queue_name_filter" {
  queue_name = local.queue_name
}
  
locals {
  queue_name = data.huaweicloud_dli_flinksql_jobs.test.jobs[0].queue_name
}
  
output "queue_name_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.queue_name_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.queue_name_filter.jobs[*].queue_name : v == local.queue_name]
  )
}
data "huaweicloud_dli_flinksql_jobs" "cu_num_filter" {
  cu_num = local.cu_num
}
  
locals {
  cu_num = data.huaweicloud_dli_flinksql_jobs.test.jobs[0].cu_num
}
  
output "cu_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.cu_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.cu_num_filter.jobs[*].cu_num : v == local.cu_num]
  )
}
data "huaweicloud_dli_flinksql_jobs" "parallel_num_filter" {
  parallel_num = local.parallel_num
}
  
locals {
  parallel_num = data.huaweicloud_dli_flinksql_jobs.test.jobs[0].parallel_num
}
  
output "parallel_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.parallel_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.parallel_num_filter.jobs[*].parallel_num : v == local.parallel_num]
  )
}
data "huaweicloud_dli_flinksql_jobs" "manager_cu_num_filter" {
  manager_cu_num = local.manager_cu_num
}
  
locals {
  manager_cu_num = data.huaweicloud_dli_flinksql_jobs.test.jobs[0].manager_cu_num
}
  
output "manager_cu_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.manager_cu_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.manager_cu_num_filter.jobs[*].manager_cu_num : v == local.manager_cu_num]
  )
}
data "huaweicloud_dli_flinksql_jobs" "tm_cu_num_filter" {
  tm_cu_num = local.tm_cu_num
}
  
locals {
  tm_cu_num = data.huaweicloud_dli_flinksql_jobs.test.jobs[0].tm_cu_num
}
  
output "tm_cu_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.tm_cu_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.tm_cu_num_filter.jobs[*].tm_cu_num: v == local.tm_cu_num]
  )
}
data "huaweicloud_dli_flinksql_jobs" "tm_slot_num_filter" {
  tm_slot_num = local.tm_slot_num
}
  
locals {
  tm_slot_num = data.huaweicloud_dli_flinksql_jobs.test.jobs[0].tm_slot_num
}
  
output "tm_slot_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.tm_slot_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.tm_slot_num_filter.jobs[*].tm_slot_num : v == local.tm_slot_num]
  )
}
`, testAccFlinkJobResource_basic(name, acceptance.HW_REGION_NAME))
}
