package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDdmLogicalSessionsKill_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDMInstanceID(t)
			acceptance.TestAccPreCheckDDMProcessId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDdmLogicalSessionsKill_basic(),
			},
		},
	})
}

func testAccDdmLogicalSessionsKill_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ddm_logical_sessions_kill" "test" {
  instance_id = "%[1]s"
  process_ids = ["%[2]s"]
}`, acceptance.HW_DDM_INSTANCE_ID, acceptance.HW_DDM_PROCESS_ID)
}
