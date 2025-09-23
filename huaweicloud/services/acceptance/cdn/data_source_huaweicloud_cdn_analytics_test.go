package cdn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running this test case, please configure the following environment variables:
// HW_CDN_DOMAIN_NAME: Configure the CDN domain name with statistical data.
// HW_CDN_START_TIME: Configure the start timestamp for querying domain name statistical data (in milliseconds).
// HW_CDN_END_TIME: Configure the end timestamp for querying domain name statistical data (in milliseconds).
// HW_CDN_STAT_TYPE: Configure the indicator type for querying statistical data. Use commas to separate multiple types.
func TestAccDataSourceCdnAnalytics_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cdn_analytics.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCDN(t)
			acceptance.TestAccPrecheckCDNAnalytics(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCdnAnalytics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "result"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCdnAnalytics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cdn_analytics" "test" {
  domain_name = "%[1]s"
  stat_type   = "%[2]s"
  start_time  = "%[3]s"
  end_time    = "%[4]s"
  action      = "detail"
  group_by    = "domain"
}
`, acceptance.HW_CDN_DOMAIN_NAME, acceptance.HW_CDN_STAT_TYPE, acceptance.HW_CDN_START_TIME, acceptance.HW_CDN_END_TIME)
}
