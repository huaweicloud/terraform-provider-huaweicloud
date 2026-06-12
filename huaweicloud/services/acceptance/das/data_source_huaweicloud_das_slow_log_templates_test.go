package das

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccSlowLogTemplates_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_das_slow_log_templates.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByTemplateId   = "data.huaweicloud_das_slow_log_templates.filter_by_template_id"
		dcFilterByTemplateId = acceptance.InitDataSourceCheck(filterByTemplateId)

		filterByNodeId   = "data.huaweicloud_das_slow_log_templates.filter_by_node_id"
		dcFilterByNodeId = acceptance.InitDataSourceCheck(filterByNodeId)

		filterByDbName   = "data.huaweicloud_das_slow_log_templates.filter_by_db_name"
		dcFilterByDbName = acceptance.InitDataSourceCheck(filterByDbName)

		filterByMinAvgExecTime   = "data.huaweicloud_das_slow_log_templates.filter_by_min_avg_execute_time"
		dcFilterByMinAvgExecTime = acceptance.InitDataSourceCheck(filterByMinAvgExecTime)

		filterByMaxAvgExecTime   = "data.huaweicloud_das_slow_log_templates.filter_by_max_avg_execute_time"
		dcFilterByMaxAvgExecTime = acceptance.InitDataSourceCheck(filterByMaxAvgExecTime)

		filterByOperation   = "data.huaweicloud_das_slow_log_templates.filter_by_operation"
		dcFilterByOperation = acceptance.InitDataSourceCheck(filterByOperation)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDasInstanceIds(t, 1)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSlowLogTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "templates.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "templates.0.template_id"),
					resource.TestCheckResourceAttrSet(all, "templates.0.template_name"),
					resource.TestCheckResourceAttrSet(all, "templates.0.execute_count"),
					resource.TestCheckResourceAttrSet(all, "templates.0.avg_execute_time"),
					resource.TestCheckResourceAttrSet(all, "templates.0.max_execute_time"),
					resource.TestCheckResourceAttrSet(all, "templates.0.db_names.#"),

					// filter by template_id
					dcFilterByTemplateId.CheckResourceExists(),
					resource.TestCheckOutput("is_template_id_filter_useful", "true"),

					// filter by node_id
					dcFilterByNodeId.CheckResourceExists(),
					resource.TestCheckOutput("is_node_id_filter_useful", "true"),

					// filter by db_name
					dcFilterByDbName.CheckResourceExists(),
					resource.TestCheckOutput("is_db_name_filter_useful", "true"),

					// filter by min_avg_execute_time
					dcFilterByMinAvgExecTime.CheckResourceExists(),
					resource.TestCheckOutput("is_min_avg_execute_time_filter_useful", "true"),

					// filter by max_avg_execute_time
					dcFilterByMaxAvgExecTime.CheckResourceExists(),
					resource.TestCheckOutput("is_max_avg_execute_time_filter_useful", "true"),

					// filter by operation
					dcFilterByOperation.CheckResourceExists(),
					resource.TestCheckOutput("is_operation_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccSlowLogTemplates_base() string {
	// 1. The maximum difference between `start_time` and `end_time` is `12` hours.
	// 2. The `end_time` must be greater than the `start_time`.
	// 3. The `start_time` and `end_time` only support year, month, day, and hour, but not minute and second.
	currentTime := time.Now().Truncate(time.Hour)
	startTime := currentTime.Add(-5 * time.Hour).Format(time.RFC3339)
	endTime := currentTime.Add(5 * time.Hour).Format(time.RFC3339)

	return fmt.Sprintf(`
locals {
  instance_ids = split(",", "%[1]s")
  start_time   = "%[2]s"
  end_time     = "%[3]s"
}
`, acceptance.HW_DAS_INSTANCE_IDS, startTime, endTime)
}

func testAccSlowLogTemplates_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_das_slow_log_templates" "all" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
}

# Filter by template_id
locals {
  template_id = try(data.huaweicloud_das_slow_log_templates.all.templates[0].template_id, "")
}

data "huaweicloud_das_slow_log_templates" "filter_by_template_id" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  template_id = local.template_id
}

