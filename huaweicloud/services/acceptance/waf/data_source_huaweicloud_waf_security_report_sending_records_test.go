package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWafSecurityReportSendingRecords_basic(t *testing.T) {
	dataSource := "data.huaweicloud_waf_security_report_sending_records.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the test data in advance before running the test cases.
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceWafSecurityReportSendingRecords_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.report_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.subscription_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.report_name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.stat_period.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.report_category"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.sending_time"),

					resource.TestCheckOutput("is_report_name_filter_useful", "true"),
					resource.TestCheckOutput("is_report_category_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceDataSourceWafSecurityReportSendingRecords_basic = `
data "huaweicloud_waf_security_report_sending_records" "test" {}

locals {
  report_name = data.huaweicloud_waf_security_report_sending_records.test.items[0].report_name
}

data "huaweicloud_waf_security_report_sending_records" "report_name_filter" {
  report_name = local.report_name
}

output "is_report_name_filter_useful" {
  value = length(data.huaweicloud_waf_security_report_sending_records.report_name_filter.items) > 0 && alltrue(
    [for item in data.huaweicloud_waf_security_report_sending_records.report_name_filter.items :
	item.report_name == local.report_name]
  )
}

locals {
  report_category = data.huaweicloud_waf_security_report_sending_records.test.items[0].report_category
}

data "huaweicloud_waf_security_report_sending_records" "report_category_filter" {
  report_category = local.report_category
}

output "is_report_category_filter_useful" {
  value = length(data.huaweicloud_waf_security_report_sending_records.report_category_filter.items) > 0 && alltrue(
    [for item in data.huaweicloud_waf_security_report_sending_records.report_category_filter.items :
	item.report_category == local.report_category]
  )
}
`
