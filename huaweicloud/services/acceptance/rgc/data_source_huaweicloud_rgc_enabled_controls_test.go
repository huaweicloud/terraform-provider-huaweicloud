package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEnabledControls_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_enabled_controls.enabled_controls"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnabledControls_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.#"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.manage_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.control_identifier"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.control_objective"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.behavior"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.owner"),
					resource.TestCheckResourceAttrSet(dataSource, "enabled_controls.0.regional_preference"),
				),
			},
		},
	})
}

func testAccDataSourceEnabledControls_basic() string {
	return `
data "huaweicloud_rgc_enabled_controls" "enabled_controls" {}
`
}
