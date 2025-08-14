package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceAppApplicationBatchAttach_basic(t *testing.T) {
	var (
		randUUID, _  = uuid.GenerateUUID()
		name         = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_workspace_app_application_batch_attach.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
			acceptance.TestAccPreCheckWorkspaceAppImageSpecCode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"null": {
				Source:            "hashicorp/null",
				VersionConstraint: "3.2.1",
			},
		},
		// This resource is a one-time action resource with no deletion logic.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceAppApplicationBatchAttach_expectError(randUUID),
				ExpectError: regexp.MustCompile(fmt.Sprintf("'%s' is a non-existing cloud application server", randUUID)),
			},
			{
				Config: testAccResourceAppApplicationBatchAttach_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "record_ids.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(resourceName, "uri"),
				),
			},
		},
	})
}

func testAccResourceAppApplicationBatchAttach_expectError(randomId string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_application_batch_attach" "expectError" {
  server_id  = "%[1]s"
  record_ids = ["%[1]s"]
}
`, randomId)
}

func testAccResourceAppApplicationBatchAttach_base(name string) string {
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

variable "application_file_name" {
  type    = string
  default = "%[1]s.exe"
}

data "huaweicloud_identity_projects" "test" {
  name = "%[2]s"
}

data "huaweicloud_obs_buckets" "test" {
  bucket = format("wks-app-%%s", data.huaweicloud_identity_projects.test.projects[0].id)
}

resource "null_resource" "test" {
  depends_on = [data.huaweicloud_obs_buckets.test]

  triggers = {
    file_name = var.application_file_name
  }

  provisioner "local-exec" {
    command = "echo '${var.script_content}' >> ${var.application_file_name}"
  }
  provisioner "local-exec" {
    command = "rm ${self.triggers.file_name}"
    when    = destroy
  }
}

resource "huaweicloud_obs_bucket_object" "test" {
  depends_on = [null_resource.test]

  bucket       = data.huaweicloud_obs_buckets.test.buckets[0].bucket
  key          = var.application_file_name
  source       = abspath(var.application_file_name)
  content_type = "application/x-msdownload"
}

resource "huaweicloud_workspace_app_warehouse_app" "test" {
  name            = "%[1]s"
  category        = "OTHER"
  os_type         = "Windows"
  version         = "1.0"
  version_name    = "terraform"
  file_store_path = var.application_file_name
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_user" "test" {
  name                   = "%[1]s"
  email                  = "%[1]s@example.com"
  password_never_expires = false
  disabled               = false
}

resource "huaweicloud_workspace_app_image_server" "test" {
  name                           = "%[1]s"
  flavor_id                      = "%[3]s"
  vpc_id                         = try(data.huaweicloud_workspace_service.test.vpc_id, null)
  subnet_id                      = try(data.huaweicloud_workspace_service.test.network_ids[0], null)
  image_id                       = "%[4]s"
  image_type                     = "gold"
  image_source_product_id        = "%[5]s"
  spec_code                      = "%[6]s"
  is_vdi                         = true
  is_delete_associated_resources = true

  authorize_accounts {
    account = huaweicloud_workspace_user.test.name
    type    = "USER"
  }

  root_volume {
    type = "SAS"
    size = 80
  }
}
`,
		name,
		acceptance.HW_REGION_NAME,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_SPEC_CODE,
	)
}

func testAccResourceAppApplicationBatchAttach_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_application_batch_attach" "test" {
  server_id  = huaweicloud_workspace_app_image_server.test.id
  record_ids = [huaweicloud_workspace_app_warehouse_app.test.record_id]
}
`, testAccResourceAppApplicationBatchAttach_base(name))
}
