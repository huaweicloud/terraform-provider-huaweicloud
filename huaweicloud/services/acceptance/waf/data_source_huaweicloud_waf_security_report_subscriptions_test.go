package waf

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityReportSubscriptions_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_security_report_subscriptions.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test, create a security report.
			acceptance.TestAccPrecheckWafInstance(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityReportSubscriptions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "items.#"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.subscription_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.report_name"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.report_category"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.report_status"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.sending_period"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.is_all_enterprise_project"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.template_eps_id"),
					resource.TestCheckResourceAttrSet(dataSource, "items.0.is_report_created"),

					resource.TestCheckOutput("is_report_name_filter_useful", "true"),
					resource.TestCheckOutput("is_report_category_filter_useful", "true"),
					resource.TestCheckOutput("is_report_status_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSecurityReportSubscriptions_basic = `
data "huaweicloud_waf_security_report_subscriptions" "test" {}

locals {
  report_name = data.huaweicloud_waf_security_report_subscriptions.test.items[0].report_name
}

data "huaweicloud_waf_security_report_subscriptions" "report_name_filter" {
  report_name = local.report_name
}

output "is_report_name_filter_useful" {
  value = length(data.huaweicloud_waf_security_report_subscriptions.report_name_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_security_report_subscriptions.report_name_filter.items[*].report_name : v == local.report_name]
  )
}

locals {
  report_category = data.huaweicloud_waf_security_report_subscriptions.test.items[0].report_category
}

data "huaweicloud_waf_security_report_subscriptions" "report_category_filter" {	
  report_category = local.report_category
}

output "is_report_category_filter_useful" {
  value = length(data.huaweicloud_waf_security_report_subscriptions.report_category_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_security_report_subscriptions.report_category_filter.items[*].report_category : v == local.report_category]
  )
}

locals {
  report_status = data.huaweicloud_waf_security_report_subscriptions.test.items[0].report_status
}

data "huaweicloud_waf_security_report_subscriptions" "report_status_filter" {
  report_status = local.report_status
}

output "is_report_status_filter_useful" {
  value = length(data.huaweicloud_waf_security_report_subscriptions.report_status_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_security_report_subscriptions.report_status_filter.items[*].report_status : v == local.report_status]
  )
}

locals {
  eps = data.huaweicloud_waf_security_report_subscriptions.test.items[0].enterprise_project_id
}

data "huaweicloud_waf_security_report_subscriptions" "eps_filter" {
  enterprise_project_id = local.eps
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_waf_security_report_subscriptions.eps_filter.items) > 0 && alltrue(
    [for v in data.huaweicloud_waf_security_report_subscriptions.eps_filter.items[*].enterprise_project_id : v == local.eps]
  )
}
`
