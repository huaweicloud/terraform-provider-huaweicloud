package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsAgentJobModify_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsAgentJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsAgentJobModify_basic(),
			},
		},
	})
}

func testAccRdsAgentJobModify_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_agent_job_modify" "test" {
  instance_id = "%[1]s"
  job_id      = "%[2]s"
  profile_id  = "4"
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_AGENT_JOB_ID)
}
