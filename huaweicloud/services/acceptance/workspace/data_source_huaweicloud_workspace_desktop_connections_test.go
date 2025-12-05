package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktopConnections_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_desktop_connections.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByUserName   = "data.huaweicloud_workspace_desktop_connections.filter_by_user_name"
		dcFilterByUserName = acceptance.InitDataSourceCheck(filterByUserName)

		filterByStatus   = "data.huaweicloud_workspace_desktop_connections.filter_by_status"
		dcFilterByStatus = acceptance.InitDataSourceCheck(filterByStatus)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceDesktopPoolImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDesktopConnections_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "desktop_connections.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "desktop_connections.0.id"),
					resource.TestCheckResourceAttrSet(all, "desktop_connections.0.name"),
					resource.TestCheckResourceAttrSet(all, "desktop_connections.0.connect_status"),
					resource.TestMatchResourceAttr(all, "desktop_connections.0.attach_users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "desktop_connections.0.attach_users.0.name"),
					resource.TestCheckResourceAttrSet(all, "desktop_connections.0.attach_users.0.user_group"),
					resource.TestCheckResourceAttrSet(all, "desktop_connections.0.attach_users.0.type"),
					// Filter by 'user_names' parameter.
					dcFilterByUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),
					// Filter by 'connect_status' parameter.
					dcFilterByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDesktopConnections_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_desktop" "test" {
  count = 2

  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id, "NOT_FOUND")
  image_type        = "market"
  image_id          = "%[1]s"
  availability_zone = try(data.huaweicloud_availability_zones.test.names[0], "NOT_FOUND")
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    try(data.huaweicloud_workspace_service.test.desktop_security_group[0].id, "NOT_FOUND"),
    try(data.huaweicloud_workspace_service.test.infrastructure_security_group[0].id, "NOT_FOUND")
  ]

  nic {
    network_id = try(data.huaweicloud_workspace_service.test.network_ids[0], "NOT_FOUND")
  }

  name       = "%[2]s-${count.index}"
  user_name  = "user-%[2]s-${count.index}"
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
`, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID, name)
}

func testAccDataDesktopConnections_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_desktop_connections" "all" {
  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  user_name      = data.huaweicloud_workspace_desktop_connections.all.desktop_connections[0].attach_users[0].name
}

# Filter by 'user_names' parameter.
data "huaweicloud_workspace_desktop_connections" "filter_by_user_name" {
  user_names = [local.user_name]

  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

locals {
  user_name_filter_result = [
    for v in data.huaweicloud_workspace_desktop_connections.filter_by_user_name.desktop_connections : contains([
      for user in v.attach_users : user.name
    ], local.user_name)
  ]
}

output "is_user_name_filter_useful" {
  value = alltrue(local.user_name_filter_result)
}

# Filter by 'connect_status' parameter.
locals {
  connect_status = data.huaweicloud_workspace_desktop_connections.all.desktop_connections[0].connect_status
}

data "huaweicloud_workspace_desktop_connections" "filter_by_status" {
  connect_status = local.connect_status

  depends_on = [
    huaweicloud_workspace_desktop.test,
  ]
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_workspace_desktop_connections.filter_by_status.desktop_connections : 
    v.connect_status == local.connect_status
  ]
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result)
}
`, testAccDataDesktopConnections_base(name))
}
