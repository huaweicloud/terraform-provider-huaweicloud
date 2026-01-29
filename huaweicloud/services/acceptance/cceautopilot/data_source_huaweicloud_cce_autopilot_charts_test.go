package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCharts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_charts.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCceChartPath(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCharts_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "charts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.values"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.instruction"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.version"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.source"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.icon_url"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.chart_url"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.create_at"),
					resource.TestCheckResourceAttrSet(dataSource, "charts.0.update_at"),
				),
			},
		},
	})
}

func testAccDataSourceCharts_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_autopilot_charts" "test" {
  depends_on = [huaweicloud_cce_autopilot_chart.test]
}
`, testAccAutopilotChart_basic())
}
