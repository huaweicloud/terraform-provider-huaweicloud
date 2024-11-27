package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbMysqlRestoredTables_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_restored_tables.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlInstanceId(t)
			acceptance.TestAccPreCheckGaussDBMysqlDatabaseName(t)
			acceptance.TestAccPreCheckGaussDBMysqlTableName(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbMysqlRestoredTables_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "databases.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.total_tables"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.tables.#"),
					resource.TestCheckResourceAttrSet(dataSource, "databases.0.tables.0.name"),

					resource.TestCheckOutput("database_name_filter_is_useful", "true"),
					resource.TestCheckOutput("table_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbMysqlRestoredTables_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_mysql_restore_time_ranges" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_mysql_restored_tables" "test" {
  instance_id     = "%[1]s"
  restore_time    = data.huaweicloud_gaussdb_mysql_restore_time_ranges.test.restore_times[0].end_time
  last_table_info = true
}

locals {
  database_name = "%[2]s"
}
data "huaweicloud_gaussdb_mysql_restored_tables" "database_name_filter" {
  instance_id     = "%[1]s"
  restore_time    = data.huaweicloud_gaussdb_mysql_restore_time_ranges.test.restore_times[0].end_time
  last_table_info = true
  database_name   = "%[2]s"
}
output "database_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_restored_tables.database_name_filter.databases) > 0
}

locals {
  table_name = "%[3]s"
}
data "huaweicloud_gaussdb_mysql_restored_tables" "table_name_filter" {
  instance_id     = "%[1]s"
  restore_time    = data.huaweicloud_gaussdb_mysql_restore_time_ranges.test.restore_times[0].end_time
  last_table_info = true
  table_name      = "%[3]s"
}
output "table_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_restored_tables.table_name_filter.databases) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_restored_tables.table_name_filter.databases[*].tables : length(v) > 0]
  )
}
`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID, acceptance.HW_GAUSSDB_MYSQL_DATABASE_NAME, acceptance.HW_GAUSSDB_MYSQL_TABLE_NAME)
}
