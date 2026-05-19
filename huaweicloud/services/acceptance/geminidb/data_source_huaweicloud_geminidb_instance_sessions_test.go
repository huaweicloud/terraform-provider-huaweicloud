package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInstanceSessions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_instance_sessions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGeminiDBNodeId(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInstanceSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.cmd"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.age"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.idle"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.db"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.addr"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.fd"),
					resource.TestCheckResourceAttrSet(dataSource, "node_sessions.0.sessions.0.multi"),

					resource.TestCheckOutput("node_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInstanceSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_instance_sessions" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}

locals {
  node_id = data.huaweicloud_geminidb_instance_sessions.test.node_sessions[0].node_id
}

data "huaweicloud_geminidb_instance_sessions" "node_id_filter" {
  instance_id = "%[1]s"
  node_id     = local.node_id
}

output "node_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_instance_sessions.node_id_filter.node_sessions.0.sessions) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_instance_sessions.node_id_filter.node_sessions[*].node_id : v == local.node_id]
  )
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, acceptance.HW_GEMINIDB_NODE_ID)
}
