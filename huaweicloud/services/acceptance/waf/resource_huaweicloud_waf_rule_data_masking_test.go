/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/datamasking_rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRuleDataMaskingResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	wafClient, err := cfg.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	policyID := state.Primary.Attributes["policy_id"]
	epsID := state.Primary.Attributes["enterprise_project_id"]
	return rules.GetWithEpsID(wafClient, policyID, state.Primary.ID, epsID).Extract()
}

func TestAccWafRuleDataMasking_basic(t *testing.T) {
	var obj interface{}

	policyName := acceptance.RandomAccResourceName()
	resourceName1 := "huaweicloud_waf_rule_data_masking.rule_1"
	resourceName2 := "huaweicloud_waf_rule_data_masking.rule_2"
	resourceName3 := "huaweicloud_waf_rule_data_masking.rule_3"
	resourceName4 := "huaweicloud_waf_rule_data_masking.rule_4"

	rc := acceptance.InitResourceCheck(
		resourceName1,
		&obj,
		getRuleDataMaskingResourceFunc,
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
				Config: testAccWafRuleDataMasking_basic(policyName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName1, "path", "/login"),
					resource.TestCheckResourceAttr(resourceName1, "subfield", "password"),
					resource.TestCheckResourceAttr(resourceName1, "field", "params"),
					resource.TestCheckResourceAttr(resourceName1, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName1, "status", "0"),
					resource.TestCheckResourceAttr(resourceName2, "field", "header"),
					resource.TestCheckResourceAttr(resourceName3, "field", "form"),
					resource.TestCheckResourceAttr(resourceName4, "field", "cookie"),
				),
			},
			{
				Config: testAccWafRuleDataMasking_update(policyName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName1, "path", "/login_new"),
					resource.TestCheckResourceAttr(resourceName1, "subfield", "secret"),
					resource.TestCheckResourceAttr(resourceName1, "field", "params"),
					resource.TestCheckResourceAttr(resourceName1, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName1, "status", "1"),
					resource.TestCheckResourceAttr(resourceName2, "field", "header"),
					resource.TestCheckResourceAttr(resourceName3, "field", "form"),
					resource.TestCheckResourceAttr(resourceName4, "field", "cookie"),
				),
			},
			{
				ResourceName:      resourceName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(resourceName1),
			},
		},
	})
}

func TestAccWafRuleDataMasking_withEpsID(t *testing.T) {
	var obj interface{}

	policyName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_rule_data_masking.rule"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getRuleDataMaskingResourceFunc,
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
				Config: testAccWafRuleDataMasking_basic_withEpsID(policyName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "path", "/login"),
					resource.TestCheckResourceAttr(resourceName, "subfield", "password"),
					resource.TestCheckResourceAttr(resourceName, "field", "params"),
				),
			},
			{
				Config: testAccWafRuleDataMasking_update_withEpsID(policyName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "path", "/login_new"),
					resource.TestCheckResourceAttr(resourceName, "subfield", "secret"),
					resource.TestCheckResourceAttr(resourceName, "field", "params"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(resourceName),
			},
		},
	})
}

func testAccWafRuleDataMasking_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_data_masking" "rule_1" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  path        = "/login"
  field       = "params"
  subfield    = "password"
  description = "test description"
  status      = 0
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
  policy_id   = huaweicloud_waf_policy.policy_1.id
  path        = "/login_new"
  field       = "params"
  subfield    = "secret"
  description = "test description update"
  status      = 1
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

func testAccWafRuleDataMasking_basic_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_data_masking" "rule" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  path                  = "/login"
  field                 = "params"
  subfield              = "password"
  enterprise_project_id = "%s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), epsID)
}

func testAccWafRuleDataMasking_update_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_data_masking" "rule" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  path                  = "/login_new"
  field                 = "params"
  subfield              = "secret"
  enterprise_project_id = "%s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), epsID)
}
