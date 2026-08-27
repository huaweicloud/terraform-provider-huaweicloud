package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsSlowLogStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_rds_slow_log_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRdsSlowLogStatistics_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_log_list.#"),
				),
			},
		},
	})
}

func testDataSourceRdsSlowLogStatistics_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_slow_log_statistics" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_START_TIME, acceptance.HW_RDS_END_TIME)
}
