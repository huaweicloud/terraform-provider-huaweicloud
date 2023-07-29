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

func getRuleKnownAttackResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		region  = acceptance.HW_REGION_NAME
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/punishment/{rule_id}"
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
		return nil, fmt.Errorf("error retrieving WAF known attack source rule: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func TestAccRuleKnownAttack_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_known_attack_source.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleKnownAttackResourceFunc,
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
				Config: testRuleKnownAttack_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "block_type", "long_ip_block"),
					resource.TestCheckResourceAttr(rName, "block_time", "500"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
				),
			},
			{
				Config: testRuleKnownAttack_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "block_time", "600"),
					resource.TestCheckResourceAttr(rName, "description", ""),
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

func TestAccRuleKnownAttack_withEpsID(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_waf_rule_known_attack_source.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getRuleKnownAttackResourceFunc,
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
				Config: testRuleKnownAttack_withEpsID(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.policy_1", "id"),
					resource.TestCheckResourceAttr(rName, "enterprise_project_id",
						acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(rName, "block_type", "long_ip_block"),
					resource.TestCheckResourceAttr(rName, "block_time", "500"),
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

func testRuleKnownAttack_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_known_attack_source" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  block_type  = "long_ip_block"
  block_time  = 500
  description = "test description"
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleKnownAttack_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_known_attack_source" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  block_type  = "long_ip_block"
  block_time  = 600
}
`, testAccWafPolicyV1_basic(name))
}

func testRuleKnownAttack_withEpsID(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_known_attack_source" "test" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  enterprise_project_id = "%s"
  block_type            = "long_ip_block"
  block_time            = 500
  description           = "test description"
}
`, testAccWafPolicyV1_basic_withEpsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
		acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
