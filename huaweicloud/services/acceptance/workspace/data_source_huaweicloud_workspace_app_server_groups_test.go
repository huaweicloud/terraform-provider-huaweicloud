package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppServerGroups_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_workspace_app_server_groups.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byId   = "data.huaweicloud_workspace_app_server_groups.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_workspace_app_server_groups.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byAppType   = "data.huaweicloud_workspace_app_server_groups.filter_by_app_type"
		dcByAppType = acceptance.InitDataSourceCheck(byAppType)

		byTags   = "data.huaweicloud_workspace_app_server_groups.filter_by_tags"
		dcByTags = acceptance.InitDataSourceCheck(byTags)

		byEpsId   = "data.huaweicloud_workspace_app_server_groups.filter_by_eps_id"
		dcByEpsId = acceptance.InitDataSourceCheck(byEpsId)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppServerGroups_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.id"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.name", name),
					resource.TestCheckResourceAttr(byId, "server_groups.0.description", "terraform script"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.app_type", "COMMON_APP"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.image_id", acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID),
					resource.TestCheckResourceAttr(byId, "server_groups.0.os_type", "Windows"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.system_disk_type", "SAS"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.system_disk_size", "80"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.is_vdi", "true"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.storage_mount_policy", "ANY"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.enterprise_project_id", "0"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.server_group_status", "true"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.site_type", "CENTER"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.app_server_flavor_count", "0"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.app_server_count", "0"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.app_group_count", "0"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.site_id"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.image_name"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.create_time"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.update_time"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.subnet_name"),
					// Check product_info fields
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.#", "1"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.flavor_id", "s6.xlarge.2"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.type", "BASE"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.architecture", "x86"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.cpu", "4"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.memory", "8192"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.system_disk_type", "SAS"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.system_disk_size", "80"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.charge_mode", "1"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.contain_data_disk", "false"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.resource_type", "hws.resource.type.workspace.appstream"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.cloud_service_type", "hws.service.type.vdi"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.volume_product_type", "workspace"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.sessions", "2"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.product_info.0.status", "abandon"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.product_info.0.product_id"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.product_info.0.descriptions"),
					resource.TestCheckResourceAttrSet(byId, "server_groups.0.product_info.0.cond_operation_az"),
					// Check tags fields
					resource.TestCheckResourceAttr(byId, "server_groups.0.tags.#", "1"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.tags.0.key", "key1"),
					resource.TestCheckResourceAttr(byId, "server_groups.0.tags.0.value", "value1"),
					// Check query parameters
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByAppType.CheckResourceExists(),
					resource.TestCheckOutput("is_app_type_filter_useful", "true"),
					dcByTags.CheckResourceExists(),
					resource.TestCheckOutput("is_tags_filter_useful", "true"),
					dcByEpsId.CheckResourceExists(),
					resource.TestCheckOutput("is_eps_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataAppServerGroups_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  description      = "terraform script"
  os_type          = "Windows"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 80
  is_vdi           = true
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
  
  tags = {
    key1 = "value1"
  }
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testAccDataAppServerGroups_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_app_server_groups" "test" {
  depends_on = [
    huaweicloud_workspace_app_server_group.test
  ]
}

# Filter by ID
locals {
  group_id = huaweicloud_workspace_app_server_group.test.id
}

data "huaweicloud_workspace_app_server_groups" "filter_by_id" {
  depends_on = [
    huaweicloud_workspace_app_server_group.test
  ]

  server_group_id = local.group_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_workspace_app_server_groups.filter_by_id.server_groups[*].id : v == local.group_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  group_name = huaweicloud_workspace_app_server_group.test.name
}

data "huaweicloud_workspace_app_server_groups" "filter_by_name" {
  depends_on = [
    huaweicloud_workspace_app_server_group.test
  ]

  server_group_name = local.group_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_workspace_app_server_groups.filter_by_name.server_groups[*].name : v == local.group_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by app_type
locals {
  app_type = huaweicloud_workspace_app_server_group.test.app_type
}

data "huaweicloud_workspace_app_server_groups" "filter_by_app_type" {
  depends_on = [
    huaweicloud_workspace_app_server_group.test
  ]

  app_type = local.app_type
}

locals {
  app_type_filter_result = [
    for v in data.huaweicloud_workspace_app_server_groups.filter_by_app_type.server_groups[*].app_type : v == local.app_type
  ]
}

output "is_app_type_filter_useful" {
  value = length(local.app_type_filter_result) > 0 && alltrue(local.app_type_filter_result)
}

# Filter by tags
locals {
  tags_str = "key1=value1"
}

data "huaweicloud_workspace_app_server_groups" "filter_by_tags" {
  depends_on = [
    huaweicloud_workspace_app_server_group.test
  ]

  tags = local.tags_str
}

locals {
  tags_filter_result = [
    for v in data.huaweicloud_workspace_app_server_groups.filter_by_tags.server_groups[*].tags : 
      v[0].key == "key1" && v[0].value == "value1"
  ]
}

output "is_tags_filter_useful" {
  value = length(local.tags_filter_result) > 0 && alltrue(local.tags_filter_result)
}

# Filter by enterprise_project_id
locals {
  eps_id = huaweicloud_workspace_app_server_group.test.enterprise_project_id
}

data "huaweicloud_workspace_app_server_groups" "filter_by_eps_id" {
  depends_on = [
    huaweicloud_workspace_app_server_group.test
  ]

  enterprise_project_id = local.eps_id
}

locals {
  eps_id_filter_result = [
    for v in data.huaweicloud_workspace_app_server_groups.filter_by_eps_id.server_groups[*].enterprise_project_id : 
      v == local.eps_id
  ]
}

output "is_eps_id_filter_useful" {
  value = length(local.eps_id_filter_result) > 0 && alltrue(local.eps_id_filter_result)
}
`, testAccDataAppServerGroups_base(name))
}
