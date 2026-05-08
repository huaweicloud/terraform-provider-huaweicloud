package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBNodeSessions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_node_sessions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceTaurusDBNodeSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "processes.#"),
					resource.TestCheckResourceAttrSet(dataSource, "processes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "processes.0.user"),
					resource.TestCheckResourceAttrSet(dataSource, "processes.0.host"),
					resource.TestCheckResourceAttrSet(dataSource, "processes.0.db"),
					resource.TestCheckResourceAttrSet(dataSource, "processes.0.command"),
					resource.TestCheckResourceAttrSet(dataSource, "processes.0.time"),
				),
			},
		},
	})
}

func testDataSourceTaurusDBNodeSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_taurusdb_node_sessions" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}
`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_NODE_ID)
}
