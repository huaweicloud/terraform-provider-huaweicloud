package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewStatusHostProtection_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_overview_status_host_protection.test"
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
				Config: testAccDataSourceAssetOverviewStatusHostProtection_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "no_risk"),
					resource.TestCheckResourceAttrSet(dataSource, "risk"),
					resource.TestCheckResourceAttrSet(dataSource, "no_protect"),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
				),
			},
		},
	})
}

const testAccDataSourceAssetOverviewStatusHostProtection_basic = `
data "huaweicloud_hss_asset_overview_status_host_protection" "test" {
  enterprise_project_id = "all_granted_eps"
}
`
