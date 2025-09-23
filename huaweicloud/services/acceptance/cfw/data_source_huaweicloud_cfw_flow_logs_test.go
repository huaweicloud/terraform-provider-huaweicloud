package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwFlowLogs_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_flow_logs.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
			acceptance.TestAccPreCheckCfwTimeRange(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCfwFlowLogs_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.app"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.direction"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_port"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.src_region_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.dst_region_name"),

					resource.TestCheckOutput("is_app_filter_useful", "true"),
					resource.TestCheckOutput("is_direction_filter_useful", "true"),
					resource.TestCheckOutput("is_src_port_filter_useful", "true"),
					resource.TestCheckOutput("is_dst_port_filter_useful", "true"),
					resource.TestCheckOutput("is_src_region_name_filter_useful", "true"),
					resource.TestCheckOutput("is_dst_region_name_filter_useful", "true"),
					resource.TestCheckOutput("is_dst_province_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCfwFlowLogs_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_cfw_flow_logs" "test" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
}

locals {
  records           = data.huaweicloud_cfw_flow_logs.test.records
  app               = local.records[0].app
  direction         = local.records[0].direction
  src_port          = local.records[0].src_port
  dst_port          = local.records[0].dst_port
  src_region_name   = local.records[0].src_region_name
  dst_region_name   = local.records[0].dst_region_name
  dst_province_name = local.records[0].dst_province_name
}

data "huaweicloud_cfw_flow_logs" "filter_by_app" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
  app            = local.app
}

data "huaweicloud_cfw_flow_logs" "filter_by_direction" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
  direction      = local.direction
}

data "huaweicloud_cfw_flow_logs" "filter_by_src_port" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
  src_port       = local.src_port
}

data "huaweicloud_cfw_flow_logs" "filter_by_dst_port" {
  fw_instance_id = "%[1]s"
  start_time     = "%[2]s"
  end_time       = "%[3]s"
  dst_port       = local.dst_port
}

data "huaweicloud_cfw_flow_logs" "filter_by_src_region_name" {
  fw_instance_id  = "%[1]s"
  start_time      = "%[2]s"
  end_time        = "%[3]s"
  src_region_name = local.src_region_name
}

data "huaweicloud_cfw_flow_logs" "filter_by_dst_region_name" {
  fw_instance_id  = "%[1]s"
  start_time      = "%[2]s"
  end_time        = "%[3]s"
  dst_region_name = local.dst_region_name
}

data "huaweicloud_cfw_flow_logs" "filter_by_dst_province_name" {
  fw_instance_id    = "%[1]s"
  start_time        = "%[2]s"
  end_time          = "%[3]s"
  dst_province_name = local.dst_province_name
}

locals {
  records_by_app               = data.huaweicloud_cfw_flow_logs.filter_by_app.records
  records_by_direction         = data.huaweicloud_cfw_flow_logs.filter_by_direction.records
  records_by_src_port          = data.huaweicloud_cfw_flow_logs.filter_by_src_port.records
  records_by_dst_port          = data.huaweicloud_cfw_flow_logs.filter_by_dst_port.records
  records_by_src_region_name   = data.huaweicloud_cfw_flow_logs.filter_by_src_region_name.records
  records_by_dst_region_name   = data.huaweicloud_cfw_flow_logs.filter_by_dst_region_name.records
  records_by_dst_province_name = data.huaweicloud_cfw_flow_logs.filter_by_dst_province_name.records
}

output "is_app_filter_useful" {
  value = length(local.records_by_app) > 0 && alltrue(
    [for v in local.records_by_app[*].app : v == local.app]
  )
}

output "is_direction_filter_useful" {
  value = length(local.records_by_direction) > 0 && alltrue(
    [for v in local.records_by_direction[*].direction : v == local.direction]
  )
}

output "is_src_port_filter_useful" {
  value = length(local.records_by_src_port) > 0 && alltrue(
    [for v in local.records_by_src_port[*].src_port : v == local.src_port]
  )
}

output "is_dst_port_filter_useful" {
  value = length(local.records_by_dst_port) > 0 && alltrue(
    [for v in local.records_by_dst_port[*].dst_port : v == local.dst_port]
  )
}

output "is_src_region_name_filter_useful" {
  value = length(local.records_by_src_region_name) > 0 && alltrue(
    [for v in local.records_by_src_region_name[*].src_region_name : v == local.src_region_name]
  )
}

output "is_dst_region_name_filter_useful" {
  value = length(local.records_by_dst_region_name) > 0 && alltrue(
    [for v in local.records_by_dst_region_name[*].dst_region_name : v == local.dst_region_name]
  )
}

output "is_dst_province_name_filter_useful" {
  value = length(local.records_by_dst_province_name) > 0 && alltrue(
    [for v in local.records_by_dst_province_name[*].dst_province_name : v == local.dst_province_name]
  )
}
`, acceptance.HW_CFW_INSTANCE_ID, acceptance.HW_CFW_START_TIME, acceptance.HW_CFW_END_TIME)
}
