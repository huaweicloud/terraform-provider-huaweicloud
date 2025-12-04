package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataUserDesktopPoolAssociations_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_user_desktop_pool_associations.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserDesktopPoolAssociations_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "associations.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "associations.0.user_id"),
					resource.TestMatchResourceAttr(all, "associations.0.desktop_pools.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "associations.0.desktop_pools.0.id"),
					resource.TestCheckResourceAttrSet(all, "associations.0.desktop_pools.0.name"),
					resource.TestCheckResourceAttrSet(all, "associations.0.desktop_pools.0.is_attached"),
					resource.TestCheckOutput("is_user_ids_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataUserDesktopPoolAssociations_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

locals {
  cpu_flavors = [for v in data.huaweicloud_workspace_flavors.test.flavors : v if v.is_gpu == false]
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_user" "test" {
  count = 2

  name  = "%[1]s${count.index}"
  email = "www.user${count.index}@test.com"
}

resource "huaweicloud_workspace_desktop_pool" "test" {
  count                         = 2
  name                          = "%[1]s-${count.index}"
  type                          = "DYNAMIC"
  size                          = 1
  product_id                    = try(local.cpu_flavors[0].id, "")
  image_type                    = "gold"
  image_id                      = "%[2]s"
  subnet_ids                    = data.huaweicloud_workspace_service.test.network_ids
  vpc_id                        = data.huaweicloud_workspace_service.test.vpc_id
  availability_zone             = data.huaweicloud_availability_zones.test.names[0]
  disconnected_retention_period = 10
  enable_autoscale              = false

  root_volume {
    type = "SAS"
    size = 80
  }

  security_groups {
    id = data.huaweicloud_workspace_service.test.desktop_security_group[0].id
  }

  dynamic "authorized_objects" {
    for_each = huaweicloud_workspace_user.test

    content {
      object_type = "USER"
      object_id   = authorized_objects.value.id
      object_name = authorized_objects.value.name
      user_group  = "administrators"
    }
  }
}
`, name, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID)
}

func testAccDataUserDesktopPoolAssociations_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  user_ids = huaweicloud_workspace_user.test[*].id
}

data "huaweicloud_workspace_user_desktop_pool_associations" "all" {
  user_ids = local.user_ids

  depends_on = [
    huaweicloud_workspace_desktop_pool.test
  ]
}

locals {
  query_result = [
    for v in data.huaweicloud_workspace_user_desktop_pool_associations.all.associations[*].user_id :
    contains(local.user_ids, v)
  ]
}

output "is_user_ids_filter_useful" {
  value = length(local.query_result) > 0 && alltrue(local.query_result)
}
`, testAccDataUserDesktopPoolAssociations_base(name))
}
