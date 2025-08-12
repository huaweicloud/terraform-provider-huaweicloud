package secmaster

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWorkflowInstanceDetail_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_secmaster_workflow_instance.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSecMasterWorkspaceID(t)
			// The workflow instance ID
			acceptance.TestAccPreCheckSecMasterInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkflowInstanceDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "workflow_instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "name"),
					resource.TestCheckResourceAttrSet(dataSource, "workflow.#"),
					resource.TestCheckResourceAttrSet(dataSource, "workflow.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "workflow.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "workflow.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "dataclass.#"),
					resource.TestCheckResourceAttrSet(dataSource, "dataclass.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "dataclass.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "playbook.#"),
					resource.TestCheckResourceAttrSet(dataSource, "playbook.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "playbook.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "trigger_type"),
					resource.TestCheckResourceAttrSet(dataSource, "status"),
					resource.TestCheckResourceAttrSet(dataSource, "start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "end_time"),
				),
			},
		},
	})
}

func testAccDataSourceWorkflowInstanceDetail_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_secmaster_workflow_instance" "test" {
  workspace_id  = "%[1]s"
  instance_id   = "%[2]s"
}
`, acceptance.HW_SECMASTER_WORKSPACE_ID, acceptance.HW_SECMASTER_INSTANCE_ID)
}
