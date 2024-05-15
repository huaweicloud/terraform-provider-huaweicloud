package rds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsPgSqlLimits_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_pg_sql_limits.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceRdsPgSqlLimits_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "sql_limits.#"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_limits.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_limits.0.query_id"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_limits.0.max_concurrency"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_limits.0.max_waiting"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_limits.0.search_path"),
					resource.TestCheckResourceAttrSet(dataSource, "sql_limits.0.is_effective"),

					resource.TestCheckOutput("sql_limit_id_filter_is_useful", "true"),
					resource.TestCheckOutput("query_id_filter_is_useful", "true"),
					resource.TestCheckOutput("max_concurrency_filter_is_useful", "true"),
					resource.TestCheckOutput("max_waiting_filter_is_useful", "true"),
					resource.TestCheckOutput("search_path_filter_is_useful", "true"),
					resource.TestCheckOutput("is_effective_filter_is_useful", "true"),
					resource.TestCheckOutput("query_string_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceRdsPgSqlLimits_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_rds_pg_sql_limit" "test" {
  depends_on = [huaweicloud_rds_pg_plugin.test]

  instance_id     = huaweicloud_rds_instance.test.id
  db_name         = "%[2]s"
  query_id        = "100"
  max_concurrency = 20
  max_waiting     = 5
  search_path     = "public"
  switch          = "open"

  lifecycle {
    ignore_changes = [query_string]
  }
}

resource "huaweicloud_rds_pg_sql_limit" "query_string" {
  depends_on = [huaweicloud_rds_pg_plugin.test]

  instance_id     = huaweicloud_rds_instance.test.id
  db_name         = "%[2]s"
  query_string    = "select"
  max_concurrency = 20
  max_waiting     = 5
  search_path     = "public"

  lifecycle {
    ignore_changes = [query_id]
  }
}

data "huaweicloud_rds_pg_sql_limits" "test" {
  depends_on = [huaweicloud_rds_pg_sql_limit.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = "%[2]s"
}

locals {
  sql_limit_id = huaweicloud_rds_pg_sql_limit.test.sql_limit_id
}
data "huaweicloud_rds_pg_sql_limits" "sql_limit_id_filter" {
  depends_on = [huaweicloud_rds_pg_sql_limit.test]

  instance_id  = huaweicloud_rds_instance.test.id
  db_name      = "%[2]s"
  sql_limit_id = huaweicloud_rds_pg_sql_limit.test.sql_limit_id
}
output "sql_limit_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_sql_limits.sql_limit_id_filter.sql_limits) > 0 && alltrue(
  [for v in data.huaweicloud_rds_pg_sql_limits.sql_limit_id_filter.sql_limits[*].id : v == local.sql_limit_id]
  )
}

locals {
  query_id = huaweicloud_rds_pg_sql_limit.test.query_id
}
data "huaweicloud_rds_pg_sql_limits" "query_id_filter" {
  depends_on = [huaweicloud_rds_pg_sql_limit.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = "%[2]s"
  query_id    = huaweicloud_rds_pg_sql_limit.test.query_id
}
output "query_id_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_sql_limits.query_id_filter.sql_limits) > 0 && alltrue(
  [for v in data.huaweicloud_rds_pg_sql_limits.query_id_filter.sql_limits[*].query_id : v == local.query_id]
  )
}

locals {
  max_concurrency = huaweicloud_rds_pg_sql_limit.test.max_concurrency
}
data "huaweicloud_rds_pg_sql_limits" "max_concurrency_filter" {
  depends_on = [huaweicloud_rds_pg_sql_limit.test]

  instance_id     = huaweicloud_rds_instance.test.id
  db_name         = "%[2]s"
  max_concurrency = huaweicloud_rds_pg_sql_limit.test.max_concurrency
}
output "max_concurrency_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_sql_limits.max_concurrency_filter.sql_limits) > 0 && alltrue(
  [for v in data.huaweicloud_rds_pg_sql_limits.max_concurrency_filter.sql_limits[*].max_concurrency : v == local.max_concurrency]
  )
}

locals {
  max_waiting = huaweicloud_rds_pg_sql_limit.test.max_waiting
}
data "huaweicloud_rds_pg_sql_limits" "max_waiting_filter" {
  depends_on = [huaweicloud_rds_pg_sql_limit.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = "%[2]s"
  max_waiting = huaweicloud_rds_pg_sql_limit.test.max_waiting
}
output "max_waiting_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_sql_limits.max_waiting_filter.sql_limits) > 0 && alltrue(
  [for v in data.huaweicloud_rds_pg_sql_limits.max_waiting_filter.sql_limits[*].max_waiting : v == local.max_waiting]
  )
}

locals {
  search_path = huaweicloud_rds_pg_sql_limit.test.search_path
}
data "huaweicloud_rds_pg_sql_limits" "search_path_filter" {
  depends_on = [huaweicloud_rds_pg_sql_limit.test]

  instance_id = huaweicloud_rds_instance.test.id
  db_name     = "%[2]s"
  search_path = huaweicloud_rds_pg_sql_limit.test.search_path
}
output "search_path_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_sql_limits.search_path_filter.sql_limits) > 0 && alltrue(
  [for v in data.huaweicloud_rds_pg_sql_limits.search_path_filter.sql_limits[*].search_path : v == local.search_path]
  )
}

locals {
  is_effective = huaweicloud_rds_pg_sql_limit.test.is_effective
}
data "huaweicloud_rds_pg_sql_limits" "is_effective_filter" {
  depends_on = [huaweicloud_rds_pg_sql_limit.test]

  instance_id  = huaweicloud_rds_instance.test.id
  db_name      = "%[2]s"
  is_effective = huaweicloud_rds_pg_sql_limit.test.is_effective
}
output "is_effective_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_sql_limits.is_effective_filter.sql_limits) > 0 && alltrue(
  [for v in data.huaweicloud_rds_pg_sql_limits.is_effective_filter.sql_limits[*].is_effective : v == local.is_effective]
  )
}

locals {
  query_string = huaweicloud_rds_pg_sql_limit.query_string.query_string
}
data "huaweicloud_rds_pg_sql_limits" "query_string_filter" {
  depends_on = [huaweicloud_rds_pg_sql_limit.query_string]

  instance_id  = huaweicloud_rds_instance.test.id
  db_name      = "%[2]s"
  query_string = huaweicloud_rds_pg_sql_limit.query_string.query_string
}
output "query_string_filter_is_useful" {
  value = length(data.huaweicloud_rds_pg_sql_limits.query_string_filter.sql_limits) > 0 && alltrue(
  [for v in data.huaweicloud_rds_pg_sql_limits.query_string_filter.sql_limits[*].query_string : v == local.query_string]
  )
}
`, testPgSqlLimit_base(name), name)
}
