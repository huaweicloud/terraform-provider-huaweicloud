package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerKubernetesClustersConfigs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_kubernetes_clusters_configs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerKubernetesClustersConfigs_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.cluster_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protect_status"),
				),
			},
		},
	})
}

const testDataSourceContainerKubernetesClustersConfigs_basic = `
data "huaweicloud_hss_container_kubernetes_clusters_configs" "test" {
  enterprise_project_id = "0"

  cluster_info_list {
    cluster_id   = "test-cluster-id"
    cluster_name = "test-cluster-name"
  }
}
`
