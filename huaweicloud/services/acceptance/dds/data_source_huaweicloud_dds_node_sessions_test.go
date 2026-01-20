package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceNodeSessions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dds_node_sessions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceNodeSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.active"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.operation"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.cost_time"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.host"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.client"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.db"),

					resource.TestCheckOutput("type_filter_useful", "true"),
					resource.TestCheckOutput("namespace_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceNodeSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_node_sessions" "test" {
  node_id = "%[1]s"
}

locals {
  type      = data.huaweicloud_dds_node_sessions.test.sessions[0].type
  namespace = data.huaweicloud_dds_node_sessions.test.sessions[0].namespace
}

data "huaweicloud_dds_node_sessions" "type_filter" {
  node_id = "%[1]s"
  type    = local.type
}

output "type_filter_useful" {
  value = length(data.huaweicloud_dds_node_sessions.type_filter.sessions) > 0 && alltrue(
    [for v in data.huaweicloud_dds_node_sessions.type_filter.sessions[*].type : v == local.type]
  )
}

data "huaweicloud_dds_node_sessions" "namespace_filter" {
  node_id   = "%[1]s"	
  namespace = local.namespace
}

output "namespace_filter_useful" {
  value = length(data.huaweicloud_dds_node_sessions.namespace_filter.sessions) > 0 && alltrue(
    [for v in data.huaweicloud_dds_node_sessions.namespace_filter.sessions[*].namespace : v == local.namespace]
  )
}
`, acceptance.HW_DDS_NODE_ID)
}
