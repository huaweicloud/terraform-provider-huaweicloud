package waf

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewsAttackTopDomains_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_waf_overviews_attack_top_domains.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
		startTime      = time.Now().Add(-24*time.Hour).UnixNano() / int64(time.Millisecond)
		endTime        = time.Now().UnixNano() / int64(time.Millisecond)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			// Before running test, please prepare a WAF domain.
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWafDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewsAttackTopDomains_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "items.0.key"),

					resource.TestCheckOutput("top_filter_is_useful", "true"),
					resource.TestCheckOutput("eps_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOverviewsAttackTopDomains_basic(startTime, endTime int64) string {
	return fmt.Sprintf(`
data "huaweicloud_waf_overviews_attack_top_domains" "test" {
  from = %[1]d
  to   = %[2]d
}

data "huaweicloud_waf_overviews_attack_top_domains" "filter_by_top" {
  from  = %[1]d
  to    = %[2]d
  top   = 1
}

output "top_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_attack_top_domains.filter_by_top.items) > 0
}

data "huaweicloud_waf_overviews_attack_top_domains" "filter_by_eps" {
  from                  = %[1]d
  to                    = %[2]d
  enterprise_project_id = "0"
}

output "eps_filter_is_useful" {
  value = length(data.huaweicloud_waf_overviews_attack_top_domains.filter_by_eps.items) > 0
}
`, startTime, endTime)
}
