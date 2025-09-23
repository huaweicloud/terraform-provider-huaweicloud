package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSwitchHoneypotPortPolicy_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID that has enabled premium edition host protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
			// The host need to bind a dynamic port honeypot policy before switching the honeypot policy.
			acceptance.TestAccPreCheckHSSPolicyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testSwitchHoneypotPortPolicy_basic(),
			},
		},
	})
}

func testSwitchHoneypotPortPolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_switch_honeypot_port_policy" "test" {
  policy_id             = "%[1]s"
  host_id               = "%[2]s"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_POLICY_ID, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
