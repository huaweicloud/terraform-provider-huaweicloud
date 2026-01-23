package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNodeSessionKill_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSNodeId(t)
			acceptance.TestAccPreCheckDDSSessionId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testNodeSessionKill_basic(),
			},
		},
	})
}

func testNodeSessionKill_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dds_node_session_kill" "test" {
  node_id  = "%[1]s"
  sessions = ["%[2]s"]
}
`, acceptance.HW_DDS_NODE_ID, acceptance.HW_DDS_SESSION_ID)
}
