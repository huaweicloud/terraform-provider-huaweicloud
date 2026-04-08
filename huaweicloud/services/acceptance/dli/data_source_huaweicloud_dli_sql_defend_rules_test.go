package dli

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSqlDefendRules_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dli_sql_defend_rules.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_dli_sql_defend_rules.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_dli_sql_defend_rules.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		byQueueName   = "data.huaweicloud_dli_sql_defend_rules.filter_by_queue_name"
		dcByQueueName = acceptance.InitDataSourceCheck(byQueueName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSqlDefendRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "rules.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "rules.0.name"),
					resource.TestCheckResourceAttrSet(all, "rules.0.uuid"),
					resource.TestCheckResourceAttrSet(all, "rules.0.id"),
					resource.TestCheckResourceAttrSet(all, "rules.0.category"),
					resource.TestCheckResourceAttrSet(all, "rules.0.description"),
					resource.TestCheckResourceAttrSet(all, "rules.0.project_id"),
					resource.TestCheckResourceAttrSet(all, "rules.0.sys_desc"),
					resource.TestMatchResourceAttr(all, "rules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(all, "rules.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),

					// Filter by name
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					// Filter by not found name
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),

					// Filter by queue_name
					dcByQueueName.CheckResourceExists(),
					resource.TestCheckOutput("is_queue_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSqlDefendRules_basic() string {
	return `
data "huaweicloud_dli_sql_defend_rules" "all" {}

# Filter by name
locals {
  first_rule_name = data.huaweicloud_dli_sql_defend_rules.all.rules[0].name
}

data "huaweicloud_dli_sql_defend_rules" "filter_by_name" {
  rule_name = local.first_rule_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dli_sql_defend_rules.filter_by_name.rules[*].name : v == local.first_rule_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by not found name
data "huaweicloud_dli_sql_defend_rules" "filter_by_not_found_name" {
  rule_name = "non-exist-rule-name"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_dli_sql_defend_rules.filter_by_not_found_name.rules) == 0
}

# Filter by queue_name
locals {
  first_queue_name = data.huaweicloud_dli_sql_defend_rules.all.rules[0].queue_names[0]
}

data "huaweicloud_dli_sql_defend_rules" "filter_by_queue_name" {
  queue_name = local.first_queue_name
}

output "is_queue_name_filter_useful" {
  value = length(data.huaweicloud_dli_sql_defend_rules.filter_by_queue_name.rules) > 0
}
`
}
