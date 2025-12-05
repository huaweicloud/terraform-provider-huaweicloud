package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOperationalReportNotification_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_operational_report_notification.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOperationalReportNotification_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "title"),
					resource.TestCheckResourceAttrSet(dataSourceName, "report_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "current_month"),
				),
			},
		},
	})
}

func testDataSourceOperationalReportNotification_basic() string {
	return `
data "huaweicloud_hss_operational_report_notification" "test" {}
`
}
