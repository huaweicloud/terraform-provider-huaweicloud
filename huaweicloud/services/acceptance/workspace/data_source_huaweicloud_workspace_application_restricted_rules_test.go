package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplicationRestrictedRules_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_application_restricted_rules.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByName   = "data.huaweicloud_workspace_application_restricted_rules.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationRestrictedRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "rules.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "rules.0.id"),
					resource.TestCheckResourceAttrSet(all, "rules.0.name"),
					resource.TestCheckResourceAttrSet(all, "rules.0.rule_source"),
					resource.TestCheckResourceAttrSet(all, "rules.0.create_time"),
					resource.TestCheckResourceAttrSet(all, "rules.0.update_time"),
					resource.TestCheckResourceAttrSet(all, "rules.0.rule.0.scope"),
					// Filter by 'name' parameter.
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataApplicationRestrictedRules_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_application_rule" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"

  detail {
    scope = "PRODUCT"

    product_rule {
      identify_condition = "process"
      publisher          = "Microsoft Corporation"
      product_name       = "Microsoft Office"
      process_name       = "WINWORD.EXE"
      support_os         = "Windows"
      version            = "1.0"
      product_version    = "2019"
    }
  }
}

resource "huaweicloud_workspace_application_rule_restriction" "test" {
  rule_ids = [huaweicloud_workspace_application_rule.test.id]
}
`, name)
}

func testAccDataApplicationRestrictedRules_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_application_restricted_rules" "all" {
  depends_on = [
    huaweicloud_workspace_application_rule_restriction.test
  ]
}

# Filter by 'name' parameter.
locals {
  rule_name = huaweicloud_workspace_application_rule.test.name
}

data "huaweicloud_workspace_application_restricted_rules" "filter_by_name" {
  name = local.rule_name

  depends_on = [
    huaweicloud_workspace_application_rule_restriction.test
  ]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_workspace_application_restricted_rules.filter_by_name.rules[*].name : v == local.rule_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`, testAccDataApplicationRestrictedRules_base(name), name)
}
