package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceGaussdbIdleSessionsQueryAndKill_basic(t *testing.T) {
	rName := "huaweicloud_gaussdb_idle_sessions_query_and_kill.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGaussdbIdleSessionsQueryAndKill_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "success", "true"),
				),
			},
		},
	})
}

func testAccResourceGaussdbIdleSessionsQueryAndKill_basic() string {
	return fmt.Sprintf(`

data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_gaussdb_idle_sessions_query_and_kill" "test" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
