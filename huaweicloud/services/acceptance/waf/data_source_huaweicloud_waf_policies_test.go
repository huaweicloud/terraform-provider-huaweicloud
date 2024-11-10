package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccDataSourcePolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dataSourceName = "data.huaweicloud_waf_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byName   = "data.huaweicloud_waf_policies.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNonExist   = "data.huaweicloud_waf_policies.non_exist_filter"
		dcByNonExist = acceptance.InitDataSourceCheck(byNonExist)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWafPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
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

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByNonExist.CheckResourceExists(),
					resource.TestCheckOutput("non_exist_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccWafPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_policies" "test" {
  enterprise_project_id = "%[2]s"
}

# Filter by name
locals {
  name = data.huaweicloud_waf_policies.test.policies.0.name
}

data "huaweicloud_waf_policies" "name_filter" {
  name                  = local.name
  enterprise_project_id = "%[2]s"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_policies.name_filter.policies[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)  
}

# Filter by non-exist name
data "huaweicloud_waf_policies" "non_exist_filter" {
  name                  = "non-exist"
  enterprise_project_id = "%[2]s"
}

locals {
  non_exist_filter_result = [
    for v in data.huaweicloud_waf_policies.non_exist_filter.policies[*].name : v == local.name
  ]
}

output "non_exist_filter_is_useful" {
  value = length(local.non_exist_filter_result) == 0
}
`, testAccWafPolicy_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
