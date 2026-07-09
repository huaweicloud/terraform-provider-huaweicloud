package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceGaussdbSpecifiedSessionsQueryAndKill_basic(t *testing.T) {
	rName := "huaweicloud_gaussdb_specified_sessions_query_and_kill.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckGaussDBSessionId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testResourceGaussdbSpecifiedSessionsQueryAndKill_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "success", "true"),
					resource.TestCheckResourceAttr(rName, "session_ids.#", "1"),
				),
			},
		},
	})
}

func testResourceGaussdbSpecifiedSessionsQueryAndKill_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_gaussdb_specified_sessions_query_and_kill" "test" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
  session_ids  = ["%[2]s"]
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_SESSION_ID)
}
