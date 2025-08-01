package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopSysprepDataSource_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		dcName = "data.huaweicloud_workspace_desktop_sysprep.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopSysprepDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "sysprep_info.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "sysprep_info.0.sysprep_version"),
					resource.TestCheckResourceAttrSet(dcName, "sysprep_info.0.support_create_image"),
				),
			},
		},
	})
}

func testAccDesktopSysprepDataSource_base(name string) string {
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
  image_id          = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id

  security_groups = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"),
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND"),
  ]
  
  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name       = "%[1]s"
  user_name  = "%[1]s-user"
  user_email = "terraform@example.com"
  user_group = "administrators"

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }

  tags = {
    foo   = "bar"
    owner = "terraform"
  }

  email_notification = true
  delete_user        = true
}
`, name)
}

func testAccDesktopSysprepDataSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_desktop_sysprep" "test" {
  desktop_id = huaweicloud_workspace_desktop.test.id
}
`, testAccDesktopSysprepDataSource_base(name))
}
