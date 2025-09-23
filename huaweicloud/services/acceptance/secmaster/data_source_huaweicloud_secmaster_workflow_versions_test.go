package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkflowVersions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_workflow_versions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			acceptance.TestAccPreCheckSecMasterWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkflowVersions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.workflow_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.creator_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.taskconfig"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.taskflow"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.taskflow_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data.0.aop_type"),

					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceWorkflowVersions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_workflow_versions" "test" {
  workspace_id = "%[1]s"
  workflow_id  = "%[2]s"
}

locals {
  status = data.huaweicloud_secmaster_workflow_versions.test.data[0].status
}

data "huaweicloud_secmaster_workflow_versions" "status_filter" {
  workspace_id = "%[1]s"
  workflow_id  = "%[2]s"
  status       = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_secmaster_workflow_versions.status_filter.data) > 0 && alltrue(
    [for v in data.huaweicloud_secmaster_workflow_versions.status_filter.data[*].status : v == local.status]
  )
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_WORKFLOW_ID)
}
