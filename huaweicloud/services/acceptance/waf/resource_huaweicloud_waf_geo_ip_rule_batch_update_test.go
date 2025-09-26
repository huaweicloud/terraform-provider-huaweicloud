package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeoIpRuleBatchUpdate_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, please ensure that there is at least one WAF instance in the current region.
			// Prepare a WAF policy with a WAF geo IP rule.
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafPolicyId(t)
			acceptance.TestAccPreCheckWafGeoRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeoIpRuleBatchUpdate_basic(),
			},
		},
	})
}

func testAccGeoIpRuleBatchUpdate_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_geo_ip_rule_batch_update" "test" {
  geoip  = "US|CA|JP"
  status = 1
  name   = "updated_geo_rule"
  white  = 2

  policy_rule_ids {
    policy_id = "%[1]s"
    rule_ids  = ["%[2]s"]
  }
}
`, acceptance.HW_WAF_POLICY_ID, acceptance.HW_WAF_GEO_RULE_ID)
}
