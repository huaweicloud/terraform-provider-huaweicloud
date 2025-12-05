package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this acceptance test, make sure that the scheduled task has been executed at least once.
func TestAccDataScheduledTaskRecordDetails_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_workspace_scheduled_task_record_details.all"
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
				Config: testAccDataScheduledTaskRecordDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "details.#", regexp.MustCompile(`^[0-9]+$`)),
					resource.TestCheckResourceAttrSet(all, "details.0.id"),
					resource.TestCheckResourceAttrSet(all, "details.0.record_id"),
					resource.TestCheckResourceAttrPair(all, "details.0.record_id",
						"data.huaweicloud_workspace_scheduled_task_records.test", "records.0.id"),
					resource.TestCheckResourceAttrSet(all, "details.0.desktop_id"),
					resource.TestCheckResourceAttrSet(all, "details.0.desktop_name"),
					resource.TestCheckResourceAttrSet(all, "details.0.exec_status"),
					resource.TestMatchResourceAttr(all, "details.0.start_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataScheduledTaskRecordDetails_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_scheduled_task_records" "test" {
  task_id = "%[1]s"
}

data "huaweicloud_workspace_scheduled_task_record_details" "all" {
  task_id   = "%[1]s"
  record_id = try(data.huaweicloud_workspace_scheduled_task_records.test.records[0].id, "")
}
`, acceptance.HW_WORKSPACE_SCHEDULED_TASK_ID)
}
