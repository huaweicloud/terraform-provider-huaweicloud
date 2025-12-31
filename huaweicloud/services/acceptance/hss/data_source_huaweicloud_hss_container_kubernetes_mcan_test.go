package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerKubernetesMCAN_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_container_kubernetes_mcan.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceContainerKubernetesMCAN_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "anp_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "region_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "agent_address"),
				),
			},
		},
	})
}

func testDataSourceContainerKubernetesMCAN_basic() string {
	return `
data "huaweicloud_hss_container_kubernetes_mcan" "test" {
  enterprise_project_id = "0"
}
`
}
