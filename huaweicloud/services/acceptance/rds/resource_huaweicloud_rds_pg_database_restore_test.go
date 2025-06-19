package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRdsPgDatabaseRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstancePgDatabaseRestoreConfig_basic(),
			},
		},
	})
}

func testAccRdsInstancePgDatabaseRestoreConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_pg_database_restore" "test" {
  instance_id  = "%[1]s"
  restore_time = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].start_time

  databases {
    old_name = "test_database"
    new_name = "test_database_update"
  }
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
