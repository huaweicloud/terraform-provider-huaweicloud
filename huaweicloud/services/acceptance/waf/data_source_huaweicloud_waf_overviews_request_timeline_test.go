package waf

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewsRequestTimeline_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_overviews_request_timeline.test"
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
				Config: testDataSourceOverviewsRequestTimeline_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "requests.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "requests.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "requests.0.timeline.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "requests.0.timeline.0.time"),

					resource.TestCheckOutput("value_is_exist", "true"),
					resource.TestCheckOutput("group_by_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOverviewsRequestTimeline_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_waf_overviews_request_timeline" "test" {
  from = %[1]d
  to   = %[2]d
}

output "value_is_exist" {
  value = length(data.huaweicloud_waf_overviews_request_timeline.test.requests) > 0
}

data "huaweicloud_waf_overviews_request_timeline" "filter_group_by" {
  from     = %[1]d
  to       = %[2]d
  group_by = "DAY"
}

output "group_by_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_request_timeline.filter_group_by.requests[0].timeline) > 0
}

data "huaweicloud_waf_overviews_request_timeline" "filter_eps" {
  from                  = %[1]d
  to                    = %[2]d
  enterprise_project_id = "0"
}

output "eps_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_request_timeline.filter_eps.requests) > 0
}
`, startTime, endTime)
}
