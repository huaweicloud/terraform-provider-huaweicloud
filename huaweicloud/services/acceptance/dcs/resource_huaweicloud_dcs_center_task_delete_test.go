package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDcsCenterTaskDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDcsCenterTaskId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDcsCenterTaskDelete_basic(),
			},
		},
	})
}

func testAccDcsCenterTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_dcs_center_task_delete" "test" {
  task_id = "%[1]s"
}
`, acceptance.HW_DCS_CENTER_TASK_ID)
}
