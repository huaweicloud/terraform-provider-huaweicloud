package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAntivirusPayPerScanFreeQuotas_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_hss_antivirus_pay_per_scan_free_quotas.test"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAntivirusPayPerScanFreeQuotas_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "free_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "occupied_free_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "antivirus_already_given"),
					resource.TestCheckResourceAttrSet(dataSourceName, "antivirus_free_quota"),
					resource.TestCheckResourceAttrSet(dataSourceName, "report_already_given"),
					resource.TestCheckResourceAttrSet(dataSourceName, "report_free_quota"),
				),
			},
		},
	})
}

func testAccDataSourceAntivirusPayPerScanFreeQuotas_basic() string {
	return `
data "huaweicloud_hss_antivirus_pay_per_scan_free_quotas" "test" {
  enterprise_project_id = "0"
}
`
}
