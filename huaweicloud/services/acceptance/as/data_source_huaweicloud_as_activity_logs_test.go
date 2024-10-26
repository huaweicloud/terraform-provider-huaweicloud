package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccActivityLogsDataSource_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_as_activity_logs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byStatus   = "data.huaweicloud_as_activity_logs.status_filter"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byStartTime   = "data.huaweicloud_as_activity_logs.start_time_filter"
		dcByStartTime = acceptance.InitDataSourceCheck(byStartTime)

		byEndTime   = "data.huaweicloud_as_activity_logs.end_time_filter"
		dcByEndTime = acceptance.InitDataSourceCheck(byEndTime)

		byStartEndTime   = "data.huaweicloud_as_activity_logs.start_end_time_filter"
		dcByStartEndTime = acceptance.InitDataSourceCheck(byStartEndTime)
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
				Config: testAccActivityLogsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.added_instances"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.changes_instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.current_instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.desire_instance_number"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.start_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.0.status"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),

					dcByStartTime.CheckResourceExists(),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),

					dcByEndTime.CheckResourceExists(),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),

					dcByStartEndTime.CheckResourceExists(),
					resource.TestCheckOutput("start_time_filter_result1_useful", "true"),
					resource.TestCheckOutput("end_time_filter_result1_useful", "true"),
				),
			},
		},
	})
}

func testAccActivityLogsDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_as_activity_logs" "test" {
  scaling_group_id = "%[1]s"
}

# Filter by status
locals {
  status = data.huaweicloud_as_activity_logs.test.activity_logs[0].status
}

data "huaweicloud_as_activity_logs" "status_filter" {
  scaling_group_id = "%[1]s"
  status           = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_as_activity_logs.status_filter.activity_logs[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

# Filter by time
locals {
  start_time = data.huaweicloud_as_activity_logs.test.activity_logs[0].start_time
  end_time   = data.huaweicloud_as_activity_logs.test.activity_logs[0].end_time

  # The attribute date field format is "yyyy-MM-dd hh:mm:ss", and the parameter date field format is "yyyy-MM-ddThh:mm:ssZ".
  # The start and end times require format conversion. 
  start_time_utc = format("%%sZ", replace(local.start_time, " ", "T"))
  end_time_utc = format("%%sZ", replace(local.end_time, " ", "T"))

  start_time_utc_before24h = timeadd(local.start_time_utc, "-24h")
  end_time_utc_after24h = timeadd(local.end_time_utc, "24h")
}

# Filter by start_time
data "huaweicloud_as_activity_logs" "start_time_filter" {
  scaling_group_id = "%[1]s"
  start_time       = local.start_time_utc_before24h
}

locals {
  start_time_filter_result = [
    for v in data.huaweicloud_as_activity_logs.start_time_filter.activity_logs[*].start_time :
    timecmp(local.start_time_utc_before24h, format("%%sZ", replace(v, " ", "T"))) == -1
  ]
}

output "is_start_time_filter_useful" {
  value = alltrue(local.start_time_filter_result) && length(local.start_time_filter_result) > 0
}

# Filter by end_time
data "huaweicloud_as_activity_logs" "end_time_filter" {
  scaling_group_id = "%[1]s"
  end_time         = local.end_time_utc_after24h
}

locals {
  end_time_filter_result = [
    for v in data.huaweicloud_as_activity_logs.end_time_filter.activity_logs[*].end_time :
    timecmp(local.end_time_utc_after24h, format("%%sZ", replace(v, " ", "T"))) == 1
  ]
}

output "is_end_time_filter_useful" {
  value = alltrue(local.end_time_filter_result) && length(local.end_time_filter_result) > 0
}

# Filter by start_time and end_time
data "huaweicloud_as_activity_logs" "start_end_time_filter" {
  scaling_group_id = "%[1]s"
  start_time       = local.start_time_utc_before24h
  end_time         = local.end_time_utc_after24h
}

locals {
  start_time_filter_result1 = [
    for v in data.huaweicloud_as_activity_logs.start_end_time_filter.activity_logs[*].start_time :
    timecmp(local.start_time_utc_before24h, format("%%sZ", replace(v, " ", "T"))) == -1
  ]

  end_time_filter_result1 = [
    for v in data.huaweicloud_as_activity_logs.start_end_time_filter.activity_logs[*].end_time :
    timecmp(local.end_time_utc_after24h, format("%%sZ", replace(v, " ", "T"))) == 1
  ]
}

output "start_time_filter_result1_useful" {
  value = alltrue(local.start_time_filter_result1) && length(local.start_time_filter_result1) > 0
}

output "end_time_filter_result1_useful" {
  value = alltrue(local.end_time_filter_result1) && length(local.end_time_filter_result1) > 0
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
