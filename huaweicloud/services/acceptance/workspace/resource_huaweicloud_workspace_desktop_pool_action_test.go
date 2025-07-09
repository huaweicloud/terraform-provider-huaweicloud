package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopPoolAction_basic(t *testing.T) {
	var (
		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_workspace_desktop_pool_action.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopPoolAction_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "FAIL"),
				),
			},
		},
	})
}

func testAccDesktopPoolAction_base(name string) string {
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

resource "huaweicloud_workspace_desktop_pool" "test" {
  name                          = "%[1]s"
  type                          = "DYNAMIC"
  size                          = 1
  product_id                    = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type                    = "gold"
  image_id                      = try(data.huaweicloud_images_images.test.images[0].id, "NOT_FOUND")
  subnet_ids                    = data.huaweicloud_workspace_service.test.network_ids
  vpc_id                        = data.huaweicloud_workspace_service.test.vpc_id
  availability_zone             = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  disconnected_retention_period = 10
  enable_autoscale              = true
  in_maintenance_mode           = true

  security_groups {
    id = data.huaweicloud_workspace_service.test.desktop_security_group.0.id
  }
  security_groups {
    id = data.huaweicloud_workspace_service.test.infrastructure_security_group.0.id
  }

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volumes {
    type = "SAS"
    size = 20
  }

  autoscale_policy {
    autoscale_type    = "AUTO_CREATED"
    min_idle          = 1
    max_auto_created  = 2
    once_auto_created = 1
  }
}
`, name)
}

func testAccDesktopPoolAction_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_desktop_pool_action" "test" {
  pool_id = huaweicloud_workspace_desktop_pool.test.id
  op_type = "os-start"
  type    = "SOFT"
}
`, testAccDesktopPoolAction_base(name))
}
