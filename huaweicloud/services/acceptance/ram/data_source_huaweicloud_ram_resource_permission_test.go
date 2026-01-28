package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceResourcePermission_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_ram_resource_permission.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMPermission(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceResourcePermission_basic(acceptance.HW_RAM_PERMISSION_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permission.#"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.content"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.is_resource_type_default"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.permission_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.permission_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.default_version"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.status"),

					resource.TestCheckOutput("is_permission_version_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceResourcePermission_basic(permissionId string) string {
	return fmt.Sprintf(`
data "huaweicloud_ram_resource_permission" "test" {
  permission_id = "%[1]s"
}

# Filter using permission_version.
locals {
  permission_version = data.huaweicloud_ram_resource_permission.test.permission[0].version
}

data "huaweicloud_ram_resource_permission" "permission_version_filter" {
  permission_id      = "%[1]s"
  permission_version = local.permission_version
}

output "is_permission_version_filter_useful" {
  value = length(data.huaweicloud_ram_resource_permission.permission_version_filter.permission) > 0 && alltrue(
    [for v in data.huaweicloud_ram_resource_permission.permission_version_filter.permission[*].version : v == local.permission_version]
  )
}
`, permissionId)
}
