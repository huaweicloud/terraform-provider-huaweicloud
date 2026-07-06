package gaussdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSingleFullSqls_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_gaussdb_single_full_sqls.test"
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
				Config: testAccDataSourceSingleFullSqls_basic(startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.#"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.db_name"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.username"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.component_id"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.schema_name"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.sql_id"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.sql_exec_id"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.query"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.db_time"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.cpu_time"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.client_addr"),
					resource.TestCheckResourceAttrSet(dataSource, "full_sqls.0.client_port"),

					resource.TestCheckOutput("node_id_filter_useful", "true"),
					resource.TestCheckOutput("sql_id_filter_useful", "true"),
					resource.TestCheckOutput("db_name_filter_useful", "true"),
					resource.TestCheckOutput("client_addr_filter_useful", "true"),
					resource.TestCheckOutput("client_port_filter_useful", "true"),
					resource.TestCheckOutput("multi_queries_filter_useful", "true"),
					resource.TestCheckOutput("compare_conditions_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceSingleFullSqls_basic(startTime, endTime string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_single_full_sqls" "test" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  node_id     = data.huaweicloud_gaussdb_single_full_sqls.test.full_sqls[0].node_id
  sql_id      = data.huaweicloud_gaussdb_single_full_sqls.test.full_sqls[0].sql_id
  db_name     = data.huaweicloud_gaussdb_single_full_sqls.test.full_sqls[0].db_name
  client_addr = data.huaweicloud_gaussdb_single_full_sqls.test.full_sqls[0].client_addr
  client_port = data.huaweicloud_gaussdb_single_full_sqls.test.full_sqls[0].client_port
}

data "huaweicloud_gaussdb_single_full_sqls" "node_id_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  node_id     = local.node_id
}

output "node_id_filter_useful" {
  value = length(data.huaweicloud_gaussdb_single_full_sqls.node_id_filter.full_sqls) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_single_full_sqls.node_id_filter.full_sqls[*].node_id : v == local.node_id]
  )
}

data "huaweicloud_gaussdb_single_full_sqls" "sql_id_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  sql_id      = local.sql_id
}

output "sql_id_filter_useful" {
  value = length(data.huaweicloud_gaussdb_single_full_sqls.sql_id_filter.full_sqls) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_single_full_sqls.sql_id_filter.full_sqls[*].sql_id : v == local.sql_id]
  )
}

data "huaweicloud_gaussdb_single_full_sqls" "db_name_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  db_name     = local.db_name
}

output "db_name_filter_useful" {
  value = length(data.huaweicloud_gaussdb_single_full_sqls.db_name_filter.full_sqls) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_single_full_sqls.db_name_filter.full_sqls[*].db_name : v == local.db_name]
  )
}

data "huaweicloud_gaussdb_single_full_sqls" "client_addr_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  client_addr = local.client_addr
}

output "client_addr_filter_useful" {
  value = length(data.huaweicloud_gaussdb_single_full_sqls.client_addr_filter.full_sqls) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_single_full_sqls.client_addr_filter.full_sqls[*].client_addr : v == local.client_addr]
  )
}

data "huaweicloud_gaussdb_single_full_sqls" "client_port_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"
  client_port = local.client_port
}

output "client_port_filter_useful" {
  value = length(data.huaweicloud_gaussdb_single_full_sqls.client_port_filter.full_sqls) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_single_full_sqls.client_port_filter.full_sqls[*].client_port : v == local.client_port]
  )
}

data "huaweicloud_gaussdb_single_full_sqls" "multi_queries_filter" {
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
  value = length(data.huaweicloud_gaussdb_single_full_sqls.multi_queries_filter.full_sqls) > 0
}

data "huaweicloud_gaussdb_single_full_sqls" "compare_conditions_filter" {
  instance_id = "%[1]s"
  begin_time  = "%[2]s"
  end_time    = "%[3]s"

  compare_conditions {
    name         = "db_time"
    enable_equal = "true"
    min          = "0"
    max          = "2"
  }
}

output "compare_conditions_filter_useful" {
  value = length(data.huaweicloud_gaussdb_single_full_sqls.compare_conditions_filter.full_sqls) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, startTime, endTime)
}
