package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationalUnitControl_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_organizational_unit_controls.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCOrganizationID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationalUnitControlConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.#"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.manage_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.control_identifier"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.control_objective"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.behavior"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.owner"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.regional_preference"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.guidance"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.service"),
					resource.TestCheckResourceAttrSet(dataSource, "control_summaries.0.implementation"),
				),
			},
		},
	})
}

func testAccDataSourceOrganizationalUnitControlConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_organizational_unit_controls" "test" {
  managed_organizational_unit_id = "%[1]s"
}
`, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
