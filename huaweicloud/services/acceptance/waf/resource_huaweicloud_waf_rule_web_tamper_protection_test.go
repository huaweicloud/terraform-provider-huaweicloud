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

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccRuleWebTamperProtection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_web_tamper_protection.test"

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
				Config: testAccDataSourceWebTamperProtection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain", "www.abc.com"),
					resource.TestCheckResourceAttr(rName, "path", "/a"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
				),
			},
			{
				Config: testAccDataSourceWebTamperProtectionn_basic_update(name),
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

func testAccDataSourceWebTamperProtection_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_web_tamper_protection" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  domain                = "www.abc.com"
  path                  = "/a"
  description           = "test description"
  status                = 0
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataSourceWebTamperProtectionn_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_web_tamper_protection" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  domain                = "www.abc.com"
  path                  = "/a"
  description           = "test description"
  status                = 1
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
