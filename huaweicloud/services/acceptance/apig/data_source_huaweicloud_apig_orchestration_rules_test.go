package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrchestrationRules_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_apig_orchestration_rules.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byId   = "data.huaweicloud_apig_orchestration_rules.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_orchestration_rules.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byNotFoundName   = "data.huaweicloud_apig_orchestration_rules.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrchestrationRules_base(name),
			},
			{
				// Update the orchestration rule's name to make sure the update time generated successfully.
				Config: testAccDataSourceOrchestrationRules_base(updateName),
			},
			{
				Config: testAccDataSourceOrchestrationRules_basic(updateName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "rules.#", regexp.MustCompile(`[1-9]\d*`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckResourceAttr(byId, "rules.#", "1"),
					resource.TestCheckResourceAttrPair(byId, "rules.0.id", "huaweicloud_apig_orchestration_rule.test", "id"),
					resource.TestCheckResourceAttr(byId, "rules.0.name", updateName),
					resource.TestCheckResourceAttr(byId, "rules.0.strategy", "hash"),
					resource.TestCheckResourceAttr(byId, "rules.0.is_preprocessing", "false"),
					resource.TestCheckResourceAttr(byId, "rules.0.mapped_param",
						"{\"mapped_param_location\":\"header\",\"mapped_param_name\":\"standard-param\",\"mapped_param_type\":\"string\"}"),
					resource.TestMatchResourceAttr(byId, "rules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(byId, "rules.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceOrchestrationRules_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_apig_orchestration_rule" "test" {
  instance_id = "%[1]s"
  name        = "%[2]s"
  strategy    = "hash"

  mapped_param = jsonencode({
    "mapped_param_name": "standard-param",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}

resource "huaweicloud_apig_orchestration_rule" "with_suffix_index" {
  count = 2

  instance_id = "%[1]s"
  name        = format("reverse_%%s_%%d", strrev("%[2]s"), count.index)
  strategy    = "hash"

  mapped_param = jsonencode({
    "mapped_param_name": "reversed-param-with-suffix-index-${count.index}",
    "mapped_param_type": "string",
    "mapped_param_location": "header"
  })
}
`, acceptance.HW_APIG_DEDICATED_INSTANCE_ID, name)
}

func testAccDataSourceOrchestrationRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_orchestration_rules" "test" {
  depends_on = [
    huaweicloud_apig_orchestration_rule.test,
    huaweicloud_apig_orchestration_rule.with_suffix_index,
  ]

  instance_id = "%[2]s"
}

# Filter by ID
locals {
  rule_id = huaweicloud_apig_orchestration_rule.test.id
}

data "huaweicloud_apig_orchestration_rules" "filter_by_id" {
  # Need to be executed after all resources are created (excluding implicit dependencies).
  depends_on = [
    huaweicloud_apig_orchestration_rule.with_suffix_index,
  ]

  instance_id = "%[2]s"
  rule_id     = local.rule_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_orchestration_rules.filter_by_id.rules[*].id : v == local.rule_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  rule_name = huaweicloud_apig_orchestration_rule.test.name
}

data "huaweicloud_apig_orchestration_rules" "filter_by_name" {
  # Since a specified name is used, there is no dependency relationship with resource attachment, and the dependency
  # needs to be manually set.
  depends_on = [
    huaweicloud_apig_orchestration_rule.test,
    huaweicloud_apig_orchestration_rule.with_suffix_index,
  ]

  instance_id = "%[2]s"
  name        = local.rule_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_orchestration_rules.filter_by_name.rules[*].name : strcontains(v, local.rule_name)
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}


# Filter by name and the name is not exist
data "huaweicloud_apig_orchestration_rules" "filter_by_not_found_name" {
  # Since there is not any parameter used and no dependency relationship with resource attachment, and the dependency
  # needs to be manually set.
  depends_on = [
    huaweicloud_apig_orchestration_rule.test,
    huaweicloud_apig_orchestration_rule.with_suffix_index,
  ]

  instance_id = "%[2]s"
  name        = "not_found_name"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_apig_orchestration_rules.filter_by_not_found_name.rules) == 0
}
`, testAccDataSourceOrchestrationRules_base(name), acceptance.HW_APIG_DEDICATED_INSTANCE_ID)
}
