package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccTaurusDBInstanceScheduledTaskDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccTaurusDBInstanceScheduledTaskDelete_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccTaurusDBInstanceScheduledTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_taurusdb_scheduled_task_delete" "test" {
  instance_id = "%[1]s"
  job_id      = "%[2]s"
}`, acceptance.HW_TAURUSDB_INSTANCE_ID, acceptance.HW_TAURUSDB_JOB_ID)
}
