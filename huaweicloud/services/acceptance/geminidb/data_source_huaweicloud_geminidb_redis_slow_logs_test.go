package geminidb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRedisSlowLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_redis_slow_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		start_time = time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05+0800")
		end_time   = time.Now().Format("2006-01-02T15:04:05+0800")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRedisSlowLogs_basic(start_time, end_time),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.whole_message"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.operate_type"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.cost_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.log_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.line_num"),

					resource.TestCheckOutput("node_id_filter_useful", "true"),
					resource.TestCheckOutput("keywords_filter_useful", "true"),
					resource.TestCheckOutput("max_cost_time_filter_useful", "true"),
					resource.TestCheckOutput("min_cost_time_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceRedisSlowLogs_basic(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_redis_slow_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  node_id = data.huaweicloud_geminidb_redis_slow_logs.test.slow_logs[0].node_id
}

data "huaweicloud_geminidb_redis_slow_logs" "node_id_filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  node_id     = local.node_id
}

output "node_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_redis_slow_logs.node_id_filter.slow_logs) > 0 && alltrue(
    [for v in data.huaweicloud_geminidb_redis_slow_logs.node_id_filter.slow_logs[*].node_id : v == local.node_id]
  )
}

data "huaweicloud_geminidb_redis_slow_logs" "keywords_filter" {    
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  keywords    = ["default"]
}

output "keywords_filter_useful" {
  value = length(data.huaweicloud_geminidb_redis_slow_logs.keywords_filter.slow_logs) > 0
}

data "huaweicloud_geminidb_redis_slow_logs" "max_cost_time_filter" {
  instance_id   = "%[1]s"
  start_time    = "%[2]s"
  end_time      = "%[3]s"
  max_cost_time = 800
}

output "max_cost_time_filter_useful" {
  value = length(data.huaweicloud_geminidb_redis_slow_logs.max_cost_time_filter.slow_logs) > 0
}

data "huaweicloud_geminidb_redis_slow_logs" "min_cost_time_filter" {
  instance_id   = "%[1]s"
  start_time    = "%[2]s"
  end_time      = "%[3]s"
  min_cost_time = 600
}

output "min_cost_time_filter_useful" {
  value = length(data.huaweicloud_geminidb_redis_slow_logs.min_cost_time_filter.slow_logs) > 0
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, startTime, endTime)
}
