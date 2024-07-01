package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcSubNetworkInterfaces_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_vpc_sub_network_interfaces.basic"
	dataSource2 := "data.huaweicloud_vpc_sub_network_interfaces.filter_by_id"
	dataSource3 := "data.huaweicloud_vpc_sub_network_interfaces.filter_by_subnet_id"
	dataSource4 := "data.huaweicloud_vpc_sub_network_interfaces.filter_by_parent_id"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)
	dc4 := acceptance.InitDataSourceCheck(dataSource4)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcSubNetworkInterfaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					dc4.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_subnet_id_filter_useful", "true"),
					resource.TestCheckOutput("is_parent_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcSubNetworkInterfaces_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_sub_network_interfaces" "basic" {
  depends_on = [huaweicloud_vpc_sub_network_interface.test]
}

data "huaweicloud_vpc_sub_network_interfaces" "filter_by_id" {
  interface_id = huaweicloud_vpc_sub_network_interface.test.id

  depends_on = [huaweicloud_vpc_sub_network_interface.test]
}

data "huaweicloud_vpc_sub_network_interfaces" "filter_by_subnet_id" {
  subnet_id = huaweicloud_vpc_subnet.test.id

  depends_on = [huaweicloud_vpc_sub_network_interface.test]
}

data "huaweicloud_vpc_sub_network_interfaces" "filter_by_parent_id" {
  parent_id = huaweicloud_compute_instance.test.network[0].port

  depends_on = [huaweicloud_vpc_sub_network_interface.test]
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_vpc_sub_network_interfaces.filter_by_id.sub_network_interfaces[*].id :
    v == huaweicloud_vpc_sub_network_interface.test.id
  ]
  subnet_id_filter_result = [
    for v in data.huaweicloud_vpc_sub_network_interfaces.filter_by_id.sub_network_interfaces[*].subnet_id :
    v == huaweicloud_vpc_sub_network_interface.test.subnet_id
  ]
  parent_id_filter_result = [
    for v in data.huaweicloud_vpc_sub_network_interfaces.filter_by_id.sub_network_interfaces[*].parent_id :
    v == huaweicloud_vpc_sub_network_interface.test.parent_id
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_sub_network_interfaces.basic.sub_network_interfaces) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}

output "is_subnet_id_filter_useful" {
  value = alltrue(local.subnet_id_filter_result) && length(local.subnet_id_filter_result) > 0
}

output "is_parent_id_filter_useful" {
  value = alltrue(local.parent_id_filter_result) && length(local.parent_id_filter_result) > 0
}
`, testSubNetworkInterface_basic(name))
}
