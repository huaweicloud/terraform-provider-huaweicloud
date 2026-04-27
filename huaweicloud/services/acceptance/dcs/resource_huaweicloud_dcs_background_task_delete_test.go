package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsBackgroundTaskDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDCSInstanceID(t)
			acceptance.TestAccPreCheckDcsBackgroundTaskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsBackgroundTaskDelete_basic(),
			},
		},
	})
}

func testAccDcsBackgroundTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dcs_background_task_delete" "test" {
  instance_id = "%[1]s"
  task_id     = "%[2]s"
}
`, acceptance.HW_DCS_INSTANCE_ID, acceptance.HW_DCS_BACKGROUND_TASK_ID)
}
