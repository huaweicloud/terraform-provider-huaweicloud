package geminidb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceInfluxdbSlowLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_geminidb_influxdb_slow_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		start_time = time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05+0800")
		end_time   = time.Now().Format("2006-01-02T15:04:05+0800")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// No data for test.
			acceptance.TestAccCheckGeminidbInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceInfluxdbSlowLogs_basic(start_time, end_time),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.#"),

					resource.TestCheckOutput("value_is_empty", "true"),
				),
			},
		},
	})
}

func testAccDataSourceInfluxdbSlowLogs_basic(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_geminidb_influxdb_slow_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

output "value_is_empty" {
  value = length(data.huaweicloud_geminidb_influxdb_slow_logs.test.slow_logs) == 0
}
`, acceptance.HW_GEMINIDB_INSATNCE_ID, startTime, endTime)
}
