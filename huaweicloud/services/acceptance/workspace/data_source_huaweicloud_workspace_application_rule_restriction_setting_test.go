package workspace

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplicationRuleRestrictionSetting_basic(t *testing.T) {
	dcName := "data.huaweicloud_workspace_application_rule_restriction_setting.test"
	dc := acceptance.InitDataSourceCheck(dcName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationRuleRestrictionSetting_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "app_restrict_rule_switch", "true"),
					resource.TestCheckResourceAttr(dcName, "app_control_mode", "0"),
					resource.TestCheckResourceAttr(dcName, "app_periodic_switch", "true"),
					resource.TestCheckResourceAttr(dcName, "app_periodic_interval", "10"),
					resource.TestCheckResourceAttr(dcName, "app_force_kill_proc_switch", "true"),
				),
			},
		},
	})
}

const testAccDataApplicationRuleRestrictionSetting_basic = `
resource "huaweicloud_workspace_application_rule_restriction_setting" "test" {
  app_restrict_rule_switch   = true
  app_control_mode           = 0
  app_periodic_switch        = true
  app_periodic_interval      = 10
  app_force_kill_proc_switch = true
}

data "huaweicloud_workspace_application_rule_restriction_setting" "test" {
  depends_on = [huaweicloud_workspace_application_rule_restriction_setting.test]
}
`
