package coc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceCocScriptApproval_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCocScriptID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testCocScriptApproval_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testCocScriptApproval_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_coc_script_approval" "test" {
  script_uuid = "%s"
  status      = "APPROVED"
  comments    = "agree comments"
}
`, acceptance.HW_COC_SCRIPT_ID)
}
