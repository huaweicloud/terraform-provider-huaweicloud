package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRulesDataMasking_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_rules_data_masking.test"
		rName          = acceptance.RandomAccResourceName()
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleId   = "data.huaweicloud_waf_rules_data_masking.filter_by_rule_id"
		dcByRuleId = acceptance.InitDataSourceCheck(byRuleId)

		byStatus   = "data.huaweicloud_waf_rules_data_masking.filter_by_status"
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
				Config: testDataSourceRulesDataMasking_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.path"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.field"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.subfield"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.0.created_at"),

					dcByRuleId.CheckResourceExists(),
					resource.TestCheckOutput("rule_id_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRulesDataMasking_base(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_waf_rule_data_masking" "test" {
  policy_id   = huaweicloud_waf_policy.policy_1.id
  path        = "/login"
  field       = "params"
  subfield    = "password"
  description = "test description"
  status      = 0
}
`, testAccWafPolicyV1_basic(name))
}

func testDataSourceRulesDataMasking_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_rules_data_masking" "test" {
  depends_on = [
    huaweicloud_waf_rule_data_masking.test
  ]

  policy_id = huaweicloud_waf_policy.policy_1.id
}

locals {
  rule_id = data.huaweicloud_waf_rules_data_masking.test.rules[0].id
}

data "huaweicloud_waf_rules_data_masking" "filter_by_rule_id" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  rule_id   = local.rule_id
}

locals {
  rule_id_filter_result = [
    for v in data.huaweicloud_waf_rules_data_masking.filter_by_rule_id.rules[*].id : v == local.rule_id
  ]
}

output "rule_id_filter_is_useful" {
  value = alltrue(local.rule_id_filter_result) && length(local.rule_id_filter_result) > 0
}

locals {
  status = data.huaweicloud_waf_rules_data_masking.test.rules[0].status
}

data "huaweicloud_waf_rules_data_masking" "filter_by_status" {
  policy_id = huaweicloud_waf_policy.policy_1.id
  status    = local.status
}

locals {
  status_filter_result = [ 
    for v in data.huaweicloud_waf_rules_data_masking.filter_by_status.rules[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, testDataSourceRulesDataMasking_base(name))
}
