package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/waf"
)

func getPolicyV2ResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	wafClient, err := cfg.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}
	return waf.GetWafPolicyV2(wafClient, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"])
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccPolicyV2_basic(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_policy_v2.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPolicyV2ResourceFunc,
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
				Config: testAccWafPolicyV2_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "false"),
					resource.TestCheckResourceAttr(resourceName, "level", "2"),
					resource.TestCheckResourceAttr(resourceName, "log_action_replaced", "true"),
					resource.TestCheckResourceAttr(resourceName, "action.0.category", "log"),
					resource.TestCheckResourceAttr(resourceName, "robot_action.0.category", "log"),
					resource.TestCheckResourceAttr(resourceName, "options.0.antileakage", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.antitamper", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.bot_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.cc", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.common", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_engine", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_other", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_scanner", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_script", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.custom", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.geoip", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.ignore", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.modulex_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.privacy", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.webattack", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.webshell", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.whiteblackip", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "extend_attribute.%"),
					resource.TestCheckResourceAttrSet(resourceName, "timestamp"),
				),
			},
			{
				Config: testAccWafPolicyV2_update1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s_update", randName)),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "true"),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "log_action_replaced", "true"),
					resource.TestCheckResourceAttr(resourceName, "action.0.category", "block"),
					resource.TestCheckResourceAttr(resourceName, "robot_action.0.category", "block"),
					resource.TestCheckResourceAttr(resourceName, "options.0.antileakage", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.antitamper", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.bot_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.cc", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.common", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_engine", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_other", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_scanner", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_script", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.custom", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.geoip", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.ignore", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.modulex_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.privacy", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.webattack", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.webshell", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.whiteblackip", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "extend_attribute.%"),
					resource.TestCheckResourceAttrSet(resourceName, "timestamp"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(resourceName),
				ImportStateVerifyIgnore: []string{
					"log_action_replaced",
					"extend",
				},
			},
		},
	})
}

func testAccWafPolicyV2_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy_v2" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "%[2]s"
  log_action_replaced   = "true"
  full_detection        = "false"
  level                 = 2

  action {
    category = "log"
  }

  robot_action {
    category = "log"
  }

  options {
    antileakage     = "false"
    antitamper      = "true"
    bot_enable      = "true"
    cc              = "true"
    crawler_other   = "false"
    crawler_scanner = "true"
    crawler_script  = "false"
    custom          = "true"
    geoip           = "true"
    ignore          = "true"
    modulex_enabled = "false"
    privacy         = "true"
    webattack       = "true"
    webshell        = "false"
    whiteblackip    = "true"
  }

  extend = {
    "extend" = jsonencode(
      {
        deep_decode             = true
        log_action_replaced     = false
        shiro_rememberMe_enable = true
      }
    )
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafPolicyV2_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy_v2" "test" {
  name                  = "%[1]s_update"
  enterprise_project_id = "%[2]s"
  log_action_replaced   = "true"
  full_detection        = "true"
  level                 = 1

  action {
    category = "block"
  }

  robot_action {
    category = "block"
  }

  options {
    antileakage     = "true"
    antitamper      = "true"
    bot_enable      = "false"
    cc              = "true"
    crawler_other   = "true"
    crawler_scanner = "true"
    crawler_script  = "true"
    custom          = "true"
    geoip           = "true"
    ignore          = "false"
    modulex_enabled = "false"
    privacy         = "true"
    webattack       = "false"
    webshell        = "false"
    whiteblackip    = "false"
  }

  extend = {
    "extend" = jsonencode(
      {
        deep_decode             = false
        log_action_replaced     = false
        shiro_rememberMe_enable = false
      }
    )
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
