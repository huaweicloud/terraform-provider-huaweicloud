package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceSmsTaskLogUpload_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSmsTaskID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testSmsTaskLogUpload_basic(),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testSmsTaskLogUpload_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_sms_task_log_upload" "test" {
  task_id    = "%s"
  log_bucket = "test"
}
`, acceptance.HW_SMS_TASK_ID)
}
