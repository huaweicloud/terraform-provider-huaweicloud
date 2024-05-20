package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRulesKnownAttackSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_rules_known_attack_source.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_waf_rules_known_attack_source.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byType   = "data.huaweicloud_waf_rules_known_attack_source.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRulesKnownAttackSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.block_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.block_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.created_at"),

					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRulesKnownAttackSource_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_rules_known_attack_source" "test" {
  depends_on = [
    huaweicloud_waf_rule_known_attack_source.test
  ]

  policy_id = huaweicloud_waf_policy.policy_1.id
}

locals {
  rule_id = data.huaweicloud_waf_rules_known_attack_source.test.rules[0].id
}

data "huaweicloud_waf_rules_known_attack_source" "filter_by_rule_id" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  rule_id   = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_waf_rules_known_attack_source.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  block_type = data.huaweicloud_waf_rules_known_attack_source.test.rules[0].block_type
}

data "huaweicloud_waf_rules_known_attack_source" "filter_by_type" {
  policy_id  = huaweicloud_waf_policy.policy_1.id
  block_type = local.block_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_waf_rules_known_attack_source.filter_by_type.rules[*].block_type : v == local.block_type
  ]
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`, testRuleKnownAttack_basic(name))
}
