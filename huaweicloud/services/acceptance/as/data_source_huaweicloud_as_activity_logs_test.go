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
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckASScalingGroupID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccActivityLogsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "scaling_group_id", acceptance.HW_AS_SCALING_GROUP_ID),
					resource.TestCheckOutput("is_start_time_filter_useful", "true"),
					resource.TestCheckOutput("is_end_time_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
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

locals {
  start_time = data.huaweicloud_as_activity_logs.test.activity_logs[0].start_time
}
data "huaweicloud_as_activity_logs" "start_time_filter" {
  scaling_group_id = "%[1]s"
  start_time       = local.start_time
}
output "is_start_time_filter_useful" {
  value = length(data.huaweicloud_as_activity_logs.start_time_filter.activity_logs) > 0 && alltrue( 
    [for v in data.huaweicloud_as_activity_logs.start_time_filter.activity_logs[*].start_time : v == local.start_time]
  )  
}

locals {
  end_time = data.huaweicloud_as_activity_logs.test.activity_logs[0].end_time
}
data "huaweicloud_as_activity_logs" "end_time_filter" {
  scaling_group_id = "%[1]s"
  end_time         = local.end_time
}
output "is_end_time_filter_useful" {
  value = length(data.huaweicloud_as_activity_logs.end_time_filter.activity_logs) >= 0 && alltrue( 
    [for v in data.huaweicloud_as_activity_logs.end_time_filter.activity_logs[*].end_time : v != local.end_time]
  )  
}

locals {
  status = data.huaweicloud_as_activity_logs.test.activity_logs[0].status
}
data "huaweicloud_as_activity_logs" "status_filter" {
  scaling_group_id = "%[1]s"
  status           = local.status
}
output "is_status_filter_useful" {
  value = length(data.huaweicloud_as_activity_logs.status_filter.activity_logs) > 0 && alltrue(
    [for v in data.huaweicloud_as_activity_logs.status_filter.activity_logs[*].status : v == local.status]
  )  
}
`, acceptance.HW_AS_SCALING_GROUP_ID)
}
