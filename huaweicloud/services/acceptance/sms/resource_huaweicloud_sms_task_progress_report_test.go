package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceSmsTaskProgressReport_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSmsTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSmsTaskProgressReport_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testSmsTaskProgressReport_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sms_task_progress_report" "test" {
  task_id             = "%s"
  subtask_name        = "DETTACH_AGENT_IMAGE"
  progress            = 100
  replicatesize       = 1000
  totalsize           = 100000
  process_trace       = "migrate details"
  migrate_speed       = 10
  compress_rate       = 50
  remain_time         = 0
  total_cpu_usage     = 50
  agent_cpu_usage     = 50
  total_mem_usage     = 500
  agent_mem_usage     = 500
  total_disk_io       = 500
  agent_disk_io       = 500
  need_migration_test = false
  agent_time          = "2025-07-07T15:30:00"
}
`, acceptance.HW_SMS_TASK_ID)
}
