package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDesktops_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()

		all = "data.huaweicloud_workspace_desktops.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByDesktopId   = "data.huaweicloud_workspace_desktops.filter_by_desktop_id"
		dcFilterByDesktopId = acceptance.InitDataSourceCheck(filterByDesktopId)

		filterByName   = "data.huaweicloud_workspace_desktops.filter_by_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByUserName   = "data.huaweicloud_workspace_desktops.filter_by_user_name"
		dcFilterByUserName = acceptance.InitDataSourceCheck(filterByUserName)

		filterByTags   = "data.huaweicloud_workspace_desktops.filter_by_tags"
		dcFilterByTags = acceptance.InitDataSourceCheck(filterByTags)

		filterByImageId   = "data.huaweicloud_workspace_desktops.filter_by_image_id"
		dcFilterByImageId = acceptance.InitDataSourceCheck(filterByImageId)

		filterByEnterpriseProjectId   = "data.huaweicloud_workspace_desktops.filter_by_eps_id"
		dcFilterByEnterpriseProjectId = acceptance.InitDataSourceCheck(filterByEnterpriseProjectId)

		filterBySubnetId   = "data.huaweicloud_workspace_desktops.filter_by_subnet_id"
		dcFilterBySubnetId = acceptance.InitDataSourceCheck(filterBySubnetId)

		filterByStatus   = "data.huaweicloud_workspace_desktops.filter_by_status"
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
				Config: testAccDataDesktops_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "desktops.#"),
					resource.TestCheckResourceAttr(all, "desktops.0.is_support_internet", "false"),

					dcFilterByDesktopId.CheckResourceExists(),
					resource.TestCheckOutput("is_desktop_id_filter_useful", "true"),

					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcFilterByUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),

					dcFilterByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),

					dcFilterByImageId.CheckResourceExists(),
					resource.TestCheckOutput("is_image_id_filter_useful", "true"),

					dcFilterByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),

					dcFilterBySubnetId.CheckResourceExists(),
					resource.TestCheckOutput("is_subnet_id_filter_useful", "true"),

					dcFilterByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataDesktops_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

data "huaweicloud_workspace_flavors" "test" {
  os_type = "Windows"
}

locals {
  cpu_flavors = [for v in data.huaweicloud_workspace_flavors.test.flavors : v if v.is_gpu == false]
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id)
  image_type        = "market"
  image_id          = "%[1]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = data.huaweicloud_workspace_service.test.vpc_id
  security_groups   = [
    data.huaweicloud_workspace_service.test.desktop_security_group.0.id,
  ]

  nic {
    network_id = data.huaweicloud_workspace_service.test.network_ids[0]
  }

  name        = "%[2]s"
  user_name   = "user-%[2]s"
  user_email  = "terraform@example.com"
  user_group  = "administrators"
  delete_user = true

  root_volume {
    type = "SAS"
    size = 80
  }

  data_volume {
    type = "SAS"
    size = 50
  }
}
`, acceptance.HW_WORKSPACE_DESKTOP_POOL_IMAGE_ID, name)
}

func testAccDataDesktops_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_desktops" "all" {
  depends_on = [
    huaweicloud_workspace_desktop.test
  ]
}

# Filter by 'desktop_id' parameter.
locals {
  desktop_id = data.huaweicloud_workspace_desktops.all.desktops[0].id
}

data "huaweicloud_workspace_desktops" "filter_by_desktop_id" {
  desktop_id = local.desktop_id
}

output "is_desktop_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_desktop_id.desktops) == 1 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_desktop_id.desktops[*].id : v == local.desktop_id]
  )
}

# Filter by 'name' parameter.
locals {
  name = data.huaweicloud_workspace_desktops.all.desktops[0].name
}

data "huaweicloud_workspace_desktops" "filter_by_name" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_name.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_name.desktops[*].name : v == local.name]
  )
}

# Filter by 'user_name' parameter.
locals {
  user_name = data.huaweicloud_workspace_desktops.all.desktops[0].attach_user_infos[0].user_name
}

data "huaweicloud_workspace_desktops" "filter_by_user_name" {
  user_name = local.user_name
}

output "is_user_name_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_user_name.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_user_name.desktops[*].attach_user_infos[0].user_name : v == local.user_name]
  )
}

# Filter by 'tags' parameter.
locals {
  tags = data.huaweicloud_workspace_desktops.all.desktops[0].tags
}

data "huaweicloud_workspace_desktops" "filter_by_tags" {
  tags = local.tags
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_workspace_desktops.filter_by_tags.desktops[*].tags : v == local.tags
  ]
}

output "is_tags_filter_useful" {
  value = alltrue(local.tags_filter_result) && length(local.tags_filter_result) > 0
}

# Filter by 'image_id' parameter.
locals {
  image_id = data.huaweicloud_workspace_desktops.all.desktops[0].image_id
}

data "huaweicloud_workspace_desktops" "filter_by_image_id" {
  image_id = local.image_id
}

output "is_image_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_image_id.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_image_id.desktops[*].image_id : v == local.image_id]
  )
}

# Filter by 'enterprise_project_id' parameter.
locals {
  eps_id = data.huaweicloud_workspace_desktops.all.desktops[0].enterprise_project_id
}

data "huaweicloud_workspace_desktops" "filter_by_eps_id" {
  enterprise_project_id = local.eps_id
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_eps_id.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_eps_id.desktops[*].enterprise_project_id : v == local.eps_id]
  )
}

# Filter by 'subnet_id' parameter.
locals {
  subnet_id = data.huaweicloud_workspace_desktops.all.desktops[0].subnet_id
}

data "huaweicloud_workspace_desktops" "filter_by_subnet_id" {
  subnet_id = local.subnet_id
}

output "is_subnet_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_subnet_id.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_subnet_id.desktops[*].subnet_id : v == local.subnet_id]
  )
}

# Filter by 'status' parameter.
locals {
  status = data.huaweicloud_workspace_desktops.all.desktops[0].status
}

data "huaweicloud_workspace_desktops" "filter_by_status" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_status.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_status.desktops[*].status : v == local.status]
  )
}
`, testAccDataDesktops_base(name))
}
