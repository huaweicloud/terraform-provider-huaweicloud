package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUpdateWorkflowInstance_basic(t *testing.T) {
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
				Config: testUpdateWorkflowInstance_basic(),
			},
		},
	})
}

func testUpdateWorkflowInstance_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_update_workflow_instance" "test" {
  workspace_id = "%[1]s"
  instance_id  = "%[2]s"
  command_type = "ActionInstanceTerminateCommand"
}`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_INSTANCE_ID)
}