locals {
  template_id_filter_result = length(local.template_id) > 0 ? [
    for v in data.huaweicloud_das_slow_log_templates.filter_by_template_id.templates :
    v.template_id == local.template_id
  ] : []
}

output "is_template_id_filter_useful" {
  value = length(local.template_id) == 0 || (length(local.template_id_filter_result) > 0 && alltrue(local.template_id_filter_result))
}

# Filter by node_id
locals {
  node_id_from_all = try(data.huaweicloud_das_slow_log_templates.all.templates[0].node_ids[0], "")
}

data "huaweicloud_das_slow_log_templates" "filter_by_node_id" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  node_id     = local.node_id_from_all
}

locals {
  node_id_filter_result = length(local.node_id_from_all) > 0 ? [
    for v in data.huaweicloud_das_slow_log_templates.filter_by_node_id.templates :
    contains(v.node_ids, local.node_id_from_all)
  ] : []
}

output "is_node_id_filter_useful" {
  value = length(local.node_id_from_all) == 0 || (length(local.node_id_filter_result) > 0 && alltrue(local.node_id_filter_result))
}

# Filter by db_name
locals {
  db_name = try(data.huaweicloud_das_slow_log_templates.all.templates[0].db_names[0], "")
}

data "huaweicloud_das_slow_log_templates" "filter_by_db_name" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  db_name     = local.db_name
}

locals {
  db_name_filter_result = length(local.db_name) > 0 ? [
    for v in data.huaweicloud_das_slow_log_templates.filter_by_db_name.templates : contains(v.db_names, local.db_name)
  ] : []
}

output "is_db_name_filter_useful" {
  value = length(local.db_name) == 0 || (length(local.db_name_filter_result) > 0 && alltrue(local.db_name_filter_result))
}

# Filter by min_avg_execute_time
locals {
  min_avg_execute_time = 3.0
}

data "huaweicloud_das_slow_log_templates" "filter_by_min_avg_execute_time" {
  instance_id          = local.instance_ids[0]
  start_time           = local.start_time
  end_time             = local.end_time
  min_avg_execute_time = local.min_avg_execute_time
}

locals {
  min_avg_execute_time_filter_result = [
    for v in data.huaweicloud_das_slow_log_templates.filter_by_min_avg_execute_time.templates :
    v.avg_execute_time >= local.min_avg_execute_time
  ]
}

output "is_min_avg_execute_time_filter_useful" {
  value = length(local.min_avg_execute_time_filter_result) > 0 && alltrue(local.min_avg_execute_time_filter_result)
}

# Filter by max_avg_execute_time
locals {
  max_avg_execute_time = 6000.0
}

data "huaweicloud_das_slow_log_templates" "filter_by_max_avg_execute_time" {
  instance_id          = local.instance_ids[0]
  start_time           = local.start_time
  end_time             = local.end_time
  max_avg_execute_time = local.max_avg_execute_time
}

locals {
  max_avg_execute_time_filter_result = [
    for v in data.huaweicloud_das_slow_log_templates.filter_by_max_avg_execute_time.templates :
    v.avg_execute_time <= local.max_avg_execute_time
  ]
}

output "is_max_avg_execute_time_filter_useful" {
  value = length(local.max_avg_execute_time_filter_result) > 0 && alltrue(local.max_avg_execute_time_filter_result)
}

# Filter by operation
locals {
  operation = try(data.huaweicloud_das_slow_log_templates.all.templates[0].template_name, "")
  operation_value = length(local.operation) > 0 ? split(" ", local.operation)[0] : ""
}

data "huaweicloud_das_slow_log_templates" "filter_by_operation" {
  instance_id = local.instance_ids[0]
  start_time  = local.start_time
  end_time    = local.end_time
  operation   = local.operation_value
}

locals {
  operation_filter_result = length(local.operation_value) > 0 ? [
    for v in data.huaweicloud_das_slow_log_templates.filter_by_operation.templates :
    strcontains(v.template_name, local.operation_value)
  ] : []
}

output "is_operation_filter_useful" {
  value = length(local.operation_value) == 0 || (length(local.operation_filter_result) > 0 && alltrue(local.operation_filter_result))
}
`, testAccSlowLogTemplates_base())
}
