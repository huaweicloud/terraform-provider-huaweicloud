package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerKubernetesClusters_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_kubernetes_clusters.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires preparing a CCE cluster and synchronizing the cluster information to HSS.
			acceptance.TestAccPreCheckHSSCCEProtection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerKubernetesClusters_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "last_update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.cluster_type"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.total_nodes_number"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.active_nodes_number"),
					resource.TestCheckResourceAttrSet(dataSource, "cluster_info_list.0.creation_timestamp"),

					resource.TestCheckOutput("is_cluster_name_filter_useful", "true"),
					resource.TestCheckOutput("is_load_agent_info_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceContainerKubernetesClusters_basic string = `
data "huaweicloud_hss_container_kubernetes_clusters" "test" {}

# Filter using cluster_name.
locals {
  cluster_name = data.huaweicloud_hss_container_kubernetes_clusters.test.cluster_info_list[0].cluster_name
}

data "huaweicloud_hss_container_kubernetes_clusters" "cluster_name_filter" {
  cluster_name = local.cluster_name
}

output "is_cluster_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_clusters.cluster_name_filter.cluster_info_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes_clusters.cluster_name_filter.cluster_info_list[*].cluster_name : v == local.cluster_name]
  )
}

# Filter using load_agent_info.
data "huaweicloud_hss_container_kubernetes_clusters" "load_agent_info_filter" {
  load_agent_info = true
}

output "is_load_agent_info_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_clusters.load_agent_info_filter.cluster_info_list) > 0
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_container_kubernetes_clusters" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_clusters.enterprise_project_id_filter.cluster_info_list) > 0
}

# Filter using non existent cluster_name.
data "huaweicloud_hss_container_kubernetes_clusters" "not_found" {
  cluster_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_container_kubernetes_clusters.not_found.cluster_info_list) == 0
}
`
