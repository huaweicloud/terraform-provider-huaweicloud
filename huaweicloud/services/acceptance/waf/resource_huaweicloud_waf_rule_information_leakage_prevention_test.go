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
		return nil, fmt.Errorf("error creating WAF client: %s", err)
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

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
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
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRuleLeakagePrevention_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(rName, "policy_id",
						"huaweicloud_waf_policy.test", "id"),
					resource.TestCheckResourceAttr(rName, "path", "/test/path"),
					resource.TestCheckResourceAttr(rName, "type", "sensitive"),
					resource.TestCheckResourceAttr(rName, "protective_action", "block"),
					resource.TestCheckResourceAttr(rName, "description", "test description"),
					resource.TestCheckResourceAttr(rName, "contents.#", "2"),
					resource.TestCheckResourceAttrSet(rName, "status"),
				),
			},
			{
				Config: testDataSourceRuleLeakagePrevention_basic_update1(name),
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
				Config: testDataSourceRuleLeakagePrevention_basic_update2(name),
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

func testDataSourceRuleLeakagePrevention_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  path                  = "/test/path"
  type                  = "sensitive"
  contents              = ["phone", "id_card"]
  protective_action     = "block"
  description           = "test description"
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceRuleLeakagePrevention_basic_update1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  path                  = "/test/val*"
  type                  = "sensitive"
  contents              = ["phone"]
  protective_action     = "log"
  description           = "test description update"
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testDataSourceRuleLeakagePrevention_basic_update2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_waf_rule_information_leakage_prevention" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  path                  = "/test/val*"
  type                  = "code"
  contents              = ["401", "405", "503"]
  protective_action     = "log"
  description           = "test description update"
  enterprise_project_id = "%[2]s"
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
