package dsc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscTemplateRules_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dsc_template_rules.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuleName   = "data.huaweicloud_dsc_template_rules.filter_by_rule_name"
		dcByRuleName = acceptance.InitDataSourceCheck(byRuleName)

		byClassificationId   = "data.huaweicloud_dsc_template_rules.filter_by_classification_id"
		dcByClassificationId = acceptance.InitDataSourceCheck(byClassificationId)

		bySecurityLevelId   = "data.huaweicloud_dsc_template_rules.filter_by_security_level_id"
		dcBySecurityLevelId = acceptance.InitDataSourceCheck(bySecurityLevelId)

		byIsUsed   = "data.huaweicloud_dsc_template_rules.filter_by_is_used"
		dcByIsUsed = acceptance.InitDataSourceCheck(byIsUsed)

		byRuleNameNotFound   = "data.huaweicloud_dsc_template_rules.not_found"
		dcByRuleNameNotFound = acceptance.InitDataSourceCheck(byRuleNameNotFound)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the scan template in advance and configure the template ID into
			// the environment variable.
			acceptance.TestAccPreCheckDSCScanTemplateID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscTemplateRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "template_rules_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "template_rules_list.0.rule_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "template_rules_list.0.rule_name"),

					dcByRuleName.CheckResourceExists(),
					resource.TestCheckOutput("rule_name_filter_is_useful", "true"),

					dcByClassificationId.CheckResourceExists(),
					resource.TestCheckOutput("classification_id_filter_is_useful", "true"),

					dcBySecurityLevelId.CheckResourceExists(),
					resource.TestCheckOutput("security_level_id_filter_is_useful", "true"),

					dcByIsUsed.CheckResourceExists(),
					resource.TestCheckOutput("is_used_filter_is_useful", "true"),

					dcByRuleNameNotFound.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDscTemplateRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dsc_template_rules" "test" {
  template_id = "%[1]s"
}

# Filter by rule_name
locals {
  rule_name = data.huaweicloud_dsc_template_rules.test.template_rules_list[0].rule_name
}

data "huaweicloud_dsc_template_rules" "filter_by_rule_name" {
  template_id = "%[1]s"
  rule_name   = local.rule_name
}

locals {
  rule_name_filter_result = [
    for v in data.huaweicloud_dsc_template_rules.filter_by_rule_name.template_rules_list[*].rule_name : v == local.rule_name
  ]
}

output "rule_name_filter_is_useful" {
  value = alltrue(local.rule_name_filter_result) && length(local.rule_name_filter_result) > 0
}

# Filter by classification_id
locals {
  classification_id = data.huaweicloud_dsc_template_rules.test.template_rules_list[0].classification_id
}

data "huaweicloud_dsc_template_rules" "filter_by_classification_id" {
  template_id        = "%[1]s"
  classification_ids = [local.classification_id]
}

locals {
  classification_id_filter_result = [
    for v in data.huaweicloud_dsc_template_rules.filter_by_classification_id.template_rules_list[*].classification_id :
    v == local.classification_id
  ]
}

output "classification_id_filter_is_useful" {
  value = alltrue(local.classification_id_filter_result) && length(local.classification_id_filter_result) > 0
}

# Filter by security_level_id
locals {
  security_level_id = data.huaweicloud_dsc_template_rules.test.template_rules_list[0].security_level_id
}

data "huaweicloud_dsc_template_rules" "filter_by_security_level_id" {
  template_id        = "%[1]s"
  security_level_ids = [local.security_level_id]
}

locals {
  security_level_id_filter_result = [
    for v in data.huaweicloud_dsc_template_rules.filter_by_security_level_id.template_rules_list[*].security_level_id :
    v == local.security_level_id
  ]
}

output "security_level_id_filter_is_useful" {
  value = alltrue(local.security_level_id_filter_result) && length(local.security_level_id_filter_result) > 0
}

# Filter by is_used
locals {
  is_used = tostring(data.huaweicloud_dsc_template_rules.test.template_rules_list[0].is_used)
}

data "huaweicloud_dsc_template_rules" "filter_by_is_used" {
  template_id = "%[1]s"
  is_used     = local.is_used
}

locals {
  is_used_filter_result = [
    for v in data.huaweicloud_dsc_template_rules.filter_by_is_used.template_rules_list[*].is_used :
    tostring(v) == local.is_used
  ]
}

output "is_used_filter_is_useful" {
  value = alltrue(local.is_used_filter_result) && length(local.is_used_filter_result) > 0
}

data "huaweicloud_dsc_template_rules" "not_found" {
  template_id = "%[1]s"
  rule_name   = "not_found"
}

output "is_not_found" {
  value = length(data.huaweicloud_dsc_template_rules.not_found.template_rules_list) == 0
}
`, acceptance.HW_DSC_SCAN_TEMPLATE_ID)
}
