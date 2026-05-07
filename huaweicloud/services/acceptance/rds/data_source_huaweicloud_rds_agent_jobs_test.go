package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsAgentJobs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_agent_jobs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsAgentJobs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.job_name"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.is_enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.run_time"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.run_status"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.failure_count"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.agent_type"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.profile_id"),
					resource.TestCheckResourceAttrSet(dataSource, "jobs.0.profile_name"),
					resource.TestCheckOutput("job_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsAgentJobs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_agent_jobs" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_agent_jobs" "job_type_filter" {
  instance_id = "%[1]s"
  job_type    = "replication"
}
output "job_type_filter_is_useful" {
  value = length(data.huaweicloud_rds_agent_jobs.job_type_filter.jobs) > 0
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
