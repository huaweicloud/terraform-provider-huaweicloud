package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, make sure that the scheduled task has been executed at least once.
func TestAccDataScheduledTaskRecords_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_scheduled_task_records.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceScheduledTaskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataScheduledTaskRecords_invalidTaskId(),
				ExpectError: regexp.MustCompile(`The scheduled task does not exist`),
			},
			{
				Config: testAccDataScheduledTaskRecords_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "records.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "records.0.id"),
					resource.TestCheckResourceAttrSet(all, "records.0.start_time"),
					resource.TestCheckResourceAttrSet(all, "records.0.task_type"),
					resource.TestCheckResourceAttrSet(all, "records.0.scheduled_type"),
					resource.TestCheckResourceAttrSet(all, "records.0.status"),
					resource.TestCheckResourceAttrSet(all, "records.0.success_num"),
					resource.TestCheckResourceAttrSet(all, "records.0.failed_num"),
					resource.TestCheckResourceAttrSet(all, "records.0.skip_num"),
					resource.TestMatchResourceAttr(all, "records.0.start_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataScheduledTaskRecords_invalidTaskId() string {
	randUUID, _ := uuid.GenerateUUID()

	return fmt.Sprintf(`
data "huaweicloud_workspace_scheduled_task_records" "invalid_task_id" {
  task_id = "%s"
}
`, randUUID)
}

func testAccDataScheduledTaskRecords_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_scheduled_task_records" "all" {
  task_id = "%s"
}
`, acceptance.HW_WORKSPACE_SCHEDULED_TASK_ID)
}
