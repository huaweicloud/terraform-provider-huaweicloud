package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafSecurityReport_basic(t *testing.T) {
	dataSource := "data.huaweicloud_waf_security_report.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPrecheckWafSecurityReportSubscription(t)
			acceptance.TestAccPrecheckWafSecurityReportId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceWafSecurityReport_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "report_id", acceptance.HW_WAF_SECURITY_REPORT_ID),
					resource.TestCheckResourceAttr(dataSource, "subscription_id", acceptance.HW_WAF_SECURITY_REPORT_SUBSCRIPTION_ID),
					resource.TestCheckResourceAttrSet(dataSource, "report_name"),
					resource.TestCheckResourceAttrSet(dataSource, "report_category"),
					resource.TestCheckResourceAttrSet(dataSource, "sending_period"),
					resource.TestCheckResourceAttrSet(dataSource, "topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "subscription_type"),
					resource.TestCheckResourceAttrSet(dataSource, "stat_period.#"),
					resource.TestCheckResourceAttrSet(dataSource, "create_time"),
				),
			},
		},
	})
}

func testDataSourceDataSourceWafSecurityReport_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_waf_security_report" "test" {
  report_id       = "%s"
  subscription_id = "%s"
}
`, acceptance.HW_WAF_SECURITY_REPORT_ID, acceptance.HW_WAF_SECURITY_REPORT_SUBSCRIPTION_ID)
}
