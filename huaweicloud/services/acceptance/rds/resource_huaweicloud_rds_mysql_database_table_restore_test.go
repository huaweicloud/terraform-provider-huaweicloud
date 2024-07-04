package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccMysqlDatabaseTableRestore_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceMysqlDatabaseTableRestoreConfig_basic(),
			},
		},
	})
}

func TestAccMysqlDatabaseTableRestore_table(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccRdsInstanceMysqlDatabaseTableRestoreConfig_table(),
			},
		},
	})
}

func testAccRdsInstanceMysqlDatabaseTableRestoreConfig_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_mysql_database_table_restore" "test" {
  restore_time    = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].start_time
  instance_id     = "%[1]s"
  is_fast_restore = true

  databases {
    old_name = "test111"
    new_name = "test111_terraform_update"
  }

  databases {
    old_name = "test333"
    new_name = "test333_terraform_update"
  }
}
`, acceptance.HW_RDS_INSTANCE_ID)
}

func testAccRdsInstanceMysqlDatabaseTableRestoreConfig_table() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

resource "huaweicloud_rds_mysql_database_table_restore" "test" {
  restore_time    = data.huaweicloud_rds_restore_time_ranges.test.restore_time[0].start_time
  instance_id     = "%[1]s"
  is_fast_restore = true

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
`, acceptance.HW_RDS_INSTANCE_ID)
}
