package waf

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewsClassification_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_overviews_classification.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		startTime      = time.Now().Add(-8*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime        = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Before running test, please prepare two WAF domain.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWafDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewsClassification_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "domain.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domain.0.items.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "domain.0.items.0.key"),

					resource.TestCheckOutput("value_is_exist", "true"),
					resource.TestCheckOutput("top_filter_is_useful", "true"),
					resource.TestCheckOutput("hosts_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOverviewsClassification_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_waf_overviews_classification" "test" {
  from = %[1]d
  to   = %[2]d
}

output "value_is_exist" {
  value = length(data.huaweicloud_waf_overviews_classification.test.domain[0].items) == 2
}

data "huaweicloud_waf_overviews_classification" "filter_by_top" {
  from = %[1]d
  to   = %[2]d
  top  = 1
}

output "top_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_classification.filter_by_top.domain[0].items) == 1
}

data "huaweicloud_waf_overviews_classification" "filter_by_hosts" {
  from  = %[1]d
  to    = %[2]d
  hosts = "%[3]s"
}

output "hosts_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_classification.filter_by_hosts.domain[0].items) == 1
}
`, startTime, endTime, acceptance.HW_WAF_DOMAIN_ID)
}
