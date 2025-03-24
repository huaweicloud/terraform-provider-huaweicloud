package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceExecutePolicy_basic(t *testing.T) {
	// Avoid CheckDestroy because this resource is a one-time action resource and there is nothing in the destroy method.
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckASScalingPolicyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceExecutePolicy_basic(),
			},
		},
	})
}

func testAccResourceExecutePolicy_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_as_execute_policy" "test" {
  scaling_policy_id = "%s"
}
`, acceptance.HW_AS_SCALING_POLICY_ID)
}
