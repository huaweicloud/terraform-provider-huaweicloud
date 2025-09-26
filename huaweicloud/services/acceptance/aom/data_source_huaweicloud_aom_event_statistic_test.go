package aom

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEventStatistic_basic(t *testing.T) {
	var (
		all           = "data.huaweicloud_aom_event_statistic.all"
		dcForAllStats = acceptance.InitDataSourceCheck(all)

		allWithStep           = "data.huaweicloud_aom_event_statistic.all_with_step"
		dcForAllWithStepStats = acceptance.InitDataSourceCheck(allWithStep)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEventStatistic_basic,
				Check: resource.ComposeTestCheckFunc(
					dcForAllStats.CheckResourceExists(),
					// Check the attributes.
					resource.TestCheckResourceAttrSet(all, "step_result"),
					resource.TestCheckResourceAttrSet(all, "timestamps.#"),
					resource.TestCheckResourceAttrSet(all, "series.#"),
					resource.TestCheckResourceAttrSet(all, "summary.%"),
					// Check series structure.
					resource.TestCheckResourceAttrSet(all, "series.0.event_severity"),
					resource.TestCheckResourceAttrSet(all, "series.0.values.#"),
					// With step.
					dcForAllWithStepStats.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(allWithStep, "step_result"),
					resource.TestCheckResourceAttrSet(allWithStep, "timestamps.#"),
					resource.TestCheckResourceAttrSet(allWithStep, "series.#"),
					resource.TestCheckResourceAttrSet(allWithStep, "summary.%"),
					// Check series structure.
					resource.TestCheckResourceAttrSet(allWithStep, "series.0.event_severity"),
					resource.TestCheckResourceAttrSet(allWithStep, "series.0.values.#"),
				),
			},
		},
	})
}

const testAccDataSourceEventStatistic_basic = `
data "huaweicloud_aom_event_statistic" "all" {
  time_range = "-1.-1.60"
}

data "huaweicloud_aom_event_statistic" "all_with_step" {
  time_range = "-1.-1.60"
  step       = 90000
}
`
