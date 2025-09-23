package ram

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssociatedPermissions_basic(t *testing.T) {
	var (
		rName = acceptance.RandomAccResourceName()

		dataSource = "data.huaweicloud_ram_resource_share_associated_permissions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byPermissionName   = "data.huaweicloud_ram_resource_share_associated_permissions.filter_by_permission_name"
		dcByPermissionName = acceptance.InitDataSourceCheck(byPermissionName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRAMResourceShare(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssociatedPermissions_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "associated_permissions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "associated_permissions.0.permission_id"),
					resource.TestCheckResourceAttrSet(dataSource, "associated_permissions.0.permission_name"),
					resource.TestCheckResourceAttrSet(dataSource, "associated_permissions.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "associated_permissions.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "associated_permissions.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "associated_permissions.0.updated_at"),

					dcByPermissionName.CheckResourceExists(),
					resource.TestCheckOutput("is_permission_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAssociatedPermissions_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ram_resource_share_associated_permissions" "test" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_share_id = huaweicloud_ram_resource_share.test.id
}

# Filter by permission_name
locals {
  permission_name = data.huaweicloud_ram_resource_share_associated_permissions.test.associated_permissions[0].permission_name
}

data "huaweicloud_ram_resource_share_associated_permissions" "filter_by_permission_name" {
  depends_on = [huaweicloud_ram_resource_share.test]

  resource_share_id = huaweicloud_ram_resource_share.test.id
  permission_name   = local.permission_name
}

locals {
  permission_name_filter_result = [
    for v in data.huaweicloud_ram_resource_share_associated_permissions.filter_by_permission_name.associated_permissions[*].
    permission_name : v == local.permission_name
  ]
}

output "is_permission_name_filter_useful" {
  value = length(local.permission_name_filter_result) > 0 && alltrue(local.permission_name_filter_result)
}
`, testRAMShare_basic(name))
}
