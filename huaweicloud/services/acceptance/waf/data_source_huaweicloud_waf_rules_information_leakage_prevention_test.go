package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccDataSourceRulesInformationLeakagePrevention_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_rules_information_leakage_prevention.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_waf_rules_information_leakage_prevention.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byStatus   = "data.huaweicloud_waf_rules_information_leakage_prevention.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byType   = "data.huaweicloud_waf_rules_information_leakage_prevention.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
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
				Config: testDataSourceRulesInformationLeakagePrevention_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.path"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.protection_action"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.contents.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.created_at"),

					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRulesInformationLeakagePrevention_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_rules_information_leakage_prevention" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_waf_rule_information_leakage_prevention.test
  ]
}

locals {
  rule_id = data.huaweicloud_waf_rules_information_leakage_prevention.test.rules[0].id
}

data "huaweicloud_waf_rules_information_leakage_prevention" "filter_by_rule_id" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  rule_id               = local.rule_id
  enterprise_project_id = "%[2]s"
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_waf_rules_information_leakage_prevention.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  status = data.huaweicloud_waf_rules_information_leakage_prevention.test.rules[0].status
}

data "huaweicloud_waf_rules_information_leakage_prevention" "filter_by_status" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  status                = local.status
  enterprise_project_id = "%[2]s"
}

locals {
  status_filter_result = [ 
    for v in data.huaweicloud_waf_rules_information_leakage_prevention.filter_by_status.rules[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  type = data.huaweicloud_waf_rules_information_leakage_prevention.test.rules[0].type
}

data "huaweicloud_waf_rules_information_leakage_prevention" "filter_by_type" {
  policy_id             = huaweicloud_waf_policy.policy_1.id
  type                  = local.type
  enterprise_project_id = "%[2]s"
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_waf_rules_information_leakage_prevention.filter_by_type.rules[*].type : v == local.type
  ]
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`, testDataSourceRuleLeakagePrevention_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
