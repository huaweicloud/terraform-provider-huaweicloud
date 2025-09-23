package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDDSSlowLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dds_slow_logs.test"
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
				Config: testDataSourceDDSSlowLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.#"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.node_id"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.node_name"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.whole_message"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.operate_type"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.cost_time"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.database"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.collection"),
					resource.TestCheckResourceAttrSet(dataSource, "slow_logs.0.log_time"),

					resource.TestCheckOutput("operate_type_filter_is_useful", "true"),
					resource.TestCheckOutput("node_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDDSSlowLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_slow_logs" "test" {
  instance_id = "%[1]s"
  start_time  = "%[2]s"
  end_time    = "%[3]s"
}

locals {
  operate_type = data.huaweicloud_dds_slow_logs.test.slow_logs.0.operate_type
  node_id      = data.huaweicloud_dds_slow_logs.test.slow_logs.0.node_id
}

data "huaweicloud_dds_slow_logs" "filter" {
  instance_id  = "%[1]s"
  start_time   = "%[2]s"
  end_time     = "%[3]s"
  operate_type = local.operate_type
  node_id      = local.node_id
}

output "operate_type_filter_is_useful" {
  value = length(data.huaweicloud_dds_slow_logs.filter.slow_logs) > 0 && alltrue(
    [for v in data.huaweicloud_dds_slow_logs.filter.slow_logs[*].operate_type : v == local.operate_type]
  )
}

output "node_id_filter_is_useful" {
  value = length(data.huaweicloud_dds_slow_logs.filter.slow_logs) > 0 && alltrue(
    [for v in data.huaweicloud_dds_slow_logs.filter.slow_logs[*].node_id : v == local.node_id]
  )
}`, acceptance.HW_DDS_INSTANCE_ID, acceptance.HW_DDS_START_TIME, acceptance.HW_DDS_END_TIME)
}
