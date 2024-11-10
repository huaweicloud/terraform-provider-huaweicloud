package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
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

		byStatus   = "data.huaweicloud_waf_rules_anti_crawler.filter_by_status"
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
  policy_id             = huaweicloud_waf_policy.test.id
  enterprise_project_id = "%[2]s"
  protection_mode       = "anticrawler_specific_url"

  depends_on = [
    huaweicloud_waf_rule_anti_crawler.test
  ]
}

locals {
  rule_id = data.huaweicloud_waf_rules_anti_crawler.test.rules[0].id
}

data "huaweicloud_waf_rules_anti_crawler" "filter_by_rule_id" {
  policy_id             = huaweicloud_waf_policy.test.id
  rule_id               = local.rule_id
  enterprise_project_id = "%[2]s"
  protection_mode       = "anticrawler_specific_url"
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
  policy_id             = huaweicloud_waf_policy.test.id
  name                  = local.name
  enterprise_project_id = "%[2]s"
  protection_mode       = "anticrawler_specific_url"
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
  policy_id             = huaweicloud_waf_policy.test.id
  name                  = "not_found"
  enterprise_project_id = "%[2]s"
  protection_mode       = "anticrawler_specific_url"
}

output "is_not_found" {
  value = length(data.huaweicloud_waf_rules_anti_crawler.not_found.rules) == 0
}

locals {
  status = data.huaweicloud_waf_rules_anti_crawler.test.rules[0].status
}

data "huaweicloud_waf_rules_anti_crawler" "filter_by_status" {
  policy_id             = huaweicloud_waf_policy.test.id
  status                = local.status
  enterprise_project_id = "%[2]s"
  protection_mode       = "anticrawler_specific_url"
}

locals {
  status_filter_result = [ 
    for v in data.huaweicloud_waf_rules_anti_crawler.filter_by_status.rules[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testDataSourceRuleAntiCrawler_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
