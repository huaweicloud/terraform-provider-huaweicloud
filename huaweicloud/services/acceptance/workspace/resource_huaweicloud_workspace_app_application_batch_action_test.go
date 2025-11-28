package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppApplicationBatchAction_basic(t *testing.T) {
	var (
		rName           = "huaweicloud_workspace_app_application_batch_action.test"
		rNameWithEnable = "huaweicloud_workspace_app_application_batch_action.enable"
		name            = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroupId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppApplicationBatchAction_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rName, "action", "disable"),
					resource.TestCheckResourceAttr(rName, "application_ids.#", "1"),
				),
			},
			{
				Config: testAccAppApplicationBatchAction_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(rNameWithEnable, "action", "enable"),
					resource.TestCheckResourceAttr(rNameWithEnable, "application_ids.#", "1"),
				),
			},
		},
	})
}

func testAccAppApplicationBatchAction_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = "%[1]s"
  name            = "%[2]s"
}

# The application is enabled by default.
resource "huaweicloud_workspace_app_publishment" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id
  name         = "%[2]s"
  type         = 3
  execute_path = "C:\\Program Files\\terraform\\terraform.exe"
}
`, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_ID, name)
}

func testAccAppApplicationBatchAction_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_application_batch_action" "test" {
  app_group_id    = huaweicloud_workspace_app_group.test.id
  action          = "disable"
  application_ids = huaweicloud_workspace_app_publishment.test[*].id
}
`, testAccAppApplicationBatchAction_base(name))
}

func testAccAppApplicationBatchAction_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_application_batch_action" "enable" {
  app_group_id    = huaweicloud_workspace_app_group.test.id
  action          = "enable"
  application_ids = huaweicloud_workspace_app_publishment.test[*].id
}
`, testAccAppApplicationBatchAction_base(name))
}
