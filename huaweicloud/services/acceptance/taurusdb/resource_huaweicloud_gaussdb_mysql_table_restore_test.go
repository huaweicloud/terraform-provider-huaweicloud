package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGaussDBTableRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccGaussDBTableRestoreonfig_basic(),
			},
		},
	})
}

func testAccGaussDBTableRestoreonfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_mysql_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_gaussdb_mysql_table_restore" "test" {
  restore_time    = data.huaweicloud_gaussdb_mysql_restore_time_ranges.test.restore_times[0].start_time
  instance_id     = "%[1]s"
  last_table_info = true

  restore_tables {
    database = "test111"
    tables {
      old_name = "table111"
      new_name = "table111_terraform_update"
    }

    tables {
      old_name = "table222"
      new_name = "table222_terraform_update"
    }
  }

  restore_tables {
    database = "test222"
    tables {
      old_name = "table111"
      new_name = "table111_terraform_update"
    }
  }
}
`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID)
}
