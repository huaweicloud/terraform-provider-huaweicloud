package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccUnlockNodeReadonlyStatus_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccUnlockNodeReadonlyStatus_basic(),
			},
		},
	})
}

func testAccUnlockNodeReadonlyStatus_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_unlock_node_readonly_status" "test" {
  instance_id              = "%s"
  status_preservation_time = 100
}`, acceptance.HW_RDS_INSTANCE_ID)
}
