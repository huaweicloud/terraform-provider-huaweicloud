package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDataMapDynamicDataInfos_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_data_map_dynamic_data_infos.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataMapDynamicDataInfos_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "vpc_dbss_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "vpc_dbss_list.0.vpc_id"),
					resource.TestCheckResourceAttrSet(dataSource, "vpc_dbss_list.0.dbss.#"),
				),
			},
		},
	})
}

const testDataSourceDataMapDynamicDataInfos_basic = `
data "huaweicloud_dsc_data_map_dynamic_data_infos" "test" {}
`
