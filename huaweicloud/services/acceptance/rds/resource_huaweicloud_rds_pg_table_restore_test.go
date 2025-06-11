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
// tonumber(data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].start_time)
func testAccRdsInstancePgTableRestoreConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_pg_table_restore" "test" {
  instances {
    instance_id  = "%[1]s"
    restore_time = 1754954459000  

    databases {
      database = "xxxx"

      schemas {
        schema = "schema-xxxx"

        tables {
          old_name = "mytable"
          new_name = "aaaaaaa"
        }
      }
    }
  }
}
`, acceptance.HW_RDS_INSTANCE_ID)
}
