package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkflowAction_basic(t *testing.T) {
	name := "CIS_Ensuring IAM Policies Are Not Created to Allow Wildcard Administrative Permissions"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testWorkflowAction_basic(name),
			},
		},
	})
}

func testWorkflowAction_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_workflows" "test" {
  workspace_id = "%[2]s"
  name         = "%[1]s"
}

resource "huaweicloud_secmaster_workflow_action" "test" {
  workspace_id = "%[2]s"
  workflow_id  = data.huaweicloud_secmaster_workflows.test.workflows[0].id
  command_type = "ActionInstanceRunCommand"
  action_type  = "workflow"
}
`, name, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
