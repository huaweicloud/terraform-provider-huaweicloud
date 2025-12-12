package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRAMPermission_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ram_permission.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMPermission(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRAMPermission_basic(acceptance.HW_RAM_PERMISSION_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permission.#"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.is_resource_type_default"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.permission_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.permission_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.default_version"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.status"),
				),
			},
			{
				Config: testDataSourceRAMPermission_with_version_basic(acceptance.HW_RAM_PERMISSION_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permission.#"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.is_resource_type_default"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.permission_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.permission_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.default_version"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "permission.0.status"),
				),
			},
		},
	})
}

func testDataSourceRAMPermission_basic(permissionId string) string {
	return fmt.Sprintf(`
data "huaweicloud_ram_resource_permission" "test" {
	permission_id = "%[1]s"
}
`, permissionId)
}

func testDataSourceRAMPermission_with_version_basic(permissionId string) string {
	return fmt.Sprintf(`
data "huaweicloud_ram_resource_permission" "test" {
	permission_id      = "%[1]s"
	permission_version = "1"
}
`, permissionId)
}
