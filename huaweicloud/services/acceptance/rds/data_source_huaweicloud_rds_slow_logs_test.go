package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSlowLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_slow_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsSlowLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.count"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.lock_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.rows_sent"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.rows_examined"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.database"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.users"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.query_sample"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.client_ip"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("database_filter_is_useful", "true"),
					resource.TestCheckOutput("users_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceRdsSlowLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_slow_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  type = "UPDATE"
}
data "huaweicloud_rds_slow_logs" "type_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  type        = "UPDATE"
}
output "type_filter_is_useful" {
  value = length(data.huaweicloud_rds_slow_logs.type_filter.slow_logs) > 0 && alltrue(
  [for v in data.huaweicloud_rds_slow_logs.type_filter.slow_logs[*].type : v == local.type]
  )
}

locals {
  database = "test1111"
}
data "huaweicloud_rds_slow_logs" "database_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  database    = "test1111"
}
output "database_filter_is_useful" {
  value = length(data.huaweicloud_rds_slow_logs.database_filter.slow_logs) > 0 && alltrue(
  [for v in data.huaweicloud_rds_slow_logs.database_filter.slow_logs[*].database : v == local.database]
  )
}

locals {
  users = "root"
}
data "huaweicloud_rds_slow_logs" "users_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  users       = "root"
}
output "users_filter_is_useful" {
  value = length(data.huaweicloud_rds_slow_logs.users_filter.slow_logs) > 0 && alltrue(
  [for v in data.huaweicloud_rds_slow_logs.users_filter.slow_logs[*].users : v == local.users]
  )
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_START_TIME, acceptance.HW_RDS_END_TIME)
}
