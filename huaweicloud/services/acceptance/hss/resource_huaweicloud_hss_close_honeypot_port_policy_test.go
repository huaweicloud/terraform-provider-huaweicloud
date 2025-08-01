package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCloseHoneypotPortPolicy_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID that has enabled premium edition host protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
			// The dynamic port honeypot policy ID is Required.
			acceptance.TestAccPreCheckHSSPolicyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCloseHoneypotPortPolicy_basic(),
			},
		},
	})
}

func testCloseHoneypotPortPolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_hss_close_honeypot_port_policy" "test" {
  policy_id             = "%[1]s"
  host_id               = "%[2]s"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_POLICY_ID, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
