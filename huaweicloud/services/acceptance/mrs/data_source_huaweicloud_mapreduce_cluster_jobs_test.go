package mrs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test, please ensure that Flink service is installed on the provided cluster.
func TestAccDataClusterJobs_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_mapreduce_cluster_jobs.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byJobId   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_id"
		dcByJobId = acceptance.InitDataSourceCheck(byJobId)

		byJobName   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_name"
		dcByJobName = acceptance.InitDataSourceCheck(byJobName)

		byUser   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_user"
		dcByUser = acceptance.InitDataSourceCheck(byUser)

		byJobType   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_type"
		dcByJobType = acceptance.InitDataSourceCheck(byJobType)

		byJobState   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_state"
		dcByJobState = acceptance.InitDataSourceCheck(byJobState)

		byJobResult   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_result"
		dcByJobResult = acceptance.InitDataSourceCheck(byJobResult)

		byQueue   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_queue"
		dcByQueue = acceptance.InitDataSourceCheck(byQueue)

		bySubmittedTimeBegin   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_submitted_time_begin"
		dcBySubmittedTimeBegin = acceptance.InitDataSourceCheck(bySubmittedTimeBegin)

		bySubmittedTimeEnd   = "data.huaweicloud_mapreduce_cluster_jobs.filter_by_submitted_time_end"
		dcBySubmittedTimeEnd = acceptance.InitDataSourceCheck(bySubmittedTimeEnd)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
			acceptance.TestAccPreCheckMrsClusterJobProgramPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataClusterJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "jobs.#", regexp.MustCompile(`^[0-9]+$`)),
					dcByJobName.CheckResourceExists(),
					resource.TestCheckOutput("is_job_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.job_progress"),
					resource.TestMatchResourceAttr(byJobId, "jobs.0.started_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byJobId, "jobs.0.submitted_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byJobId, "jobs.0.finished_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.elapsed_time"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.arguments"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.launcher_id"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.properties"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.app_id"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.tracking_url"),
					resource.TestCheckResourceAttrSet(byJobId, "jobs.0.queue"),
					resource.TestCheckOutput("is_job_name_filter_useful", "true"),
					dcByJobId.CheckResourceExists(),
					dcByUser.CheckResourceExists(),
					resource.TestCheckOutput("is_user_filter_useful", "true"),
					dcByJobType.CheckResourceExists(),
					resource.TestCheckOutput("is_job_type_filter_useful", "true"),
					dcByJobState.CheckResourceExists(),
					resource.TestCheckOutput("is_job_state_filter_useful", "true"),
					dcByJobResult.CheckResourceExists(),
					resource.TestCheckOutput("is_job_result_filter_useful", "true"),
					dcByQueue.CheckResourceExists(),
					resource.TestCheckOutput("is_queue_filter_useful", "true"),
					dcBySubmittedTimeBegin.CheckResourceExists(),
					resource.TestCheckOutput("is_submitted_time_begin_filter_useful", "true"),
					dcBySubmittedTimeEnd.CheckResourceExists(),
					resource.TestCheckOutput("is_submitted_time_end_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataClusterJobs_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_mapreduce_job" "test" {
  cluster_id   = "%[1]s"
  name         = "%[2]s"
  type         = "Flink"
  program_path = "%[3]s"

  service_parameters = {
    "mrs.cluster.is.user-agency" = "false"
  }
}
`, acceptance.HW_MRS_CLUSTER_ID, name, acceptance.HW_MRS_CLUSTER_JOB_PROGRAM_PATH)
}

func testAccDataClusterJobs_basic() string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_mapreduce_cluster_jobs" "test" {
  cluster_id = "%[2]s"
  depends_on = [huaweicloud_mapreduce_job.test]
}

# Filter by 'job_id' parameters.
locals {
  job_id = huaweicloud_mapreduce_job.test.id
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_job_id" {
  cluster_id = "%[2]s"
  job_id     = local.job_id
}

locals {
  job_id_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_id.jobs[*].job_id :
  v == local.job_id]
}

output "is_job_id_filter_useful" {
  value = length(local.job_id_filter_result) > 0 && alltrue(local.job_id_filter_result)
}

# Filter by 'job_name' parameters.
locals {
  job_name = huaweicloud_mapreduce_job.test.name
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_job_name" {
  cluster_id = "%[2]s"
  job_name   = local.job_name

  depends_on = [huaweicloud_mapreduce_job.test]
}

locals {
  job_name_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_name.jobs[*].job_name :
  strcontains(v, local.job_name)]
}

output "is_job_name_filter_useful" {
  value = length(local.job_name_filter_result) > 0 && alltrue(local.job_name_filter_result)
}

# Filter by 'user' parameters.
locals {
  user = try(data.huaweicloud_mapreduce_cluster_jobs.test.jobs[0].user, null)
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_user" {
  cluster_id = "%[2]s"
  user       = local.user
}

locals {
  user_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_user.jobs[*].user :
  v == local.user]
}

output "is_user_filter_useful" {
  value = length(local.user_filter_result) > 0 && alltrue(local.user_filter_result)
}

# Filter by 'job_type' parameters.
locals {
  job_type = huaweicloud_mapreduce_job.test.type
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_job_type" {
  cluster_id = "%[2]s"
  job_type   = local.job_type

  depends_on = [huaweicloud_mapreduce_job.test]
}

locals {
  job_type_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_type.jobs[*].job_type :
  v == local.job_type]
}

output "is_job_type_filter_useful" {
  value = length(local.job_type_filter_result) > 0 && alltrue(local.job_type_filter_result)
}

# Filter by 'job_state' parameters.
locals {
  job_state = huaweicloud_mapreduce_job.test.status
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_job_state" {
  cluster_id = "%[2]s"
  job_state  = local.job_state
}

locals {
  job_state_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_state.jobs[*].job_state :
  v == local.job_state]
}

output "is_job_state_filter_useful" {
  value = length(local.job_state_filter_result) > 0 && alltrue(local.job_state_filter_result)
}

# Filter by 'job_result' parameters.
locals {
  job_result = try(data.huaweicloud_mapreduce_cluster_jobs.test.jobs[0].job_result, null)
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_job_result" {
  cluster_id = "%[2]s"
  job_result = local.job_result
}

locals {
  job_result_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_job_result.jobs[*].job_result :
  v == local.job_result]
}

output "is_job_result_filter_useful" {
  value = length(local.job_result_filter_result) > 0 && alltrue(local.job_result_filter_result)
}

# Filter by 'queue' parameters.
locals {
  queue = try(data.huaweicloud_mapreduce_cluster_jobs.test.jobs[0].queue, null)
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_queue" {
  cluster_id = "%[2]s"
  queue      = local.queue
}

locals {
  queue_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_queue.jobs[*].queue :
  v == local.queue]
}

output "is_queue_filter_useful" {
  value = length(local.queue_filter_result) > 0 && alltrue(local.queue_filter_result)
}

# Filter by 'submitted_time_begin' parameters.
locals {
  submitted_time = try(data.huaweicloud_mapreduce_cluster_jobs.test.jobs[0].submitted_time, null)
}

data "huaweicloud_mapreduce_cluster_jobs" "filter_by_submitted_time_begin" {
  cluster_id           = "%[2]s"
  submitted_time_begin = local.submitted_time
}

locals {
  submitted_time_begin_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_submitted_time_begin.jobs[*].submitted_time :
  timecmp(v, local.submitted_time) >= 0]
}

output "is_submitted_time_begin_filter_useful" {
  value = length(local.submitted_time_begin_filter_result) > 0 && alltrue(local.submitted_time_begin_filter_result)
}

# Filter by 'submitted_time_end' parameters.
data "huaweicloud_mapreduce_cluster_jobs" "filter_by_submitted_time_end" {
  cluster_id         = "%[2]s"
  submitted_time_end = local.submitted_time
}

locals {
  submitted_time_end_filter_result = [for v in data.huaweicloud_mapreduce_cluster_jobs.filter_by_submitted_time_end.jobs[*].submitted_time :
  timecmp(local.submitted_time, v) >= 0]
}

output "is_submitted_time_end_filter_useful" {
  value = length(local.submitted_time_end_filter_result) > 0 && alltrue(local.submitted_time_end_filter_result)
}
`, testAccDataClusterJobs_base(), acceptance.HW_MRS_CLUSTER_ID)
}
