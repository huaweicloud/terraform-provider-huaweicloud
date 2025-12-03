package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOverviewHotInformation_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_overview_hot_information.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOverviewHotInformation_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "total_num"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.title"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.update_time"),
				),
			},
		},
	})
}

func testDataSourceOverviewHotInformation_basic() string {
	return `
data "huaweicloud_hss_overview_hot_information" "test" {}
`
}
