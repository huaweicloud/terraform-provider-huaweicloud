package dbss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOperationLogs_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_dbss_operation_logs.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byUserName   = "data.huaweicloud_dbss_operation_logs.filter_by_user_name"
		dcByUserName = acceptance.InitDataSourceCheck(byUserName)

		byOperateName   = "data.huaweicloud_dbss_operation_logs.filter_by_operate_name"
		dcByOperateName = acceptance.InitDataSourceCheck(byOperateName)

		byResult   = "data.huaweicloud_dbss_operation_logs.filter_by_result"
		dcByResult = acceptance.InitDataSourceCheck(byResult)

		byTimeRange   = "data.huaweicloud_dbss_operation_logs.filter_by_time_range"
		dcByTimeRange = acceptance.InitDataSourceCheck(byTimeRange)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running this test case, please prepare an audit instance with user operation logs.
			acceptance.TestAccPrecheckDbssInstanceId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOperationLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.description"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.result"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.action"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.function"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.user"),
					resource.TestCheckResourceAttrSet(dataSourceName, "logs.0.time"),

					dcByUserName.CheckResourceExists(),
					resource.TestCheckOutput("user_name_filter_useful", "true"),

					dcByOperateName.CheckResourceExists(),
					resource.TestCheckOutput("operate_name_filter_useful", "true"),

					dcByResult.CheckResourceExists(),
					resource.TestCheckOutput("result_filter_useful", "true"),

					dcByTimeRange.CheckResourceExists(),
					resource.TestCheckOutput("time_range_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOperationLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dbss_operation_logs" "test" {
  instance_id = "%[1]s"
}

locals {
  user_name = data.huaweicloud_dbss_operation_logs.test.logs[0].user
}

data "huaweicloud_dbss_operation_logs" "filter_by_user_name" {
  instance_id = "%[1]s"
  user_name   = local.user_name
}

output "user_name_filter_useful" {
  value = length(data.huaweicloud_dbss_operation_logs.filter_by_user_name.logs) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_operation_logs.filter_by_user_name.logs[*].user : v == local.user_name]
  )
}

locals {
  operate_name = data.huaweicloud_dbss_operation_logs.test.logs[0].name
}

data "huaweicloud_dbss_operation_logs" "filter_by_operate_name" {
  instance_id  = "%[1]s"
  operate_name = local.operate_name
}

output "operate_name_filter_useful" {
  value = length(data.huaweicloud_dbss_operation_logs.filter_by_operate_name.logs) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_operation_logs.filter_by_operate_name.logs[*].name : v == local.operate_name]
  )
}

locals {
  result = data.huaweicloud_dbss_operation_logs.test.logs[0].result
}

data "huaweicloud_dbss_operation_logs" "filter_by_result" {
  instance_id = "%[1]s"
  result      = local.result
}

output "result_filter_useful" {
  value = length(data.huaweicloud_dbss_operation_logs.filter_by_result.logs) > 0 && alltrue(
    [for v in data.huaweicloud_dbss_operation_logs.filter_by_result.logs[*].result : v == local.result]
  )
}

data "huaweicloud_dbss_operation_logs" "filter_by_time_range" {
  instance_id = "%[1]s"
  time_range  = "DAY"
}

output "time_range_filter_useful" {
  value = length(data.huaweicloud_dbss_operation_logs.filter_by_time_range.logs) > 0
}
`, acceptance.HW_DBSS_INSATNCE_ID)
}
