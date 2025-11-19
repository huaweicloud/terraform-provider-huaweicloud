package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccApplicationRuleRestrictionSetting_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_application_rule_restriction_setting.test"
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
				Config: testAccApplicationRuleRestrictionSetting_basic_step1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app_restrict_rule_switch", "false"),
				),
			},
			{
				Config: testAccApplicationRuleRestrictionSetting_basic_step2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "app_restrict_rule_switch", "true"),
					resource.TestCheckResourceAttr(resourceName, "app_control_mode", "0"),
					resource.TestCheckResourceAttr(resourceName, "app_periodic_switch", "true"),
					resource.TestCheckResourceAttr(resourceName, "app_periodic_interval", "10"),
					resource.TestCheckResourceAttr(resourceName, "app_force_kill_proc_switch", "true"),
				),
			},
		},
	})
}

const testAccApplicationRuleRestrictionSetting_basic_step1 = `
resource "huaweicloud_workspace_application_rule_restriction_setting" "test" {
  app_restrict_rule_switch = false

  enable_force_new = "true"
}
`

const testAccApplicationRuleRestrictionSetting_basic_step2 = `
resource "huaweicloud_workspace_application_rule_restriction_setting" "test" {
  app_restrict_rule_switch   = true
  app_control_mode           = 0
  app_periodic_switch        = true
  app_periodic_interval      = 10
  app_force_kill_proc_switch = true

  enable_force_new = "true"
}
`
