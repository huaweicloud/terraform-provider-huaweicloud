package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `service_info_list`.
func TestAccDataSourceKubernetesServices_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_kubernetes_services.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKubernetesServices_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_info_list.#"),
				),
			},
		},
	})
}

const testAccDataSourceKubernetesServices_basic = `
data "huaweicloud_hss_kubernetes_services" "test" {
  enterprise_project_id = "0"
}
`
