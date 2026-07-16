package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbRealTimeTransactions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_real_time_transactions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbRealTimeTransactions_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rows.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.sessionid"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.pid"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.query_id"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.datname"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.client_addr"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.client_port"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.usename"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.query"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.backend_start"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.xact_start"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.application_name"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.state_change"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.query_start"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.waiting"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.unique_sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.top_xid"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.current_xid"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.exec_time"),
					resource.TestCheckResourceAttrSet(dataSource, "rows.0.xlog_quantity"),
				),
			},
		},
	})
}

func testDataSourceGaussdbRealTimeTransactions_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_instance_real_time_sessions" "test" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}

data "huaweicloud_gaussdb_real_time_transactions" "test" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}

data "huaweicloud_gaussdb_real_time_transactions" "query_info_filter" {
  instance_id  = "%[1]s"
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id

  transaction_query_info {
    xlog_quantity = "0"
    client_addrs  = ["100.125.31.5"]
    datnames      = ["postgres"]
    usenames      = ["root"]
  }
}

output "query_info_filter" {
  value = length(data.huaweicloud_gaussdb_real_time_transactions.query_info_filter.rows) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID)
}
