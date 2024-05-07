package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliSqlJobs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dli_sql_jobs.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliSqlJobs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.queue_name"),

					resource.TestCheckOutput("job_id_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
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

data "huaweicloud_dli_sql_jobs" "job_id_filter" {
  job_id = local.job_id
}
  
locals {
  job_id = data.huaweicloud_dli_sql_jobs.test.jobs[0].id
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
  type = data.huaweicloud_dli_sql_jobs.test.jobs[0].type
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
  status = data.huaweicloud_dli_sql_jobs.test.jobs[0].status
}
  
output "status_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_jobs.status_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_jobs.status_filter.jobs[*].status : v == local.status]
  )
}

data "huaweicloud_dli_sql_jobs" "queue_name_filter" {
  queue_name = local.queue_name
}
  
locals {
  queue_name = data.huaweicloud_dli_sql_jobs.test.jobs[0].queue_name
}
  
output "queue_name_filter_is_useful" {
  value = length(data.huaweicloud_dli_sql_jobs.queue_name_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_dli_sql_jobs.queue_name_filter.jobs[*].queue_name : v == local.queue_name]
  )
}
`, testAccSqlJobBaseResource_basic(name))
}
