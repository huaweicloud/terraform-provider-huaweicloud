package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSqlExecutionPlanStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_sql_execution_plan_status.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceSqlExecutionPlanStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sql_plan_bind_state_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_plan_bind_state_list.0.outline"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_plan_bind_state_list.0.cost"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_plan_bind_state_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_plan_bind_state_list.0.sql_hash"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_plan_bind_state_list.0.plan_id"),
				),
			},
		},
	})
}

func testDataSourceSqlExecutionPlanStatus_basic() string {
	startTime := time.Now().UTC().Add(-8 * time.Hour).UnixMilli()
	endTime := time.Now().UTC().Add(8 * time.Hour).UnixMilli()
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_nodes" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_slow_sql_list" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]d"
  end_time    = "%[3]d"
  threshold   = 1
  node_ids    = data.huaweicloud_gaussdb_instance_nodes.test.nodes[*].id
}

data "huaweicloud_gaussdb_sql_execution_plan_status" "test" {
  instance_id = "%[1]s"
  node_id     = data.huaweicloud_gaussdb_slow_sql_list.test.slow_sql_infos[0].node_id
  sql_id      = data.huaweicloud_gaussdb_slow_sql_list.test.slow_sql_infos[0].sql_id
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
