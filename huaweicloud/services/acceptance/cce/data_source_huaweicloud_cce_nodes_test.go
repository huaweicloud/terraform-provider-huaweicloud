package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccNodesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_cce_nodes.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccNodesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "ids.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "nodes.0.name", rName),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("nodeId_filter_is_useful", "true"),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccNodesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_nodes" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id

  depends_on = [huaweicloud_cce_node.test]
}

data "huaweicloud_cce_nodes" "name_filter" {
  cluster_id = huaweicloud_cce_cluster.test.id
  name       = huaweicloud_cce_node.test.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_cce_nodes.name_filter.nodes) > 0 && alltrue(
    [for v in data.huaweicloud_cce_nodes.name_filter.nodes : v.name == huaweicloud_cce_node.test.name]
  )
}

data "huaweicloud_cce_nodes" "nodeId_filter" {
  cluster_id = huaweicloud_cce_cluster.test.id
  node_id    = huaweicloud_cce_node.test.id
}

output "nodeId_filter_is_useful" {
  value = length(data.huaweicloud_cce_nodes.nodeId_filter.nodes) > 0 && alltrue(
    [for v in data.huaweicloud_cce_nodes.nodeId_filter.nodes : v.id == huaweicloud_cce_node.test.id]
  )
}

data "huaweicloud_cce_nodes" "status_filter" {
  cluster_id = huaweicloud_cce_cluster.test.id
  status     = huaweicloud_cce_node.test.status
}

output "status_filter_is_useful" {
  value = length(data.huaweicloud_cce_nodes.status_filter.nodes) > 0 && alltrue(
    [for v in data.huaweicloud_cce_nodes.status_filter.nodes : v.status == huaweicloud_cce_node.test.status]
  )
}
`, testAccCceCluster_config(rName))
}
