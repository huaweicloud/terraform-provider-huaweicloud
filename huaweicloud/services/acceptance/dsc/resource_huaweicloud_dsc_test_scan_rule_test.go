package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTestScanRule_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a scan rule ID, and the values of each parameter in the script must be
			// consistent with the scan rule.
			acceptance.TestAccPreCheckDscScanRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testTestScanRule_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("huaweicloud_dsc_test_scan_rule.test", "is_match"),
				),
			},
		},
	})
}

func testTestScanRule_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dsc_test_scan_rule" "test" {
  category       = "BUILT_SELF"
  data           = "test"
  effective_mode = "IN"
  location       = "NAME"
  rule_content   = ["location"]
  rule_id        = "%s"
}
`, acceptance.HW_DSC_SCAN_RULE_ID)
}
