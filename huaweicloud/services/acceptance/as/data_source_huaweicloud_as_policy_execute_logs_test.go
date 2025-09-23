package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPolicyExecuteLogsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_policy_execute_logs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byLogID   = "data.huaweicloud_as_policy_execute_logs.log_id_filter"
		dcByLogID = acceptance.InitDataSourceCheck(byLogID)

		byResourceID   = "data.huaweicloud_as_policy_execute_logs.resource_id_filter"
		dcByResourceID = acceptance.InitDataSourceCheck(byResourceID)

		byResourceType   = "data.huaweicloud_as_policy_execute_logs.resource_type_filter"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)

		byExecuteType   = "data.huaweicloud_as_policy_execute_logs.execute_type_filter"
		dcByExecuteType = acceptance.InitDataSourceCheck(byExecuteType)

		byStatus   = "data.huaweicloud_as_policy_execute_logs.status_filter"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byStartTime   = "data.huaweicloud_as_policy_execute_logs.start_time_filter"
		dcByStartTime = acceptance.InitDataSourceCheck(byStartTime)

		byEndTime   = "data.huaweicloud_as_policy_execute_logs.end_time_filter"
		dcByEndTime = acceptance.InitDataSourceCheck(byEndTime)

		byStartEndTime   = "data.huaweicloud_as_policy_execute_logs.start_end_time_filter"
		dcByStartEndTime = acceptance.InitDataSourceCheck(byStartEndTime)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the AS policy containing the execute logs in advance and configure the AS policy ID into
			// the environment variable.
			acceptance.TestAccPreCheckASScalingPolicyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyExecuteLogsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.desire_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.execute_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.execute_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.job_records.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.limit_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.old_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.scaling_policy_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.scaling_resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.scaling_resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "execute_logs.0.type"),

					dcByLogID.CheckResourceExists(),
					resource.TestCheckOutput("is_log_id_filter_useful", "true"),

					dcByResourceID.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),

					dcByExecuteType.CheckResourceExists(),
					resource.TestCheckOutput("is_execute_type_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					dcByStartTime.CheckResourceExists(),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),

					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),

					dcByStartEndTime.CheckResourceExists(),
					resource.TestCheckOutput("start_end_time_filter_result", "true"),
				),
			},
		},
	})
}

func testAccPolicyExecuteLogsDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_policy_execute_logs" "test" {
  scaling_policy_id = "%[1]s"
}

// Filter using log_id.
locals {
  log_id = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].id
}

data "huaweicloud_as_policy_execute_logs" "log_id_filter" {
  scaling_policy_id = "%[1]s"
  log_id            = local.log_id
}

locals {
  log_id_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.log_id_filter.execute_logs[*].id : v == local.log_id
  ]
}

output "is_log_id_filter_useful" {
  value = length(local.log_id_filter_result) > 0 && alltrue(local.log_id_filter_result)
}

// Filter using scaling_resource_id.
locals {
  scaling_resource_id = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].scaling_resource_id
}

data "huaweicloud_as_policy_execute_logs" "resource_id_filter" {
  scaling_policy_id   = "%[1]s"
  scaling_resource_id = local.scaling_resource_id
}

locals {
  scaling_resource_id_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.resource_id_filter.execute_logs[*].scaling_resource_id :
    v == local.scaling_resource_id
  ]
}

output "is_resource_id_filter_useful" {
  value = length(local.scaling_resource_id_filter_result) > 0 && alltrue(local.scaling_resource_id_filter_result)
}

// Filter using scaling_resource_type.
locals {
  scaling_resource_type = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].scaling_resource_type
}

data "huaweicloud_as_policy_execute_logs" "resource_type_filter" {
  scaling_policy_id     = "%[1]s"
  scaling_resource_type = local.scaling_resource_type
}

locals {
  scaling_resource_type_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.resource_type_filter.execute_logs[*].scaling_resource_type :
    v == local.scaling_resource_type
  ]
}

output "is_resource_type_filter_useful" {
  value = length(local.scaling_resource_type_filter_result) > 0 && alltrue(local.scaling_resource_type_filter_result)
}

// Filter using execute_type.
locals {
  execute_type = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].execute_type
}

data "huaweicloud_as_policy_execute_logs" "execute_type_filter" {
  scaling_policy_id = "%[1]s"
  execute_type      = local.execute_type
}

locals {
  execute_type_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.execute_type_filter.execute_logs[*].execute_type : v == local.execute_type
  ]
}

output "is_execute_type_filter_useful" {
  value = length(local.execute_type_filter_result) > 0 && alltrue(local.execute_type_filter_result)
}

// Filter using status.
locals {
  status = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].status
}

data "huaweicloud_as_policy_execute_logs" "status_filter" {
  scaling_policy_id = "%[1]s"
  status            = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.status_filter.execute_logs[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}

// Filter using start_time and end_time.
locals {
  execute_time         = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].execute_time
  start_time_before24h = timeadd(local.execute_time, "-24h")
  end_time_after24h    = timeadd(local.execute_time, "24h")
}

data "huaweicloud_as_policy_execute_logs" "start_time_filter" {
  scaling_policy_id = "%[1]s"
  start_time        = local.start_time_before24h
}

locals {
  start_time_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.start_time_filter.execute_logs[*].execute_time :
    timecmp(local.start_time_before24h, v) == -1
  ]
}

output "is_start_time_filter_useful" {
  value = alltrue(local.start_time_filter_result) && length(local.start_time_filter_result) > 0
}

data "huaweicloud_as_policy_execute_logs" "end_time_filter" {
  scaling_policy_id = "%[1]s"
  end_time          = local.end_time_after24h
}

locals {
  end_time_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.end_time_filter.execute_logs[*].execute_time :
    timecmp(local.end_time_after24h, v) == 1
  ]
}

output "is_end_time_filter_useful" {
  value = alltrue(local.end_time_filter_result) && length(local.end_time_filter_result) > 0
}

data "huaweicloud_as_policy_execute_logs" "start_end_time_filter" {
  scaling_policy_id = "%[1]s"
  start_time        = local.start_time_before24h
  end_time          = local.end_time_after24h
}

locals {
  start_end_time_filter_result = [
    for v in data.huaweicloud_as_policy_execute_logs.start_end_time_filter.execute_logs[*].execute_time :
    timecmp(local.start_time_before24h, v) == -1 && timecmp(local.end_time_after24h, v) == 1
  ]
}

output "start_end_time_filter_result" {
  value = alltrue(local.start_end_time_filter_result) && length(local.start_end_time_filter_result) > 0
}
`, acceptance.HW_AS_SCALING_POLICY_ID)
}
