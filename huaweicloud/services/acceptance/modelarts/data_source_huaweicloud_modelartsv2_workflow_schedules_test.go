package modelarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceV2WorkflowSchedules_basic(t *testing.T) {
	var (
		dcName = "data.huaweicloud_modelartsv2_workflow_schedules.test"
		dc     = acceptance.InitDataSourceCheck(dcName)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckModelArtsWorkflowId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceV2WorkflowSchedules_nonExistentWorkflow(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dcName, "schedules.#", "0"),
				),
			},
			{
				Config: testAccDatasourceV2WorkflowSchedules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_created_schedule_found", "true"),
					resource.TestCheckResourceAttrSet(dcName, "schedules.0.id"),
					resource.TestCheckResourceAttr(dcName, "schedules.0.workflow_id", acceptance.HW_MODELARTS_WORKFLOW_ID),
					resource.TestCheckResourceAttrSet(dcName, "schedules.0.content"),
					resource.TestCheckResourceAttrSet(dcName, "schedules.0.enable"),
					resource.TestCheckResourceAttrSet(dcName, "schedules.0.type"),
					resource.TestCheckResourceAttr(dcName, "schedules.0.action", "run"),
					resource.TestCheckResourceAttr(dcName, "schedules.0.policies.#", "1"),
					resource.TestCheckResourceAttrSet(dcName, "schedules.0.policies.0.on_failure"),
					resource.TestCheckResourceAttrSet(dcName, "schedules.0.policies.0.on_running"),
					resource.TestCheckResourceAttrSet(dcName, "schedules.0.user_id"),
					resource.TestMatchResourceAttr(dcName, "schedules.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDatasourceV2WorkflowSchedules_nonExistentWorkflow() string {
	return fmt.Sprintf(`
data "huaweicloud_modelartsv2_workflow_schedules" "test" {
  workflow_id = "%[1]s"
}
`, acceptance.HW_MODELARTS_WORKFLOW_ID)
}

func testAccDatasourceV2WorkflowSchedules_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_modelartsv2_workflow_schedule" "test" {
  workflow_id = "%[1]s"
  content     = jsonencode({
    cron   = "0 0 12 * * ?"
    method = "fixed"
  })
}

data "huaweicloud_modelartsv2_workflow_schedules" "test" {
  depends_on = [huaweicloud_modelartsv2_workflow_schedule.test]

  workflow_id = "%[1]s"
}

output "is_created_schedule_found" {
  value = length([for v in data.huaweicloud_modelartsv2_workflow_schedules.test.schedules : v if
    v.id == huaweicloud_modelartsv2_workflow_schedule.test.id]) > 0
}
`, acceptance.HW_MODELARTS_WORKFLOW_ID)
}
