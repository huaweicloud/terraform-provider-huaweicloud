package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWebTamperRefresh_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, please ensure that there is at least one WAF instance in the current region.
			// Prepare a WAF policy and a web tamper protection rule.
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafPolicyId(t)
			acceptance.TestAccPreCheckWafWebTamperRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRuleWebTamperRefresh_basic(),
			},
		},
	})
}

func testAccRuleWebTamperRefresh_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_rule_web_tamper_protection_refresh" "test" {
  policy_id = "%[1]s"
  rule_id   = "%[2]s"
}
`, acceptance.HW_WAF_POLICY_ID, acceptance.HW_WAF_WEB_TAMPER_RULE_ID)
}
