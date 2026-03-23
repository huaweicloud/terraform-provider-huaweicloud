package mrs

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataClusterNodes_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_mapreduce_cluster_nodes.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byNodeGroup = "data.huaweicloud_mapreduce_cluster_nodes.filter_by_node_group"
		dcByGroup   = acceptance.InitDataSourceCheck(byNodeGroup)

		byNodeName = "data.huaweicloud_mapreduce_cluster_nodes.filter_by_node_name"
		dcByName   = acceptance.InitDataSourceCheck(byNodeName)

		byQueryNodeDetail   = "data.huaweicloud_mapreduce_cluster_nodes.filter_by_query_node_detail"
		dcByQueryNodeDetail = acceptance.InitDataSourceCheck(byQueryNodeDetail)

		byQueryEcsDetail   = "data.huaweicloud_mapreduce_cluster_nodes.filter_by_query_ecs_detail"
		dcByQueryEcsDetail = acceptance.InitDataSourceCheck(byQueryEcsDetail)

		byInternalIp   = "data.huaweicloud_mapreduce_cluster_nodes.filter_by_internal_ip"
		dcByInternalIp = acceptance.InitDataSourceCheck(byInternalIp)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMrsClusterID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceClusterNodes_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "nodes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "nodes.0.node_name"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.resource_id"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.node_group_name"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.node_type"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.charging_mode"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.deployment_type"),
					resource.TestCheckResourceAttr(all, "nodes.0.server_info.#", "1"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.server_info.0.server_id"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.server_info.0.server_name"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.server_info.0.server_type"),
					resource.TestMatchResourceAttr(all, "nodes.0.server_info.0.data_volumes.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "nodes.0.server_info.0.data_volumes.0.type"),
					resource.TestMatchResourceAttr(all, "nodes.0.server_info.0.data_volumes.0.size", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "nodes.0.server_info.0.data_volumes.0.count", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttr(all, "nodes.0.server_info.0.root_volume.#", "1"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.server_info.0.root_volume.0.type"),
					resource.TestMatchResourceAttr(all, "nodes.0.server_info.0.root_volume.0.size", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(all, "nodes.0.server_info.0.root_volume.0.count", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "nodes.0.server_info.0.cpu_type"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.server_info.0.internal_ip"),
					resource.TestCheckResourceAttrSet(all, "nodes.0.node_status"),
					// Filter by 'node_group' parameter.
					dcByGroup.CheckResourceExists(),
					resource.TestCheckOutput("is_node_group_filter_useful", "true"),
					// Filter by 'node_name' parameter.
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_node_name_filter_useful", "true"),
					// Filter by 'query_node_detail' parameter.
					dcByQueryNodeDetail.CheckResourceExists(),
					resource.TestCheckResourceAttr(byQueryNodeDetail, "nodes.0.node_detail.#", "1"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.running_status"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.cpu_usage"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.memory_usage"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.disk_usage"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.total_memory"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.available_memory"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.total_hard_disk_space"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.available_hard_disk_space"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.network_read"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.node_detail.0.network_write"),
					resource.TestMatchResourceAttr(byQueryNodeDetail, "nodes.0.component_infos.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.id"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.name"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.instance_group_name"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.running_status"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.ha_status"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.config_status"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.role_name"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.role_short_name"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.role_type"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.service_name"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.support_decom"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.support_reinstall"),
					resource.TestCheckResourceAttrSet(byQueryNodeDetail, "nodes.0.component_infos.0.support_collect_stack_info"),
					// some components not have 'pair_name' and 'relation_pairs', so we don't check them.
					// Filter by 'query_ecs_detail' parameter.
					dcByQueryEcsDetail.CheckResourceExists(),
					resource.TestMatchResourceAttr(byQueryEcsDetail, "nodes.0.server_info.0.cpu", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(byQueryEcsDetail, "nodes.0.server_info.0.mem", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'internal_ip' parameter.
					dcByInternalIp.CheckResourceExists(),
					resource.TestCheckOutput("is_internal_ip_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceClusterNodes_basic() string {
	return fmt.Sprintf(`
# Without any filter parameters.
data "huaweicloud_mapreduce_cluster_nodes" "test" {
  cluster_id = "%[1]s"
}

# Filter by 'node_group' parameter.
# For a cluster, the master node group name is always 'master_node_default_group'
locals {
  node_group_name = "master_node_default_group"
}

data "huaweicloud_mapreduce_cluster_nodes" "filter_by_node_group" {
  cluster_id = "%[1]s"
  node_group = local.node_group_name
}

locals {
  node_group_filter_result = [for v in data.huaweicloud_mapreduce_cluster_nodes.filter_by_node_group.nodes[*].node_group_name :
  v == local.node_group_name]
}

output "is_node_group_filter_useful" {
  value = length(local.node_group_filter_result) > 0 && alltrue(local.node_group_filter_result)
}

# Filter by 'node_name' parameter.
# Fuzzy matching. The node name must contain the 'master' string.
locals {
  node_name = "master"
}

data "huaweicloud_mapreduce_cluster_nodes" "filter_by_node_name" {
  cluster_id = "%[1]s"
  node_name  = local.node_name
}

locals {
  node_name_filter_result = [for v in data.huaweicloud_mapreduce_cluster_nodes.filter_by_node_name.nodes[*].node_name :
  strcontains(v, local.node_name)]
}

output "is_node_name_filter_useful" {
  value = length(local.node_name_filter_result) > 0 && alltrue(local.node_name_filter_result)
}

# Filter by 'query_node_detail' parameters.
data "huaweicloud_mapreduce_cluster_nodes" "filter_by_query_node_detail" {
  cluster_id        = "%[1]s"
  query_node_detail = true
}

# Filter by 'query_ecs_detail' parameters.
data "huaweicloud_mapreduce_cluster_nodes" "filter_by_query_ecs_detail" {
  cluster_id       = "%[1]s"
  query_ecs_detail = true
}

# Filter by 'internal_ip' parameters.
locals {
  internal_ip = try(data.huaweicloud_mapreduce_cluster_nodes.test.nodes[0].server_info[0].internal_ip, null)
}

data "huaweicloud_mapreduce_cluster_nodes" "filter_by_internal_ip" {
  cluster_id  = "%[1]s"
  internal_ip = local.internal_ip
}

locals {
  internal_ip_filter_result = [
    for v in flatten(data.huaweicloud_mapreduce_cluster_nodes.filter_by_internal_ip.nodes[*].server_info[*].internal_ip) :
    v == local.internal_ip
  ]
}

output "is_internal_ip_filter_useful" {
  value = length(local.internal_ip_filter_result) > 0 && alltrue(local.internal_ip_filter_result)
}
`, acceptance.HW_MRS_CLUSTER_ID)
}
