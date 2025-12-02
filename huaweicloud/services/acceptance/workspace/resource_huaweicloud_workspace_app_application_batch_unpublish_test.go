package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAppApplicationBatchUnpublish_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		rName          = "huaweicloud_workspace_app_application_batch_publish.test"
		unpublishRName = "huaweicloud_workspace_app_application_batch_unpublish.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccAppApplicationBatchUnpublish_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "app_group_id"),
					resource.TestCheckResourceAttr(rName, "applications.#", "2"),
					// Check attribute(s).
					resource.TestCheckResourceAttrSet(rName, "applications.0.id"),
				),
			},
			{
				Config: testAccAppApplicationBatchUnpublish_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(unpublishRName, "app_group_id"),
					resource.TestCheckResourceAttr(unpublishRName, "application_ids.#", "2"),
				),
			},
		},
	})
}

func testAccAppApplicationBatchUnpublish_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  app_type         = "COMMON_APP"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 90
  is_vdi           = true
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
}

# A server group can only be associated with one application group.
resource "huaweicloud_workspace_app_group" "test" {
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  name            = "%[1]s"
}
`, name,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testAccAppApplicationBatchUnpublish_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_application_batch_publish" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id

  applications {
    name           = "%[2]s-1"
    execute_path   = "C:\\Program Files\\7-Zip\\7zFM.exe"
    source_type    = 3
    sandbox_enable = true
    version        = "1.0.0"
    publisher      = "terraform"
    work_path      = "C:\\Program Files\\Sandboxie"
    command_param  = "/box:DefaultBox \"7567\""
    description    = "Created APP by script"
    icon_path      = "C:\\Program Files\\7-Zip\\7zFM.exe"
    icon_index     = 0
  }
  applications {
    name         = "%[2]s-2"
    execute_path = "C:\\Program Files\\7-Zip\\7zFM.exe"
    source_type  = 3
  }
}
`, testAccAppApplicationBatchUnpublish_base(name), name)
}

func testAccAppApplicationBatchUnpublish_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_app_application_batch_unpublish" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id

  application_ids = huaweicloud_workspace_app_application_batch_publish.test.applications[*].id
}
`, testAccAppApplicationBatchUnpublish_basic_step1(name))
}
