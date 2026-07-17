package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbTopSqlStatements_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_top_sql_statements.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbTopSqlStatements_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.user_name"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.sql_text"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.calls_percent"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.cpu_percent"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.io_percent"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.calls"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.returned_rows"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.tuple_read"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.avg_elapse_time"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.total_elapse_time"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.cpu_time"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.io_time"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.min_elapse_time"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.max_elapse_time"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.sql_hit_ratio"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "top_sql_infos.0.db_name"),
					resource.TestCheckOutput("query_filter", "true"),
					resource.TestCheckOutput("multi_queries_filter", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbTopSqlStatements_basic() string {
	now := time.Now().UTC()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	startTime := startOfDay.UnixMilli()

	endOfDay := startOfDay.Add(24 * time.Hour)
	endTime := endOfDay.UnixMilli()

	return fmt.Sprintf(`
data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = "%[1]s"
}

data "huaweicloud_gaussdb_top_sql_statements" "test" {
  instance_id = "%[1]s"
  node_ids    = [data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id]
  start_time  = "%[2]d"
  end_time    = "%[3]d"
}

data "huaweicloud_gaussdb_top_sql_statements" "query_filter" {
  instance_id    = "%[1]s"
  node_ids       = [data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id]
  start_time     = "%[2]d"
  end_time       = "%[3]d"
  support_system = true
  db_name        = "postgres"
}

output "query_filter" {
  value = length(data.huaweicloud_gaussdb_top_sql_statements.query_filter.top_sql_infos) > 0
}

data "huaweicloud_gaussdb_top_sql_statements" "multi_queries_filter" {
  instance_id    = "%[1]s"
  node_ids       = [data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id]
  start_time     = "%[2]d"
  end_time       = "%[3]d"
  support_system = true
  db_name        = "postgres"

  multi_queries {
    name      = "query"
    condition = "and"
    values    = ["SELECT"]
    is_fuzzy  = "true"
  }
}

output "multi_queries_filter" {
  value = length(data.huaweicloud_gaussdb_top_sql_statements.multi_queries_filter.top_sql_infos) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
