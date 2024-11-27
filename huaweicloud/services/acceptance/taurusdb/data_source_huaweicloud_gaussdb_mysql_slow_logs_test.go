package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBMysqlSlowLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_slow_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBMysqlInstanceId(t)
			acceptance.TestAccPreCheckGaussDBMysqlNodeId(t)
			acceptance.TestAccPreCheckGaussDBMysqlTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDBMysqlSlowLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.count"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.lock_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.rows_sent"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.rows_examined"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.database"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.users"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.query_sample"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.0.client_ip"),

					resource.TestCheckOutput("operate_type_filter_is_useful", "true"),
					resource.TestCheckOutput("database_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussDBMysqlSlowLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_mysql_slow_logs" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
}

locals {
  operate_type = "UPDATE"
}
data "huaweicloud_gaussdb_mysql_slow_logs" "operate_type_filter" {
  instance_id  = "%[1]s"
  node_id      = "%[2]s"
  start_time   = "%[3]s"
  end_time     = "%[4]s"
  operate_type = "UPDATE"
}
output "operate_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_slow_logs.operate_type_filter.slow_log_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_slow_logs.operate_type_filter.slow_log_list[*].type : v == local.operate_type]
  )
}

locals {
  database = "test_db_1"
}
data "huaweicloud_gaussdb_mysql_slow_logs" "database_filter" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
  database    = "test_db_1"
}
output "database_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_slow_logs.database_filter.slow_log_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_slow_logs.database_filter.slow_log_list[*].database : v == local.database]
  )
}
`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID, acceptance.HW_GAUSSDB_MYSQL_NODE_ID, acceptance.HW_GAUSSDB_MYSQL_START_TIME,
		acceptance.HW_GAUSSDB_MYSQL_END_TIME)
}
