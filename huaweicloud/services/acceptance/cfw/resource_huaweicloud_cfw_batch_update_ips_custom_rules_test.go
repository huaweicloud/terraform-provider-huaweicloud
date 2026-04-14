package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccBatchUpdateIpsCustomRulesAction_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting the firewall instance ID and IPS custom rule ID for CFW.
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwIpsCustomRuleId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testBatchUpdateIpsCustomRulesAction_basic(),
			},
		},
	})
}

func testBatchUpdateIpsCustomRulesAction_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_cfw_batch_update_ips_custom_rules" "test" {
  fw_instance_id = "%s"
  action_type    = 0
  ips_ids        = ["%s"]
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_IPS_CUSTOM_RULE_ID)
}
