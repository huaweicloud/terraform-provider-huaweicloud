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

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccRuleDataMasking_basic(t *testing.T) {
	var obj interface{}

	policyName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_rule_data_masking.test"

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
				Config: testAccWafRuleDataSourceMasking_basic(policyName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "path", "/login"),
					resource.TestCheckResourceAttr(resourceName, "subfield", "password"),
					resource.TestCheckResourceAttr(resourceName, "field", "params"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "status", "0"),
				),
			},
			{
				Config: testAccWafRuleDataSourceMasking_update(policyName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "path", "/login_new"),
					resource.TestCheckResourceAttr(resourceName, "subfield", "secret"),
					resource.TestCheckResourceAttr(resourceName, "field", "params"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "status", "1"),
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

func testAccWafRuleDataSourceMasking_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_data_masking" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  path                  = "/login"
  field                 = "params"
  subfield              = "password"
  description           = "test description"
  status                = 0
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafRuleDataSourceMasking_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_data_masking" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  path                  = "/login_new"
  field                 = "params"
  subfield              = "secret"
  description           = "test description update"
  status                = 1
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
