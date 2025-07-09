package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopRemoteConsoleDataSource_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_desktop_remote_console.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopRemoteConsoleDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "remote_console.#", "1"),
					resource.TestCheckResourceAttrSet(dcName, "remote_console.0.type"),
					resource.TestCheckResourceAttrSet(dcName, "remote_console.0.url"),
					resource.TestCheckResourceAttrSet(dcName, "remote_console.0.protocol"),
				),
			},
		},
	})
}

func testAccDesktopRemoteConsoleDataSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_desktop_remote_console" "test" {
  desktop_id = huaweicloud_workspace_desktop.test.id

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}
`, testAccDesktopRemoteConsoleDataSource_base(name))
}

func testAccDesktopRemoteConsoleDataSource_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_images" "test" {
  name_regex = "WORKSPACE"
  visibility = "market"
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[1].id, "NOT_FOUND")
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"),
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name       = "%[1]s"
  user_name  = "user-%[1]s"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 10
  }

  email_notification = true
  delete_user        = true
}
`, name)
}
