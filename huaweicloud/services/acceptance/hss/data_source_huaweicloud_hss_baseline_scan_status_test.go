package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceBaselineScanStatus_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_baseline_scan_status.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBaselineScanStatus_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "scan_status"),
					resource.TestCheckResourceAttrSet(dataSource, "scanned_time"),
				),
			},
		},
	})
}

const testAccDataSourceBaselineScanStatus_basic = `
data "huaweicloud_hss_baseline_scan_status" "test" {
  enterprise_project_id = "all_granted_eps"
}
`
