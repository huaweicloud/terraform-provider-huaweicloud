package geminidb

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to the lack of testing conditions, only the error situations of API calls were verified.
func TestAccEnlargeFailNodeDelete_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
			acceptance.TestAccPreCheckGeminiDBNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testEnlargeFailNodeDelete_basic(),
				ExpectError: regexp.MustCompile("error deleting GeminiDB instance enlarge failed node"),
			},
		},
	})
}

func testEnlargeFailNodeDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_enlarge_fail_node_delete" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, acceptance.HW_GEMINIDB_NODE_ID)
}
