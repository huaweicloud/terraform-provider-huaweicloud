package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRiskRuleAction_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckDbssInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				// The ont-time resource do not need to be checked and no processing is performed in the Read method.
				Config: testAccRiskRuleAction_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("result_is_success", "true"),
				),
			},
		},
	})
}

func testAccRiskRuleAction_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dbss_audit_risk_rules" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_dbss_audit_risk_rule_action" "test" {
  instance_id = "%[1]s"
  risk_ids    = data.huaweicloud_dbss_audit_risk_rules.test.rules[0].id
  action      = "OFF"
}

output "result_is_success" {
  value = (huaweicloud_dbss_audit_risk_rule_action.test.result == "SUCCESS")
}
`, acceptance.HW_DBSS_INSATNCE_ID)
}
