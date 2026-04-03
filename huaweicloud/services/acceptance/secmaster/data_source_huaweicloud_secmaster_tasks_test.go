package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecmasterTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_secmaster_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a secmaster workspace.
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecmasterTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.approver"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.business_type"),
					resource.TestCheckResourceAttrSet(dataSource, "tasks.0.status"),

					resource.TestCheckOutput("is_business_type_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSecmasterTasks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_tasks" "test" {
  workspace_id = "%[1]s"
}

# Filter by business_type
locals {
  business_type = data.huaweicloud_secmaster_tasks.test.tasks[0].business_type
}

data "huaweicloud_secmaster_tasks" "filter_by_business_type" {
  workspace_id  = "%[1]s"
  business_type = local.business_type
}

output "is_business_type_filter_useful" {
  value = length(data.huaweicloud_secmaster_tasks.filter_by_business_type.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_tasks.filter_by_business_type.tasks[*].business_type :
    v == local.business_type]
  )
}

# Filter by name
locals {
  name = data.huaweicloud_secmaster_tasks.test.tasks[0].name
}

data "huaweicloud_secmaster_tasks" "filter_by_name" {
  workspace_id = "%[1]s"
  name         = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_secmaster_tasks.filter_by_name.tasks) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_tasks.filter_by_name.tasks[*].name : v == local.name]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID)
}
