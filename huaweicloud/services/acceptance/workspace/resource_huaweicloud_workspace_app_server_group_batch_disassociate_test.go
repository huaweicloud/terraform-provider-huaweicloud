package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test, please make sure the server group have no application group.
func TestAccResourceAppServerGroupBatchDisassociate_basic(t *testing.T) {
	var (
		resNameForServerGroupBatchDisassociate = "huaweicloud_workspace_app_server_group_batch_disassociate.test"
		resNameForApplicationGroup             = "huaweicloud_workspace_app_group.test"
		name                                   = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time batch disassociate application groups from server group resource and there is no
		// logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAppServerGroupBatchDisassociate_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resNameForApplicationGroup, "name", name),
					resource.TestCheckResourceAttr(resNameForApplicationGroup, "type", "SESSION_DESKTOP_APP"),
					resource.TestCheckResourceAttr(resNameForApplicationGroup, "server_group_id",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID),
				),
			},
			{
				Config: testAccResourceAppServerGroupBatchDisassociate_step2(name),
				Check: resource.ComposeTestCheckFunc(
					// Check that the batch disassociate resource is created successfully
					resource.TestCheckResourceAttr(resNameForServerGroupBatchDisassociate, "server_group_id",
						acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID),
				),
			},
			{
				Config: testAccResourceAppServerGroupBatchDisassociate_step3(name),
				Check: resource.ComposeTestCheckFunc(
					// Refresh and check the application group now have empty server_group_id after disassociation
					resource.TestCheckResourceAttr(resNameForApplicationGroup, "server_group_id", ""),
				),
			},
		},
	})
}

func testAccResourceAppServerGroupBatchDisassociate_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = "%[1]s"
  name            = "%[2]s"
  type            = "SESSION_DESKTOP_APP"

  lifecycle {
    ignore_changes = [
      server_group_id
    ]
  }
}
`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID, name)
}

func testAccResourceAppServerGroupBatchDisassociate_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_server_group_batch_disassociate" "test" {
  server_group_id = "%[2]s"
}
`, testAccResourceAppServerGroupBatchDisassociate_step1(name), acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID)
}

func testAccResourceAppServerGroupBatchDisassociate_step3(name string) string {
	// This step is just for checking the updated values, no additional resources needed
	return testAccResourceAppServerGroupBatchDisassociate_step2(name)
}
