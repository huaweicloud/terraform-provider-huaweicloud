package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/secmaster"
)

func getResourceWorkflowVersionFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("secmaster", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SecMaster client: %s", err)
	}

	return secmaster.GetWorkflowVersionInfo(client, state.Primary.Attributes["workspace_id"], state.Primary.ID)
}

func TestAccResourceWorkflowVersion_basic(t *testing.T) {
	var (
		rName      = "huaweicloud_secmaster_workflow_version.test"
		name       = acceptance.RandomAccResourceName()
		updateName = acceptance.RandomAccResourceName()

		taskFlow1 = "eyJlbGVtZW50VHlwZSI6Int9Iiwid29ya2Zsb3ciOnsiZWxlbWVudHMiOlt7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImJwbW4" +
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

		taskFlow2 = "eyJlbGVtZW50VHlwZSI6Int9Iiwid29ya2Zsb3ciOnsiZWxlbWVudHMiOlt7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImJwbW4yOm" +
			"RlZmluaXRpb25zIiwiYXR0cmlidXRlcyI6eyJ4bWxuczp4c2kiOiJodHRwOi8vd3d3LnczLm9yZy8yMDAxL1hNTFNjaGVtYS1pbnN0YW5jZSIsInhtbG" +
			"5zOmJwbW4yIjoiaHR0cDovL3d3dy5vbWcub3JnL3NwZWMvQlBNTi8yMDEwMDUyNC9NT0RFTCIsInhtbG5zOmJwbW5kaSI6Imh0dHA6Ly93d3cub21nLm" +
			"9yZy9zcGVjL0JQTU4vMjAxMDA1MjQvREkiLCJ4bWxuczpkYyI6Imh0dHA6Ly93d3cub21nLm9yZy9zcGVjL0RELzIwMTAwNTI0L0RDIiwieG1sbnM6ZG" +
			"kiOiJodHRwOi8vd3d3Lm9tZy5vcmcvc3BlYy9ERC8yMDEwMDUyNC9ESSIsInhtbG5zOmZsb3dhYmxlIjoiaHR0cDovL2Zsb3dhYmxlLm9yZy9icG1uIi" +
			"wieG1sbnM6c29hciI6Imh0dHBzOi8vc29jLmNsb3Vkc3JlLmNvbS9zb2FyL3YyIiwiaWQiOiJkaWFncmFtX1Byb2Nlc3NfMTc1NjcxMjM0MzgyNyIsIn" +
			"RhcmdldE5hbWVzcGFjZSI6Imh0dHA6Ly9mbG93YWJsZS5vcmcvYnBtbiIsInhzaTpzY2hlbWFMb2NhdGlvbiI6Imh0dHA6Ly93d3cub21nLm9yZy9zcG" +
			"VjL0JQTU4vMjAxMDA1MjQvTU9ERUwgQlBNTjIwLnhzZCJ9LCJlbGVtZW50cyI6W3sidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbjI6cHJvY2Vzcy" +
			"IsImF0dHJpYnV0ZXMiOnsiaWQiOiJQcm9jZXNzXzE3NTY3MTIzNDM4MjciLCJuYW1lIjoiU2VydmljZV9Qcm9jZXNzMTc1NjcxMjM0MzgyNyIsImlzRX" +
			"hlY3V0YWJsZSI6InRydWUifSwiZWxlbWVudHMiOlt7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImJwbW4yOnN0YXJ0RXZlbnQiLCJhdHRyaWJ1dGVzIj" +
			"p7ImlkIjoiRXZlbnRfMDh4MHoxbyJ9LCJlbGVtZW50cyI6W3sidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbjI6b3V0Z29pbmciLCJlbGVtZW50cy" +
			"I6W3sidHlwZSI6InRleHQiLCJ0ZXh0IjoiRmxvd18wYmt5ZGh5In1dfV19LHsidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbjI6ZW5kRXZlbnQiLC" +
			"JhdHRyaWJ1dGVzIjp7ImlkIjoiRXZlbnRfMG5mbXhsbiJ9LCJlbGVtZW50cyI6W3sidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbjI6aW5jb21pbm" +
			"ciLCJlbGVtZW50cyI6W3sidHlwZSI6InRleHQiLCJ0ZXh0IjoiRmxvd18wZ3lwbTRxIn1dfV19LHsidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbj" +
			"I6dXNlclRhc2siLCJhdHRyaWJ1dGVzIjp7ImlkIjoiQWN0aXZpdHlfMXlocnlrZCIsIm5hbWUiOiJ0ZXN0IiwiZmxvd2FibGU6ZHVlRGF0ZSI6IiIsIn" +
			"NvYXI6ZHVlX2hhbmRsZSI6ImludGVycnVwdCIsInNvYXI6aGFuZGxlcnMiOiJbXSJ9LCJlbGVtZW50cyI6W3sidHlwZSI6ImVsZW1lbnQiLCJuYW1lIj" +
			"oiYnBtbjI6ZG9jdW1lbnRhdGlvbiJ9LHsidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbjI6aW5jb21pbmciLCJlbGVtZW50cyI6W3sidHlwZSI6In" +
			"RleHQiLCJ0ZXh0IjoiRmxvd18wYmt5ZGh5In1dfSx7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImJwbW4yOm91dGdvaW5nIiwiZWxlbWVudHMiOlt7In" +
			"R5cGUiOiJ0ZXh0IiwidGV4dCI6IkZsb3dfMGd5cG00cSJ9XX1dfSx7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImJwbW4yOnNlcXVlbmNlRmxvdyIsIm" +
			"F0dHJpYnV0ZXMiOnsiaWQiOiJGbG93XzBia3lkaHkiLCJzb3VyY2VSZWYiOiJFdmVudF8wOHgwejFvIiwidGFyZ2V0UmVmIjoiQWN0aXZpdHlfMXlocn" +
			"lrZCJ9fSx7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImJwbW4yOnNlcXVlbmNlRmxvdyIsImF0dHJpYnV0ZXMiOnsiaWQiOiJGbG93XzBneXBtNHEiLC" +
			"Jzb3VyY2VSZWYiOiJBY3Rpdml0eV8xeWhyeWtkIiwidGFyZ2V0UmVmIjoiRXZlbnRfMG5mbXhsbiJ9fV19LHsidHlwZSI6ImVsZW1lbnQiLCJuYW1lIj" +
			"oiYnBtbmRpOkJQTU5EaWFncmFtIiwiYXR0cmlidXRlcyI6eyJpZCI6IkJQTU5EaWFncmFtXzEifSwiZWxlbWVudHMiOlt7InR5cGUiOiJlbGVtZW50Ii" +
			"wibmFtZSI6ImJwbW5kaTpCUE1OUGxhbmUiLCJhdHRyaWJ1dGVzIjp7ImlkIjoiQlBNTlBsYW5lXzEiLCJicG1uRWxlbWVudCI6IlByb2Nlc3NfMTc1Nj" +
			"cxMjM0MzgyNyJ9LCJlbGVtZW50cyI6W3sidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbmRpOkJQTU5TaGFwZSIsImF0dHJpYnV0ZXMiOnsiaWQiOi" +
			"JBY3Rpdml0eV8xeWhyeWtkX2RpIiwiYnBtbkVsZW1lbnQiOiJBY3Rpdml0eV8xeWhyeWtkIn0sImVsZW1lbnRzIjpbeyJ0eXBlIjoiZWxlbWVudCIsIm" +
			"5hbWUiOiJkYzpCb3VuZHMiLCJhdHRyaWJ1dGVzIjp7IngiOiItNDcwIiwieSI6Ii01MjAiLCJ3aWR0aCI6IjEwMCIsImhlaWdodCI6IjgwIn19LHsidH" +
			"lwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbmRpOkJQTU5MYWJlbCJ9XX0seyJ0eXBlIjoiZWxlbWVudCIsIm5hbWUiOiJicG1uZGk6QlBNTlNoYXBlIi" +
			"wiYXR0cmlidXRlcyI6eyJpZCI6IkV2ZW50XzBuZm14bG5fZGkiLCJicG1uRWxlbWVudCI6IkV2ZW50XzBuZm14bG4ifSwiZWxlbWVudHMiOlt7InR5cG" +
			"UiOiJlbGVtZW50IiwibmFtZSI6ImRjOkJvdW5kcyIsImF0dHJpYnV0ZXMiOnsieCI6Ii0zMTgiLCJ5IjoiLTQ5OCIsIndpZHRoIjoiMzYiLCJoZWlnaH" +
			"QiOiIzNiJ9fV19LHsidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbmRpOkJQTU5TaGFwZSIsImF0dHJpYnV0ZXMiOnsiaWQiOiJFdmVudF8wOHgwej" +
			"FvX2RpIiwiYnBtbkVsZW1lbnQiOiJFdmVudF8wOHgwejFvIn0sImVsZW1lbnRzIjpbeyJ0eXBlIjoiZWxlbWVudCIsIm5hbWUiOiJkYzpCb3VuZHMiLC" +
			"JhdHRyaWJ1dGVzIjp7IngiOiItNTg4IiwieSI6Ii00OTgiLCJ3aWR0aCI6IjM2IiwiaGVpZ2h0IjoiMzYifX1dfSx7InR5cGUiOiJlbGVtZW50Iiwibm" +
			"FtZSI6ImJwbW5kaTpCUE1ORWRnZSIsImF0dHJpYnV0ZXMiOnsiaWQiOiJGbG93XzBia3lkaHlfZGkiLCJicG1uRWxlbWVudCI6IkZsb3dfMGJreWRoeS" +
			"J9LCJlbGVtZW50cyI6W3sidHlwZSI6ImVsZW1lbnQiLCJuYW1lIjoiZGk6d2F5cG9pbnQiLCJhdHRyaWJ1dGVzIjp7IngiOiItNTUyIiwieSI6Ii00OD" +
			"AifX0seyJ0eXBlIjoiZWxlbWVudCIsIm5hbWUiOiJkaTp3YXlwb2ludCIsImF0dHJpYnV0ZXMiOnsieCI6Ii00NzAiLCJ5IjoiLTQ4MCJ9fV19LHsidH" +
			"lwZSI6ImVsZW1lbnQiLCJuYW1lIjoiYnBtbmRpOkJQTU5FZGdlIiwiYXR0cmlidXRlcyI6eyJpZCI6IkZsb3dfMGd5cG00cV9kaSIsImJwbW5FbGVtZW" +
			"50IjoiRmxvd18wZ3lwbTRxIn0sImVsZW1lbnRzIjpbeyJ0eXBlIjoiZWxlbWVudCIsIm5hbWUiOiJkaTp3YXlwb2ludCIsImF0dHJpYnV0ZXMiOnsieC" +
			"I6Ii0zNzAiLCJ5IjoiLTQ4MCJ9fSx7InR5cGUiOiJlbGVtZW50IiwibmFtZSI6ImRpOndheXBvaW50IiwiYXR0cmlidXRlcyI6eyJ4IjoiLTMxOCIsIn" +
			"kiOiItNDgwIn19XX1dfV19XX1dfX0="

		object interface{}
		rc     = acceptance.InitResourceCheck(
			rName,
			&object,
			getResourceWorkflowVersionFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWorkflowVersion_basic(name, taskFlow1),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "workspace_id", acceptance.HW_SECMASTER_WORKSPACE_ID),
					resource.TestCheckResourceAttr(rName, "workflow_id", acceptance.HW_SECMASTER_WORKFLOW_ID),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "taskflow", taskFlow1),
					resource.TestCheckResourceAttr(rName, "taskconfig", "{\"node_info\":{},\"usertask_info\":{}}"),
					resource.TestCheckResourceAttr(rName, "taskflow_type", "JSON"),
					resource.TestCheckResourceAttr(rName, "aop_type", "NORMAL"),
					resource.TestCheckResourceAttr(rName, "description", "terraform test"),
					resource.TestCheckResourceAttrSet(rName, "status"),
					resource.TestCheckResourceAttrSet(rName, "creator_id"),
				),
			},
			{
				Config: testAccWorkflowVersion_update(updateName, taskFlow2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", updateName),
					resource.TestCheckResourceAttr(rName, "taskflow", taskFlow2),
					resource.TestCheckResourceAttr(rName, "taskconfig", "{\"test\":\"abc\"}"),
					resource.TestCheckResourceAttr(rName, "aop_type", "SURVEY"),
					resource.TestCheckResourceAttr(rName, "description", ""),
				),
			},
			{
				Config: testAccWorkflowVersion_update_status(updateName, taskFlow2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "status", "pending_approval"),
					resource.TestCheckResourceAttrSet(rName, "version"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccWorkflowVersionImportStateFunc(rName),
			},
		},
	})
}

