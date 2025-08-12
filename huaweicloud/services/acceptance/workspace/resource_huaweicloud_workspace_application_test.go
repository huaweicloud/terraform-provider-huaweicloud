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

func getApplicationFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("workspace", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	return workspace.GetApplicationById(client, state.Primary.ID)
}

func TestAccResourceApplication_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		withTest    = "huaweicloud_workspace_application.test"
		application interface{}
		rcWithTest  = acceptance.InitResourceCheck(withTest, &application, getApplicationFunc)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckProjectID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			rcWithTest.CheckResourceDestroy(),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceApplication_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithTest.CheckResourceExists(),
					resource.TestCheckResourceAttr(withTest, "name", name),
					resource.TestCheckResourceAttr(withTest, "version", "1.0.0"),
					resource.TestCheckResourceAttr(withTest, "description", "Created by terraform script"),
					resource.TestCheckResourceAttr(withTest, "authorization_type", "ALL_USER"),
					resource.TestCheckResourceAttr(withTest, "install_type", "QUIET_INSTALL"),
					resource.TestCheckResourceAttr(withTest, "support_os", "Windows"),
					resource.TestCheckResourceAttr(withTest, "catalog_id", "1"),
					resource.TestCheckResourceAttr(withTest, "application_file_store.#", "1"),
					resource.TestCheckResourceAttr(withTest, "application_file_store.0.store_type", "LINK"),
					resource.TestCheckResourceAttr(withTest, "application_file_store.0.file_link", "https://www.huaweicloud.com/TerraformTest.msi"),
				),
			},
			{
				Config: testAccResourceApplication_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rcWithTest.CheckResourceExists(),
					resource.TestCheckResourceAttr(withTest, "support_os", "Linux"),
					resource.TestCheckResourceAttr(withTest, "catalog_id", "2"),
					resource.TestCheckResourceAttr(withTest, "application_file_store.#", "1"),
					resource.TestCheckResourceAttr(withTest, "application_file_store.0.store_type", "OBS"),
					resource.TestCheckResourceAttr(withTest, "application_file_store.0.bucket_store.#", "1"),
					resource.TestCheckResourceAttr(withTest, "application_file_store.0.bucket_store.0.bucket_file_path", "dir1/TerraformTest.apk"),
				),
			},
			{
				ResourceName:      withTest,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccResourceApplication_base = `
data "huaweicloud_workspace_application_catalogs" "test" {}
`

func testAccResourceApplication_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_application" "test" {
  name               = "%[2]s"
  version            = "1.0.0"
  description        = "Created by terraform script"
  authorization_type = "ALL_USER"
  install_type       = "QUIET_INSTALL"
  support_os         = "Windows"
  catalog_id         = try(data.huaweicloud_workspace_application_catalogs.test.catalogs[0].id, "NOT_FOUND")

  application_file_store {
    store_type = "LINK"
    file_link  = "https://www.huaweicloud.com/TerraformTest.msi"
  }
}
`, testAccResourceApplication_base, name)
}

func testAccResourceApplication_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_application" "test" {
  name               = "%[2]s"
  version            = "1.0.0"
  description        = "created by terraform script."
  authorization_type = "ALL_USER"
  install_type       = "QUIET_INSTALL"
  support_os         = "Linux"
  catalog_id         = try(data.huaweicloud_workspace_application_catalogs.test.catalogs[1].id, "NOT_FOUND")

  application_file_store {
    store_type = "OBS"
	
    bucket_store {
      bucket_name      = "app-center-%[3]s"
      bucket_file_path = "dir1/TerraformTest.apk"
    }
  }
}
`, testAccResourceApplication_base, name, acceptance.HW_PROJECT_ID)
}
