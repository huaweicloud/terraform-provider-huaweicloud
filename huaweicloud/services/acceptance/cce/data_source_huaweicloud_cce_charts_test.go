package cce

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccChartsDataSource_basic(t *testing.T) {
	datasourceName := "data.huaweicloud_cce_charts.test"
	dc := acceptance.InitDataSourceCheck(datasourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      dc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testCharts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
				),
			},
		},
	})
}

const testCharts_basic = `
data "huaweicloud_cce_charts" "test" {}

output "is_results_not_empty" {
  value = length(data.huaweicloud_cce_charts.test.charts) > 0
}
`
