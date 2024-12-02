package taurusdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBMysqlErrorLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_mysql_error_logs.test"
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
				Config: testDataSourceGaussDBMysqlErrorLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "error_log_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "error_log_list.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "error_log_list.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "error_log_list.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "error_log_list.0.content"),

					resource.TestCheckOutput("level_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussDBMysqlErrorLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_mysql_error_logs" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
}

locals {
  level = "WARNING"
}
data "huaweicloud_gaussdb_mysql_error_logs" "level_filter" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
  level       = "WARNING"
}
output "level_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_mysql_error_logs.level_filter.error_log_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_mysql_error_logs.level_filter.error_log_list[*].level : v == local.level]
  )
}
`, acceptance.HW_GAUSSDB_MYSQL_INSTANCE_ID, acceptance.HW_GAUSSDB_MYSQL_NODE_ID, acceptance.HW_GAUSSDB_MYSQL_START_TIME,
		acceptance.HW_GAUSSDB_MYSQL_END_TIME)
}
