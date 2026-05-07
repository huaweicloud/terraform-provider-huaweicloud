package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsAgentJobHistorySteps_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_agent_job_history_steps.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsAgentJobHistorySteps_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "steps.#"),
					resource.TestCheckResourceAttrSet(dataSource, "steps.0.step_id"),
					resource.TestCheckResourceAttrSet(dataSource, "steps.0.step_name"),
					resource.TestCheckResourceAttrSet(dataSource, "steps.0.run_status"),
					resource.TestCheckResourceAttrSet(dataSource, "steps.0.run_time"),
					resource.TestCheckResourceAttrSet(dataSource, "steps.0.run_duration"),
					resource.TestCheckResourceAttrSet(dataSource, "steps.0.message"),
				),
			},
		},
	})
}

func testDataSourceRdsAgentJobHistorySteps_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_agent_jobs" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_agent_job_histories" "test" {
  instance_id = "%[1]s"
  job_id      = data.huaweicloud_rds_agent_jobs.test.jobs[1].job_id
}

data "huaweicloud_rds_agent_job_history_steps" "test" {
  instance_id = "%[1]s"
  history_id  = data.huaweicloud_rds_agent_job_histories.test.histories[0].history_id
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
