package taurusdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksSlowLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_slow_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckTaurusDBHtapInstanceId(t)
			acceptance.TestAccPreCheckTaurusDBHtapNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTaurusDBHtapStarrocksSlowLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "instance_id", acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID),
					resource.TestCheckResourceAttr(dataSource, "node_id", acceptance.HW_TAURUSDB_HTAP_NODE_ID),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.#"),
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

func testAccDataSourceTaurusDBHtapStarrocksSlowLogs_basic() string {
	cst := time.FixedZone("CST", 8*3600)
	endTime := time.Now().In(cst).Format("2006-01-02T15:04:05-0700")
	startTime := time.Now().In(cst).Add(-24 * time.Hour).Format("2006-01-02T15:04:05-0700")

	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_slow_logs" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
}

locals {
  operate_type = data.huaweicloud_taurusdb_htap_starrocks_slow_logs.test.slow_log_list[0].type
  database     = data.huaweicloud_taurusdb_htap_starrocks_slow_logs.test.slow_log_list[0].database
}

data "huaweicloud_taurusdb_htap_starrocks_slow_logs" "operate_type_filter" {
  instance_id  = "%[1]s"
  node_id      = "%[2]s"
  start_time   = "%[3]s"
  end_time     = "%[4]s"
  operate_type = local.operate_type
}

output "operate_type_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_htap_starrocks_slow_logs.operate_type_filter.slow_log_list) > 0 && alltrue(
  [for v in data.huaweicloud_taurusdb_htap_starrocks_slow_logs.operate_type_filter.slow_log_list[*].type : v == local.operate_type]
  )
}

data "huaweicloud_taurusdb_htap_starrocks_slow_logs" "database_filter" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
  database    = local.database
}
output "database_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_htap_starrocks_slow_logs.database_filter.slow_log_list) > 0 && alltrue(
  [for v in data.huaweicloud_taurusdb_htap_starrocks_slow_logs.database_filter.slow_log_list[*].database : v == local.database]
  )
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_NODE_ID, startTime, endTime)
}
