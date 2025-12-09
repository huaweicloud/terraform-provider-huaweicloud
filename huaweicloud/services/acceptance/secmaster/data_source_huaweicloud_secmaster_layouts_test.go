package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceLayouts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_layouts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceLayouts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.boa_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_built_in"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.layout_cfg"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.layout_json"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.sections_sum"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.tabs_sum"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.used_by"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.workspace_id"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_used_by_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceLayouts_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_layouts" "test" {
  workspace_id = "%[1]s"
}

# Filter by name
locals {
  name = data.huaweicloud_secmaster_layouts.test.data[0].name
}

data "huaweicloud_secmaster_layouts" "name_filter" {
  workspace_id = "%[1]s"
  name         = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_layouts.name_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_layouts.name_filter.data[*].name : v == local.name]
  )
}

# Filter by used_by
locals {
  used_by = data.huaweicloud_secmaster_layouts.test.data[0].used_by
}

data "huaweicloud_secmaster_layouts" "used_by_filter" {
  workspace_id = "%[1]s"
  used_by      = local.used_by
}

output "is_used_by_filter_useful" {
  value = length(data.huaweicloud_secmaster_layouts.used_by_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_layouts.used_by_filter.data[*].used_by : v == local.used_by]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
