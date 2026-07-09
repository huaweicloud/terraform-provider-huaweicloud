package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSqlTraces_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_sql_traces.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckGaussDBSqlExecId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSqlTraces_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "traces.#"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.sql_exec_id"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.finish_time"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.all_time"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.client_addr"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.client_port"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.is_slow_sql"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.execution_time_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.execution_time_details.0.resource_time.#"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.execution_time_details.0.kernel_time.#"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.execution_time_details.0.kernel_execution_time.#"),
					resource.TestCheckResourceAttrSet(dataSource, "traces.0.execution_time_details.0.wait_event_time.#"),

					resource.TestCheckOutput("sql_exec_id_filter_useful", "true"),
					resource.TestCheckOutput("transaction_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSqlTraces_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_sql_traces" "test" {
  instance_id = "%[1]s"
  sql_exec_id = "%[2]s"
}

locals {
  transaction_id = data.huaweicloud_gaussdb_sql_traces.test.traces[0].transaction_id
}

data "huaweicloud_gaussdb_sql_traces" "transaction_id_filter" {
  instance_id    = "%[1]s"
  transaction_id = local.transaction_id
}

output "sql_exec_id_filter_useful" {
  value = length(data.huaweicloud_gaussdb_sql_traces.test.traces) > 0
}

output "transaction_id_filter_useful" {
  value = length(data.huaweicloud_gaussdb_sql_traces.transaction_id_filter.traces) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_SQL_EXEC_ID)
}
