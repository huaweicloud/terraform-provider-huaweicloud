package aom

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUniAgentBatchUpgrade_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckUniAgentId(t)
			acceptance.TestAccPreCheckUniAgentInnerIp(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccUniAgentBatchUpgrade_basic(),
			},
		},
	})
}

func testAccUniAgentBatchUpgrade_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_aom_uniagent_batch_upgrade" "test" {
  version = "%[1]s"

  agent_list {
    agent_id = "%[2]s"
    inner_ip = "%[3]s"
  }
}
`, acceptance.HW_AOM_UNIAGENT_VERSION, acceptance.HW_AOM_UNIAGENT_AGENT_ID, acceptance.HW_AOM_UNIAGENT_INNER_IP)
}
