package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceVpcNetworkInterfaces_basic(t *testing.T) {
	dataSource1 := "data.huaweicloud_vpc_network_interfaces.basic"
	dataSource2 := "data.huaweicloud_vpc_network_interfaces.filter_by_id"
	dataSource3 := "data.huaweicloud_vpc_network_interfaces.filter_by_name"
	rName := acceptance.RandomAccResourceName()
	dc1 := acceptance.InitDataSourceCheck(dataSource1)
	dc2 := acceptance.InitDataSourceCheck(dataSource2)
	dc3 := acceptance.InitDataSourceCheck(dataSource3)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceVpcNetworkInterfaces_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc1.CheckResourceExists(),
					dc2.CheckResourceExists(),
					dc3.CheckResourceExists(),
					resource.TestCheckOutput("is_results_not_empty", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceVpcNetworkInterfaces_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_vpc_network_interfaces" "basic" {
  depends_on = [huaweicloud_vpc_network_interface.test]
}

data "huaweicloud_vpc_network_interfaces" "filter_by_id" {
  interface_id = [huaweicloud_vpc_network_interface.test.id]

  depends_on = [huaweicloud_vpc_network_interface.test]
}

data "huaweicloud_vpc_network_interfaces" "filter_by_name" {
  name = "%[2]s"

  depends_on = [huaweicloud_vpc_network_interface.test]
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_vpc_network_interfaces.filter_by_id.ports[*].id :
    v == huaweicloud_vpc_network_interface.test.id
  ]
  name_filter_result = [
    for v in data.huaweicloud_vpc_network_interfaces.filter_by_id.ports[*].name : v == "%[2]s"
  ]
}

output "is_results_not_empty" {
  value = length(data.huaweicloud_vpc_network_interfaces.basic.ports) > 0
}

output "is_id_filter_useful" {
  value = alltrue(local.id_filter_result) && length(local.id_filter_result) > 0
}

output "is_name_filter_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}
`, testAccNetworkInterface_basic(name), name)
}
