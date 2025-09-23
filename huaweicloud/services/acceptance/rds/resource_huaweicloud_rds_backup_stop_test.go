package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccRdsBackupStop_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsBackupStop_basic(),
			},
		},
	})
}

func testAccRdsBackupStop_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_rds_backup_stop" "test" {
  instance_id = "%s"
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
