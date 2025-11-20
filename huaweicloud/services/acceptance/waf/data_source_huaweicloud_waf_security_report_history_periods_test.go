package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityReportHistoryPeriods_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_security_report_history_periods.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPrecheckWafSecurityReportSubscription(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityReportHistoryPeriods_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
				),
			},
		},
	})
}

func testAccDataSourceSecurityReportHistoryPeriods_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_waf_security_report_history_periods" "test" {
  subscription_id = "%s"
}
`, acceptance.HW_WAF_SECURITY_REPORT_SUBSCRIPTION_ID)
}
