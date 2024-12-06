package workspace

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getAppPersonalFoldersFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	assignments, err := workspace.ListAppPersonalFolders(client, state.Primary.ID)
	if len(assignments) < 1 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Body: []byte("All personal folders have been removed from the Workspace APP service"),
			},
		}
	}
	return assignments, err
}

func TestAccAppPersonalFolders_basic(t *testing.T) {
	var (
		resourceName        = "huaweicloud_workspace_app_personal_folders.test"
		storageResourceName = "huaweicloud_workspace_app_nas_storage.test"
		name                = acceptance.RandomAccResourceName()

		personalFolders interface{}
		rc              = acceptance.InitResourceCheck(resourceName, &personalFolders, getAppPersonalFoldersFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSfsFileSystemNames(t, 1)
			acceptance.TestAccPrecheckWorkspaceUserNames(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppPersonalFolders_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "storage_id",
						"huaweicloud_workspace_app_nas_storage.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "assignments.#",
						strconv.Itoa(len(strings.Split(acceptance.HW_WORKSPACE_USER_NAMES, ",")))),
				),
			},
			// Import by ID (also NAS storage ID).
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Import by NAS storage name.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppNasStorageImportStateIdFunc(storageResourceName),
			},
		},
	})
}

func testAccAppPersonalFolders_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_nas_storage" "test" {
  name = "%[1]s"

  storage_metadata {
    storage_handle = element(split(",", "%[2]s"), 0)
    storage_class  = "sfs"
  }
}

resource "huaweicloud_workspace_app_storage_policy" "test" {
  server_actions = ["GetObject", "PutObject", "DeleteObject"]
  client_actions = ["GetObject", "PutObject", "DeleteObject"]
}

resource "huaweicloud_workspace_app_personal_folders" "test" {
  storage_id = huaweicloud_workspace_app_nas_storage.test.id

  dynamic "assignments" {
    for_each = split(",", "%[3]s")

    content {
      policy_statement_id = huaweicloud_workspace_app_storage_policy.test.id
      attach              = assignments.value
      attach_type         = "USER"
    }
  }
}
`, name, acceptance.HW_SFS_FILE_SYSTEM_NAMES, acceptance.HW_WORKSPACE_USER_NAMES)
}
