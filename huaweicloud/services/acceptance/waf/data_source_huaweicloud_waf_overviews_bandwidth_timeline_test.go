package waf

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewsBandwidthTimeline_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_overviews_bandwidth_timeline.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		startTime      = time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime        = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Before running test, please ensure have data on console.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewsBandwidthTimeline_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.timeline.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "bandwidths.0.timeline.0.time"),

					resource.TestCheckOutput("value_is_exist", "true"),
					resource.TestCheckOutput("group_by_filter_is_useful", "true"),
					resource.TestCheckOutput("display_option_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOverviewsBandwidthTimeline_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_waf_overviews_bandwidth_timeline" "test" {
  from = %[1]d
  to   = %[2]d
}

output "value_is_exist" {
  value = length(data.huaweicloud_waf_overviews_bandwidth_timeline.test.bandwidths) > 0
}

data "huaweicloud_waf_overviews_bandwidth_timeline" "filter_group_by" {
  from     = %[1]d
  to       = %[2]d
  group_by = "DAY"
}

output "group_by_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_bandwidth_timeline.filter_group_by.bandwidths[0].timeline) > 0
}

data "huaweicloud_waf_overviews_bandwidth_timeline" "filter_display_option" {
  from           = %[1]d
  to             = %[2]d
  display_option = "1"
}

output "display_option_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_bandwidth_timeline.filter_display_option.bandwidths[0].timeline) > 0
}
`, startTime, endTime)
}
