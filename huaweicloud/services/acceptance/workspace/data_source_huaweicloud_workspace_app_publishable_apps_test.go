package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAppPublishableApps_basic(t *testing.T) {
	var (
		notFounddataSource = "data.huaweicloud_workspace_app_publishable_apps.not_found"
		notFound           = acceptance.InitDataSourceCheck(notFounddataSource)
		dataSource         = "data.huaweicloud_workspace_app_publishable_apps.test"
		rName              = acceptance.RandomAccResourceName()
		dc                 = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAppPublishableApps_notFound(),
				Check: resource.ComposeTestCheckFunc(
					notFound.CheckResourceExists(),
					resource.TestCheckResourceAttr(notFounddataSource, "group_images.#", "0"),
					resource.TestCheckResourceAttr(notFounddataSource, "apps.#", "0"),
				),
			},
			{
				Config: testDataSourceAppPublishableApps_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "group_images.#"),
					resource.TestMatchResourceAttr(dataSource, "apps.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.execute_path"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.publisher"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.publishable"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.icon_path"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.icon_index"),
					resource.TestCheckResourceAttrSet(dataSource, "apps.0.source_image_ids.#"),
					// The `description`, `work_path` and `command_param` may be empty, so do not assert them.
				),
			},
		},
	})
}

func testDataSourceAppPublishableApps_notFound() string {
	uuid, _ := uuid.GenerateUUID()
	return fmt.Sprintf(`
data "huaweicloud_workspace_app_publishable_apps" "not_found" {
  app_group_id = "%s"
}
`, uuid)
}

func testDataSourceAppPublishableApps_base(name string) string {
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

func testDataSourceAppPublishableApps_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_workspace_app_publishable_apps" "test" {
  app_group_id = huaweicloud_workspace_app_group.test.id
}
`, testDataSourceAppPublishableApps_base(name))
}
