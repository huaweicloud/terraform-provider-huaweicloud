package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAdvancedIpsRule_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwIpsAdvancedRuleID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAdvancedIpsRule_basic(),
			},
		},
	})
}

func testAdvancedIpsRule_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_firewalls" "test" {
  fw_instance_id = "%[1]s"
}

resource "huaweicloud_cfw_advanced_ips_rule" "test" {
  fw_instance_id = "%[1]s"
  ips_rule_id    = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  param          = "{\"threshold\":11}"
  action         = 0
  ips_rule_type  = 1
  status         = 0
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_IPS_ADVANCED_RULE_ID)
}
