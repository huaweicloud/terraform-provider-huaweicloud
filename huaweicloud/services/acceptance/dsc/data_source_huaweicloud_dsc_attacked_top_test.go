package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscAttackedTop_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_attacked_top.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDscEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDscAttackedTop_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "attacked_api_num"),
					resource.TestCheckResourceAttrSet(dataSource, "attacked_api_top.#"),
					resource.TestCheckResourceAttrSet(dataSource, "attacked_asset_num"),
					resource.TestCheckResourceAttrSet(dataSource, "attacked_asset_top.#"),
				),
			},
		},
	})
}

const testDataSourceDscAttackedTop_basic = `
data "huaweicloud_dsc_attacked_top" "test" {}
`
