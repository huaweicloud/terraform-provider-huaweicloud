package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationalUnitComplianceStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_organizational_unit_compliance_status.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCOrganizationID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationalUnitComplianceStatusConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_status"),
				),
			},
		},
	})
}

func testAccDataSourceOrganizationalUnitComplianceStatusConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_organizational_unit_compliance_status" "test" {
  managed_organizational_unit_id = "%[1]s"
}
`, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
