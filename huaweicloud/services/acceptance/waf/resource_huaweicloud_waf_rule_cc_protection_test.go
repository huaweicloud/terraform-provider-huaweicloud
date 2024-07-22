package waf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getRuleCCProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		getRuleCCProtectionHttpUrl = "v1/{project_id}/waf/policy/{policy_id}/cc/{rule_id}"
		getRuleCCProtectionProduct = "waf"
	)
	getRuleCCProtectionClient, err := cfg.NewServiceClient(getRuleCCProtectionProduct, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF Client: %s", err)
	}

	getPath := getRuleCCProtectionClient.Endpoint + getRuleCCProtectionHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", getRuleCCProtectionClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", state.Primary.Attributes["policy_id"])
	getPath = strings.ReplaceAll(getPath, "{rule_id}", state.Primary.ID)

	queryParam := ""
	if epsID := state.Primary.Attributes["enterprise_project_id"]; epsID != "" {
		queryParam = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	getPath += queryParam

	getRuleCCProtectionOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRuleCCProtectionResp, err := getRuleCCProtectionClient.Request("GET", getPath, &getRuleCCProtectionOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RuleCCProtection: %s", err)
	}
	return utils.FlattenResponse(getRuleCCProtectionResp)
}

func TestAccRuleCCProtection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_cc_protection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleCCProtectionResourceFunc,
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
				Config: testRuleCCProtection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "protective_action", "block"),
					resource.TestCheckResourceAttr(rName, "rate_limit_mode", "cookie"),
					resource.TestCheckResourceAttr(rName, "block_page_type", "application/json"),
					resource.TestCheckResourceAttr(rName, "page_content", "test page content"),
					resource.TestCheckResourceAttr(rName, "user_identifier", "test_identifier"),
					resource.TestCheckResourceAttr(rName, "limit_num", "10"),
					resource.TestCheckResourceAttr(rName, "limit_period", "60"),
					resource.TestCheckResourceAttr(rName, "lock_time", "5"),
					resource.TestCheckResourceAttr(rName, "request_aggregation", "true"),
					resource.TestCheckResourceAttr(rName, "all_waf_instances", "true"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "params"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "test content"),
					resource.TestCheckResourceAttr(rName, "conditions.0.subfield", "test_subfield"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "ip"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "equal"),
					resource.TestCheckResourceAttr(rName, "conditions.1.content", "192.168.0.1"),
				),
			},
			{
				Config: testRuleCCProtection_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "protective_action", "dynamic_block"),
					resource.TestCheckResourceAttr(rName, "rate_limit_mode", "policy"),
					resource.TestCheckResourceAttr(rName, "limit_num", "20"),
					resource.TestCheckResourceAttr(rName, "limit_period", "100"),
					resource.TestCheckResourceAttr(rName, "unlock_num", "15"),
					resource.TestCheckResourceAttr(rName, "all_waf_instances", "false"),
					resource.TestCheckResourceAttr(rName, "description", ""),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "response_code"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "equal"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "200"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "header"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "equal_any"),
					resource.TestCheckResourceAttr(rName, "conditions.1.subfield", "test_subfield"),
					resource.TestCheckResourceAttrPair(rName, "conditions.1.reference_table_id", "huaweicloud_waf_reference_table.ref_table", "id"),
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

func TestAccRuleCCProtection_withEpsID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_cc_protection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleCCProtectionResourceFunc,
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
				Config: testRuleCCProtection_basicWithEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "protective_action", "block"),
					resource.TestCheckResourceAttr(rName, "rate_limit_mode", "other"),
					resource.TestCheckResourceAttr(rName, "other_user_identifier", "test_referer"),
					resource.TestCheckResourceAttr(rName, "limit_num", "20"),
					resource.TestCheckResourceAttr(rName, "limit_period", "70"),
					resource.TestCheckResourceAttr(rName, "request_aggregation", "true"),
					resource.TestCheckResourceAttr(rName, "all_waf_instances", "true"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "url"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "test_url"),
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

func testRuleCCProtection_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_cc_protection" "test" {
  policy_id           = huaweicloud_waf_policy.policy_1.id
  name                = "%s"
  protective_action   = "block"
  rate_limit_mode     = "cookie"
  block_page_type     = "application/json"
  page_content        = "test page content"
  user_identifier     = "test_identifier"
  limit_num           = 10
  limit_period        = 60
  lock_time           = 5
  request_aggregation = true
  all_waf_instances   = true
  description         = "test description"
  status              = 0

  conditions {
    field    = "params"
    logic    = "contain"
    content  = "test content"
    subfield = "test_subfield"
  }

  conditions {
    field   = "ip"
    logic   = "equal"
    content = "192.168.0.1"
  }
}
`, testAccWafPolicyV1_basic(name), name)
}

func testRuleCCProtection_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_reference_table" "ref_table" {
  name        = "%s"
  type        = "header"
  description = "tf acc"

  conditions = [
    "test_table"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}

resource "huaweicloud_waf_rule_cc_protection" "test" {
  policy_id         = huaweicloud_waf_policy.policy_1.id
  name              = "%s_update"
  protective_action = "dynamic_block"
  rate_limit_mode   = "policy"
  limit_num         = 20
  limit_period      = 100
  unlock_num        = 15
  all_waf_instances = false
  status            = 1

  conditions {
    field    = "response_code"
    logic    = "equal"
    content  = "200"
  }

  conditions {
    field              = "header"
    logic              = "equal_any"
    subfield           = "test_subfield"
    reference_table_id = huaweicloud_waf_reference_table.ref_table.id
  }
}
`, testAccWafPolicyV1_basic(name), name, name)
}

func testRuleCCProtection_basicWithEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_cc_protection" "test" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  name                  = "%s"
  enterprise_project_id = "%s"
  protective_action     = "block"
  rate_limit_mode       = "other"
  other_user_identifier = "test_referer"
  limit_num             = 20
  limit_period          = 70
  request_aggregation   = true
  all_waf_instances     = true
  description           = "test description"

  conditions {
    field    = "url"
    logic    = "contain"
    content  = "test_url"
  }
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), name, epsID)
}
