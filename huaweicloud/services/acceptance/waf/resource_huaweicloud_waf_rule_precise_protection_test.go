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

func getRulePreciseProtectionResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		preciseProtectionHttpUrl = "v1/{project_id}/waf/policy/{policy_id}/custom/{rule_id}"
		product                  = "waf"
	)
	preciseProtectionClient, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF Client: %s", err)
	}

	getRulePath := preciseProtectionClient.Endpoint + preciseProtectionHttpUrl
	getRulePath = strings.ReplaceAll(getRulePath, "{project_id}", preciseProtectionClient.ProjectID)
	getRulePath = strings.ReplaceAll(getRulePath, "{policy_id}", state.Primary.Attributes["policy_id"])
	getRulePath = strings.ReplaceAll(getRulePath, "{rule_id}", state.Primary.ID)

	queryParam := ""
	if epsID := state.Primary.Attributes["enterprise_project_id"]; epsID != "" {
		queryParam = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	getRulePath += queryParam

	getRuleOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRuleResp, err := preciseProtectionClient.Request("GET", getRulePath, &getRuleOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving RulePreciseProtection: %s", err)
	}
	return utils.FlattenResponse(getRuleResp)
}

func TestAccRulePreciseProtection_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_precise_protection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRulePreciseProtectionResourceFunc,
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
				Config: testRulePreciseProtection_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "priority", "10"),
					resource.TestCheckResourceAttr(rName, "action", "block"),
					resource.TestCheckResourceAttr(rName, "start_time", "2023-05-01 13:01:20"),
					resource.TestCheckResourceAttr(rName, "end_time", "2023-05-10 14:10:30"),
					resource.TestCheckResourceAttr(rName, "description", "description information"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "url"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "login"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "params"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.1.subfield", "param_info"),
					resource.TestCheckResourceAttr(rName, "conditions.1.content", "register"),
				),
			},
			{
				Config: testRulePreciseProtection_basicUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id", "huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttrPair(rName, "conditions.1.reference_table_id", "huaweicloud_waf_reference_table.ref_table", "id"),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "priority", "20"),
					resource.TestCheckResourceAttr(rName, "action", "pass"),
					resource.TestCheckResourceAttr(rName, "start_time", ""),
					resource.TestCheckResourceAttr(rName, "end_time", ""),
					resource.TestCheckResourceAttr(rName, "description", "description information update"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "method"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "equal"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "GET"),
					resource.TestCheckResourceAttr(rName, "conditions.1.field", "header"),
					resource.TestCheckResourceAttr(rName, "conditions.1.logic", "prefix_any"),
					resource.TestCheckResourceAttr(rName, "conditions.1.subfield", "test_sub"),
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

func TestAccRulePreciseProtection_knownAttackSourceId(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_precise_protection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRulePreciseProtectionResourceFunc,
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
				Config: testRulePreciseProtection_knownAttackSourceId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "priority", "20"),
					resource.TestCheckResourceAttr(rName, "action", "block"),
					resource.TestCheckResourceAttrPair(rName, "known_attack_source_id",
						"huaweicloud_waf_rule_known_attack_source.test", "id"),
					resource.TestCheckResourceAttr(rName, "description", "description information"),
				),
			},
			{
				Config: testRulePreciseProtection_updateKnownAttackSourceId(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "priority", "20"),
					resource.TestCheckResourceAttr(rName, "action", "log"),
					resource.TestCheckResourceAttr(rName, "known_attack_source_id", ""),
					resource.TestCheckResourceAttr(rName, "description", "description information"),
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

func TestAccRulePreciseProtection_WithEpsID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_precise_protection.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRulePreciseProtectionResourceFunc,
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
				Config: testRulePreciseProtection_basicWithEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "priority", "10"),
					resource.TestCheckResourceAttr(rName, "action", "block"),
					resource.TestCheckResourceAttr(rName, "status", "1"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "url"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "contain"),
					resource.TestCheckResourceAttr(rName, "conditions.0.content", "login"),
				),
			},
			{
				Config: testRulePreciseProtection_updateWithEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "name", fmt.Sprintf("%s_update", name)),
					resource.TestCheckResourceAttr(rName, "priority", "20"),
					resource.TestCheckResourceAttr(rName, "status", "0"),
					resource.TestCheckResourceAttr(rName, "conditions.0.field", "header"),
					resource.TestCheckResourceAttr(rName, "conditions.0.logic", "prefix_any"),
					resource.TestCheckResourceAttr(rName, "conditions.0.subfield", "test_sub"),
					resource.TestCheckResourceAttrPair(rName, "conditions.0.reference_table_id", "huaweicloud_waf_reference_table.ref_table", "id"),
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

