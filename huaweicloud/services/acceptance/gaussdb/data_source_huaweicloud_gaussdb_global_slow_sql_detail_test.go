package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbGlobalSlowSqlDetail_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_global_slow_sql_detail.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussDbGlobalSlowSqlDetail_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.#"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.client_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.client_port"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.sql_text"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.query_plan"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.finish_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.returned_rows"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.fetched_rows"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.fetched_pages"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.hit_pages"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.total_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.cpu_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.plan_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.io_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.lock_count"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_sql_details.0.lock_time"),
				),
			},
		},
	})
}

func testDataSourceGaussDbGlobalSlowSqlDetail_basic() string {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startTime := startOfDay.UnixMilli()

	endOfDay := startOfDay.Add(24 * time.Hour)
	endTime := endOfDay.UnixMilli()

	return fmt.Sprintf(`
data "huaweicloud_gaussdb_instance_nodes" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_slow_sql_list" "test_list" {
  instance_id = "%[1]s"
  start_time  = "%[2]d"
  end_time    = "%[3]d"
  threshold   = 5
  node_ids    = data.huaweicloud_gaussdb_instance_nodes.test.nodes[*].id
}

data "huaweicloud_gaussdb_slow_sql_detail" "test_detail" {
  instance_id = "%[1]s"
  start_time  = "%[2]d"
  end_time    = "%[3]d"
  sql_id      = data.huaweicloud_gaussdb_slow_sql_list.test_list.slow_sql_infos.0.sql_id
  node_ids    = data.huaweicloud_gaussdb_instance_nodes.test.nodes[*].id
}

data "huaweicloud_gaussdb_global_slow_sql_detail" "test" {
  instance_id    = "%[1]s"
  start_time     = "%[2]d"
  end_time       = "%[3]d"
  sql_id         = data.huaweicloud_gaussdb_slow_sql_detail.test_detail.slow_sql_details.0.sql_id
  node_ids       = data.huaweicloud_gaussdb_instance_nodes.test.nodes[*].id
  component_type = "cn"
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
