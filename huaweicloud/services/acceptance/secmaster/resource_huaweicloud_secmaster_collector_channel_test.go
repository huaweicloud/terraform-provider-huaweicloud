package secmaster

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to the lack of a test environment, only the expected failure scenario can be tested at present.
// The values of `group_id` and `parser_id` in the test script are mock data.
func TestAccResourceCollectorChannel_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceCollectorChannel_basic(),
				ExpectError: regexp.MustCompile(`参数不合法`),
			},
		},
	})
}

func testAccResourceCollectorChannel_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_collector_channel" "test" {
  workspace_id = "%[1]s"
  title        = "tf_test_channel"
  group_id     = "00000000-0000-0000-0000-000000000000"
  parser_id    = "00000000-0000-0000-0000-000000000000"

  input {
  }

  output {
  }

  nodes {
  }
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
