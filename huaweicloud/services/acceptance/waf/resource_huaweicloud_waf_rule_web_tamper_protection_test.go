/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/webtamperprotection_rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRuleWebTamperProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	wafClient, err := cfg.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	policyID := state.Primary.Attributes["policy_id"]
	epsID := state.Primary.Attributes["enterprise_project_id"]
	return rules.GetWithEpsID(wafClient, policyID, state.Primary.ID, epsID).Extract()
}

func TestAccWafRuleWebTamperProtection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_web_tamper_protection.rule_1"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleWebTamperProtectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleWebTamperProtection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain", "www.abc.com"),
					resource.TestCheckResourceAttr(rName, "path", "/a"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
				),
			},
			{
				Config: testAccWafRuleWebTamperProtection_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "1"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName),
			},
		},
	})
}

func TestAccWafRuleWebTamperProtection_withEpsID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_web_tamper_protection.rule_1"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleWebTamperProtectionResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWafRuleWebTamperProtection_basic_withEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "domain", "www.abc.com"),
					resource.TestCheckResourceAttr(rName, "path", "/a"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName),
			},
		},
	})
}

func testAccWafRuleWebTamperProtection_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_web_tamper_protection" "rule_1" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  domain      = "www.abc.com"
  path        = "/a"
  description = "test description"
  status      = 0
}
`, testAccWafPolicyV1_basic(name))
}

func testAccWafRuleWebTamperProtection_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_web_tamper_protection" "rule_1" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  domain      = "www.abc.com"
  path        = "/a"
  description = "test description"
  status      = 1
}
`, testAccWafPolicyV1_basic(name))
}

func testAccWafRuleWebTamperProtection_basic_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_web_tamper_protection" "rule_1" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  domain                = "www.abc.com"
  path                  = "/a"
  description           = "test description"
  enterprise_project_id = "%s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), epsID)
}
