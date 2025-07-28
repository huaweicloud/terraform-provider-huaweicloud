package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetManualCollect_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_manual_collect.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetManualCollect_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "scan_status"),
					resource.TestCheckResourceAttrSet(dataSource, "scanned_time"),
				),
			},
		},
	})
}

func testAccDataSourceAssetManualCollect_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_hss_asset_manual_collect" "test" {
  type                  = "kernel-module"
  host_id               = "%s"
  enterprise_project_id = "0"
}
`, acceptance.HW_HSS_HOST_PROTECTION_HOST_ID)
}
