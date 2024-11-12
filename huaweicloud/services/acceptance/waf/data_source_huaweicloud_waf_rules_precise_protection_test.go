package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccDataSourceRulesPreciseProtection_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_rules_precise_protection.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_waf_rules_precise_protection.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byName           = "data.huaweicloud_waf_rules_precise_protection.filter_by_rule_id"
		dcByName         = acceptance.InitDataSourceCheck(byName)
		byNameNotFound   = "data.huaweicloud_waf_rules_precise_protection.not_found"
		dcByNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)

		byStatus   = "data.huaweicloud_waf_rules_precise_protection.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
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
				Config: testDataSourceRulesPreciseProtection_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.priority"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.conditions.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.created_at"),

					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRulesPreciseProtection_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_rules_precise_protection" "test" {
  policy_id             = huaweicloud_waf_policy.test.id
  enterprise_project_id = "%[2]s"

  depends_on = [
    huaweicloud_waf_rule_precise_protection.test
  ]
}

locals {
  rule_id = data.huaweicloud_waf_rules_precise_protection.test.rules[0].id
}

data "huaweicloud_waf_rules_precise_protection" "filter_by_rule_id" {
  policy_id             = huaweicloud_waf_policy.test.id
  rule_id               = local.rule_id
  enterprise_project_id = "%[2]s"
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_waf_rules_precise_protection.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  name = data.huaweicloud_waf_rules_precise_protection.test.rules[0].name
}

data "huaweicloud_waf_rules_precise_protection" "filter_by_name" {
  policy_id             = huaweicloud_waf_policy.test.id
  name                  = local.name
  enterprise_project_id = "%[2]s"
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_rules_precise_protection.filter_by_name.rules[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

data "huaweicloud_waf_rules_precise_protection" "not_found" {
  policy_id             = huaweicloud_waf_policy.test.id
  name                  = "not_found"
  enterprise_project_id = "%[2]s"
}
  
output "is_not_found" {
  value = length(data.huaweicloud_waf_rules_precise_protection.not_found.rules) == 0
}

locals {
  status = data.huaweicloud_waf_rules_precise_protection.test.rules[0].status
}

data "huaweicloud_waf_rules_precise_protection" "filter_by_status" {
  policy_id             = huaweicloud_waf_policy.test.id
  status                = local.status
  enterprise_project_id = "%[2]s"
}

locals {
  status_filter_result = [ 
    for v in data.huaweicloud_waf_rules_precise_protection.filter_by_status.rules[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testDataSourceRulePreciseProtection_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
