package rfs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRfsPrivateModuleVersions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rfs_private_module_versions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the private module data in the environment before running this test case.
			acceptance.TestAccPreCheckRfsModuleName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRfsPrivateModuleVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "versions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.module_name"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.module_id"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.module_version"),
					resource.TestCheckResourceAttrSet(dataSource, "versions.0.create_time"),

					resource.TestCheckOutput("is_module_id_filter_result", "true"),
				),
			},
		},
	})
}

func testDataSourceRfsPrivateModuleVersions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rfs_private_module_versions" "test" {
  module_name = "%[1]s"
  sort_key    = "create_time"
  sort_dir    = "asc"
}

locals {
  module_id = data.huaweicloud_rfs_private_module_versions.test.versions.0.module_id
}

# Filter by module_id
data "huaweicloud_rfs_private_module_versions" "module_id_filter" {
  module_name = "%[1]s"
  module_id   = local.module_id
}

locals {
  module_id_filter_result = [
    for v in data.huaweicloud_rfs_private_module_versions.module_id_filter.versions[*].module_id : v == local.module_id
  ]
}

output "is_module_id_filter_result" {
  value = length(local.module_id_filter_result) > 0 && alltrue(local.module_id_filter_result)
}
`, acceptance.HW_RFS_MODULE_NAME)
}
