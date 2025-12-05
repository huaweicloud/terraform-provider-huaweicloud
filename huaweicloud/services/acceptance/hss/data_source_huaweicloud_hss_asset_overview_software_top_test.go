package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetOverviewSoftwareTop_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_asset_overview_software_top.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case needs to ensure the existence of a host with host protection enabled,
			// and the host is under the default enterprise project.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAssetOverviewSoftwareTop_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.percentage"),
				),
			},
		},
	})
}

func testDataSourceAssetOverviewSoftwareTop_basic() string {
	return `
data "huaweicloud_hss_asset_overview_software_top" "test" {
  enterprise_project_id = "0"
}
`
}
