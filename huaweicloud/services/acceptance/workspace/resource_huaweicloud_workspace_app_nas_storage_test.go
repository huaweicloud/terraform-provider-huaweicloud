package workspace

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getAppNasStorageFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	return workspace.GetAppNasStorageById(client, state.Primary.ID)
}

func getFirstSfsFileSystemName(names []string) string {
	if len(names) < 1 {
		return ""
	}
	return names[0]
}

func TestAccAppNasStorage_basic(t *testing.T) {
	var (
		resourceName      = "huaweicloud_workspace_app_nas_storage.test"
		name              = acceptance.RandomAccResourceName()
		sfsFileSystemName = getFirstSfsFileSystemName(strings.Split(acceptance.HW_SFS_FILE_SYSTEM_NAMES, ","))

		appGroup interface{}
		rc       = acceptance.InitResourceCheck(resourceName, &appGroup, getAppNasStorageFunc)
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
				Config: testAccAppNasStorage_basic(name, sfsFileSystemName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "storage_metadata.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "storage_metadata.0.storage_handle", sfsFileSystemName),
					resource.TestCheckResourceAttr(resourceName, "storage_metadata.0.storage_class", "sfs"),
					resource.TestCheckResourceAttrSet(resourceName, "storage_metadata.0.export_location"),
					resource.TestMatchResourceAttr(resourceName, "created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
			// Import by ID.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Import by name.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccAppNasStorageImportStateIdFunc(resourceName),
			},
		},
	})
}

func testAccAppNasStorageImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var storageName string
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("the resource (%s) is not found in the tfstate", resourceName)
		}
		storageName = rs.Primary.Attributes["name"]
		if storageName == "" {
			return "", fmt.Errorf("the NAS storage name is missing")
		}
		return storageName, nil
	}
}

func testAccAppNasStorage_basic(name, sfsFileSystemName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_nas_storage" "test" {
  name = "%[1]s"

  storage_metadata {
    storage_handle = "%[2]s"
    storage_class  = "sfs"
  }
}
`, name, sfsFileSystemName)
}
