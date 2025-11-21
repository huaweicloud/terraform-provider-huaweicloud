package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityReportSubscription_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_waf_security_report_subscription.test"
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
				Config: testAccDataSourceSecurityReportSubscription_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sending_period"),
					resource.TestCheckResourceAttrSet(dataSource, "report_name"),
					resource.TestCheckResourceAttrSet(dataSource, "report_category"),
					resource.TestCheckResourceAttrSet(dataSource, "topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "subscription_type"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.overview_statistics_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.group_by_day_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.request_statistics_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.qps_statistics_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.bandwidth_statistics_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.response_code_statistics_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.attack_type_distribution_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.top_attacked_domains_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.top_attack_source_ips_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.top_attacked_urls_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.top_attack_source_locations_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "report_content_subscription.0.top_abnormal_urls_enable"),
					resource.TestCheckResourceAttrSet(dataSource, "is_all_enterprise_project"),
					resource.TestCheckResourceAttrSet(dataSource, "enterprise_project_id"),
				),
			},
		},
	})
}

func testAccDataSourceSecurityReportSubscription_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_waf_security_report_subscription" "test" {
  subscription_id = "%s"
}
`, acceptance.HW_WAF_SECURITY_REPORT_SUBSCRIPTION_ID)
}
