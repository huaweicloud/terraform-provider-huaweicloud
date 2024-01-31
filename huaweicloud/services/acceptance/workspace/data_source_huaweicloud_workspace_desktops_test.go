package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDesktops_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_workspace_desktops.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byDesktopId   = "data.huaweicloud_workspace_desktops.filter_by_desktop_id"
		dcByDesktopId = acceptance.InitDataSourceCheck(byDesktopId)

		byName   = "data.huaweicloud_workspace_desktops.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byUserName = "data.huaweicloud_workspace_desktops.filter_by_user_name"
		dcUserName = acceptance.InitDataSourceCheck(byUserName)

		byTags = "data.huaweicloud_workspace_desktops.filter_by_tags"
		dcTags = acceptance.InitDataSourceCheck(byTags)

		byImageId = "data.huaweicloud_workspace_desktops.filter_by_image_id"
		dcImageId = acceptance.InitDataSourceCheck(byImageId)

		byEnterpriseProjectId = "data.huaweicloud_workspace_desktops.filter_by_eps_id"
		dcEnterpriseProjectId = acceptance.InitDataSourceCheck(byEnterpriseProjectId)

		bySubnetId = "data.huaweicloud_workspace_desktops.filter_by_subnet_id"
		dcSubnetId = acceptance.InitDataSourceCheck(bySubnetId)

		byStatus = "data.huaweicloud_workspace_desktops.filter_by_status"
		dcStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDesktops_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "desktops.#"),
					resource.TestCheckResourceAttr(dataSourceName, "desktops.0.is_support_internet", "false"),
					dcByDesktopId.CheckResourceExists(),
					resource.TestCheckOutput("is_desktop_id_filter_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),

					dcTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),

					dcImageId.CheckResourceExists(),
					resource.TestCheckOutput("is_image_id_filter_useful", "true"),

					dcEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),

					dcSubnetId.CheckResourceExists(),
					resource.TestCheckOutput("is_subnet_id_filter_useful", "true"),

					dcStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceDesktops_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s
		
data "huaweicloud_workspace_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  os_type           = "Windows"
}

resource "huaweicloud_workspace_desktop" "test" {
  flavor_id         = try(data.huaweicloud_workspace_flavors.test.flavors[0].id)
  image_type        = "market"
  image_id          = try(data.huaweicloud_images_images.test.images[0].id)
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  vpc_id            = huaweicloud_vpc.test.id
  security_groups   = [
    huaweicloud_workspace_service.test.desktop_security_group.0.id,
    huaweicloud_networking_secgroup.test.id,
  ]

  nic {
    network_id = huaweicloud_vpc_subnet.test.id
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

data "huaweicloud_workspace_desktops" "test" {
  depends_on = [
    huaweicloud_workspace_desktop.test
  ] 
}

// By desktop ID filter
locals {
  desktop_id = data.huaweicloud_workspace_desktops.test.desktops[0].id
}

data "huaweicloud_workspace_desktops" "filter_by_desktop_id" {
  desktop_id = local.desktop_id
}

output "is_desktop_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_desktop_id.desktops) == 1 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_desktop_id.desktops[*].id : v == local.desktop_id]
  )
}

// By desktop name filter
locals {
  name = data.huaweicloud_workspace_desktops.test.desktops[0].name
}

data "huaweicloud_workspace_desktops" "filter_by_name" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_name.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_name.desktops[*].name : v == local.name]
  )
}

// By user name filter
locals {
  user_name = data.huaweicloud_workspace_desktops.test.desktops[0].attach_user_infos[0].user_name
}

data "huaweicloud_workspace_desktops" "filter_by_user_name" {
  user_name = local.user_name
}

output "is_user_name_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_user_name.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_user_name.desktops[*].attach_user_infos[0].user_name : v == local.user_name]
  )
}

// By tags filter
locals {
  tags = data.huaweicloud_workspace_desktops.test.desktops[0].tags
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

// By image ID filter
locals {
  image_id = data.huaweicloud_workspace_desktops.test.desktops[0].image_id
}

data "huaweicloud_workspace_desktops" "filter_by_image_id" {
  image_id = local.image_id
}

output "is_image_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_image_id.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_image_id.desktops[*].image_id : v == local.image_id]
  )
}

// By enperprise project ID filter
locals {
  eps_id = data.huaweicloud_workspace_desktops.test.desktops[0].enterprise_project_id
}

data "huaweicloud_workspace_desktops" "filter_by_eps_id" {
  enterprise_project_id = local.eps_id
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_eps_id.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_eps_id.desktops[*].enterprise_project_id : v == local.eps_id]
  )
}

// By subnet ID filter
locals {
  subnet_id = data.huaweicloud_workspace_desktops.test.desktops[0].subnet_id
}

data "huaweicloud_workspace_desktops" "filter_by_subnet_id" {
  subnet_id = local.subnet_id
}

output "is_subnet_id_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_subnet_id.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_subnet_id.desktops[*].subnet_id : v == local.subnet_id]
  )
}

// By desktop status filter
locals {
  status = data.huaweicloud_workspace_desktops.test.desktops[0].status
}

data "huaweicloud_workspace_desktops" "filter_by_status" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_workspace_desktops.filter_by_status.desktops) > 0 && alltrue(
    [for v in data.huaweicloud_workspace_desktops.filter_by_status.desktops[*].status : v == local.status]
  )
}
`, testAccDesktop_base(rName), rName)
}
