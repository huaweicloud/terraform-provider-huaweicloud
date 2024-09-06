package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRulesAntiCrawler_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_rules_anti_crawler.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_waf_rules_anti_crawler.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byName           = "data.huaweicloud_waf_rules_anti_crawler.filter_by_name"
		dcByName         = acceptance.InitDataSourceCheck(byName)
		byNameNotFound   = "data.huaweicloud_waf_rules_anti_crawler.not_found"
		dcByNameNotFound = acceptance.InitDataSourceCheck(byNameNotFound)

		byProtectionMode   = "data.huaweicloud_waf_rules_anti_crawler.filter_by_protection_mode"
		dcByProtectionMode = acceptance.InitDataSourceCheck(byProtectionMode)

		byStatus   = "data.huaweicloud_waf_rules_anti_crawler.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRulesAntiCrawler_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.protection_mode"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.priority"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.conditions.0.field"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.conditions.0.logic"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.created_at"),

					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found", "true"),

					dcByProtectionMode.CheckResourceExists(),
					resource.TestCheckOutput("protection_mode_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRulesAntiCrawler_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_rules_anti_crawler" "test" {
  depends_on = [
    huaweicloud_waf_rule_anti_crawler.test
  ]

  policy_id = huaweicloud_waf_policy.test.id
}

locals {
  rule_id = data.huaweicloud_waf_rules_anti_crawler.test.rules[0].id
}

data "huaweicloud_waf_rules_anti_crawler" "filter_by_rule_id" {
  policy_id = huaweicloud_waf_policy.test.id
  rule_id   = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_waf_rules_anti_crawler.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  name = data.huaweicloud_waf_rules_anti_crawler.test.rules[0].name
}

data "huaweicloud_waf_rules_anti_crawler" "filter_by_name" {
  policy_id = huaweicloud_waf_policy.test.id
  name      = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_rules_anti_crawler.filter_by_name.rules[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

data "huaweicloud_waf_rules_anti_crawler" "not_found" {
  policy_id = huaweicloud_waf_policy.test.id
  name      = "not_found"
}

output "is_not_found" {
  value = length(data.huaweicloud_waf_rules_anti_crawler.not_found.rules) == 0
}

locals {
  protection_mode = data.huaweicloud_waf_rules_anti_crawler.test.rules[0].protection_mode
}

data "huaweicloud_waf_rules_anti_crawler" "filter_by_protection_mode" {
  policy_id       = huaweicloud_waf_policy.test.id
  protection_mode = local.protection_mode
}

locals {
  protection_mode_filter_result = [
    for v in data.huaweicloud_waf_rules_anti_crawler.filter_by_protection_mode.rules[*].protection_mode : 
    v == local.protection_mode
  ]
}

output "protection_mode_filter_is_useful" {
  value = alltrue(local.protection_mode_filter_result) && length(local.protection_mode_filter_result) > 0
}

locals {
  status = data.huaweicloud_waf_rules_anti_crawler.test.rules[0].status
}

data "huaweicloud_waf_rules_anti_crawler" "filter_by_status" {
  policy_id = huaweicloud_waf_policy.test.id
  status    = local.status
}

locals {
  status_filter_result = [ 
    for v in data.huaweicloud_waf_rules_anti_crawler.filter_by_status.rules[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testRuleAntiCrawler_excepProtectionMode(name))
}
