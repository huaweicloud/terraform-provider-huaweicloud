package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAsv2ActivityLogs_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_asv2_activity_logs.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byLogID   = "data.huaweicloud_asv2_activity_logs.log_id_filter"
		dcByLogID = acceptance.InitDataSourceCheck(byLogID)

		byStartTime   = "data.huaweicloud_asv2_activity_logs.start_time_filter"
		dcByStartTime = acceptance.InitDataSourceCheck(byStartTime)

		byEndTime   = "data.huaweicloud_asv2_activity_logs.end_time_filter"
		dcByEndTime = acceptance.InitDataSourceCheck(byEndTime)

		byStartEndTime   = "data.huaweicloud_asv2_activity_logs.start_end_time_filter"
		dcByStartEndTime = acceptance.InitDataSourceCheck(byStartEndTime)

		byType   = "data.huaweicloud_asv2_activity_logs.type_filter"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byStatus   = "data.huaweicloud_asv2_activity_logs.status_filter"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Please prepare the scaling group ID with scaling activity logs in advance.
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceAsv2ActivityLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.#"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.description"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.desire_value"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.instance_value"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.scaling_value"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "scaling_activity_log.0.type"),

					dcByLogID.CheckResourceExists(),
					resource.TestCheckOutput("is_log_id_filter_useful", "true"),

					dcByStartTime.CheckResourceExists(),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),

					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),

					dcByStartEndTime.CheckResourceExists(),
					resource.TestCheckOutput("start_time_filter_result1_useful", "true"),
					resource.TestCheckOutput("end_time_filter_result1_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceAsv2ActivityLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_asv2_activity_logs" "test" {
  scaling_group_id = "%[1]s"
}

# Filter by log_id
locals {
  log_id = data.huaweicloud_asv2_activity_logs.test.scaling_activity_log[0].id
}

data "huaweicloud_asv2_activity_logs" "log_id_filter" {
  scaling_group_id = "%[1]s"
  log_id           = local.log_id
}

locals {
  log_id_filter_result = [
    for v in data.huaweicloud_asv2_activity_logs.log_id_filter.scaling_activity_log[*].id : v == local.log_id
  ]
}

output "is_log_id_filter_useful" {
  value = alltrue(local.log_id_filter_result) && length(local.log_id_filter_result) > 0
}

# Filter by time
locals {
  start_time_utc = data.huaweicloud_asv2_activity_logs.test.scaling_activity_log[0].start_time
  end_time_utc   = data.huaweicloud_asv2_activity_logs.test.scaling_activity_log[0].end_time

  start_time_utc_before24h = timeadd(local.start_time_utc, "-24h")
  end_time_utc_after24h = timeadd(local.end_time_utc, "24h")
}

# Filter by start_time
data "huaweicloud_asv2_activity_logs" "start_time_filter" {
  scaling_group_id = "%[1]s"
  start_time       = local.start_time_utc_before24h
}

locals {
  start_time_filter_result = [
    for v in data.huaweicloud_asv2_activity_logs.start_time_filter.scaling_activity_log[*].start_time :
    timecmp(local.start_time_utc_before24h, v) == -1
  ]
}

output "is_start_time_filter_useful" {
  value = alltrue(local.start_time_filter_result) && length(local.start_time_filter_result) > 0
}

# Filter by end_time
data "huaweicloud_asv2_activity_logs" "end_time_filter" {
  scaling_group_id = "%[1]s"
  end_time         = local.end_time_utc_after24h
}

locals {
  end_time_filter_result = [
    for v in data.huaweicloud_asv2_activity_logs.end_time_filter.scaling_activity_log[*].end_time :
    timecmp(local.end_time_utc_after24h, v) == 1
  ]
}

output "is_end_time_filter_useful" {
  value = alltrue(local.end_time_filter_result) && length(local.end_time_filter_result) > 0
}

# Filter by start_time and end_time
data "huaweicloud_asv2_activity_logs" "start_end_time_filter" {
  scaling_group_id = "%[1]s"
  start_time       = local.start_time_utc_before24h
  end_time         = local.end_time_utc_after24h
}

locals {
  start_time_filter_result1 = [
    for v in data.huaweicloud_asv2_activity_logs.start_end_time_filter.scaling_activity_log[*].start_time :
    timecmp(local.start_time_utc_before24h, v) == -1
  ]

  end_time_filter_result1 = [
    for v in data.huaweicloud_asv2_activity_logs.start_end_time_filter.scaling_activity_log[*].end_time :
    timecmp(local.end_time_utc_after24h, v) == 1
  ]
}

output "start_time_filter_result1_useful" {
  value = alltrue(local.start_time_filter_result1) && length(local.start_time_filter_result1) > 0
}

output "end_time_filter_result1_useful" {
  value = alltrue(local.end_time_filter_result1) && length(local.end_time_filter_result1) > 0
}

# Filter by type
locals {
  type = data.huaweicloud_asv2_activity_logs.test.scaling_activity_log[0].type
}

data "huaweicloud_asv2_activity_logs" "type_filter" {
  scaling_group_id = "%[1]s"
  type             = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_asv2_activity_logs.type_filter.scaling_activity_log[*].type : v == local.type
  ]
}

output "is_type_filter_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}

# Filter by status
locals {
  status = data.huaweicloud_asv2_activity_logs.test.scaling_activity_log[0].status
}

data "huaweicloud_asv2_activity_logs" "status_filter" {
  scaling_group_id = "%[1]s"
  status           = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_asv2_activity_logs.status_filter.scaling_activity_log[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
