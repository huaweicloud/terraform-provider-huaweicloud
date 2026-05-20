package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNodeSessions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_node_sessions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGeminiDBNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNodeSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.cmd"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.age"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.idle"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.db"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.addr"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.fd"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.multi"),

					resource.TestCheckOutput("addr_prefix_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceNodeSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_node_sessions" "test" {
  node_id = "%[1]s"
}

locals {
  addr_prefix = data.huaweicloud_geminidb_node_sessions.test.sessions[0].addr
}

data "huaweicloud_geminidb_node_sessions" "addr_prefix_filter" {
  node_id     = "%[1]s"
  addr_prefix = local.addr_prefix
}

output "addr_prefix_filter_useful" {
  value = length(data.huaweicloud_geminidb_node_sessions.addr_prefix_filter.sessions) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_node_sessions.addr_prefix_filter.sessions[*].addr : v == local.addr_prefix]
  )
}
`, acceptance.HW_GEMINIDB_NODE_ID)
}
