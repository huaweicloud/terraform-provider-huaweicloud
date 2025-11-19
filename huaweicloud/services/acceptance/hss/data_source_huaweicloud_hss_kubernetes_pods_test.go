package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations,
// this test case can only test the scenario with empty `data_list` and `pod_info_list`.
func TestAccDataSourceKubernetesPods_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_kubernetes_pods.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKubernetesPods_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "pod_info_list.#"),
				),
			},
		},
	})
}

func testDataSourceKubernetesPods_basic() string {
	return `
data "huaweicloud_hss_kubernetes_pods" "test" {
  pod_name              = "non-exist"
  enterprise_project_id = "0"
}
`
}
