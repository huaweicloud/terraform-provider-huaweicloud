package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceRAMPermissions_basic(t *testing.T) {
	var (
		rName = "data.huaweicloud_ram_resource_permissions.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byResourceType   = "data.huaweicloud_ram_resource_permissions.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)

		byName   = "data.huaweicloud_ram_resource_permissions.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceRAMPermissions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.id"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.name"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.resource_type"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.is_resource_type_default"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "permissions.0.updated_at"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceRAMPermissions_basic() string {
	return `
data "huaweicloud_ram_resource_permissions" "test" {
}

# Filter by resource_type
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

# Filter by name
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
