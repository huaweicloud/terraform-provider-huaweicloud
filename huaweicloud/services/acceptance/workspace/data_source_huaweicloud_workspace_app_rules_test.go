package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppRulesDataSource_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_app_rules.test"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterByName   = "data.huaweicloud_workspace_app_rules.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAppRulesDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "app_rules.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(dcName, "app_rules.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "app_rules.0.name"),
					resource.TestCheckResourceAttrSet(dcName, "app_rules.0.rule_source"),
					resource.TestCheckResourceAttrSet(dcName, "app_rules.0.create_time"),
					resource.TestCheckResourceAttrSet(dcName, "app_rules.0.update_time"),
					resource.TestCheckResourceAttrSet(dcName, "app_rules.0.rule.0.scope"),
					dcFilterByName.CheckResourceExists(),
					resource.TestMatchResourceAttr(filterByName, "app_rules.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccAppRulesDataSource_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_rule" "with_product_rule" {
  name        = "%[1]s"
  description = "Created by terraform script"

  rule {
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

resource "huaweicloud_workspace_app_rule" "with_path_rule" {
  name        = "%[1]s_path"
  description = "Created by terraform script for path rule"

  rule {
    scope = "PATH"

    path_rule {
      path = "C:\\Program Files\\Microsoft Office\\root\\Office16\\WINWORD.EXE"
    }
  }
}
`, name)
}

func testAccAppRulesDataSource_basic(name string) string {
	// the name filter case need validate the context is contain the filter parameter?
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_rules" "test" {
  depends_on = [
    huaweicloud_workspace_app_rule.with_product_rule,
    huaweicloud_workspace_app_rule.with_path_rule
  ]
}

data "huaweicloud_workspace_app_rules" "filter_by_name" {
  name = "%[2]s"

  depends_on = [
    huaweicloud_workspace_app_rule.with_product_rule,
    huaweicloud_workspace_app_rule.with_path_rule
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_workspace_app_rules.filter_by_name.app_rules) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_app_rules.filter_by_name.app_rules[*].name : strcontains(v, "%[2]s") == true]
  )
}
`, testAccAppRulesDataSource_base(name), name)
}
