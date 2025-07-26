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
				),
			},
		},
	})
}

func testAccAppRulesDataSource_base(name string) string {
	// waiting for the implement of app rule resource
	return fmt.Sprintf("%[1]s", name)
}

func testAccAppRulesDataSource_basic(name string) string {
	// the name filter case need validate the context is contain the filter parameter?
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_rules" "test" {}

data "huaweicloud_workspace_app_rules" "filter_by_name" {
  name = "test-rule"
}
`, testAccAppRulesDataSource_base(name))
}