func testAccWorkflowVersion_basic(name, taskFlow string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_workflow_version" "test" {
  workspace_id  = "%[1]s"
  workflow_id   = "%[2]s"
  name          = "%[3]s"
  taskflow      = "%[4]s"
  taskconfig    = "{\"node_info\":{},\"usertask_info\":{}}"
  taskflow_type = "JSON"
  aop_type      = "NORMAL"
  description   = "terraform test"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_WORKFLOW_ID, name, taskFlow)
}

func testAccWorkflowVersion_update(name, taskFlow string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_workflow_version" "test" {
  workspace_id  = "%[1]s"
  workflow_id   = "%[2]s"
  name          = "%[3]s"
  taskflow      = "%[4]s"
  taskconfig    = "{\"test\":\"abc\"}"
  taskflow_type = "JSON"
  aop_type      = "SURVEY"
  description   = ""
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_WORKFLOW_ID, name, taskFlow)
}

func testAccWorkflowVersion_update_status(name, taskFlow string) string {
	return fmt.Sprintf(`
resource "huaweicloud_secmaster_workflow_version" "test" {
  workspace_id  = "%[1]s"
  workflow_id   = "%[2]s"
  name          = "%[3]s"
  taskflow      = "%[4]s"
  taskconfig    = "{\"test\":\"abc\"}"
  taskflow_type = "JSON"
  aop_type      = "SURVEY"
  description   = ""
  status        = "pending_approval"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_WORKFLOW_ID, name, taskFlow)
}

func testAccWorkflowVersionImportStateFunc(rName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var workspaceId, versionId string
		rs, ok := s.RootModule().Resources[rName]
		if !ok {
			return "", fmt.Errorf("resource (%s) not found: %s", rName, rs)
		}

		workspaceId = rs.Primary.Attributes["workspace_id"]
		versionId = rs.Primary.ID

		if workspaceId == "" || versionId == "" {
			return "", fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s/%s'",
				workspaceId, versionId)
		}

		return fmt.Sprintf("%s/%s", workspaceId, versionId), nil
	}
}
