package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTemplateDeployParams_Basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_template_deploy_params.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTemplateDeployParams_Basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.#"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.nullable"),
					resource.TestCheckResourceAttrSet(dataSource, "variables.0.sensitive"),
				),
			},
		},
	})
}

const testAccDataSourceTemplateDeployParams_Basic = `
data "huaweicloud_rgc_template_deploy_params" "test" {
  template_name    = "VPC"
  template_version = "V1"
}
`
