package ddm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDdmPhysicalSessionsKill_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckDDMProcessId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccDdmPhysicalSessionsKill_basic(),
			},
		},
	})
}

func testAccDdmPhysicalSessionsKill_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_ddm_physical_sessions_kill" "test" {
  instance_id = "%[1]s"
  process_ids = ["%[2]s"]
}`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_DDM_PROCESS_ID)
}
