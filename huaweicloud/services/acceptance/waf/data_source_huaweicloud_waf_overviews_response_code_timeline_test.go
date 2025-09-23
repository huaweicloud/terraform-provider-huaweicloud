package waf

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Due to testing environment limitations, this test case can only test the scenario with empty `response_codes`.
func TestAccDataSourceOverviewsResponseCodeTimeline_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_overviews_response_code_timeline.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		startTime      = time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime        = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewsResponseCodeTimeline_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "response_codes.#"),
				),
			},
		},
	})
}

func testDataSourceOverviewsResponseCodeTimeline_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_waf_overviews_response_code_timeline" "test" {
  from                  = %[1]d
  to                    = %[2]d
  group_by              = "DAY"
  response_source       = "WAF"
  enterprise_project_id = "0"
}
`, startTime, endTime)
}
