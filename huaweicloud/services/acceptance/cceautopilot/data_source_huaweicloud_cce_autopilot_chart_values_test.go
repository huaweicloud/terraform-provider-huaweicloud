package cceautopilot

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCCEAutopilotShowChartValues_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cce_autopilot_chart_values.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCCEAutopilotShowChartValues_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "values.%"),
				),
			},
		},
	})
}

func testAccDataSourceCCEAutopilotShowChartValues_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_cce_autopilot_chart_values" "test" {
 chart_id = huaweicloud_cce_autopilot_chart.test.id
}
`, testAccAutopilotChart_basic())
}
