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

func TestAccResourceAppWarehouseApplication_basic(t *testing.T) {
	var (
		application  interface{}
		resourceName = "huaweicloud_workspace_app_warehouse_application.test"
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
			acceptance.TestAccPreCheckWorkspaceAppFileName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		CheckDestroy: rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWarehouseApplication_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "category", "OTHER"),
					resource.TestCheckResourceAttr(resourceName, "os_type", "Windows"),
					resource.TestCheckResourceAttr(resourceName, "version", "1.0"),
					resource.TestCheckResourceAttr(resourceName, "version_name", "terraform"),
					resource.TestCheckResourceAttr(resourceName, "file_store_path", acceptance.HW_WORKSPACE_APP_FILE_NAME),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by script"),
					resource.TestCheckResourceAttrSet(resourceName, "record_id"),
				),
			},
			{
				Config: testAccResourceWarehouseApplication_basic_step2(updateName),
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

func executionFileUploadResourcesConfig() string {
	return fmt.Sprintf(`
variable "script_content" {
  type    = string
  default = <<EOT
def main():  
    print("Hello, World!")  

if __name__ == "__main__":  
    main()
EOT
}

data "huaweicloud_identity_projects" "test" {
  name = "%[1]s"
}

data "huaweicloud_obs_buckets" "test" {
  bucket = format("wks-app-%%s", data.huaweicloud_identity_projects.test.projects[0].id)
}

resource "null_resource" "test" {
  depends_on = [data.huaweicloud_obs_buckets.test]

  provisioner "local-exec" {
    command = "echo '${var.script_content}' >> %[2]s"
  }
  provisioner "local-exec" {
    command = "rm %[2]s"
    when    = destroy
  }
}

resource "huaweicloud_obs_bucket_object" "test" {
  depends_on = [null_resource.test]

  bucket       = data.huaweicloud_obs_buckets.test.buckets[0].bucket
  key          = "%[2]s"
  source       = abspath("%[2]s")
  content_type = "application/x-msdownload"
}`, acceptance.HW_REGION_NAME, acceptance.HW_WORKSPACE_APP_FILE_NAME)
}

func testAccResourceWarehouseApplication_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_warehouse_application" "test" {
  name            = "%[2]s"
  category        = "OTHER"
  os_type         = "Windows"
  version         = "1.0"
  version_name    = "terraform"
  file_store_path = "%[3]s"
  description     = "Created by script"
}
`, executionFileUploadResourcesConfig(), name, acceptance.HW_WORKSPACE_APP_FILE_NAME)
}

func testAccResourceWarehouseApplication_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_warehouse_application" "test" {
  name            = "%[2]s"
  category        = "PRODUCTIVITY_AND_COLLABORATION"
  os_type         = "Linux"
  version         = "2.0"
  version_name    = "terraform_update"
  file_store_path = "%[3]s"
  description     = "Updated by script"
}
`, executionFileUploadResourcesConfig(), name, acceptance.HW_WORKSPACE_APP_FILE_NAME)
}
