package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppPublishableApplications_basic(t *testing.T) {
	var (
		notFound   = "data.huaweicloud_workspace_app_publishable_applications.not_found"
		dcNotFound = acceptance.InitDataSourceCheck(notFound)
		all        = "data.huaweicloud_workspace_app_publishable_applications.test"
		dcAll      = acceptance.InitDataSourceCheck(all)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataAppPublishableApplications_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcAll.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "group_images.#"),
					// Deprecated attribute `apps` is still supported, so do not remove it.
					resource.TestMatchResourceAttr(all, "apps.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "apps.0.name"),
					resource.TestCheckResourceAttrSet(all, "apps.0.execute_path"),
					resource.TestCheckResourceAttrSet(all, "apps.0.version"),
					resource.TestCheckResourceAttrSet(all, "apps.0.publisher"),
					resource.TestCheckResourceAttrSet(all, "apps.0.publishable"),
					resource.TestCheckResourceAttrSet(all, "apps.0.icon_path"),
					resource.TestCheckResourceAttrSet(all, "apps.0.icon_index"),
					resource.TestCheckResourceAttrSet(all, "apps.0.source_image_ids.#"),
					// The instead of the deprecated attribute `apps`, we should use the new attribute `applications`.
					resource.TestMatchResourceAttr(all, "applications.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "applications.0.name"),
					resource.TestCheckResourceAttrSet(all, "applications.0.execute_path"),
					resource.TestCheckResourceAttrSet(all, "applications.0.version"),
					resource.TestCheckResourceAttrSet(all, "applications.0.publisher"),
					resource.TestCheckResourceAttrSet(all, "applications.0.publishable"),
					resource.TestCheckResourceAttrSet(all, "applications.0.icon_path"),
					resource.TestCheckResourceAttrSet(all, "applications.0.icon_index"),
					// The `description`, `work_path` and `command_param` may be empty, so do not assert them.
				),
			},
			// If the incorrect test step is executed first, it will cause subsequent test steps failed.
			{
				Config: testDataAppPublishableApplications_notFound(),
				Check: resource.ComposeTestCheckFunc(
					dcNotFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(notFound, "group_images.#", "0"),
					// Deprecated attribute `apps` is still supported, so do not remove it.
					resource.TestCheckResourceAttr(notFound, "apps.#", "0"),
					// The instead of the deprecated attribute `apps`, we should use the new attribute `applications`.
					resource.TestCheckResourceAttr(notFound, "applications.#", "0"),
				),
			},
		},
	})
}

func testDataAppPublishableApplications_notFound() string {
	uuid, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_publishable_applications" "not_found" {
  app_group_id = "%s"
}
`, uuid)
}

func testDataAppPublishableApplications_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = true
  app_type         = "COMMON_APP"
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
}

resource "huaweicloud_workspace_app_server" "test" {
  name            = "%[1]s" 
  server_group_id = huaweicloud_workspace_app_server_group.test.id
  type            = "createApps"
  flavor_id       = huaweicloud_workspace_app_server_group.test.flavor_id

  root_volume {
    type = huaweicloud_workspace_app_server_group.test.system_disk_type
    size = huaweicloud_workspace_app_server_group.test.system_disk_size
  }

  vpc_id              = huaweicloud_workspace_app_server_group.test.vpc_id
  subnet_id           = huaweicloud_workspace_app_server_group.test.subnet_id
  update_access_agent = false
}

resource "huaweicloud_workspace_app_group" "test" {
  depends_on   = [huaweicloud_workspace_app_server.test]

  name            = "%[1]s"
  server_group_id = huaweicloud_workspace_app_server_group.test.id
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testDataAppPublishableApplications_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_workspace_app_publishable_applications" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id
}
`, testDataAppPublishableApplications_base(name))
}
