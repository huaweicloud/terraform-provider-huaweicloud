package gaussdb

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBParameterSettingTemplate_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_parameter_setting_template.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBParameterSettingTemplate_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.#"),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.0.restart_required"),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.0.value_range"),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "expansion_parameters.0.risk_description"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBParameterSettingTemplate_basic() string {
	return `
data "huaweicloud_gaussdb_parameter_setting_template" "test" {}
`
}
