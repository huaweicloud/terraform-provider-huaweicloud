package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_cdn_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		// set a query start time
		// make sure the CDN log file is generated on the query time
		timeStamp = acceptance.HW_CDN_TIMESTAMP

		byStartTime   = "data.huaweicloud_cdn_logs.filter_by_start_time"
		dcByStartTime = acceptance.InitDataSourceCheck(byStartTime)

		byEndTime   = "data.huaweicloud_cdn_logs.filter_by_end_time"
		dcByEndTime = acceptance.InitDataSourceCheck(byEndTime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCdnDomainName(t)
			acceptance.TestAccPrecheckTimeStamp(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataLogs_basic(timeStamp),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.size"),
					resource.TestCheckResourceAttrSet(dataSource, "logs.0.link"),

					dcByStartTime.CheckResourceExists(),
					resource.TestCheckOutput("start_time_filter_is_useful", "true"),

					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataLogs_basic(time string) string {
	return fmt.Sprintf(`
data "huaweicloud_cdn_logs" "test" {
  domain_name = "%[1]s"
  start_time  = "%[2]s"
}

locals {
  start_time = data.huaweicloud_cdn_logs.test.logs[0].start_time
}

data "huaweicloud_cdn_logs" "filter_by_start_time" {
  domain_name = "%[1]s"
  start_time  = local.start_time
}

locals {
  start_time_filter_result = [
    for v in data.huaweicloud_cdn_logs.filter_by_start_time.logs[*].start_time : v == local.start_time
  ]
}

output "start_time_filter_is_useful" {
  value = alltrue(local.start_time_filter_result) && length(local.start_time_filter_result) > 0
}

locals {
  end_time = data.huaweicloud_cdn_logs.test.logs[0].end_time
}

data "huaweicloud_cdn_logs" "filter_by_end_time" {
  domain_name = "%[1]s"
  start_time  = "%[2]s"
  end_time    = local.end_time
}

locals {
  end_time_filter_result = [
    for v in data.huaweicloud_cdn_logs.filter_by_end_time.logs[*].end_time : v == local.end_time
  ]
}

output "end_time_filter_is_useful" {
  value = alltrue(local.end_time_filter_result) && length(local.end_time_filter_result) > 0
}
`, acceptance.HW_CDN_DOMAIN_NAME, time)
}
