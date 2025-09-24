package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCcProtectionRuleBatchDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, please ensure that there is at least one WAF instance in the current region.
			// Prepare a WAF policy with a WAF CC protection rule.
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafPolicyId(t)
			acceptance.TestAccPreCheckWafRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccCcProtectionRuleBatchDelete_basic(),
			},
		},
	})
}

func testAccCcProtectionRuleBatchDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_cc_protection_rule_batch_delete" "test" {
  policy_rule_ids {
    policy_id = "%[1]s"
    rule_ids  = ["%[2]s"]
  }
}
`, acceptance.HW_WAF_POLICY_ID, acceptance.HW_WAF_RULE_ID)
}
