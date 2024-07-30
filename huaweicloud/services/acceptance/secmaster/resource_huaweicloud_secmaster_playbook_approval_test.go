package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPlaybookApproval_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterPlaybookApproval(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testPlaybookApproval_basic(),
			},
		},
	})
}

func testPlaybookApproval_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_playbook_approval" "test" {
  workspace_id = "%[1]s"
  version_id   = "%[2]s"
  result       = "PASS"
  content      = "ok"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_UNAUDITED_VERSION_ID)
}
