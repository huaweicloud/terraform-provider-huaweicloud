package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceControls_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_controls.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceControls_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "controls.#"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.identifier"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.guidance"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.resource.#"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.framework.#"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.service"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.implementation"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.behavior"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.owner"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.control_objective"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "controls.0.release_date"),
				),
			},
		},
	})
}

const testAccDataSourceControls_basic = `
data "huaweicloud_rgc_controls" "test" {}
`
