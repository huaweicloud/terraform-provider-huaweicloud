package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataStudioWorkspaceUsers_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dataarts_studio_workspace_users.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckUserName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataStudioWorkspaceUsers_nonExistentWorkspace(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "users.#", "0"),
				),
			},
			{
				Config: testAccDataStudioWorkspaceUsers_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "workspace_id", acceptance.HW_DATAARTS_WORKSPACE_ID),
					resource.TestMatchResourceAttr(dataSource, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.name"),
					resource.TestMatchResourceAttr(dataSource, "users.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "users.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "users.0.roles.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.roles.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.roles.0.code"),
					resource.TestCheckResourceAttrSet(dataSource, "users.0.roles.0.name"),
				),
			},
		},
	})
}

func testAccDataStudioWorkspaceUsers_nonExistentWorkspace() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_users" "test" {
  workspace_id = "%[1]s"
}
`, randUUID)
}

func testAccDataStudioWorkspaceUsers_basic_base() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_identity_users" "test" {
  name = "%[2]s"
}

resource "huaweicloud_dataarts_studio_workspace_user" "test" {
  workspace_id = "%[1]s"
  user_id      = try(data.huaweicloud_identity_users.test.users[0].id, "NOT_FOUND")

  dynamic "roles" {
    for_each = slice(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles, 0, 1)

    content {
      id = roles.value.id
    }
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_USER_NAME)
}

func testAccDataStudioWorkspaceUsers_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_dataarts_studio_workspace_users" "test" {
  depends_on = [
    huaweicloud_dataarts_studio_workspace_user.test,
  ]

  workspace_id = "%[2]s"
}
`, testAccDataStudioWorkspaceUsers_basic_base(), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
