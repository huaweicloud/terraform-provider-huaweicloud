package antiddos

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceWeeklyProtectionStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_antiddos_weekly_protection_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byPeriodStartDate   = "data.huaweicloud_antiddos_weekly_protection_statistics.test_period_start_date"
		dcByPeriodStartDate = acceptance.InitDataSourceCheck(byPeriodStartDate)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceWeeklyProtectionStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "ddos_intercept_times"),
					resource.TestCheckResourceAttrSet(dataSource, "weekdata.#"),
					resource.TestCheckResourceAttrSet(dataSource, "weekdata.0.ddos_blackhole_times"),
					resource.TestCheckResourceAttrSet(dataSource, "weekdata.0.ddos_intercept_times"),
					resource.TestCheckResourceAttrSet(dataSource, "weekdata.0.max_attack_bps"),
					resource.TestCheckResourceAttrSet(dataSource, "weekdata.0.max_attack_conns"),
					resource.TestCheckResourceAttrSet(dataSource, "weekdata.0.period_start_date"),

					dcByPeriodStartDate.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byPeriodStartDate, "ddos_intercept_times"),
					resource.TestCheckResourceAttrSet(byPeriodStartDate, "weekdata.#"),
					resource.TestCheckResourceAttrSet(byPeriodStartDate, "weekdata.0.ddos_blackhole_times"),
					resource.TestCheckResourceAttrSet(byPeriodStartDate, "weekdata.0.ddos_intercept_times"),
					resource.TestCheckResourceAttrSet(byPeriodStartDate, "weekdata.0.max_attack_bps"),
					resource.TestCheckResourceAttrSet(byPeriodStartDate, "weekdata.0.max_attack_conns"),
					resource.TestCheckResourceAttrSet(byPeriodStartDate, "weekdata.0.period_start_date"),
				),
			},
		},
	})
}

func testDataSourceWeeklyProtectionStatistics_basic() string {
	now := time.Now().UnixMilli()

	return fmt.Sprintf(`
data "huaweicloud_antiddos_weekly_protection_statistics" "test" {
}

data "huaweicloud_antiddos_weekly_protection_statistics" "test_period_start_date" {
  period_start_date = "%d"
}
`, now)
}
