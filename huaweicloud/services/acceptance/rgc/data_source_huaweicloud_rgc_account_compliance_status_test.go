package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccountComplianceStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_account_compliance_status.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAccountComplianceStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "compliance_status"),
				),
			},
		},
	})
}

func testDataSourceDataSourceAccountComplianceStatus_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_account_compliance_status" "test" {
  managed_account_id = "%[1]s"
}
`, acceptance.HW_RGC_ACCOUNT_ID)
}
