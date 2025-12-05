package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopUserBatchAttach_basic(t *testing.T) {
	var (
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_workspace_desktop_user_batch_attach.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopUserBatchAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desktops.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "desktops.0.user_name", name),
				),
			},
		},
	})
}

func testAccDesktopUserBatchAttach_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "test" {
  name        = "%[2]s"
  email       = "basic@example.com"
  description = "Created by acc test"
}

resource "huaweicloud_workspace_desktop_user_batch_attach" "test" {
  desktops {
    desktop_id = "%[1]s"
    user_name  = huaweicloud_workspace_user.test.name
    user_email = huaweicloud_workspace_user.test.email
  }

  enable_force_new = "true"
}
`, acceptance.HW_WORKSPACE_DESKTOP_IDS, name)
}

func TestAccDesktopUserBatchAttach_mulity(t *testing.T) {
	var (
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_workspace_desktop_user_batch_attach.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopUserBatchAttach_mulity(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "desktops.#", "1"),
				),
			},
		},
	})
}

func testAccDesktopUserBatchAttach_mulity(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "test1" {
  name        = "%[2]s_1"
  email       = "basic@example.com"
  description = "Created by acc test"
}

resource "huaweicloud_workspace_user" "test2" {
  name        = "%[2]s_2"
  email       = "basic@example.com"
  description = "Created by acc test"
}

resource "huaweicloud_workspace_desktop_user_batch_attach" "test" {
  desktops {
    desktop_id = "%[1]s"

    attach_user_infos {
      user_name  = huaweicloud_workspace_user.test1.name
      user_id    = huaweicloud_workspace_user.test1.id
      user_group = "users"
    }

    attach_user_infos {
      user_name  = huaweicloud_workspace_user.test2.name
      user_id    = huaweicloud_workspace_user.test2.id
      user_group = "users"
    }
  }
}
`, acceptance.HW_WORKSPACE_DESKTOP_IDS, name)
}
