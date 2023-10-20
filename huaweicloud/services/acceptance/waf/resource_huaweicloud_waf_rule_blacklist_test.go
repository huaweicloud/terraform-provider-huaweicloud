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

func TestAccWafRuleBlackList_basic(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	rName1 := "huaweicloud_waf_rule_blacklist.rule_1"
	rName2 := "huaweicloud_waf_rule_blacklist.rule_2"
	rName3 := "huaweicloud_waf_rule_blacklist.rule_3"

	rc := acceptance.InitResourceCheck(
		rName1,
		&obj,
		getRuleBlackListResourceFunc,
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
				Config: testAccWafRuleBlackList_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.0/24"),
					resource.TestCheckResourceAttr(rName1, "action", "0"),
					resource.TestCheckResourceAttr(rName1, "status", "0"),

					resource.TestCheckResourceAttr(rName2, "ip_address", "192.165.0.0/24"),
					resource.TestCheckResourceAttr(rName2, "action", "1"),

					resource.TestCheckResourceAttr(rName3, "ip_address", "192.160.0.0/24"),
					resource.TestCheckResourceAttr(rName3, "action", "0"),
					resource.TestCheckResourceAttr(rName3, "name", randName),
					resource.TestCheckResourceAttr(rName3, "description", "test description"),
				),
			},
			{
				Config: testAccWafRuleBlackList_update(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName1, "ip_address", "192.168.0.125"),
					resource.TestCheckResourceAttr(rName1, "action", "2"),
					resource.TestCheckResourceAttr(rName1, "status", "1"),

					resource.TestCheckResourceAttr(rName2, "ip_address", "192.150.0.0/24"),
					resource.TestCheckResourceAttr(rName2, "action", "0"),

					resource.TestCheckResourceAttrPair(rName3, "address_group_id",
						"huaweicloud_waf_address_group.test", "id"),
					resource.TestCheckResourceAttr(rName3, "action", "2"),
					resource.TestCheckResourceAttr(rName3, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(rName3, "description", ""),
					resource.TestCheckResourceAttrSet(rName3, "address_group_name"),
					resource.TestCheckResourceAttrSet(rName3, "address_group_size"),
				),
			},
			{
				ResourceName:      rName1,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFRuleImportState(rName1),
			},
		},
	})
}

func TestAccWafRuleBlackList_withEpsID(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_blacklist.rule"

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
				Config: testAccWafRuleBlackList_basic_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "ip_address", "192.160.0.0/24"),
					resource.TestCheckResourceAttr(rName, "action", "0"),
					resource.TestCheckResourceAttr(rName, "name", randName),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
				),
			},
			{
				Config: testAccWafRuleBlackList_update_withEpsID(randName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "address_group_id",
						"huaweicloud_waf_address_group.test", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "action", "2"),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttrSet(rName, "address_group_name"),
					resource.TestCheckResourceAttrSet(rName, "address_group_size"),
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

func testAccWafRuleBlackList_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_blacklist" "rule_1" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.168.0.0/24"
  status     = 0
}

resource "huaweicloud_waf_rule_blacklist" "rule_2" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.165.0.0/24"
  action     = 1
}

resource "huaweicloud_waf_rule_blacklist" "rule_3" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  ip_address  = "192.160.0.0/24"
  name        = "%s"
  description = "test description"
}
`, testAccWafPolicyV1_basic(name), name)
}

func testAccWafRuleBlackList_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_address_group" "test" {
  name         = "%s"
  description  = "example_description"
  ip_addresses = ["192.168.1.0/24"]

  depends_on   = [huaweicloud_waf_dedicated_instance.instance_1]
}

resource "huaweicloud_waf_rule_blacklist" "rule_1" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.168.0.125"
  action     = 2
  status     = 1
}

resource "huaweicloud_waf_rule_blacklist" "rule_2" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  ip_address = "192.150.0.0/24"
  action     = 0
}

resource "huaweicloud_waf_rule_blacklist" "rule_3" {
  policy_id        = huaweicloud_waf_policy.policy_1.id
  address_group_id = huaweicloud_waf_address_group.test.id
  action           = 2
  name             = "%s_update"
  description      = ""
}
`, testAccWafPolicyV1_basic(name), name, name)
}

func testAccWafRuleBlackList_basic_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_blacklist" "rule" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  ip_address            = "192.160.0.0/24"
  name                  = "%s"
  description           = "test description"
  enterprise_project_id = "%s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), name, epsID)
}

func testAccWafRuleBlackList_update_withEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_address_group" "test" {
  name                  = "%[2]s"
  description           = "example_description"
  ip_addresses          = ["192.168.1.0/24"]
  enterprise_project_id = "%[3]s"

  depends_on   = [huaweicloud_waf_dedicated_instance.instance_1]
}

resource "huaweicloud_waf_rule_blacklist" "rule" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  address_group_id      = huaweicloud_waf_address_group.test.id
  action                = 2
  name                  = "%[2]s_update"
  description           = ""
  enterprise_project_id = "%[3]s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), name, epsID)
}
