package taurusdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceTaurusDBHtapStarrocksErrorLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_taurusdb_htap_starrocks_error_logs.test"
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
				Config: testAccDataSourceTaurusDBHtapStarrocksErrorLogs_basic(),
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

func testAccDataSourceTaurusDBHtapStarrocksErrorLogs_basic() string {
	cst := time.FixedZone("CST", 8*3600)
	endTime := time.Now().In(cst).Format("2006-01-02T15:04:05-0700")
	startTime := time.Now().In(cst).Add(-10 * time.Hour).Format("2006-01-02T15:04:05-0700")

	return fmt.Sprintf(`
data "huaweicloud_taurusdb_htap_starrocks_error_logs" "test" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
  level       = "ALL"
}

locals {
  level = data.huaweicloud_taurusdb_htap_starrocks_error_logs.test.error_log_list[0].level
}

data "huaweicloud_taurusdb_htap_starrocks_error_logs" "level_filter" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
  level       = local.level
}

output "level_filter_is_useful" {
  value = length(data.huaweicloud_taurusdb_htap_starrocks_error_logs.level_filter.error_log_list) > 0 && alltrue(
  [for v in data.huaweicloud_taurusdb_htap_starrocks_error_logs.level_filter.error_log_list[*].level : v == local.level]
  )
}
`, acceptance.HW_TAURUSDB_HTAP_INSTANCE_ID, acceptance.HW_TAURUSDB_HTAP_NODE_ID, startTime, endTime)
}
