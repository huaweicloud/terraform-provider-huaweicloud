package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBMysqlInstantTaskDelete_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlJobId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBMysqlInstantTaskDelete_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccGaussDBMysqlInstantTaskDelete_basic() string {
	return fmt.Sprintf(`
resource "huaweicloud_gaussdb_mysql_instant_task_delete" "test" {
  job_id = "%[1]s"
}`, acceptance.HW_GAUSSDB_MYSQL_JOB_ID)
}
