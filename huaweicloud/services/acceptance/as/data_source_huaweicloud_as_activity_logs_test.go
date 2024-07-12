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
					resource.TestCheckResourceAttrSet(dataSourceName, "activity_logs.#"),
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
