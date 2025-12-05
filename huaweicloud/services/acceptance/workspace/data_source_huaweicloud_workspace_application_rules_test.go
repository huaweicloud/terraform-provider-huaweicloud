package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplicationRules_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_application_rules.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterByName   = "data.huaweicloud_workspace_application_rules.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationRules_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "rules.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dcName, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "rules.0.name"),
					resource.TestCheckResourceAttrSet(dcName, "rules.0.rule_source"),
					resource.TestCheckResourceAttrSet(dcName, "rules.0.create_time"),
					resource.TestCheckResourceAttrSet(dcName, "rules.0.update_time"),
					resource.TestCheckResourceAttrSet(dcName, "rules.0.detail.0.scope"),
					// Filter by 'name' parameter.
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataApplicationRules_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_application_rule" "with_product_rule" {
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

resource "huaweicloud_workspace_application_rule" "with_path_rule" {
  name        = "%[1]s_path"
  description = "Created by terraform script for path rule"

  detail {
    scope = "PATH"

    path_rule {
      path = "C:\\Program Files\\Microsoft Office\\root\\Office16\\WINWORD.EXE"
    }
  }
}
`, name)
}

func testAccDataApplicationRules_basic(name string) string {
	// the name filter case need validate the context is contain the filter parameter?
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_application_rules" "test" {
  depends_on = [
    huaweicloud_workspace_application_rule.with_product_rule,
    huaweicloud_workspace_application_rule.with_path_rule
  ]
}

# Filter by 'name' parameter.
data "huaweicloud_workspace_application_rules" "filter_by_name" {
  name = "%[2]s"

  depends_on = [
    huaweicloud_workspace_application_rule.with_product_rule,
    huaweicloud_workspace_application_rule.with_path_rule
  ]
}

locals {
  name_filter_result = [for v in data.huaweicloud_workspace_application_rules.filter_by_name.rules[*].name : strcontains(v, "%[2]s")]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`, testAccDataApplicationRules_base(name), name)
}
