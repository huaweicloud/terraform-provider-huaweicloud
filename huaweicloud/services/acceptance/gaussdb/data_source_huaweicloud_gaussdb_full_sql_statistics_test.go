package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceFullSqlStatistics_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_full_sql_statistics.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
		startTime  = time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05+0800")
		endTime    = time.Now().Format("2006-01-02T15:04:05+0800")
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceFullSqlStatistics_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.#"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.template"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.sql_count"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.total_sql_time"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.avg_db_time"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.start_time_stamp"),
					resource.TestCheckResourceAttrSet(dataSource, "statistics.0.end_time_stamp"),

					resource.TestCheckOutput("sql_id_filter_useful", "true"),
					resource.TestCheckOutput("username_filter_useful", "true"),
					resource.TestCheckOutput("db_name_filter_useful", "true"),
					resource.TestCheckOutput("transaction_id_filter_useful", "true"),
					resource.TestCheckOutput("order_filter_useful", "true"),
					resource.TestCheckOutput("multi_queries_filter_useful", "true"),
					resource.TestCheckOutput("compare_conditions_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceFullSqlStatistics_basic(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_full_sql_statistics" "test" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  sql_id = data.huaweicloud_gaussdb_full_sql_statistics.test.statistics[0].sql_id
}

data "huaweicloud_gaussdb_full_sql_statistics" "sql_id_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  sql_id      = local.sql_id
}

output "sql_id_filter_useful" {
  value = length(data.huaweicloud_gaussdb_full_sql_statistics.sql_id_filter.statistics) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_full_sql_statistics.sql_id_filter.statistics[*].sql_id : v == local.sql_id]
  )
}

data "huaweicloud_gaussdb_full_sql_statistics" "username_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  username    = "root"
}

output "username_filter_useful" {
  value = length(data.huaweicloud_gaussdb_full_sql_statistics.username_filter.statistics) > 0
}

data "huaweicloud_gaussdb_full_sql_statistics" "db_name_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  db_name     = "postgres"
}

output "db_name_filter_useful" {
  value = length(data.huaweicloud_gaussdb_full_sql_statistics.db_name_filter.statistics) > 0
}

data "huaweicloud_gaussdb_full_sql_statistics" "transaction_id_filter" {
  instance_id    = "%[1]s"
  begin_time     = "%[2]s"
  end_time       = "%[3]s"
  transaction_id = "0"
}

output "transaction_id_filter_useful" {
  value = length(data.huaweicloud_gaussdb_full_sql_statistics.transaction_id_filter.statistics) > 0
}

data "huaweicloud_gaussdb_full_sql_statistics" "order_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  order_by    = "sql_count"
  order       = "ASC"
}

locals {
  len     = length(data.huaweicloud_gaussdb_full_sql_statistics.test.statistics)
  value_1 = data.huaweicloud_gaussdb_full_sql_statistics.test.statistics[local.len - 1].sql_count
  value_2 = data.huaweicloud_gaussdb_full_sql_statistics.order_filter.statistics[0].sql_count
}

output "order_filter_useful" {
  value = (local.value_1 == local.value_2)
}

data "huaweicloud_gaussdb_full_sql_statistics" "multi_queries_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"

  multi_queries {
    name      = "query"
    condition = "OR"
    values    = ["select"]
    is_fuzzy  = "true"
  }
}

output "multi_queries_filter_useful" {
  value = length(data.huaweicloud_gaussdb_full_sql_statistics.multi_queries_filter.statistics) > 0
}

data "huaweicloud_gaussdb_full_sql_statistics" "compare_conditions_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"

  compare_conditions {
    name         = "sql_count"
    enable_equal = "true"
    min          = "10"
    max          = "15"
  }
}

output "compare_conditions_filter_useful" {
  value = length(data.huaweicloud_gaussdb_full_sql_statistics.compare_conditions_filter.statistics) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
