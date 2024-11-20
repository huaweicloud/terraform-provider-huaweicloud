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

func getWarehouseApplicationFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	return workspace.GetWarehouseApplicationById(client, state.Primary.ID)
}

func TestAccAppWarehouseApp_basic(t *testing.T) {
	var (
		application  interface{}
		resourceName = "huaweicloud_workspace_app_warehouse_app.test"
		name         = acceptance.RandomAccResourceName()
		updateName   = acceptance.RandomAccResourceName()
		rc           = acceptance.InitResourceCheck(
			resourceName,
			&application,
			getWarehouseApplicationFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceFileStorePath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccWarehouseApp_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "category", "OTHER"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "version_name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "file_store_path", acceptance.HW_WORKSPACE_APP_FILE_STRORE_OBS_PATH),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
				),
			},
			{
				Config: testAccWarehouseApp_basic_step2(updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "category", "PRODUCTIVITY_AND_COLLABORATION"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Linux"),
					resource.TestCheckResourceAttr(resourceName, "version", "2.0"),
					resource.TestCheckResourceAttr(resourceName, "version_name", "terraform_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated by script"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"icon",
				},
			},
		},
	})
}

func testAccWarehouseApp_basic_step1(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_warehouse_app" "test" {
  name            = "%[1]s"
  category        = "OTHER"
  os_type         = "Windows"
  version         = "1.0"
  version_name    = "terraform"
  file_store_path = "%[2]s"
  description     = "Created by script"
}
`, name, acceptance.HW_WORKSPACE_APP_FILE_STRORE_OBS_PATH)
}

func testAccWarehouseApp_basic_step2(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_warehouse_app" "test" {
  name            = "%[1]s"
  category        = "PRODUCTIVITY_AND_COLLABORATION"
  os_type         = "Linux"
  version         = "2.0"
  version_name    = "terraform_update"
  file_store_path = "%[2]s"
  description     = "Updated by script"
}
`, name, acceptance.HW_WORKSPACE_APP_FILE_STRORE_OBS_PATH)
}
