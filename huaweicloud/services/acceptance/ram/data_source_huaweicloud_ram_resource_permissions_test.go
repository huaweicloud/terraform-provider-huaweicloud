package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourcePermissions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_resource_permissions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byResourceType   = "data.huaweicloud_ram_resource_permissions.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)

		byPermissionType   = "data.huaweicloud_ram_resource_permissions.filter_by_permission_type"
		dcByPermissionType = acceptance.InitDataSourceCheck(byPermissionType)

		byName   = "data.huaweicloud_ram_resource_permissions.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourcePermissions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.is_resource_type_default"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.permission_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.permission_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.default_version"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.status"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),

					dcByPermissionType.CheckResourceExists(),
					resource.TestCheckOutput("is_permission_type_filter_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceResourcePermissions_basic() string {
	return `
data "huaweicloud_ram_resource_permissions" "test" {}

# Filter by resource_type.
locals {
  resource_type = data.huaweicloud_ram_resource_permissions.test.permissions[0].resource_type
}

data "huaweicloud_ram_resource_permissions" "filter_by_resource_type" {
  resource_type = local.resource_type
}

locals {
  resource_type_filter_result = [
    for v in data.huaweicloud_ram_resource_permissions.filter_by_resource_type.permissions[*].resource_type : v == local.resource_type
  ]
}

output "is_resource_type_filter_useful" {
  value = length(local.resource_type_filter_result) > 0 && alltrue(local.resource_type_filter_result)
}

# Filter by permission_type.
locals {
  permission_type = data.huaweicloud_ram_resource_permissions.test.permissions[0].permission_type
}

data "huaweicloud_ram_resource_permissions" "filter_by_permission_type" {
  permission_type = local.permission_type
}

locals {
  permission_type_filter_result = [
    for v in data.huaweicloud_ram_resource_permissions.filter_by_permission_type.permissions[*].permission_type : v == local.permission_type
  ]
}

output "is_permission_type_filter_useful" {
  value = length(local.permission_type_filter_result) > 0 && alltrue(local.permission_type_filter_result)
}

# Filter by name.
locals {
  name = data.huaweicloud_ram_resource_permissions.test.permissions[0].name
}

data "huaweicloud_ram_resource_permissions" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_ram_resource_permissions.filter_by_name.permissions[*].name : v == local.name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}
`
}
