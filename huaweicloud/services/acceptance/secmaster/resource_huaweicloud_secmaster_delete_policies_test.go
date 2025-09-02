package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Note: Due to limited test conditions, this test case only verifies the expected error scenario.
func TestAccDeletePolicies_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDeletePolicies_basic(),
				ExpectError: regexp.MustCompile(`无效的参数`),
			},
		},
	})
}

// Due to the lack of a test environment, only expected failure scenarios can be tested at present.
// The value of `batch_ids` in the test script is mock data.
func testAccDeletePolicies_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_delete_policies" "test" {
  workspace_id = "%s"
  batch_ids    = ["0f60a2f3-4de7-4ed1-8b08-388765e1ffd1"]
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