func testRulePreciseProtection_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_precise_protection" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  name        = "%s"
  priority    = 10
  action      = "block"
  start_time  = "2023-05-01 13:01:20"
  end_time    = "2023-05-10 14:10:30"
  description = "description information"
  status      = 0

  conditions {
    field   = "url"
    logic   = "contain"
    content = "login"
  }

  conditions {
    field    = "params"
    logic    = "contain"
    subfield = "param_info"
    content  = "register"
  }
}
`, testAccWafPolicyV1_basic(name), name)
}

func testRulePreciseProtection_basicUpdate(name string) string {
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

resource "huaweicloud_waf_rule_precise_protection" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  name        = "%s_update"
  priority    = 20
  action      = "pass"
  description = "description information update"
  status      = 1

  conditions {
    field   = "method"
    logic   = "equal"
    content = "GET"
  }

  conditions {
    field              = "header"
    logic              = "prefix_any"
    subfield           = "test_sub"
    reference_table_id = huaweicloud_waf_reference_table.ref_table.id
  }
}
`, testAccWafPolicyV1_basic(name), name, name)
}

func testRulePreciseProtection_basicWithEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_precise_protection" "test" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  name                  = "%s"
  priority              = 10
  enterprise_project_id = "%s"

  conditions {
    field   = "url"
    logic   = "contain"
    content = "login"
  }
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), name, epsID)
}

func testRulePreciseProtection_updateWithEpsID(name, epsID string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_reference_table" "ref_table" {
  name                  = "%[2]s"
  type                  = "header"
  description           = "tf acc"
  enterprise_project_id = "%[3]s"

  conditions = [
    "test_table"
  ]

  depends_on = [
    huaweicloud_waf_dedicated_instance.instance_1
  ]
}

resource "huaweicloud_waf_rule_precise_protection" "test" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  name                  = "%[2]s_update"
  priority              = 20
  status                = 0
  enterprise_project_id = "%[3]s"

  conditions {
    field              = "header"
    logic              = "prefix_any"
    subfield           = "test_sub"
    reference_table_id = huaweicloud_waf_reference_table.ref_table.id
  }
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), name, epsID)
}

func testRulePreciseProtection_knownAttackSourceId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_known_attack_source" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  block_type  = "long_ip_block"
  block_time  = 500
  description = "test description"
}

resource "huaweicloud_waf_rule_precise_protection" "test" {
  policy_id              = huaweicloud_waf_policy.policy_1.id
  name                   = "%s"
  priority               = 20
  action                 = "block"
  known_attack_source_id = huaweicloud_waf_rule_known_attack_source.test.id
  description            = "description information"

  conditions {
    field   = "method"
    logic   = "equal"
    content = "GET"
  }

  depends_on = [
    huaweicloud_waf_dedicated_domain.domain_1
  ]
}
`, testAccWafDedicatedDomainV1_policy(name), name)
}

func testRulePreciseProtection_updateKnownAttackSourceId(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_known_attack_source" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  block_type  = "long_ip_block"
  block_time  = 500
  description = "test description"
}

resource "huaweicloud_waf_rule_precise_protection" "test" {
  policy_id              = huaweicloud_waf_policy.policy_1.id
  name                   = "%s"
  priority               = 20
  action                 = "log"
  description            = "description information"

  conditions {
    field   = "method"
    logic   = "equal"
    content = "GET"
  }

  depends_on = [
    huaweicloud_waf_dedicated_domain.domain_1
  ]
}
`, testAccWafDedicatedDomainV1_policy(name), name)
}

// testWAFRuleImportState use to return an id with format <policy_id>/<id> or <policy_id>/<id>/<enterprise_project_id>
func testWAFRuleImportState(name string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", name, rs)
		}

		policyID := rs.Primary.Attributes["policy_id"]
		if policyID == "" {
			return "", fmt.Errorf("attribute (policy_id) of Resource (%s) not found: %s", name, rs)
		}

		epsID := rs.Primary.Attributes["enterprise_project_id"]
		if epsID == "" {
			return policyID + "/" + rs.Primary.ID, nil
		}
		return policyID + "/" + rs.Primary.ID + "/" + epsID, nil
	}
}
