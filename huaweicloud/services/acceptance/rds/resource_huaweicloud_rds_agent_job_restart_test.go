package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsAgentJobRestart_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsAgentJobRestart_basic(),
			},
		},
	})
}

func testAccRdsAgentJobRestart_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_agent_jobs" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_agent_job_restart" "test" {
  instance_id = "%[1]s"
  job_id      = data.huaweicloud_rds_agent_jobs.test.jobs[0].job_id
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
