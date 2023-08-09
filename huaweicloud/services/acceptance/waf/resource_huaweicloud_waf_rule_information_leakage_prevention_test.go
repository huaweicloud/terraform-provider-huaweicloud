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

func getRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/antileakage/{rule_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", state.Primary.Attributes["policy_id"])
	getPath = strings.ReplaceAll(getPath, "{rule_id}", state.Primary.ID)

	queryParam := ""
	if epsID := state.Primary.Attributes["enterprise_project_id"]; epsID != "" {
		queryParam = fmt.Sprintf("?enterprise_project_id=%s", epsID)
	}
	getPath += queryParam

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving WAF information leakage prevention rule: %s", err)
	}
	return utils.FlattenResponse(getResp)
}

func TestAccRuleLeakagePrevention_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_information_leakage_prevention.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleResourceFunc,
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
				Config: testRuleLeakagePrevention_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "path", "/test/path"),
					resource.TestCheckResourceAttr(rName, "type", "sensitive"),
					resource.TestCheckResourceAttr(rName, "protective_action", "block"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "contents.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testRuleLeakagePrevention_basic_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "path", "/test/val*"),
					resource.TestCheckResourceAttr(rName, "type", "sensitive"),
					resource.TestCheckResourceAttr(rName, "contents.0", "phone"),
					resource.TestCheckResourceAttr(rName, "protective_action", "log"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
				),
			},
			{
				Config: testRuleLeakagePrevention_basic_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "code"),
					resource.TestCheckResourceAttr(rName, "contents.#", "3"),
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

func TestAccRuleLeakagePrevention_typeCode(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_information_leakage_prevention.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleResourceFunc,
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
				Config: testRuleLeakagePrevention_typeCode(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "path", "/test/val*"),
					resource.TestCheckResourceAttr(rName, "type", "code"),
					resource.TestCheckResourceAttr(rName, "protective_action", "log"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "contents.#", "3"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testRuleLeakagePrevention_typeCode_update1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "path", "/test/path"),
					resource.TestCheckResourceAttr(rName, "type", "code"),
					resource.TestCheckResourceAttr(rName, "contents.0", "507"),
					resource.TestCheckResourceAttr(rName, "protective_action", "block"),
					resource.TestCheckResourceAttr(rName, "description", "test description update"),
				),
			},
			{
				Config: testRuleLeakagePrevention_typeCode_update2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "type", "sensitive"),
					resource.TestCheckResourceAttr(rName, "contents.#", "3"),
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

func TestAccRuleLeakagePrevention_withEpsID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_information_leakage_prevention.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleResourceFunc,
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
				Config: testRuleLeakagePrevention_withEpsID(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "path", "/test/val*"),
					resource.TestCheckResourceAttr(rName, "type", "code"),
					resource.TestCheckResourceAttr(rName, "protective_action", "log"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "contents.#", "3"),
					resource.TestCheckResourceAttrSet(rName, "status"),
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

func testRuleLeakagePrevention_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id         = huaweicloud_waf_policy.policy_1.id
  path              = "/test/path"
  type              = "sensitive"
  contents          = ["phone", "id_card"]
  protective_action = "block"
  description       = "test description"
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleLeakagePrevention_basic_update1(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id         = huaweicloud_waf_policy.policy_1.id
  path              = "/test/val*"
  type              = "sensitive"
  contents          = ["phone"]
  protective_action = "log"
  description       = "test description update"
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleLeakagePrevention_basic_update2(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id         = huaweicloud_waf_policy.policy_1.id
  path              = "/test/val*"
  type              = "code"
  contents          = ["401", "405", "503"]
  protective_action = "log"
  description       = "test description update"
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleLeakagePrevention_typeCode(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id         = huaweicloud_waf_policy.policy_1.id
  path              = "/test/val*"
  type              = "code"
  contents          = ["401", "405", "503"]
  protective_action = "log"
  description       = "test description"
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleLeakagePrevention_typeCode_update1(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id         = huaweicloud_waf_policy.policy_1.id
  path              = "/test/path"
  type              = "code"
  contents          = ["507"]
  protective_action = "block"
  description       = "test description update"
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleLeakagePrevention_typeCode_update2(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id         = huaweicloud_waf_policy.policy_1.id
  path              = "/test/path"
  type              = "sensitive"
  contents          = ["phone", "id_card", "email"]
  protective_action = "block"
  description       = "test description update"
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleLeakagePrevention_withEpsID(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  path                  = "/test/val*"
  type                  = "code"
  contents              = ["401", "405", "503"]
  protective_action     = "log"
  description           = "test description"
  enterprise_project_id = "%s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
