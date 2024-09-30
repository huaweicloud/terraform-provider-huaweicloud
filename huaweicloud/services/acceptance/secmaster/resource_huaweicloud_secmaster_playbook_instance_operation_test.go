package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPlaybookInstanceOperation_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				// Currently, no test conditions can be constructed.
				Config: testAccPlaybookInstanceOperation_basic(),
			},
		},
	})
}

func testAccPlaybookInstanceOperation_basic() string {
	return fmt.Sprintf(`

resource "huaweicloud_secmaster_playbook_instance_operation" "test" {
  workspace_id = "%[1]s"
  instance_id  = "%[2]s"
  operation    = "TERMINATE"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_INSTANCE_ID)
}
