package dli

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSparkJobs_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dli_spark_jobs.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byJobId   = "data.huaweicloud_dli_spark_jobs.job_id_filter"
		dcByJobId = acceptance.InitDataSourceCheck(byJobId)

		byState   = "data.huaweicloud_dli_spark_jobs.state_filter"
		dcByState = acceptance.InitDataSourceCheck(byState)

		byOwner   = "data.huaweicloud_dli_spark_jobs.owner_filter"
		dcByOwner = acceptance.InitDataSourceCheck(byOwner)

		byJobName   = "data.huaweicloud_dli_spark_jobs.job_name_filter"
		dcByJobName = acceptance.InitDataSourceCheck(byJobName)

		byQueueName   = "data.huaweicloud_dli_spark_jobs.queue_name_filter"
		dcByQueueName = acceptance.InitDataSourceCheck(byQueueName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDliGeneralQueueName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSparkJobs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "jobs.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "jobs.0.id"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.owner"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.queue"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.cluster_name"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.state"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.duration"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.sc_type"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.req_body"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "jobs.0.updated_at"),

					// Filter by job ID
					dcByJobId.CheckResourceExists(),
					resource.TestCheckOutput("job_id_filter_is_useful", "true"),

					// Filter by state
					dcByState.CheckResourceExists(),
					resource.TestCheckOutput("state_filter_is_useful", "true"),

					// Filter by owner
					dcByOwner.CheckResourceExists(),
					resource.TestCheckOutput("owner_filter_is_useful", "true"),

					// Filter by job name
					dcByJobName.CheckResourceExists(),
					resource.TestCheckOutput("job_name_filter_is_useful", "true"),

					// Filter by queue name
					dcByQueueName.CheckResourceExists(),
					resource.TestCheckOutput("queue_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSparkJobs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dli_spark_jobs" "all" {
  depends_on = [huaweicloud_dli_spark_job.test]
}

# Filter by job ID
locals {
  job_id = huaweicloud_dli_spark_job.test.id
}

data "huaweicloud_dli_spark_jobs" "job_id_filter" {
  job_id = local.job_id
}

locals {
  job_id_filter_result = [
    for v in data.huaweicloud_dli_spark_jobs.job_id_filter.jobs[*].id : v == local.job_id
  ]
}

output "job_id_filter_is_useful" {
  value = length(local.job_id_filter_result) > 0 && alltrue(local.job_id_filter_result)
}

# Filter by state
locals {
  state = data.huaweicloud_dli_spark_jobs.all.jobs[0].state
}

data "huaweicloud_dli_spark_jobs" "state_filter" {
  state = local.state
}

locals {
  state_filter_result = [
    for v in data.huaweicloud_dli_spark_jobs.state_filter.jobs[*].state : v == local.state
  ]
}

output "state_filter_is_useful" {
  value = length(local.state_filter_result) > 0 && alltrue(local.state_filter_result)
}

# Filter by owner
locals {
  owner = data.huaweicloud_dli_spark_jobs.all.jobs[0].owner
}

data "huaweicloud_dli_spark_jobs" "owner_filter" {
  owner = local.owner
}

locals {
  owner_filter_result = [
    for v in data.huaweicloud_dli_spark_jobs.owner_filter.jobs[*].owner : v == local.owner
  ]
}

output "owner_filter_is_useful" {
  value = length(local.owner_filter_result) > 0 && alltrue(local.owner_filter_result)
}

# Filter by job name
locals {
  job_name = data.huaweicloud_dli_spark_jobs.all.jobs[0].name
}

data "huaweicloud_dli_spark_jobs" "job_name_filter" {
  job_name = local.job_name
}

locals {
  job_name_filter_result = [
    for v in data.huaweicloud_dli_spark_jobs.job_name_filter.jobs[*].name : v == local.job_name
  ]
}

output "job_name_filter_is_useful" {
  value = length(local.job_name_filter_result) > 0 && alltrue(local.job_name_filter_result)
}

# Filter by queue name
locals {
  queue_name = data.huaweicloud_dli_spark_jobs.all.jobs[0].queue
}

data "huaweicloud_dli_spark_jobs" "queue_name_filter" {
  queue_name = local.queue_name
}

locals {
  queue_name_filter_result = [
    for v in data.huaweicloud_dli_spark_jobs.queue_name_filter.jobs[*].queue : v == local.queue_name
  ]
}

output "queue_name_filter_is_useful" {
  value = length(local.queue_name_filter_result) > 0 && alltrue(local.queue_name_filter_result)
}
`, testAccDliSparkJob_basic(name, acceptance.RandomAccResourceNameWithDash()))
}
