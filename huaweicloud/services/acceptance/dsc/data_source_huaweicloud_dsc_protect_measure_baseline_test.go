package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceProtectMeasureBaseline_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_protect_measure_baseline.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceProtectMeasureBaseline_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.#"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.key"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.detail.#"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.detail.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.detail.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.detail.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.detail.0.protect_level"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.detail.0.data_type_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "protect_measure_baseline.0.detail.0.measure_info.#"),
				),
			},
		},
	})
}

const testDataSourceProtectMeasureBaseline_basic = `
data "huaweicloud_dsc_protect_measure_baseline" "test" {}
`
