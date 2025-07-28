package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAntivirusPayPerScanSwitchStatus_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_antivirus_pay_per_scan_switch_status.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAntivirusPayPerScanSwitchStatus_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "enabled"),
				),
			},
		},
	})
}

const testAccDataSourceAntivirusPayPerScanSwitchStatus_basic string = `
data "huaweicloud_hss_antivirus_pay_per_scan_switch_status" "test" {
  enterprise_project_id = "all_granted_eps"
}
`
