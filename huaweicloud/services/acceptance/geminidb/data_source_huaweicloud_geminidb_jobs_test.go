package geminidb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceJobs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_jobs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceJobs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.instance.0.name"),

					resource.TestCheckOutput("job_id_filter_useful", "true"),
					resource.TestCheckOutput("name_filter_useful", "true"),
					resource.TestCheckOutput("status_filter_useful", "true"),
					resource.TestCheckOutput("start_time_filter_useful", "true"),
					resource.TestCheckOutput("end_time_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceJobs_basic = `
data "huaweicloud_geminidb_jobs" "test" {}

locals {
  job_id     = data.huaweicloud_geminidb_jobs.test.jobs[0].id
  name       = data.huaweicloud_geminidb_jobs.test.jobs[0].name
  status     = data.huaweicloud_geminidb_jobs.test.jobs[0].status
  start_time = data.huaweicloud_geminidb_jobs.test.jobs[0].start_time
  end_time   = data.huaweicloud_geminidb_jobs.test.jobs[0].end_time
}

data "huaweicloud_geminidb_jobs" "job_id_filter" {
  job_id = local.job_id
}

output "job_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_jobs.job_id_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_jobs.job_id_filter.jobs[*].id : v == local.job_id]
  )
}

data "huaweicloud_geminidb_jobs" "name_filter" {	
  name = local.name
}

output "name_filter_useful" {
  value = length(data.huaweicloud_geminidb_jobs.name_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_jobs.name_filter.jobs[*].name : v == local.name]
  )
}

data "huaweicloud_geminidb_jobs" "status_filter" {
  status = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_geminidb_jobs.status_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_jobs.status_filter.jobs[*].status : v == local.status]
  )
}

data "huaweicloud_geminidb_jobs" "start_time_filter" {
  start_time = local.start_time
}

output "start_time_filter_useful" {
  value = length(data.huaweicloud_geminidb_jobs.start_time_filter.jobs) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_jobs.start_time_filter.jobs[*].start_time : v == local.start_time]
  )
}

data "huaweicloud_geminidb_jobs" "end_time_filter" {
  end_time = local.end_time
}

output "end_time_filter_useful" {
  value = length(data.huaweicloud_geminidb_jobs.end_time_filter.jobs) > 0
}
`
