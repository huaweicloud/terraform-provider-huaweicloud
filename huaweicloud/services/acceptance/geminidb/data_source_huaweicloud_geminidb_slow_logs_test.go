package geminidb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccGeminiDBSlowLogs_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_geminidb_slow_logs.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccGeminiDBSlowLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "slow_log_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slow_log_list.0.time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slow_log_list.0.database"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slow_log_list.0.query_sample"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slow_log_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "slow_log_list.0.start_time"),

					resource.TestCheckOutput("node_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccGeminiDBSlowLogs_basic() string {
	beginTime := time.Now().UTC().Add(-30 * 24 * time.Hour)
	beginTimeString := beginTime.Format("2006-01-02T15:04:05+0000")
	endTime := time.Now().UTC()
	endTimeString := endTime.Format("2006-01-02T15:04:05+0000")
	return fmt.Sprintf(`
data "huaweicloud_geminidb_slow_logs" "test" {
  instance_id = "%[1]s"
  start_date  = "%[2]s"
  end_date    = "%[3]s"
}

data "huaweicloud_geminidb_slow_logs" "node_id_filter" {
  instance_id = "%[1]s"
  start_date  = "%[2]s"
  end_date    = "%[3]s"
  node_id     = "%[4]s"
}

output "node_id_filter_useful" {
  value = length(data.huaweicloud_geminidb_slow_logs.node_id_filter.slow_log_list) > 0
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, beginTimeString, endTimeString, acceptance.HW_GEMINIDB_NODE_ID)
}
