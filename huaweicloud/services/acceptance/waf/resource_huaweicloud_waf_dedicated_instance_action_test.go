package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// The operation tested in this test case is to modify the security group of the WAF dedicated instance.
func TestAccDedicatedInstanceAction_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// WAF dedicated instance is only supported in some regions, `cn-east-5` is currently supported.
			// Prepare a WAF dedicated instance and a security group before executing this test case.
			acceptance.TestAccPreCheckWafDedicatedInstanceAction(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInstanceAction_basic(),
			},
		},
	})
}

func testAccDedicatedInstanceAction_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_waf_dedicated_instance_action" "test" {
  instance_id = "%[1]s"
  action      = "security_groups"
  params      = ["%[2]s"]
}
`, acceptance.HW_WAF_DEDICATED_INSTANCE_ID, acceptance.HW_WAF_DEDICATED_INSTANCE_SECURITY_GROUP_ID)
}
