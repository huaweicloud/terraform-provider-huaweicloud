package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccResourceRdsPgTableRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstancePgTableRestoreConfig_basic(),
			},
		},
	})
}

func testAccRdsInstancePgTableRestoreConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_pg_table_restore" "test" {
  instance_id  = "%[1]s"
  restore_time = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].start_time

  databases {
    database = "test1"

    schemas {
      schema = "test1"

      tables {
        old_name = "table1"
        new_name = "table1_test_update"
      }
    }
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
