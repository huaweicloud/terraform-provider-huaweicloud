package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceTaskConsistencyResultReport_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSmsTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testTaskConsistencyResultReport_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testTaskConsistencyResultReport_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sms_task_consistency_result_report" "test" {
  task_id = "%s"

  consistency_result {
    dir_check             = "/root/test"
    num_total_files       = 1
    num_different_files   = 1
    num_target_miss_files = 1
    num_target_more_files = 1
  }
}
`, acceptance.HW_SMS_TASK_ID)
}
