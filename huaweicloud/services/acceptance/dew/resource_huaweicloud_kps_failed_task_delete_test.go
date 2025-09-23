package dew

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccKpsFailedTaskDelete_basic(t *testing.T) {
	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare one failed taskId, then config it to the environment variable.
			acceptance.TestAccPreCheckKpsTaskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKpsFailedTaskDelete_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("is_delete_success", "true"),
				),
			},
		},
	})
}

func testAccKpsFailedTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_kps_failed_task_delete" "test" {
  task_id = "%[1]s"
}

data "huaweicloud_kps_failed_tasks" "after_delete" {
  depends_on = [huaweicloud_kps_failed_task_delete.test]
}

output "is_delete_success" {
  value = alltrue([for v in data.huaweicloud_kps_failed_tasks.after_delete.tasks[*].id : v != "%[1]s"])
}
`, acceptance.HW_KPS_FAILED_TASK_ID)
}
