package dli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliSqlJobs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dli_sql_jobs.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byJobId   = "data.huaweicloud_dli_sql_jobs.job_id_filter"
		dcByJobId = acceptance.InitDataSourceCheck(byJobId)

		byType   = "data.huaweicloud_dli_sql_jobs.type_filter"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byStatus   = "data.huaweicloud_dli_sql_jobs.status_filter"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byQueueName   = "data.huaweicloud_dli_sql_jobs.queue_name_filter"
		dcByQueueName = acceptance.InitDataSourceCheck(byQueueName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliSQLQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliSqlJobs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "jobs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByJobId.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.type"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.queue_name"),
					resource.TestCheckOutput("job_id_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					dcByQueueName.CheckResourceExists(),
					resource.TestCheckOutput("queue_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDliSqlJobs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dli_sql_jobs" "test" {
  depends_on = [
    huaweicloud_dli_sql_job.test
  ]
}

locals {
  job_id = huaweicloud_dli_sql_job.test.id
}

data "huaweicloud_dli_sql_jobs" "job_id_filter" {
  job_id = local.job_id
}

output "job_id_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_jobs.job_id_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_jobs.job_id_filter.jobs[*].id : v == local.job_id]
  )
}

data "huaweicloud_dli_sql_jobs" "type_filter" {
  type = local.type
}

locals {
  type = data.huaweicloud_dli_sql_jobs.job_id_filter.jobs[0].type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_jobs.type_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_jobs.type_filter.jobs[*].type : v == local.type]
  )
}

data "huaweicloud_dli_sql_jobs" "status_filter" {
  status = local.status
}

locals {
  status = huaweicloud_dli_sql_job.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_jobs.status_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_jobs.status_filter.jobs[*].status : v == local.status]
  )
}

locals {
  queue_name = huaweicloud_dli_sql_job.test.queue_name
}

data "huaweicloud_dli_sql_jobs" "queue_name_filter" {
  queue_name = local.queue_name
}

output "queue_name_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_jobs.queue_name_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_jobs.queue_name_filter.jobs[*].queue_name : v == local.queue_name]
  )
}
`, testAccSqlJobBaseResource_basic(name))
}
