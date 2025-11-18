package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataRuleEngineRules_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_cdn_rule_engine_rules.test"
		dc    = acceptance.InitDataSourceCheck(rName)
		name  = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataRuleEngineRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "domain_name", acceptance.HW_CDN_DOMAIN_NAME),
					resource.TestCheckResourceAttrSet(rName, "rules.#"),
					resource.TestCheckOutput("is_rule_id_set_and_valid", "true"),
					resource.TestCheckOutput("is_rule_name_set_and_valid", "true"),
					resource.TestCheckOutput("is_rule_status_set_and_valid", "true"),
					resource.TestCheckOutput("is_rule_priority_set_and_valid", "true"),
					resource.TestCheckOutput("is_rule_conditions_set_and_valid", "true"),
					resource.TestCheckOutput("is_rule_actions_set_and_valid", "true"),
				),
			},
		},
	})
}

func testAccDataRuleEngineRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cdn_rule_engine_rules" "test" {
  depends_on = [huaweicloud_cdn_rule_engine_rule.test]

  domain_name = "%[2]s"
}

locals {
  rule_name    = huaweicloud_cdn_rule_engine_rule.test.name
  queried_rule = [for rule in data.huaweicloud_cdn_rule_engine_rules.test.rules : rule if rule.name == local.rule_name]
}

output "is_rule_id_set_and_valid" {
  value = length(local.queried_rule) >= 1 && local.queried_rule[0].id == huaweicloud_cdn_rule_engine_rule.test.id
}

output "is_rule_name_set_and_valid" {
  value = length(local.queried_rule) >= 1 && local.queried_rule[0].name == huaweicloud_cdn_rule_engine_rule.test.name
}

output "is_rule_status_set_and_valid" {
  value = length(local.queried_rule) >= 1 && local.queried_rule[0].status == huaweicloud_cdn_rule_engine_rule.test.status
}

output "is_rule_priority_set_and_valid" {
  value = length(local.queried_rule) >= 1 && local.queried_rule[0].priority == huaweicloud_cdn_rule_engine_rule.test.priority
}

output "is_rule_conditions_set_and_valid" {
  value = length(local.queried_rule) >= 1 && local.queried_rule[0].conditions != ""
}

output "is_rule_actions_set_and_valid" {
  value = length(local.queried_rule) >= 1 && length(local.queried_rule[0].actions) == length(huaweicloud_cdn_rule_engine_rule.test.actions)
}
`, testAccRuleEngineRule_basic_step1(name), acceptance.HW_CDN_DOMAIN_NAME)
}
