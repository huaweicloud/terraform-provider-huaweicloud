package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsAgentJobHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_agent_job_histories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsAgentJobHistories_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "histories.#"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.history_id"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.run_status"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.run_time"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.run_duration"),
					resource.TestCheckResourceAttrSet(dataSource, "histories.0.message"),
					resource.TestCheckOutput("run_status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsAgentJobHistories_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_agent_jobs" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_agent_job_histories" "test" {
  instance_id = "%[1]s"
  job_id      = data.huaweicloud_rds_agent_jobs.test.jobs[1].job_id
}

data "huaweicloud_rds_agent_job_histories" "run_status_filter" {
  instance_id = "%[1]s"
  job_id      = data.huaweicloud_rds_agent_jobs.test.jobs[1].job_id
  run_status  = data.huaweicloud_rds_agent_job_histories.test.histories[0].run_status
}
locals {
  run_status = data.huaweicloud_rds_agent_job_histories.test.histories[0].run_status
}
output "run_status_filter_is_useful" {
  value = length(data.huaweicloud_rds_agent_job_histories.run_status_filter.histories) > 0 && alltrue(
  [for v in data.huaweicloud_rds_agent_job_histories.run_status_filter.histories[*].run_status : v == local.run_status]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
