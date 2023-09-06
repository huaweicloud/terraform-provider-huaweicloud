package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccDataSourceWafPoliciesV1_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_policies.policies_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPoliciesV1_conf(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPoliciesID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "policies.0.name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.full_detection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.protection_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.robot_action"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.level"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.deep_inspection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.header_inspection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.shiro_decryption_check"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.basic_web_protection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.general_check"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.crawler_engine"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.crawler_scanner"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.crawler_script"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.crawler_other"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.webshell"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.cc_attack_protection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.precise_protection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.blacklist"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.data_masking"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.false_alarm_masking"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.web_tamper_protection"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.geolocation_access_control"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.information_leakage_prevention"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.bot_enable"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.known_attack_source"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.anti_crawler"),
				),
			},
		},
	})
}

func TestAccDataSourceWafPoliciesV1_withEpsID(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	dataSourceName := "data.huaweicloud_waf_policies.policies_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPoliciesV1_conf_epsID(name, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWafPoliciesID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
					resource.TestCheckResourceAttr(dataSourceName, "name", name),
					resource.TestCheckResourceAttr(dataSourceName, "policies.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "policies.0.name", name),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.options.0.blacklist"),
				),
			},
		},
	})
}

func testAccCheckWafPoliciesID(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmtp.Errorf("Can't find WAF policies data source: %s.", r)
		}
		if rs.Primary.ID == "" {
			return fmtp.Errorf("The WAF policies data source ID does not set.")
		}
		return nil
	}
}

func testAccWafPoliciesV1_conf(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_policies" "policies_1" {
  name = huaweicloud_waf_policy.policy_1.name
}
`, testAccWafPolicyV1_basic(name))
}

func testAccWafPoliciesV1_conf_epsID(name, epsID string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_policies" "policies_1" {
  name                  = huaweicloud_waf_policy.policy_1.name
  enterprise_project_id = "%s"
}
`, testAccWafPolicyV1_basic_withEpsID(name, epsID), epsID)
}
