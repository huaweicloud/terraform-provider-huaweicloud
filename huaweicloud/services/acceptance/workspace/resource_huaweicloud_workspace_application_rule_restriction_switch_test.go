package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApplicationRuleRestrictionSwitch_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_application_rule_restriction_switch.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationRuleRestrictionSwitch_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "enable"),
				),
			},
			{
				Config: testAccApplicationRuleRestrictionSwitch_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "action", "disable"),
				),
			},
		},
	})
}

const testAccApplicationRuleRestrictionSwitch_basic_step1 = `
resource "huaweicloud_workspace_application_rule_restriction_switch" "test" {
  action = "enable"

  enable_force_new = "true"
}
`

const testAccApplicationRuleRestrictionSwitch_basic_step2 = `
resource "huaweicloud_workspace_application_rule_restriction_switch" "test" {
  action = "disable"

  enable_force_new = "true"
}
`
