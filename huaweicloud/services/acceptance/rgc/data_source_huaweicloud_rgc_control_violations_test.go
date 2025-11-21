package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceControlViolations_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_control_violations.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCControlViolations(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceControlViolations_account(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.control_id"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.parent_organizational_unit_id"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.parent_organizational_unit_name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.resource"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.service"),
				),
			},
			{
				Config: testDataSourceDataSourceRgcControlViolations_ou(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.#"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.account_name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.control_id"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.display_name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.parent_organizational_unit_id"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.parent_organizational_unit_name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.resource"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.resource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "control_violations.0.service"),
				),
			},
		},
	})
}

func testDataSourceDataSourceControlViolations_account() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_control_violations" "test" {
  account_id = "%[1]s"
}
`, acceptance.HW_RGC_ACCOUNT_ID)
}

func testDataSourceDataSourceRgcControlViolations_ou() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_control_violations" "test" {
  organizational_unit_id = "%[1]s"
}
`, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
