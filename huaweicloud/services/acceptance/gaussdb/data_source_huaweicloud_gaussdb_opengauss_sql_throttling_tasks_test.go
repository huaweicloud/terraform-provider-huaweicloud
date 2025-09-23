package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceGaussdbOpengaussSqlThrottlingTasks_basic(t *testing.T) {
	dataSource := "data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussInstanceId(t)
			acceptance.TestAccPreCheckGaussDBOpenGaussTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceGaussdbOpengaussSqlThrottlingTasks_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.task_id"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.task_name"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.task_scope"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.limit_type"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.limit_type_value"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.sql_model"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.rule_name"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.parallel_size"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.creator"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.node_infos.#"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.node_infos.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "limit_task_list.0.node_infos.0.sql_id"),
					resource.TestCheckOutput("task_scope_filter_is_useful", "true"),
					resource.TestCheckOutput("limit_type_filter_is_useful", "true"),
					resource.TestCheckOutput("limit_type_value_filter_is_useful", "true"),
					resource.TestCheckOutput("task_name_filter_is_useful", "true"),
					resource.TestCheckOutput("sql_model_filter_is_useful", "true"),
					resource.TestCheckOutput("rule_name_filter_is_useful", "true"),
					resource.TestCheckOutput("start_time_filter_is_useful", "true"),
					resource.TestCheckOutput("end_time_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceGaussdbOpengaussSqlThrottlingTasks_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_gaussdb_opengauss_instances" "test" {}

locals {
  instance       = [for v in  data.huaweicloud_gaussdb_opengauss_instances.test.instances : v if v.id == "%[1]s"][0]
  master_node_id = [for v in local.instance.nodes : v if v.role == "master"][0].id
}

data "huaweicloud_gaussdb_opengauss_sql_templates" "test" {
  instance_id = "%[1]s"
  node_id     = local.master_node_id
}

resource "huaweicloud_gaussdb_opengauss_sql_throttling_task" "test" {
  instance_id      = "%[1]s"
  task_scope       = "SQL"
  limit_type       = "SQL_ID"
  limit_type_value = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
  task_name        = "%[2]s"
  parallel_size    = 10
  start_time       = "%[3]s"
  end_time         = "%[4]s"
  sql_model        = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_model
  node_infos {
    node_id = local.master_node_id
    sql_id  = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
  }
}

`, acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, name, acceptance.HW_GAUSSDB_OPENGAUSS_START_TIME, acceptance.HW_GAUSSDB_OPENGAUSS_END_TIME)
}

func testDataSourceGaussdbOpengaussSqlThrottlingTasks_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "test" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
}

locals {
  task_scope = "SQL"
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "task_scope_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
  task_scope  = "SQL"
}

output "task_scope_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.task_scope_filter.limit_task_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.task_scope_filter.limit_task_list[*].task_scope :
  v == local.task_scope]
  )
}

locals {
  limit_type = "SQL_ID"
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "limit_type_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
  limit_type  = "SQL_ID"
}

output "limit_type_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.limit_type_filter.limit_task_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.limit_type_filter.limit_task_list[*].limit_type :
  v == local.limit_type]
  )
}

locals {
  limit_type_value = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "limit_type_value_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id      = "%[2]s"
  limit_type_value = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_id
}

output "limit_type_value_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.limit_type_value_filter.limit_task_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.limit_type_value_filter.limit_task_list[*].limit_type_value :
  v == local.limit_type_value]
  )
}

locals {
  task_name = "%[3]s"
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "task_name_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
  task_name   = "%[3]s"
}

output "task_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.task_name_filter.limit_task_list) > 0
}

locals {
  sql_model = data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_model
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "sql_model_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
  sql_model   =  data.huaweicloud_gaussdb_opengauss_sql_templates.test.node_limit_sql_model_list[1].sql_model
}

output "sql_model_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.sql_model_filter.limit_task_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.sql_model_filter.limit_task_list[*].sql_model :
  v == local.sql_model]
  )
}

locals {
  rule_name = data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.test.limit_task_list[0].rule_name
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "rule_name_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
  rule_name   = data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.test.limit_task_list[0].rule_name
}

output "rule_name_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.rule_name_filter.limit_task_list) > 0 && alltrue(
  [for v in data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.rule_name_filter.limit_task_list[*].rule_name :
  v == local.rule_name]
  )
}

locals {
  start_time = "%[4]s"
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "start_time_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
  start_time  = "%[4]s"
}

output "start_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.start_time_filter.limit_task_list) > 0
}

locals {
  end_time = "%[5]s"
}

data "huaweicloud_gaussdb_opengauss_sql_throttling_tasks" "end_time_filter" {
  depends_on = [huaweicloud_gaussdb_opengauss_sql_throttling_task.test]

  instance_id = "%[2]s"
  end_time    = "%[5]s"
}

output "end_time_filter_is_useful" {
  value = length(data.huaweicloud_gaussdb_opengauss_sql_throttling_tasks.end_time_filter.limit_task_list) > 0
}
`, testDataSourceGaussdbOpengaussSqlThrottlingTasks_base(name), acceptance.HW_GAUSSDB_OPENGAUSS_INSTANCE_ID, name,
		acceptance.HW_GAUSSDB_OPENGAUSS_START_TIME, acceptance.HW_GAUSSDB_OPENGAUSS_END_TIME)
}
