package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHssContainerKubernetesTemplate_basic(t *testing.T) {
	dataSource := "data.huaweicloud_hss_container_kubernetes_template.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceHssContainerKubernetesTemplate_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "is_default"),
					resource.TestCheckResourceAttrSet(dataSource, "runtime_info.0.runtime_name"),
					resource.TestCheckResourceAttrSet(dataSource, "runtime_info.0.runtime_path"),
				),
			},
		},
	})
}

const testDataSourceDataSourceHssContainerKubernetesTemplate_basic = `
data "huaweicloud_hss_container_kubernetes_template" "test" {
  enterprise_project_id = "all_granted_eps"
}`
