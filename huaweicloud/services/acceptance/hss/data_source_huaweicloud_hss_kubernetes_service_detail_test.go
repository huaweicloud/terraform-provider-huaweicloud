package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `service_info_list`.
func TestAccDataSourceKubernetesServiceDetail_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_kubernetes_service_detail.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKubernetesServiceDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "service_port_list.#"),
				),
			},
		},
	})
}

// The `service_id` used is dummy data for testing.
func testDataSourceKubernetesServiceDetail_basic() string {
	return `
data "huaweicloud_hss_kubernetes_service_detail" "test" {
  service_id = "f69812ba-bf72-11f0-a281-0255ac10024f"
}
`
}
