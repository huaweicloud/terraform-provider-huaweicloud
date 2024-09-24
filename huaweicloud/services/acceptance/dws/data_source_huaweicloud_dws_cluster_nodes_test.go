package dws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceClusterNodes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dws_cluster_nodes.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		byId       = "data.huaweicloud_dws_cluster_nodes.filter_by_id"
		dcById     = acceptance.InitDataSourceCheck(byId)
		byStatus   = "data.huaweicloud_dws_cluster_nodes.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
		byType     = "data.huaweicloud_dws_cluster_nodes.filter_by_resource_type"
		dcByType   = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDwsClusterId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceClusterNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "nodes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.spec"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "nodes.0.sub_status"),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("id_filter_is_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
				),
			},
		},
	})
}

func testDataSourceClusterNodes_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dws_cluster_nodes" "test" {
  cluster_id = "%[1]s"
}

data "huaweicloud_dws_cluster_nodes" "filter_by_id" {
  cluster_id = "%[1]s"
  node_id    = local.node_id
}

locals {
  node_id          = data.huaweicloud_dws_cluster_nodes.test.nodes[0].id
  id_fliter_nodes  = data.huaweicloud_dws_cluster_nodes.filter_by_id.nodes
}

# The query result is a ring containing the node, and a ring consists of at least three nodes.
output "id_filter_is_useful" {
  value = length(local.id_fliter_nodes) >= 3 && contains(local.id_fliter_nodes[*].id, local.node_id)
}

locals {
  status = data.huaweicloud_dws_cluster_nodes.test.nodes[0].status
}

data "huaweicloud_dws_cluster_nodes" "filter_by_status" {
  cluster_id = "%[1]s"
  filter_by  = "status"
  filter     = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_dws_cluster_nodes.filter_by_status.nodes[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

data "huaweicloud_dws_cluster_nodes" "filter_by_resource_type" {
  cluster_id = "%[1]s"
  filter_by  = "instCreateType"
  filter     = "INST"
}

# For a cluster, there are at least three nodes.
output "resource_type_filter_is_useful" {
  value = length(data.huaweicloud_dws_cluster_nodes.filter_by_resource_type.nodes[*]) >= 3
}
`, acceptance.HW_DWS_CLUSTER_ID)
}
