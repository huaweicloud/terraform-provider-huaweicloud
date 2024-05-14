package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcFlowLogs_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceName()
	dataSource1 := "data.huaweicloud_vpc_flow_logs.basic"
	dataSource2 := "data.huaweicloud_vpc_flow_logs.filter_by_name"
	dataSource3 := "data.huaweicloud_vpc_flow_logs.filter_by_id"
	dataSource4 := "data.huaweicloud_vpc_flow_logs.filter_by_resource_type"
	dataSource5 := "data.huaweicloud_vpc_flow_logs.filter_by_group_id"
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)
	dc5 := acceptance.InitDataSourceCheck(dataSource5)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcFlowLogs_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					dc5.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_resource_type_filter_useful", "true"),
					resource.TestCheckOutput("is_group_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcFlowLogs_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_flow_logs" "basic" {
  depends_on = [huaweicloud_vpc_flow_log.flow_log]
}

data "huaweicloud_vpc_flow_logs" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_vpc_flow_log.flow_log]
}

data "huaweicloud_vpc_flow_logs" "filter_by_id" {
  flow_log_id = huaweicloud_vpc_flow_log.flow_log.id

  depends_on = [huaweicloud_vpc_flow_log.flow_log]
}

data "huaweicloud_vpc_flow_logs" "filter_by_resource_type" {
  resource_type = "network"

  depends_on = [huaweicloud_vpc_flow_log.flow_log]
}

data "huaweicloud_vpc_flow_logs" "filter_by_group_id" {
  log_group_id = huaweicloud_lts_group.acc_group.id

  depends_on = [huaweicloud_vpc_flow_log.flow_log]
}

locals {
  name_filter_result   = [for v in data.huaweicloud_vpc_flow_logs.filter_by_name.flow_logs[*].name : v == "%[2]s"]
  id_filter_result     = [
    for v in data.huaweicloud_vpc_flow_logs.filter_by_id.flow_logs[*].id : v == huaweicloud_vpc_flow_log.flow_log.id
  ]
  resource_type_filter_result = [
    for v in data.huaweicloud_vpc_flow_logs.filter_by_resource_type.flow_logs[*].resource_type : v == "network"
  ]
  group_id_filter_result = [
    for v in data.huaweicloud_vpc_flow_logs.filter_by_group_id.flow_logs[*].log_group_id : v == huaweicloud_lts_group.acc_group.id
  ]
  
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_flow_logs.basic.flow_logs) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}

output "is_resource_type_filter_useful" {
  value = alltrue(local.resource_type_filter_result) && length(local.resource_type_filter_result) > 0
}

output "is_group_id_filter_useful" {
  value = alltrue(local.group_id_filter_result) && length(local.group_id_filter_result) > 0
}
`, testAccFlowLog_basic(name, name, "Created by Terraform"), name)
}
