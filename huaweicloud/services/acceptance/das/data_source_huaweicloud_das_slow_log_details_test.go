package das

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSlowLogDetails_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_slow_log_details.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByNodeIDs   = "data.huaweicloud_das_slow_log_details.filter_by_node_ids"
		dcFilterByNodeIDs = acceptance.InitDataSourceCheck(filterByNodeIDs)

		filterByDbName   = "data.huaweicloud_das_slow_log_details.filter_by_db_name"
		dcFilterByDbName = acceptance.InitDataSourceCheck(filterByDbName)

		filterByClientIPAddress   = "data.huaweicloud_das_slow_log_details.filter_by_client_ip_address"
		dcFilterByClientIPAddress = acceptance.InitDataSourceCheck(filterByClientIPAddress)

		filterByUserName   = "data.huaweicloud_das_slow_log_details.filter_by_user_name"
		dcFilterByUserName = acceptance.InitDataSourceCheck(filterByUserName)

		filterByKilled   = "data.huaweicloud_das_slow_log_details.filter_by_killed"
		dcFilterByKilled = acceptance.InitDataSourceCheck(filterByKilled)

		filterByExecTimeMin   = "data.huaweicloud_das_slow_log_details.filter_by_execute_time_min"
		dcFilterByExecTimeMin = acceptance.InitDataSourceCheck(filterByExecTimeMin)

		filterByExecTimeMax   = "data.huaweicloud_das_slow_log_details.filter_by_execute_time_max"
		dcFilterByExecTimeMax = acceptance.InitDataSourceCheck(filterByExecTimeMax)

		filterByRowsMaxExamined   = "data.huaweicloud_das_slow_log_details.filter_by_rows_max_examined"
		dcFilterByRowsMaxExamined = acceptance.InitDataSourceCheck(filterByRowsMaxExamined)

		filterByRowsMinExamined   = "data.huaweicloud_das_slow_log_details.filter_by_rows_min_examined"
		dcFilterByRowsMinExamined = acceptance.InitDataSourceCheck(filterByRowsMinExamined)

		filterByFuzzySql   = "data.huaweicloud_das_slow_log_details.filter_by_fuzzy_sql"
		dcFilterByFuzzySql = acceptance.InitDataSourceCheck(filterByFuzzySql)

		filterByOperation   = "data.huaweicloud_das_slow_log_details.filter_by_operation"
		dcFilterByOperation = acceptance.InitDataSourceCheck(filterByOperation)

		filterBySortField   = "data.huaweicloud_das_slow_log_details.filter_by_sort_field"
		dcFilterBySortField = acceptance.InitDataSourceCheck(filterBySortField)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSlowLogDetails_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "details.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "details.0.sql_template_id"),
					resource.TestCheckResourceAttrSet(all, "details.0.original_sql"),
					resource.TestCheckResourceAttrSet(all, "details.0.occurrence_time"),
					resource.TestCheckResourceAttrSet(all, "details.0.db_name"),
					resource.TestCheckResourceAttrSet(all, "details.0.client_ip_address"),
					resource.TestCheckResourceAttrSet(all, "details.0.user_name"),
					resource.TestCheckResourceAttrSet(all, "details.0.killed"),
					resource.TestCheckResourceAttrSet(all, "details.0.execute_time"),
					resource.TestCheckResourceAttrSet(all, "details.0.rows_examined"),
					resource.TestCheckResourceAttrSet(all, "details.0.rows_sent"),
					resource.TestCheckResourceAttrSet(all, "details.0.sql_type"),
					resource.TestCheckResourceAttrSet(all, "details.0.tunable"),

					// filter by node_ids
					dcFilterByNodeIDs.CheckResourceExists(),
					resource.TestCheckOutput("is_node_ids_filter_useful", "true"),

					// filter by db_name
					dcFilterByDbName.CheckResourceExists(),
					resource.TestCheckOutput("is_db_name_filter_useful", "true"),

					// filter by client_ip_address
					dcFilterByClientIPAddress.CheckResourceExists(),
					resource.TestCheckOutput("is_client_ip_address_filter_useful", "true"),

					// filter by user_name
					dcFilterByUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),

					// filter by killed
					dcFilterByKilled.CheckResourceExists(),
					resource.TestCheckOutput("is_killed_filter_useful", "true"),

					// filter by execute_time_min
					dcFilterByExecTimeMin.CheckResourceExists(),
					resource.TestCheckOutput("is_execute_time_min_filter_useful", "true"),

					// filter by execute_time_max
					dcFilterByExecTimeMax.CheckResourceExists(),
					resource.TestCheckOutput("is_execute_time_max_filter_useful", "true"),

					// filter by rows_max_examined
					dcFilterByRowsMaxExamined.CheckResourceExists(),
					resource.TestCheckOutput("is_rows_max_examined_filter_useful", "true"),

					// filter by rows_min_examined
					dcFilterByRowsMinExamined.CheckResourceExists(),
					resource.TestCheckOutput("is_rows_min_examined_filter_useful", "true"),

					// filter by fuzzy_sql
					dcFilterByFuzzySql.CheckResourceExists(),
					resource.TestCheckOutput("is_fuzzy_sql_filter_useful", "true"),

					// filter by operation
					dcFilterByOperation.CheckResourceExists(),
					resource.TestCheckOutput("is_operation_filter_useful", "true"),

					// filter by sort_field
					dcFilterBySortField.CheckResourceExists(),
					resource.TestCheckOutput("is_sort_field_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSlowLogDetails_base() string {
	// 1.The earliest `start_time` is at most two days earlier than the current time
	// 2.The latest `end_time` is at most one day later than the current time
	// 3.The `end_time` must be greater than the `start_time`
	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, -1).Format(time.RFC3339)
	endTime := currentTime.AddDate(0, 0, 1).Format(time.RFC3339)

	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
  start_time   = "%[2]s"
  end_time     = "%[3]s"
}
`, acceptance.HW_DAS_INSTANCE_IDS, startTime, endTime)
}

func testAccDataSlowLogDetails_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_slow_log_details" "all" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
}

# Filter by node_ids
locals {
  node_id_from_all = try(data.huaweicloud_das_slow_log_details.all.details[0].node_id, "")
  node_ids         = local.node_id_from_all != "" ? [local.node_id_from_all] : []
}

data "huaweicloud_das_slow_log_details" "filter_by_node_ids" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  node_ids    = local.node_ids
}

locals {
  node_ids_filter_result = length(local.node_ids) > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_node_ids.details : contains(local.node_ids, v.node_id)
  ] : []
}

output "is_node_ids_filter_useful" {
  value = length(local.node_ids) == 0 || (length(local.node_ids_filter_result) > 0 && alltrue(local.node_ids_filter_result))
}

# Filter by db_name
locals {
  db_name = try(data.huaweicloud_das_slow_log_details.all.details[0].db_name, "")
}

data "huaweicloud_das_slow_log_details" "filter_by_db_name" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  db_name     = local.db_name
}

locals {
  db_name_filter_result = [
    for v in data.huaweicloud_das_slow_log_details.filter_by_db_name.details : v.db_name == local.db_name
  ]
}

output "is_db_name_filter_useful" {
  value = length(local.db_name_filter_result) > 0 && alltrue(local.db_name_filter_result)
}

# Filter by client_ip_address
locals {
  client_ip_address = try(data.huaweicloud_das_slow_log_details.all.details[0].client_ip_address, "")
}

data "huaweicloud_das_slow_log_details" "filter_by_client_ip_address" {
  instance_id       = local.instance_ids[0]
  start_time        = local.start_time
  end_time          = local.end_time
  client_ip_address = local.client_ip_address
}

locals {
  client_ip_address_filter_result = length(local.client_ip_address) > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_client_ip_address.details : v.client_ip_address == local.client_ip_address
  ] : []
}

output "is_client_ip_address_filter_useful" {
  value = length(local.client_ip_address) == 0 || (
    length(local.client_ip_address_filter_result) > 0 && alltrue(local.client_ip_address_filter_result)
  )
}

# Filter by user_name
locals {
  user_name = try(data.huaweicloud_das_slow_log_details.all.details[0].user_name, "")
}

data "huaweicloud_das_slow_log_details" "filter_by_user_name" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  user_name   = local.user_name
}

locals {
  user_name_filter_result = [
    for v in data.huaweicloud_das_slow_log_details.filter_by_user_name.details : v.user_name == local.user_name
  ]
}

output "is_user_name_filter_useful" {
  value = length(local.user_name_filter_result) > 0 && alltrue(local.user_name_filter_result)
}

# Filter by killed
locals {
  killed = try(data.huaweicloud_das_slow_log_details.all.details[0].killed, "")
}

data "huaweicloud_das_slow_log_details" "filter_by_killed" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  killed      = local.killed
}

locals {
  killed_filter_result = length(local.killed) > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_killed.details : v.killed == local.killed
  ] : []
}

output "is_killed_filter_useful" {
  value = length(local.killed) == 0 || (length(local.killed_filter_result) > 0 && alltrue(local.killed_filter_result))
}

# Filter by execute_time_min
locals {
  execute_time_min = 3
}

data "huaweicloud_das_slow_log_details" "filter_by_execute_time_min" {
  instance_id      = local.instance_ids[0]
  start_time       = local.start_time
  end_time         = local.end_time
  execute_time_min = local.execute_time_min
}

locals {
  execute_time_min_filter_result = local.execute_time_min > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_execute_time_min.details : v.execute_time * 1000 >= local.execute_time_min
  ] : []
}

output "is_execute_time_min_filter_useful" {
  value = local.execute_time_min == 0 || (length(local.execute_time_min_filter_result) > 0 && alltrue(local.execute_time_min_filter_result))
}

# Filter by execute_time_max
locals {
  execute_time_max = 6000
}

data "huaweicloud_das_slow_log_details" "filter_by_execute_time_max" {
  instance_id      = local.instance_ids[0]
  start_time       = local.start_time
  end_time         = local.end_time
  execute_time_max = local.execute_time_max
}

locals {
  execute_time_max_filter_result = local.execute_time_max > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_execute_time_max.details : v.execute_time * 1000 <= local.execute_time_max
  ] : []
}

output "is_execute_time_max_filter_useful" {
  value = local.execute_time_max == 0 || (length(local.execute_time_max_filter_result) > 0 && alltrue(local.execute_time_max_filter_result))
}

# Filter by rows_max_examined
locals {
  rows_max_examined = try(data.huaweicloud_das_slow_log_details.all.details[0].rows_examined, 0)
}

data "huaweicloud_das_slow_log_details" "filter_by_rows_max_examined" {
  instance_id       = local.instance_ids[0]
  start_time        = local.start_time
  end_time          = local.end_time
  rows_max_examined = local.rows_max_examined
}

locals {
  rows_max_examined_filter_result = local.rows_max_examined > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_rows_max_examined.details : v.rows_examined <= local.rows_max_examined
  ] : []
}

output "is_rows_max_examined_filter_useful" {
  value = local.rows_max_examined == 0 || (length(local.rows_max_examined_filter_result) > 0 && alltrue(local.rows_max_examined_filter_result))
}

# Filter by rows_min_examined
locals {
  rows_min_examined = try(data.huaweicloud_das_slow_log_details.all.details[0].rows_examined, 0)
}

data "huaweicloud_das_slow_log_details" "filter_by_rows_min_examined" {
  instance_id       = local.instance_ids[0]
  start_time        = local.start_time
  end_time          = local.end_time
  rows_min_examined = local.rows_min_examined
}

locals {
  rows_min_examined_filter_result = local.rows_min_examined > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_rows_min_examined.details : v.rows_examined >= local.rows_min_examined
  ] : []
}

output "is_rows_min_examined_filter_useful" {
  value = local.rows_min_examined == 0 || (length(local.rows_min_examined_filter_result) > 0 && alltrue(local.rows_min_examined_filter_result))
}

# Filter by fuzzy_sql
locals {
  fuzzy_sql = try(split(" ", data.huaweicloud_das_slow_log_details.all.details[0].original_sql)[0], "")
}

data "huaweicloud_das_slow_log_details" "filter_by_fuzzy_sql" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  fuzzy_sql   = local.fuzzy_sql
}

locals {
  fuzzy_sql_filter_result = length(local.fuzzy_sql) > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_fuzzy_sql.details : strcontains(v.original_sql, local.fuzzy_sql)
  ] : []
}

output "is_fuzzy_sql_filter_useful" {
  value = length(local.fuzzy_sql) == 0 || (length(local.fuzzy_sql_filter_result) > 0 && alltrue(local.fuzzy_sql_filter_result))
}

# Filter by operation
locals {
  operation = try(data.huaweicloud_das_slow_log_details.all.details[0].sql_type, "")
}

data "huaweicloud_das_slow_log_details" "filter_by_operation" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  operation   = local.operation
}

locals {
  operation_filter_result = length(local.operation) > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_operation.details : contains(split(",", local.operation), v.sql_type)
  ] : []
}

output "is_operation_filter_useful" {
  value = length(local.operation) == 0 || (length(local.operation_filter_result) > 0 && alltrue(local.operation_filter_result))
}

# Filter by sort_field
locals {
  sort_field = "executeTime"
}

data "huaweicloud_das_slow_log_details" "filter_by_sort_field" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  sort_field  = local.sort_field
  sort_asc    = false
}

locals {
  sort_field_filter_result = length(local.sort_field) > 0 ? [
    for v in data.huaweicloud_das_slow_log_details.filter_by_sort_field.details : v.execute_time != null
  ] : []
}

output "is_sort_field_filter_useful" {
  value = length(local.sort_field) == 0 || (length(local.sort_field_filter_result) > 0 && alltrue(local.sort_field_filter_result))
}
`, testAccDataSlowLogDetails_base())
}
