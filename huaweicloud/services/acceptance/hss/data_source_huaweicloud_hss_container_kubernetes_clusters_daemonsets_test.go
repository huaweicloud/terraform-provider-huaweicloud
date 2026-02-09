package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerKubernetesClustersDaemonsets_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_kubernetes_clusters_daemonsets.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSCCEProtection(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerKubernetesClustersDaemonsets_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "upgradeful_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "err_running_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "err_access_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.latest_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.agent_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.namespace"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.node_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.ds_info.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cluster_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.installed_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.access_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.combined_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.registry_info.#"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceContainerKubernetesClustersDaemonsets_base() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_container_kubernetes_cluster_daemonset" "test" {
  cluster_id   = "%[1]s"
  cluster_name = "%[2]s"
  auto_upgrade = true

  runtime_info {
    runtime_name = "crio_endpoint"
    runtime_path = "user/test"
  }

  schedule_info {
    node_selector = ["test=test"]
  }
}
`, acceptance.HW_CCE_CLUSTER_ID, acceptance.HW_CCE_CLUSTER_NAME)
}

func testDataSourceContainerKubernetesClustersDaemonsets_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_container_kubernetes_clusters_daemonsets" "test" {
  depends_on = [huaweicloud_hss_container_kubernetes_cluster_daemonset.test]

  enterprise_project_id   = "all_granted_eps"
  show_cluster_log_status = true
  access_status           = true
}

# Filter using type.
locals {
  type = data.huaweicloud_hss_container_kubernetes_clusters_daemonsets.test.data_list[0].cluster_type
}

data "huaweicloud_hss_container_kubernetes_clusters_daemonsets" "type_filter" {
  enterprise_project_id = "all_granted_eps"
  type                  = local.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_hss_container_kubernetes_clusters_daemonsets.type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_kubernetes_clusters_daemonsets.type_filter.data_list[*].cluster_type : v == local.type]
  )
}
`, testDataSourceContainerKubernetesClustersDaemonsets_base())
}
