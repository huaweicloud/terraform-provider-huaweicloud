/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/datamasking_rules"
)

func TestAccWafRuleDataMasking_basic(t *testing.T) {
	var rule rules.DataMasking
	policyName := acceptance.RandomAccResourceName()
	resourceName1 := "huaweicloud_waf_rule_data_masking.rule_1"
	resourceName2 := "huaweicloud_waf_rule_data_masking.rule_2"
	resourceName3 := "huaweicloud_waf_rule_data_masking.rule_3"
	resourceName4 := "huaweicloud_waf_rule_data_masking.rule_4"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckWafRuleDataMaskingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleDataMasking_basic(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleDataMaskingExists(resourceName1, &rule),
					resource.TestCheckResourceAttr(resourceName1, "path", "/login"),
					resource.TestCheckResourceAttr(resourceName1, "subfield", "password"),
					resource.TestCheckResourceAttr(resourceName1, "field", "params"),
					resource.TestCheckResourceAttr(resourceName2, "field", "header"),
					resource.TestCheckResourceAttr(resourceName3, "field", "form"),
					resource.TestCheckResourceAttr(resourceName4, "field", "cookie"),
				),
			},
			{
				Config: testAccWafRuleDataMasking_update(policyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafRuleDataMaskingExists(resourceName1, &rule),
					resource.TestCheckResourceAttr(resourceName1, "path", "/login_new"),
					resource.TestCheckResourceAttr(resourceName1, "subfield", "secret"),
					resource.TestCheckResourceAttr(resourceName1, "field", "params"),
					resource.TestCheckResourceAttr(resourceName2, "field", "header"),
					resource.TestCheckResourceAttr(resourceName3, "field", "form"),
					resource.TestCheckResourceAttr(resourceName4, "field", "cookie"),
				),
			},
			{
				ResourceName:      resourceName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccWafRuleImportStateIdFunc(resourceName1),
			},
		},
	})
}

func testAccCheckWafRuleDataMaskingDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	wafClient, err := config.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_waf_rule_data_masking" {
			continue
		}

		policyID := rs.Primary.Attributes["policy_id"]
		_, err := rules.Get(wafClient, policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("WAF data masking rule still exists")
		}
	}

	return nil
}

func testAccCheckWafRuleDataMaskingExists(n string, rule *rules.DataMasking) resource.TestCheckFunc {
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
			return fmt.Errorf("WAF data masking rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccWafRuleDataMasking_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_data_masking" "rule_1" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login"
  field     = "params"
  subfield  = "password"
}
resource "huaweicloud_waf_rule_data_masking" "rule_2" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login"
  field     = "header"
  subfield  = "password"
}
resource "huaweicloud_waf_rule_data_masking" "rule_3" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login"
  field     = "form"
  subfield  = "password"
}
resource "huaweicloud_waf_rule_data_masking" "rule_4" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login"
  field     = "cookie"
  subfield  = "password"
}
`, testAccWafPolicyV1_basic(name))
}

func testAccWafRuleDataMasking_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_data_masking" "rule_1" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login_new"
  field     = "params"
  subfield  = "secret"
}
resource "huaweicloud_waf_rule_data_masking" "rule_2" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login"
  field     = "header"
  subfield  = "secret"
}
resource "huaweicloud_waf_rule_data_masking" "rule_3" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login"
  field     = "form"
  subfield  = "secret"
}
resource "huaweicloud_waf_rule_data_masking" "rule_4" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path      = "/login"
  field     = "cookie"
  subfield  = "secret"
}
`, testAccWafPolicyV1_basic(name))
}
