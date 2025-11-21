package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRAMListPermissionVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ram_permission_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMPermissionVersions(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRAMListPermissionVersions_basic(acceptance.HW_RAM_PERMISSION_ID),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.default_version"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.is_resource_type_default"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.permission_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.permission_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "permissions.0.version"),
				),
			},
		},
	})
}

func testDataSourceRAMListPermissionVersions_basic(permissionId string) string {
	return fmt.Sprintf(`
data "huaweicloud_ram_permission_versions" "test" {
	permission_id = "%[1]s"
}
`, permissionId)
}
