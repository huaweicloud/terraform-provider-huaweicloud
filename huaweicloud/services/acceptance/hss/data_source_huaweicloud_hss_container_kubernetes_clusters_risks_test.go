package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `data_list`.
func TestAccDataSourceContainerKubernetesClustersRisks_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_kubernetes_clusters_risks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerKubernetesClustersRisks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
				),
			},
		},
	})
}

// The `cluster_id_list` used is dummy data for testing.
func testDataSourceContainerKubernetesClustersRisks_basic() string {
	return `
data "huaweicloud_hss_container_kubernetes_clusters_risks" "test" {
  cluster_id_list       = ["f69812ba-bf72-11f0-a281-0255ac10024f"]
  detect_type           = "image"
  enterprise_project_id = "0"
}
`
}
