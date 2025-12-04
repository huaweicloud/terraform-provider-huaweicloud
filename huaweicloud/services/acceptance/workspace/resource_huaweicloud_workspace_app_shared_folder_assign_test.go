package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAppSharedFolderAssign_basic(t *testing.T) {
	var name = acceptance.RandomAccResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSfsFileSystemNames(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppSharedFolderAssign_step1(name),
			},
			{
				Config: testAccAppSharedFolderAssign_step2(name),
			},
		},
	})
}

func testAccAppSharedFolderAssign_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_nas_storage" "test" {
  name = "%[1]s"

  storage_metadata {
    storage_handle = element(split(",", "%[2]s"), 0)
    storage_class  = "sfs"
  }
}

resource "huaweicloud_workspace_app_shared_folder" "test" {
  storage_id = huaweicloud_workspace_app_nas_storage.test.id
  name       = "%[1]s"
}

resource "huaweicloud_workspace_user" "test" {
  name  = "%[1]s"
  email = "tf@example.com"
}
`, name, acceptance.HW_SFS_FILE_SYSTEM_NAMES)
}

func testAccAppSharedFolderAssign_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_shared_folder_assign" "test_with_add_items" {
  storage_id       = huaweicloud_workspace_app_nas_storage.test.id
  storage_claim_id = huaweicloud_workspace_app_shared_folder.test.id

  add_items {
    policy_statement_id = "DEFAULT_1"
    attach              = huaweicloud_workspace_user.test.name
    attach_type         = "USER"
  }
}
`, testAccAppSharedFolderAssign_base(name))
}

func testAccAppSharedFolderAssign_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_shared_folder_assign" "test_with_del_items" {
  storage_id       = huaweicloud_workspace_app_nas_storage.test.id
  storage_claim_id = huaweicloud_workspace_app_shared_folder.test.id

  del_items {
    attach      = huaweicloud_workspace_user.test.name
    attach_type = "USER"
  }
}
`, testAccAppSharedFolderAssign_base(name))
}
