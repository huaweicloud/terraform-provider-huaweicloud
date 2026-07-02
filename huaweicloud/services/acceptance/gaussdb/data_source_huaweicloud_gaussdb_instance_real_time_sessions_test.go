package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDBInstanceRealTimeSessions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_instance_real_time_sessions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGaussDBInstanceRealTimeSessions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.session_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.pid"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.unique_sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.database_name"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.client_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.wait"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.wait_event"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.query_runtime"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.query"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.back_end_start"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.query_start"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.application_name"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.exec_time"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.trans_num"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.rollback_num"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.sql_num"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.client_port"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.query_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.transaction_time_cost"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.top_transaction_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.current_transaction_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.xlog_quantity_pretty"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.wait_status"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.lwt_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.thread_name"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.smp_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sessions.0.component_name"),
					resource.TestCheckOutput("name_filter", "true"),
				),
			},
		},
	})
}

func testAccDataSourceGaussDBInstanceRealTimeSessions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_instance_real_time_sessions" "test" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}

data "huaweicloud_gaussdb_instance_real_time_sessions" "name_filter" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id

  query_info {
    database_name = "postgres"
  }
}

output "name_filter" {
  value = length(data.huaweicloud_gaussdb_instance_real_time_sessions.name_filter.sessions) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
