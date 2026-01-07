package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccWorkflowVersionValidation_basic(t *testing.T) {
	taskFlow := "eyJlbGVtZW50VHlwZSI6Int9Iiwid29ya2Zsb3ciOnsiZWxlbWVudHMiOlt7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImJwbW4" +
		"yOmRlZmluaXRpb25zIiwiYXR0cmlidXRlcyI6eyJ4bWxuczp4c2kiOiJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYS1pbnN0YW5jZSI" +
		"sInhtbG5zOmJwbW4yIjoiaHR0cDovL3d3dy5vbWcub3JnL3NwZWMvQlBNTi8yMDEwMDUyNC9NT0RFTCIsInhtbG5zOmJwbW5kaSI6Imh0dHA6Ly93" +
		"d3cub21nLm9yZy9zcGVjL0JQTU4vMjAxMDA1MjQvREkiLCJ4bWxuczpzb2FyIjoiaHR0cHM6Ly9zb2MuY2xvdWRzcmUuY29tL3NvYXIvdjIiLCJpZ" +
		"CI6ImRpYWdyYW1fUHJvY2Vzc18xNzU0NTU1ODMxODc2IiwidGFyZ2V0TmFtZXNwYWNlIjoiaHR0cDovL2Zsb3dhYmxlLm9yZy9icG1uIiwieHNpOn" +
		"NjaGVtYUxvY2F0aW9uIjoiaHR0cDovL3d3dy5vbWcub3JnL3NwZWMvQlBNTi8yMDEwMDUyNC9NT0RFTCBCUE1OMjAueHNkIn0sImVsZW1lbnRzIjp" +
		"beyJ0eXBlIjoiZWxlbWVudCIsIm5hbWUiOiJicG1uMjpwcm9jZXNzIiwiYXR0cmlidXRlcyI6eyJpZCI6IlByb2Nlc3NfMTc1NDU1NTgzMTg3NiIs" +
		"Im5hbWUiOiJTZXJ2aWNlX1Byb2Nlc3MxNzU0NTU1ODMxODc2IiwiaXNFeGVjdXRhYmxlIjoidHJ1ZSJ9fSx7InR5cGUiOiJlbGVtZW50IiwibmFtZ" +
		"SI6ImJwbW5kaTpCUE1ORGlhZ3JhbSIsImF0dHJpYnV0ZXMiOnsiaWQiOiJCUE1ORGlhZ3JhbV8xIn0sImVsZW1lbnRzIjpbeyJ0eXBlIjoiZWxlbW" +
		"VudCIsIm5hbWUiOiJicG1uZGk6QlBNTlBsYW5lIiwiYXR0cmlidXRlcyI6eyJpZCI6IkJQTU5QbGFuZV8xIiwiYnBtbkVsZW1lbnQiOiJQcm9jZXN" +
		"zXzE3NTQ1NTU4MzE4NzYifX1dfV19XX19"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testWorkflowVersionValidation_basic(taskFlow),
			},
		},
	})
}

func testWorkflowVersionValidation_basic(taskflow string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_workflow_version_validation" "test" {
  workspace_id   = "%[1]s"
  aopworkflow_id = "%[2]s"
  mode           = "BASIC"
  taskconfig     = "{\"node_info\":{},\"usertask_info\":{}}"
  taskflow       = "%[3]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_WORKFLOW_ID, taskflow)
}
