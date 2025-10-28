package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAntivirusPayPerScanSwitchStatus_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAntivirusPayPerScanSwitchStatus_basic(),
			},
		},
	})
}

func testAntivirusPayPerScanSwitchStatus_basic() string {
	return `
resource "huaweicloud_hss_antivirus_pay_per_scan_switch_status" "test" {
  enabled               = true
  enterprise_project_id = "0"
}
`
}
