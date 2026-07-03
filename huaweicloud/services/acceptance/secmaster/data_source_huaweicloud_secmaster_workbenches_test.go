package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkbenches_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_workbenches.test"
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
				Config: testDataSourceWorkbenches_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.url"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.url_openwith_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.icon"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_deleted"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.is_favorite"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),

					resource.TestCheckOutput("filter_by_type_is_useful", "true"),
					resource.TestCheckOutput("filter_by_status_is_useful", "true"),
					resource.TestCheckOutput("filter_by_creator_type_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceWorkbenches_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_workbenches" "test" {
  workspace_id = "%[1]s"
}

locals {
  workbench_type = data.huaweicloud_secmaster_workbenches.test.data[0].type
}

# Filter by type
data "huaweicloud_secmaster_workbenches" "filter_by_type" {
  workspace_id = "%[1]s"
  type         = local.workbench_type
}

output "filter_by_type_is_useful" {
  value = length(data.huaweicloud_secmaster_workbenches.filter_by_type.data) > 0 && alltrue(
    [for info in data.huaweicloud_secmaster_workbenches.filter_by_type.data :
      info.type == local.workbench_type]
  )
}

locals {
  workbench_status = data.huaweicloud_secmaster_workbenches.test.data[0].status
}

# Filter by status
data "huaweicloud_secmaster_workbenches" "filter_by_status" {
  workspace_id = "%[1]s"
  status       = local.workbench_status
}

output "filter_by_status_is_useful" {
  value = length(data.huaweicloud_secmaster_workbenches.filter_by_status.data) > 0 && alltrue(
    [for info in data.huaweicloud_secmaster_workbenches.filter_by_status.data :
      info.status == local.workbench_status]
  )
}

locals {
  workbench_creator_type = data.huaweicloud_secmaster_workbenches.test.data[0].creator_id
}

# Filter by creator_type
data "huaweicloud_secmaster_workbenches" "filter_by_creator_type" {
  workspace_id = "%[1]s"
  creator_type = "system"
}

output "filter_by_creator_type_is_useful" {
  value = length(data.huaweicloud_secmaster_workbenches.filter_by_creator_type.data) > 0
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
