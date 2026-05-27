package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBHtapStarrocksNodeRestart_basic(t *testing.T) {
	resourceName := "huaweicloud_taurusdb_htap_starrocks_node_restart.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBHtapStarrocksNodeRestart_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", acceptance.HW_TAURUSDB_HTAP_NODE_ID),
				),
			},
		},
	})
}

func testAccTaurusDBHtapStarrocksNodeRestart_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_htap_starrocks_node_restart" "test" {
  taurusdb_instance_id  = "%[1]s"
  starrocks_instance_id = "%[2]s"
  starrocks_node_id     = "%[3]s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_NODE_ID)
}
