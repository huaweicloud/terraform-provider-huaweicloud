package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBMysqlInstanceScheduledTaskDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlInstanceId(t)
			acceptance.TestAccPreCheckGaussDBMysqlJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBMysqlInstanceScheduledTaskDelete_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccGaussDBMysqlInstanceScheduledTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_scheduled_task_delete" "test" {
  instance_id = "%[1]s"
  job_id      = "%[2]s"
}`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID, acceptance.HW_GAUSSDB_MYSQL_JOB_ID)
}
