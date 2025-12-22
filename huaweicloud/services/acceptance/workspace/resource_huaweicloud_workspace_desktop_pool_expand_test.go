package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDesktopPoolExpand_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceNameWithDash()
		rName = "huaweicloud_workspace_desktop_pool_expand.test"
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		// This resource is a one-time action resource and there is no logic in the delete method.
		// lintignore:AT001
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDesktopPoolExpand_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(rName, "pool_id"),
					resource.TestCheckResourceAttrSet(rName, "size"),
				),
			},
		},
	})
}

func testAccDesktopPoolExpand_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  os_type           = "Windows"
}

data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_desktop_pool" "test" {
  name       = "%[1]s"
  type       = "STATIC"
  size       = 1
  product_id = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "")
  image_type = "gold"
  image_id   = "%[2]s"
  subnet_ids = try(slice(data.huaweicloud_workspace_service.test.network_ids, 0, 1), [])

  root_volume {
    type = "SAS"
    size = 80
  }

  lifecycle {
    ignore_changes = [size]
  }
}
`, name, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID)
}

func testAccDesktopPoolExpand_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_workspace_desktop_pool_expand" "test" {
  pool_id = huaweicloud_workspace_desktop_pool.test.id
  size    = 1
}
`, testAccDesktopPoolExpand_base(name))
}
