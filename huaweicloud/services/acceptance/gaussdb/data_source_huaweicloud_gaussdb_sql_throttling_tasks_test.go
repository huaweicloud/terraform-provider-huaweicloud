package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussDbSqlThrottlingTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_sql_throttling_tasks.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBInstanceId(t)
			acceptance.TestAccPreCheckGaussDBSqlId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussSqlThrottlingTasks_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.task_name"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.task_scope"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.limit_type"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.limit_type_value"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.sql_model"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.parallel_size"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.node_infos.#"),

					resource.TestCheckOutput("task_scope_filter_is_useful", "true"),
					resource.TestCheckOutput("limit_type_filter_is_useful", "true"),
					resource.TestCheckOutput("task_name_filter_is_useful", "true"),
					resource.TestCheckOutput("sql_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussSqlThrottlingTasks_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_sql_throttling_tasks" "test" {
  instance_id = "%[1]s"
}

locals {
  task_scope = data.huaweicloud_gaussdb_sql_throttling_tasks.test.limit_task_list[0].task_scope
}

data "huaweicloud_gaussdb_sql_throttling_tasks" "task_scope_filter" {
  instance_id = "%[1]s"
  task_scope  = local.task_scope
}

output "task_scope_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_sql_throttling_tasks.task_scope_filter.limit_task_list) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_sql_throttling_tasks.task_scope_filter.limit_task_list[*].task_scope :
    v == local.task_scope]
  )
}

locals {
  limit_type = data.huaweicloud_gaussdb_sql_throttling_tasks.test.limit_task_list[0].limit_type
}

data "huaweicloud_gaussdb_sql_throttling_tasks" "limit_type_filter" {
  instance_id = "%[1]s"
  limit_type  = local.limit_type
}

output "limit_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_sql_throttling_tasks.limit_type_filter.limit_task_list) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_sql_throttling_tasks.limit_type_filter.limit_task_list[*].limit_type :
    v == local.limit_type]
  )
}

locals {
  task_name = data.huaweicloud_gaussdb_sql_throttling_tasks.test.limit_task_list[0].task_name
}

data "huaweicloud_gaussdb_sql_throttling_tasks" "task_name_filter" {
  instance_id = "%[1]s"
  task_name   = local.task_name
}

output "task_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_sql_throttling_tasks.task_name_filter.limit_task_list) > 0 && alltrue(
    [for v in data.huaweicloud_gaussdb_sql_throttling_tasks.task_name_filter.limit_task_list[*].task_name :
    v == local.task_name]
  )
}

data "huaweicloud_gaussdb_sql_throttling_tasks" "sql_id_filter" {
  instance_id = "%[1]s"
  sql_id      = "%[2]s"
}

output "sql_id_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_sql_throttling_tasks.sql_id_filter.limit_task_list) > 0
}
`, acceptance.HW_GAUSSDB_INSTANCE_ID, acceptance.HW_GAUSSDB_SQL_ID)
}
