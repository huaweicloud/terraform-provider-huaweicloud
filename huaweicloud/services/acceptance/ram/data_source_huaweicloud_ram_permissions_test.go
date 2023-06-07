package ram

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceRAMPermissions_basic(t *testing.T) {
	rName := "data.huaweicloud_ram_permissions.test"
	dc := acceptance.InitDataSourceCheck(rName)

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

					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceRAMPermissions_basic() string {
	return `
data "huaweicloud_ram_permissions" "test" {
}

data "huaweicloud_ram_permissions" "resource_type_filter" {
  resource_type = "vpc:subnets"
}
output "resource_type_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_ram_permissions.resource_type_filter.permissions[*].resource_type : v == "vpc:subnets"])
}

data "huaweicloud_ram_permissions" "name_filter" {
  name = "default vpc subnets statement"
}
output "name_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_ram_permissions.name_filter.permissions[*].name : v == "default vpc subnets statement"])
}
`
}
