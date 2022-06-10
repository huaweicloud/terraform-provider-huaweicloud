package waf

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/whiteblackip_rules"
)

func TestAccWafRuleBlackList_basic(t *testing.T) {
	var rule rules.WhiteBlackIP
	randName := acceptance.RandomAccResourceName()
	rName1 := "huaweicloud_waf_rule_blacklist.rule_1"
	rName2 := "huaweicloud_waf_rule_blacklist.rule_2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafRuleBlackListDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleBlackList_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(rName1, &rule),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(rName1, "action", "0"),

					testAccCheckWafRuleBlackListExists(rName2, &rule),
					resource.TestCheckResourceAttr(rName2, "ip_address", "192.165.0.0/24"),
					resource.TestCheckResourceAttr(rName2, "action", "1"),
				),
			},
			{
				Config: testAccWafRuleBlackList_update(randName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleBlackListExists(rName1, &rule),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.125"),
					resource.TestCheckResourceAttr(rName1, "action", "2"),

					testAccCheckWafRuleBlackListExists(rName2, &rule),
					resource.TestCheckResourceAttr(rName2, "ip_address", "192.150.0.0/24"),
					resource.TestCheckResourceAttr(rName2, "action", "0"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccWafRuleImportStateIdFunc(rName1),
			},
		},
	})
}

func testAccCheckWafRuleBlackListDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_rule_blacklist" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Waf rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleBlackListExists(n string, rule *rules.WhiteBlackIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
		}

		policyID := rs.Primary.Attributes["policy_id"]
		found, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Id != rs.Primary.ID {
			return fmt.Errorf("WAF black list rule not found")
		}

		*rule = *found

		return nil
	}
}

// testAccWafRuleImportStateIdFunc is used to test exporting rule information from the HuaweiCloud to terraform.
// It is also called by other rules unit tests.
func testAccWafRuleImportStateIdFunc(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policy, ok := s.RootModule().Resources["huaweicloud_waf_policy.policy_1"]
		if !ok {
			return "", fmt.Errorf("WAF policy not found")
		}
		rule, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("WAF rule not found")
		}

		if policy.Primary.ID == "" || rule.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", policy.Primary.ID, rule.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", policy.Primary.ID, rule.Primary.ID), nil
	}
}

func testAccWafRuleBlackList_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_blacklist" "rule_1" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.168.0.0/24"
}

resource "huaweicloud_waf_rule_blacklist" "rule_2" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.165.0.0/24"
  action     = 1
}
`, testAccWafPolicyV1_basic(name))
}

func testAccWafRuleBlackList_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_blacklist" "rule_1" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.168.0.125"
  action     = 2
}

resource "huaweicloud_waf_rule_blacklist" "rule_2" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.150.0.0/24"
  action     = 0
}
`, testAccWafPolicyV1_basic(name))
}
