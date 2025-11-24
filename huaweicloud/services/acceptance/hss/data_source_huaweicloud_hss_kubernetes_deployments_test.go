package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `data_list`.
func TestAccDataSourceKubernetesDeployments_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_kubernetes_deployments.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKubernetesDeployments_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "resources_info_list.#"),
				),
			},
		},
	})
}

func testDataSourceKubernetesDeployments_basic() string {
	return `
data "huaweicloud_hss_kubernetes_deployments" "test" {
  enterprise_project_id = "0"
}
`
}
