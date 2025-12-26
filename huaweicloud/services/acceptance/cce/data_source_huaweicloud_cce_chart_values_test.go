package cce

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHuaweiCloudCceShowChartValues_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_chart_values.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHuaweiCloudCceShowChartValues_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "id"),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
					resource.TestCheckResourceAttrSet(dataSource, "values"),
				),
			},
		},
	})
}

const testDataSourceHuaweiCloudCceShowChartValues_basic = `

data "huaweicloud_cce_chart_values" "test" {
  chart_id = "399882cb-4d7a-4450-a8ca-112f278db59c"
}
`
