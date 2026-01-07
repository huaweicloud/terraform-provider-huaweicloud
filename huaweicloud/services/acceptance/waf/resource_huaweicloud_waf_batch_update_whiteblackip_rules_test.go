package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchUpdateWhiteBlackIpRules_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	cidr := acceptance.RandomCidr()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckWafPolicyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccBatchUpdateWhiteBlackIpRules_basic(name, cidr),
			},
		},
	})
}

func testAccBatchUpdateWhiteBlackIpRules_basic(name, cidr string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_rule_blacklist" "test" {
  name       = "%[1]s"
  action     = 1
  policy_id  = "%[2]s"
  ip_address = "%[3]s"
}

resource "huaweicloud_waf_batch_update_whiteblackip_rules" "test" {
  name  = "%[1]s"
  white = 1
  addr  = "%[3]s"

  policy_rule_ids {
    policy_id = "%[2]s"
    rule_ids  = [huaweicloud_waf_rule_blacklist.test.id]
  }
}
`, name, acceptance.HW_WAF_POLICY_ID, cidr)
}
