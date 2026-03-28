package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceIntelligentSessionKillStatistic_basic(t *testing.T) {
	rName := "data.huaweicloud_rds_intelligent_session_kill_statistic.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRdsInstanceId(t)
			acceptance.TestAccPreCheckRdsNodeId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceIntelligentSessionKillStatistic_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "statistics.#"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.keyword"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.raw_sql_text"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.ids.#"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.count"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.total_time"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.avg_time"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.max_time"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.strategy"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.advice_concurrency"),
					resource.TestCheckResourceAttrSet(rName, "statistics.0.type"),
					resource.TestCheckOutput("node_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceIntelligentSessionKillStatistic_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_rds_intelligent_session_kill_statistic" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_rds_intelligent_session_kill_statistic" "node_id_filter" {
  instance_id = "%[1]s"
  node_id     = "%[2]s"
}
output "node_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_intelligent_session_kill_statistic.node_id_filter.statistics) > 0
}
`, acceptance.HW_RDS_INSTANCE_ID, acceptance.HW_RDS_NODE_ID)
}
