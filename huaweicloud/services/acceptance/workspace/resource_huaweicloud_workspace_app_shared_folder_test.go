package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getAppSharedFolderFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	return workspace.GetAppSharedFolderById(client, state.Primary.Attributes["storage_id"], state.Primary.ID)
}

func TestAccAppSharedFolder_basic(t *testing.T) {
	var (
		resourceName = "huaweicloud_workspace_app_shared_folder.test"
		name         = acceptance.RandomAccResourceName()

		sharedFolder interface{}
		rc           = acceptance.InitResourceCheck(resourceName, &sharedFolder, getAppSharedFolderFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckSfsFileSystemNames(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppSharedFolder_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "storage_id",
						"huaweicloud_workspace_app_nas_storage.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "path", fmt.Sprintf("shares/%s/", name)),
					resource.TestCheckResourceAttr(resourceName, "delimiter", "/"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The format of import ID is <storage_id>/<shared_folder_id>.
				ImportStateIdFunc: testAccAppSharedFolderImportStateIdFunc(resourceName),
				ImportStateVerifyIgnore: []string{
					"name",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The format of import ID is <storage_name>/<shared_folder_id>.
				ImportStateIdFunc: testAccAppSharedFolderImportStateIdFunc(resourceName, name),
				ImportStateVerifyIgnore: []string{
					"name",
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// The format of import ID is <storage_name>/<shared_folder_name>.
				ImportStateIdFunc: testAccAppSharedFolderImportStateIdFunc(resourceName, name, name),
				ImportStateVerifyIgnore: []string{
					"name",
				},
			},
		},
	})
}

func testAccAppSharedFolderImportStateIdFunc(resourceName string, names ...string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) is not found in the tfstate", resourceName)
		}

		switch len(names) {
		case 0:
			return fmt.Sprintf("%s/%s", rs.Primary.Attributes["storage_id"], rs.Primary.ID), nil
		case 1:
			// Related NAS storage ID is unknow, using storage name instead of storage ID.
			// The ID format is: <storage_name>/<shared_folder_id>
			return fmt.Sprintf("%s/%s", names[0], rs.Primary.ID), nil
		case 2:
			// Bosth related NAS storage ID and shared folder ID are unknow, using their names instead of IDs.
			// The ID format is: <storage_name>/<shared_folder_name>
			return fmt.Sprintf("%s/%s", names[0], names[1]), nil
		}

		return "", fmt.Errorf("invalid length of the names input, most two IDs are allowed")
	}
}

func testAccAppSharedFolder_basic(name string) string {
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
`, name, acceptance.HW_SFS_FILE_SYSTEM_NAMES)
}
