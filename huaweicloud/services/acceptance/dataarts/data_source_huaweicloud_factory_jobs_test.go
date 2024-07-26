package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceFactoryJobs_basic(t *testing.T) {
	var (
		rName      = acceptance.RandomAccResourceName()
		dataSource = "data.huaweicloud_dataarts_factory_jobs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byName   = "data.huaweicloud_dataarts_factory_jobs.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byBatchType   = "data.huaweicloud_dataarts_factory_jobs.filter_by_batch_type"
		dcByBatchType = acceptance.InitDataSourceCheck(byBatchType)

		byRealTimeType   = "data.huaweicloud_dataarts_factory_jobs.filter_by_real_time_type"
		dcByRealTimeType = acceptance.InitDataSourceCheck(byRealTimeType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDatasourceFactoryJobs_workspace_id_not_found,
				ExpectError: regexp.MustCompile("detail msg Workspace does not exists"),
			},
			{
				Config: testAccDatasourceFactoryJobs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.process_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.priority"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.directory"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.created_by"),
					resource.TestMatchResourceAttr(dataSource, "jobs.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.updated_by"),
					resource.TestMatchResourceAttr(dataSource, "jobs.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					dcByBatchType.CheckResourceExists(),
					dcByRealTimeType.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_fuzzy_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_name", "true"),
					resource.TestCheckOutput("is_batch_type_filter_useful", "true"),
					resource.TestCheckOutput("is_real_time_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDatasourceFactoryJobs_workspace_id_not_found = `
data "huaweicloud_dataarts_factory_jobs" "test" {
  workspace_id = "not_found"
}
`

func testAccDatasourceFactoryJobs_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_smn_topic" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dataarts_factory_job" "batch_job" {
  name         = "%[2]s_batch_job"
  workspace_id = "%[3]s"
  process_type = "BATCH"

  nodes {
    name = "SMN_%[2]s_batch_job"
    type = "SMN"

    location {
      x = 10
      y = 11
    }

    properties {
      name  = "topic"
      value = huaweicloud_smn_topic.test.topic_urn
    }

    properties {
      name  = "messageType"
      value = "NORMAL"
    }

    properties {
      name  = "message"
      value = "terraform acceptance test"
    }
  }

  schedule {
    type = "EXECUTE_ONCE"
  }
}
`, testFactoryJob_basic(name), name, acceptance.HW_DATAARTS_WORKSPACE_ID)
}

func testAccDatasourceFactoryJobs_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dataarts_factory_jobs" "test" {
  depends_on = [
    huaweicloud_dataarts_factory_job.batch_job
  ]

  workspace_id = "%[2]s"
}

locals {
  job_name       = huaweicloud_dataarts_factory_job.batch_job.name
  batch_type     = huaweicloud_dataarts_factory_job.batch_job.process_type
  real_time_type = huaweicloud_dataarts_factory_job.test.process_type
}

# Filter by name (Exact match)
data "huaweicloud_dataarts_factory_jobs" "filter_by_name" {
  depends_on = [
    huaweicloud_dataarts_factory_job.batch_job
  ]

  workspace_id = "%[2]s"
  name         = local.job_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_factory_jobs.filter_by_name.jobs[*].name : v == local.job_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by name (Fuzzy search)
data "huaweicloud_dataarts_factory_jobs" "filter_by_fuzzy_name" {
  depends_on = [
    huaweicloud_dataarts_factory_job.batch_job
  ]

  workspace_id = "%[2]s"
  name         = "tf_test"
}

output "is_fuzzy_name_filter_useful" {
  value = length(data.huaweicloud_dataarts_factory_jobs.filter_by_fuzzy_name.jobs) >= 1
}

# Filter by name (Not found)
data "huaweicloud_dataarts_factory_jobs" "not_found_name" {
  depends_on = [
    huaweicloud_dataarts_factory_job.batch_job
  ]

  workspace_id = "%[2]s"
  name         = "not_found_name"
}

output "not_found_name" {
  value = length(data.huaweicloud_dataarts_factory_jobs.not_found_name.jobs) == 0
}

# Filter by "BATCH" type
data "huaweicloud_dataarts_factory_jobs" "filter_by_batch_type" {
  depends_on = [
    huaweicloud_dataarts_factory_job.batch_job
  ]

  workspace_id = "%[2]s"
  process_type = local.batch_type
}

locals {
  batch_type_filter_result = [
    for v in data.huaweicloud_dataarts_factory_jobs.filter_by_batch_type.jobs[*].process_type : v == local.batch_type
  ]
}

output "is_batch_type_filter_useful" {
  value = length(local.batch_type_filter_result) > 0 && alltrue(local.batch_type_filter_result)
}

# Filter by "REAL_TIME" type
data "huaweicloud_dataarts_factory_jobs" "filter_by_real_time_type" {
  depends_on = [
    huaweicloud_dataarts_factory_job.test
  ]

  workspace_id = "%[2]s"
  process_type = local.real_time_type
}

locals {
  real_time_filter_result = [
    for v in data.huaweicloud_dataarts_factory_jobs.filter_by_real_time_type.jobs[*].process_type : v == local.real_time_type
  ]
}

output "is_real_time_type_filter_useful" {
  value = length(local.real_time_filter_result) > 0 && alltrue(local.real_time_filter_result)
}
`, testAccDatasourceFactoryJobs_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
