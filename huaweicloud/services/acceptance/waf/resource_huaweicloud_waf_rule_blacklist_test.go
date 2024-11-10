package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/whiteblackip_rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRuleBlackListResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	wafClient, err := cfg.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}

	policyID := state.Primary.Attributes["policy_id"]
	epsID := state.Primary.Attributes["enterprise_project_id"]
	return rules.GetWithEpsId(wafClient, policyID, state.Primary.ID, epsID).Extract()
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccRuleBlackList_basic(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_blacklist.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleBlackListResourceFunc,
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
				Config: testAccDataSourceWafRuleBlackList_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(rName, "action", "0"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
				),
			},
			{
				Config: testAccDataSourceWafRuleBlackList_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "ip_address", "192.168.0.125"),
					resource.TestCheckResourceAttr(rName, "action", "2"),
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

func testAccDataSourceWafRuleBlackList_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_blacklist" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  ip_address            = "192.168.0.0/24"
  status                = 0
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccDataSourceWafRuleBlackList_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_blacklist" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  ip_address            = "192.168.0.125"
  action                = 2
  status                = 1
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
