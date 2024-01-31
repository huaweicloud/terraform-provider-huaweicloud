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
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckASScalingPolicyID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyExecuteLogsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "scaling_policy_id", acceptance.HW_AS_SCALING_POLICY_ID),
					resource.TestCheckOutput("is_log_id_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_id_filter_useful", "true"),
					resource.TestCheckOutput("is_execute_type_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
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

// Filter using log ID.
locals {
  log_id = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].id
}

data "huaweicloud_as_policy_execute_logs" "log_id_filter" {
  scaling_policy_id = "%[1]s"
  log_id            = local.log_id
}

output "is_log_id_filter_useful" {
  value = length(data.huaweicloud_as_policy_execute_logs.log_id_filter.execute_logs) > 0 && alltrue(
    [for v in data.huaweicloud_as_policy_execute_logs.log_id_filter.execute_logs[*].id : v == local.log_id]
  )
}

// Filter using scaling resource ID.
locals {
  scaling_resource_id = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].scaling_resource_id
}

data "huaweicloud_as_policy_execute_logs" "resource_id_filter" {
  scaling_policy_id   = "%[1]s"
  scaling_resource_id = local.scaling_resource_id
}

output "is_resource_id_filter_useful" {
  value = length(data.huaweicloud_as_policy_execute_logs.resource_id_filter.execute_logs) > 0 && alltrue(
    [for v in data.huaweicloud_as_policy_execute_logs.resource_id_filter.execute_logs[*].scaling_resource_id : v == local.scaling_resource_id]
  )
}

// Filter using execute type.
locals {
  execute_type = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].execute_type
}

data "huaweicloud_as_policy_execute_logs" "execute_type_filter" {
  scaling_policy_id = "%[1]s"
  execute_type      = local.execute_type
}

output "is_execute_type_filter_useful" {
  value = length(data.huaweicloud_as_policy_execute_logs.execute_type_filter.execute_logs) > 0 && alltrue(
    [for v in data.huaweicloud_as_policy_execute_logs.execute_type_filter.execute_logs[*].execute_type : v == local.execute_type]
  )
}

// Filter using status.
locals {
  status = data.huaweicloud_as_policy_execute_logs.test.execute_logs[0].status
}

data "huaweicloud_as_policy_execute_logs" "status_filter" {
  scaling_policy_id = "%[1]s"
  status            = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_as_policy_execute_logs.status_filter.execute_logs) > 0 && alltrue(
    [for v in data.huaweicloud_as_policy_execute_logs.status_filter.execute_logs[*].status : v == local.status]
  )
}
`, acceptance.HW_AS_SCALING_POLICY_ID)
}
