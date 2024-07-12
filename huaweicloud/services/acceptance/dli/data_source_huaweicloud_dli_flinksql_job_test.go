package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliFlinkSQLJobs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dli_flinksql_jobs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		rName      = acceptance.RandomAccResourceName()

		byId   = "data.huaweicloud_dli_flinksql_jobs.job_id_filter"
		dcById = acceptance.InitDataSourceCheck(byId)

		byQueueName   = "data.huaweicloud_dli_flinksql_jobs.queue_name_filter"
		dcByQueuename = acceptance.InitDataSourceCheck(byQueueName)

		byCuNum   = "data.huaweicloud_dli_flinksql_jobs.cu_num_filter"
		dcByCuNum = acceptance.InitDataSourceCheck(byCuNum)

		byParallelNum   = "data.huaweicloud_dli_flinksql_jobs.parallel_num_filter"
		dcByParallelNum = acceptance.InitDataSourceCheck(byParallelNum)

		byManageCuNum   = "data.huaweicloud_dli_flinksql_jobs.manager_cu_num_filter"
		dcByManageCuNum = acceptance.InitDataSourceCheck(byManageCuNum)

		byTmCuNum   = "data.huaweicloud_dli_flinksql_jobs.tm_cu_num_filter"
		dcByTmCuNum = acceptance.InitDataSourceCheck(byTmCuNum)

		byTmSlotNum   = "data.huaweicloud_dli_flinksql_jobs.tm_slot_num_filter"
		dcByTmSlotNum = acceptance.InitDataSourceCheck(byTmSlotNum)

		byTags   = "data.huaweicloud_dli_flinksql_jobs.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliGenaralQueueName(t)
			acceptance.TestAccPreCheckDliFlinkVersion(t)
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
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("job_id_filter_is_useful", "true"),
					dcByQueuename.CheckResourceExists(),
					resource.TestCheckOutput("queue_name_filter_is_useful", "true"),
					dcByCuNum.CheckResourceExists(),
					resource.TestCheckOutput("cu_num_filter_is_useful", "true"),
					dcByParallelNum.CheckResourceExists(),
					resource.TestCheckOutput("parallel_num_filter_is_useful", "true"),
					dcByManageCuNum.CheckResourceExists(),
					resource.TestCheckOutput("manager_cu_num_filter_is_useful", "true"),
					dcByTmCuNum.CheckResourceExists(),
					resource.TestCheckOutput("tm_cu_num_filter_is_useful", "true"),
					dcByTmSlotNum.CheckResourceExists(),
					resource.TestCheckOutput("tm_slot_num_filter_is_useful", "true"),
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("tags_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDliFlinkSQLJobs_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_dli_flinksql_job" "test" {
  name          = "%s"
  type          = "flink_opensource_sql_job"
  run_mode      = "exclusive_cluster"
  sql           = local.opensourceSql
  queue_name    = "%s"
  flink_version = "%s"
  restart_when_exception = true
  cu_number         = 4
  resume_checkpoint = true
  manager_cu_number = 2
  tm_cus            = 2
  tm_slot_num       = 3
  parallel_number   = 2

  tags = {
    foo = "bar"
  }
}

data "huaweicloud_dli_flinksql_jobs" "test" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]
}

data "huaweicloud_dli_flinksql_jobs" "job_id_filter" {
  job_id = local.job_id
}

locals {
  job_id = huaweicloud_dli_flinksql_job.test.id
}

output "job_id_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.job_id_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.job_id_filter.jobs[*].id : v == local.job_id]
  )
}

locals {
  queue_name = huaweicloud_dli_flinksql_job.test.queue_name
}

data "huaweicloud_dli_flinksql_jobs" "queue_name_filter" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]

  queue_name = local.queue_name
}

output "queue_name_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.queue_name_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.queue_name_filter.jobs[*].queue_name : v == local.queue_name]
  )
}

locals {
  cu_num = huaweicloud_dli_flinksql_job.test.cu_number
}


data "huaweicloud_dli_flinksql_jobs" "cu_num_filter" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]

  cu_num = local.cu_num
}

output "cu_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.cu_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.cu_num_filter.jobs[*].cu_num : v == local.cu_num]
  )
}

locals {
  parallel_num = huaweicloud_dli_flinksql_job.test.parallel_number
}

data "huaweicloud_dli_flinksql_jobs" "parallel_num_filter" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]

  parallel_num = local.parallel_num
}

output "parallel_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.parallel_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.parallel_num_filter.jobs[*].parallel_num : v == local.parallel_num]
  )
}

locals {
  manager_cu_num = huaweicloud_dli_flinksql_job.test.manager_cu_number
}

data "huaweicloud_dli_flinksql_jobs" "manager_cu_num_filter" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]

  manager_cu_num = local.manager_cu_num
}

output "manager_cu_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.manager_cu_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.manager_cu_num_filter.jobs[*].manager_cu_num : v == local.manager_cu_num]
  )
}

locals {
  tm_cu_num = huaweicloud_dli_flinksql_job.test.tm_cus
}

data "huaweicloud_dli_flinksql_jobs" "tm_cu_num_filter" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]

  tm_cu_num = local.tm_cu_num
}

output "tm_cu_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.tm_cu_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.tm_cu_num_filter.jobs[*].tm_cu_num : v == local.tm_cu_num]
  )
}

locals {
  tm_slot_num = huaweicloud_dli_flinksql_job.test.tm_slot_num
}

data "huaweicloud_dli_flinksql_jobs" "tm_slot_num_filter" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]

  tm_slot_num = local.tm_slot_num
}

output "tm_slot_num_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.tm_slot_num_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flinksql_jobs.tm_slot_num_filter.jobs[*].tm_slot_num : v == local.tm_slot_num]
  )
}

data "huaweicloud_dli_flinksql_jobs" "filter_by_tags" {
  depends_on = [
    huaweicloud_dli_flinksql_job.test
  ]

  tags = huaweicloud_dli_flinksql_job.test.tags
}

output "tags_filter_is_useful" {
  value = length(data.huaweicloud_dli_flinksql_jobs.filter_by_tags.jobs) > 0
}
`, testAccFlinkJobResource_base(), name, acceptance.HW_DLI_GENERAL_QUEUE_NAME, acceptance.HW_DLI_FLINK_VERSION)
}
