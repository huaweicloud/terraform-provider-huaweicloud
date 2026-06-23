package taurusdb

import (
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBNodeSessionsKill_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBNodeId(t)
			acceptance.TestAccPreCheckTaurusDBNodeSessionId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBNodeSessionsKill_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("huaweicloud_taurusdb_node_sessions_kill.test", "processes_killed.#"),
					resource.TestCheckOutput("killed_process_id", acceptance.HW_TAURUSDB_NODE_SESSION_ID),
					resource.TestCheckResourceAttrSet("huaweicloud_taurusdb_node_sessions_kill.test", "processes_not_found.#"),
					resource.TestCheckOutput("not_found_process_id", "999999999"),
				),
			},
		},
	})
}

func testAccTaurusDBNodeSessionsKill_basic() string {
	processId, err := strconv.Atoi(acceptance.HW_TAURUSDB_NODE_SESSION_ID)
	if err != nil {
		log.Printf("[ERROR] failed to parse TaurusDB node session ID: %s", err)
	}
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_node_sessions_kill" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  processes   = [%[3]d, 999999999]
}

output killed_process_id {
  value = huaweicloud_taurusdb_node_sessions_kill.test.processes_killed[0]
}

output not_found_process_id {
  value = huaweicloud_taurusdb_node_sessions_kill.test.processes_not_found[0]
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_NODE_ID, processId)
}
