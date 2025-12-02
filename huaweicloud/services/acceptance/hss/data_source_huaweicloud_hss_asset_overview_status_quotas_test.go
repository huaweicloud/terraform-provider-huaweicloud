package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewStatusQuotas_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_overview_status_quotas.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case need a host with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetOverviewStatusQuotas_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.used_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.total_num"),
				),
			},
		},
	})
}

const testAccDataSourceAssetOverviewStatusQuotas_basic = `
data "huaweicloud_hss_asset_overview_status_quotas" "test" {
  enterprise_project_id = "all_granted_eps"
}
`
