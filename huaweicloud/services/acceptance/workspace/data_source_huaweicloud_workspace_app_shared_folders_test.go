package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppSharedFolders_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_workspace_app_shared_folders.test"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByStorageClaimId   = "data.huaweicloud_workspace_app_shared_folders.filter_by_storage_claim_id"
		dcFilterByStorageClaimId = acceptance.InitDataSourceCheck(filterByStorageClaimId)

		filterByPath   = "data.huaweicloud_workspace_app_shared_folders.filter_by_path"
		dcFilterByPath = acceptance.InitDataSourceCheck(filterByPath)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSfsFileSystemNames(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAppSharedFolders_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "shared_folders.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'storage_claim_id' parameter.
					dcFilterByStorageClaimId.CheckResourceExists(),
					resource.TestCheckOutput("is_storage_claim_id_filter_useful", "true"),
					// Filter by 'path' parameter.
					dcFilterByPath.CheckResourceExists(),
					resource.TestCheckOutput("is_path_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAppSharedFolders_base(name string) string {
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
}`, name, acceptance.HW_SFS_FILE_SYSTEM_NAMES)
}

func testAccDataSourceAppSharedFolders_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter
data "huaweicloud_workspace_app_shared_folders" "test" {
  storage_id = huaweicloud_workspace_app_nas_storage.test.id

  depends_on = [
    huaweicloud_workspace_app_shared_folder.test,
  ]
}

# Filter by 'storage_claim_id' parameter
locals {
  storage_claim_id = huaweicloud_workspace_app_shared_folder.test.id
}

data "huaweicloud_workspace_app_shared_folders" "filter_by_storage_claim_id" {
  storage_id       = huaweicloud_workspace_app_nas_storage.test.id
  storage_claim_id = local.storage_claim_id

  depends_on = [
    huaweicloud_workspace_app_shared_folder.test,
  ]
}

locals {
  storage_claim_id_filter_result = [
    for o in data.huaweicloud_workspace_app_shared_folders.filter_by_storage_claim_id.shared_folders:
    o.storage_claim_id == local.storage_claim_id
  ]
}

output "is_storage_claim_id_filter_useful" {
  value = length(local.storage_claim_id_filter_result) == 1 && alltrue(local.storage_claim_id_filter_result)
}

# Filter by 'path' parameter
locals {
  folder_path = "shares/%[2]s/"
}

data "huaweicloud_workspace_app_shared_folders" "filter_by_path" {
  storage_id = huaweicloud_workspace_app_nas_storage.test.id
  path       = local.folder_path

  depends_on = [
    huaweicloud_workspace_app_shared_folder.test,
  ]
}

locals {
  path_filter_result = [for o in data.huaweicloud_workspace_app_shared_folders.filter_by_path.shared_folders: o.folder_path == local.folder_path]
}

output "is_path_filter_useful" {
  value = length(local.path_filter_result) >= 1 && alltrue(local.path_filter_result)
}
`, testAccDataSourceAppSharedFolders_base(name), name)
}
