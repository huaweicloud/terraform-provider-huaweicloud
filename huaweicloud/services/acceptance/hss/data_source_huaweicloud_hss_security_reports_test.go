package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityReports_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_security_reports.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSecurityReports_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.report_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.report_sub_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.default_report"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.latest_create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.report_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.report_category"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.report_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.report_create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.sending_period"),

					resource.TestCheckOutput("is_report_category_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceSecurityReports_basic() string {
	return `
data "huaweicloud_hss_security_reports" "test" {
  enterprise_project_id = "all_granted_eps"
}

# Filter by report_category
locals {
  report_category = data.huaweicloud_hss_security_reports.test.data_list[0].report_category
}

data "huaweicloud_hss_security_reports" "report_category_filter" {
  enterprise_project_id = "all_granted_eps"
  report_category       = local.report_category
}

output "is_report_category_filter_useful" {
  value = length(data.huaweicloud_hss_security_reports.report_category_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_security_reports.report_category_filter.data_list[*].report_category : v == local.report_category]
  )
}
`
}
