package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationalUnitEnabledControl_basic(t *testing.T) {
	rName := "data.huaweicloud_rgc_organizational_unit_enabled_control.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOrganizationalUnitEnabledControl(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationalUnitEnabledControl_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "control.#"),
					resource.TestCheckResourceAttrSet(rName, "control.0.control_identifier"),
					resource.TestCheckResourceAttrSet(rName, "regions.#"),
					resource.TestCheckResourceAttrSet(rName, "state"),
					resource.TestCheckResourceAttrSet(rName, "version"),
				),
			},
		},
	})
}

func testAccDataSourceOrganizationalUnitEnabledControl_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_organizational_unit_enabled_control" "test" {
  control_id                     = "%[1]s"
  managed_organizational_unit_id = "%[2]s"
}
`, acceptance.HW_RGC_CONTROL_ID, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
