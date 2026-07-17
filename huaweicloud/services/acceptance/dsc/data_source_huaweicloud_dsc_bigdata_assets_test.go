package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscBigdataAssets_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_bigdata_assets.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDscBigdataAssets_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "assets.#"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.asset_name"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.authorized"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.default"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ds_address"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ds_authorized"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ds_name"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ds_port"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ds_type"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ds_version"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ins_id"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ins_name"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.ins_type"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "assets.0.source_type"),
				),
			},
		},
	})
}

const testDataSourceDscBigdataAssets_basic = `
data "huaweicloud_dsc_bigdata_assets" "test" {}
`
