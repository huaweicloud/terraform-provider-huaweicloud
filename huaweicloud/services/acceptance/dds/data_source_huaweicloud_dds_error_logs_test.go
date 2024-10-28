package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDDSErrorLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_error_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
			acceptance.TestAccPreCheckDDSTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDDSErrorLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.log_time"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.severity"),
					resource.TestCheckResourceAttrSet(dataSource, "error_logs.0.raw_message"),

					resource.TestCheckOutput("severity_filter_is_useful", "true"),
					resource.TestCheckOutput("node_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDDSErrorLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_error_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  severity = data.huaweicloud_dds_error_logs.test.error_logs.0.severity
  node_id  = data.huaweicloud_dds_error_logs.test.error_logs.0.node_id
}

data "huaweicloud_dds_error_logs" "filter" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
  severity    = local.severity
  node_id     = local.node_id
}

output "severity_filter_is_useful" {
  value = length(data.huaweicloud_dds_error_logs.filter.error_logs) > 0 && alltrue(
    [for v in data.huaweicloud_dds_error_logs.filter.error_logs[*].severity : v == local.severity]
  )
}

output "node_id_filter_is_useful" {
  value = length(data.huaweicloud_dds_error_logs.filter.error_logs) > 0 && alltrue(
    [for v in data.huaweicloud_dds_error_logs.filter.error_logs[*].node_id : v == local.node_id]
  )
}`, acceptance.HW_DDS_INSTANCE_ID, acceptance.HW_DDS_START_TIME, acceptance.HW_DDS_END_TIME)
}
