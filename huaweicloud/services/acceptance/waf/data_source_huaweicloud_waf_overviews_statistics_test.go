package waf

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewsStatistics_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_overviews_statistics.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		startTime      = time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime        = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Before running test, please prepare a WAF domain and ensure that has access records.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWafDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewsStatistics_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "statistics.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "statistics.0.key"),
					resource.TestCheckResourceAttrSet(dataSourceName, "statistics.0.num"),

					resource.TestCheckOutput("value_is_exist", "true"),
					resource.TestCheckOutput("hosts_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOverviewsStatistics_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_waf_overviews_statistics" "test" {
  from = %[1]d
  to   = %[2]d
}

output "value_is_exist" {
  value = data.huaweicloud_waf_overviews_statistics.test.statistics[0].num > 0
}

data "huaweicloud_waf_overviews_statistics" "filter_by_hosts" {
  from  = %[1]d
  to    = %[2]d
  hosts = "%[3]s"
}

output "hosts_filter_is_useful" {
  value = data.huaweicloud_waf_overviews_statistics.filter_by_hosts.statistics[0].num > 0
}
`, startTime, endTime, acceptance.HW_WAF_DOMAIN_ID)
}
