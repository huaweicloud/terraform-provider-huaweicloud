package geminidb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBScheduledTaskCancel_basic(t *testing.T) {
	resourceName := "huaweicloud_geminidb_scheduled_task_cancel.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbJobID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBScheduledTaskCancel_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "job_id", acceptance.HW_GEMINIDB_JOB_ID),
				),
			},
		},
	})
}

func testAccGeminiDBScheduledTaskCancel_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_geminidb_scheduled_task_cancel" "test" {
  job_id = "%s"
}
`, acceptance.HW_GEMINIDB_JOB_ID)
}
