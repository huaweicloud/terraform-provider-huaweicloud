package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPolicyResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	wafClient, err := cfg.WafV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating WAF client: %s", err)
	}
	return policies.GetWithEpsID(wafClient, state.Primary.ID, state.Primary.Attributes["enterprise_project_id"]).Extract()
}

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccPolicy_basic(t *testing.T) {
	var obj interface{}

	randName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_waf_policy.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getPolicyResourceFunc,
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
				Config: testAccWafPolicy_basic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "level", "1"),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "log"),
					resource.TestCheckResourceAttr(resourceName, "robot_action", "log"),
					resource.TestCheckResourceAttr(resourceName, "options.0.basic_web_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.general_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_engine", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_scanner", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_script", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_other", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.webshell", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.cc_attack_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.precise_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.blacklist", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.data_masking", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.false_alarm_masking", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.web_tamper_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.geolocation_access_control", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.information_leakage_prevention", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.bot_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.known_attack_source", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.anti_crawler", "false"),
					resource.TestCheckResourceAttr(resourceName, "deep_inspection", "false"),
					resource.TestCheckResourceAttr(resourceName, "header_inspection", "false"),
					resource.TestCheckResourceAttr(resourceName, "shiro_decryption_check", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "bind_hosts.#"),
				),
			},
			{
				Config: testAccWafPolicy_update1(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"_updated"),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "true"),
					resource.TestCheckResourceAttr(resourceName, "protection_mode", "block"),
					resource.TestCheckResourceAttr(resourceName, "level", "3"),
					resource.TestCheckResourceAttr(resourceName, "robot_action", "block"),
					resource.TestCheckResourceAttr(resourceName, "options.0.basic_web_protection", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.general_check", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_engine", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_scanner", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_script", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_other", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.webshell", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.cc_attack_protection", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.precise_protection", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.blacklist", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.data_masking", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.false_alarm_masking", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.web_tamper_protection", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.geolocation_access_control", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.information_leakage_prevention", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.bot_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.known_attack_source", "true"),
					resource.TestCheckResourceAttr(resourceName, "options.0.anti_crawler", "true"),
					resource.TestCheckResourceAttr(resourceName, "deep_inspection", "true"),
					resource.TestCheckResourceAttr(resourceName, "header_inspection", "true"),
					resource.TestCheckResourceAttr(resourceName, "shiro_decryption_check", "true"),
				),
			},
			{
				Config: testAccWafPolicy_update2(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "full_detection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.basic_web_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.general_check", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_engine", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_scanner", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_script", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.crawler_other", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.webshell", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.cc_attack_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.precise_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.blacklist", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.data_masking", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.false_alarm_masking", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.web_tamper_protection", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.geolocation_access_control", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.information_leakage_prevention", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.bot_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.known_attack_source", "false"),
					resource.TestCheckResourceAttr(resourceName, "options.0.anti_crawler", "false"),
					resource.TestCheckResourceAttr(resourceName, "deep_inspection", "false"),
					resource.TestCheckResourceAttr(resourceName, "header_inspection", "true"),
					resource.TestCheckResourceAttr(resourceName, "shiro_decryption_check", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testWAFResourceImportState(resourceName),
			},
		},
	})
}

func testAccWafPolicy_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                  = "%[1]s"
  level                 = 1
  enterprise_project_id = "%[2]s"
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafPolicy_update1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                   = "%[1]s_updated"
  full_detection         = true
  protection_mode        = "block"
  level                  = 3
  robot_action           = "block"
  deep_inspection        = true
  header_inspection      = true
  shiro_decryption_check = true
  enterprise_project_id  = "%[2]s"

  options {
    anti_crawler                   = true
    basic_web_protection           = true
    blacklist                      = true
    bot_enable                     = true
    cc_attack_protection           = true
    crawler_engine                 = true
    crawler_other                  = true
    crawler_scanner                = true
    false_alarm_masking            = true
    general_check                  = true
    geolocation_access_control     = true
    information_leakage_prevention = true
    known_attack_source            = true
    precise_protection             = true
    web_tamper_protection          = true
    webshell                       = true
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testAccWafPolicy_update2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_policy" "test" {
  name                   = "%[1]s_updated"
  full_detection         = false
  protection_mode        = "block"
  level                  = 3
  robot_action           = "block"
  deep_inspection        = false
  header_inspection      = true
  shiro_decryption_check = false
  enterprise_project_id  = "%[2]s"

  options {
    
  }
}
`, name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
