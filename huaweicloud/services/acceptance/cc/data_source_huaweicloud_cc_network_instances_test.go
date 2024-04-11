package cc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCcNetworkInstances_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cc_network_instances.filter_by_id"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)
	baseConfig := testDataSourceCcNetworkInstances_base(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCcNetworkInstances_basic(baseConfig, rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "network_instances.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "network_instances.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "network_instances.0.cloud_connection_id"),
					resource.TestCheckOutput("is_id_useful", "true"),
					resource.TestCheckOutput("is_cloud_connection_id_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCcNetworkInstances_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  count = 2
	  
  name = "%[1]s_${count.index}"
  cidr = cidrsubnet("10.12.0.0/16", 4, count.index)
}
	  
resource "huaweicloud_vpc_subnet" "test" {
  count = 2
	  
  name       = "%[1]s_${count.index}"
  vpc_id     = huaweicloud_vpc.test[count.index].id
  cidr       = cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 4, 1), 1)
}
  
resource "huaweicloud_cc_connection" "test" {
  name                  = "%[1]s"
  enterprise_project_id = "0"
  description           = "accDemo"
}

resource "huaweicloud_cc_network_instance" "test1" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = huaweicloud_vpc.test[0].id
  project_id          = "%[2]s"
  region_id           = huaweicloud_vpc.test[0].region
  description         = "desc 1"
  
  cidrs = [
    huaweicloud_vpc_subnet.test[0].cidr,
  ]
}

resource "huaweicloud_cc_network_instance" "test2" {
  type                = "vpc"
  cloud_connection_id = huaweicloud_cc_connection.test.id
  instance_id         = huaweicloud_vpc.test[1].id
  project_id          = "%[2]s"
  region_id           = huaweicloud_vpc.test[1].region
  description         = "desc 2"
	
  cidrs = [
    huaweicloud_vpc_subnet.test[1].cidr,
  ]

  depends_on = [
    huaweicloud_cc_network_instance.test1,
  ]
}
  `, name, acceptance.HW_PROJECT_ID)
}

func testDataSourceCcNetworkInstances_basic(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  id = huaweicloud_cc_network_instance.test1.id
}

data "huaweicloud_cc_network_instances" "filter_by_id" {
  network_instance_id = local.id
	
  depends_on = [
    huaweicloud_cc_network_instance.test1,
    huaweicloud_cc_network_instance.test2,
  ]
}
  
output "is_id_useful" {
  value = length(data.huaweicloud_cc_network_instances.filter_by_id.network_instances) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_network_instances.filter_by_id.network_instances[*].id : v == local.id]
  )
}
  
locals {
  cloud_connection_id = huaweicloud_cc_network_instance.test1.cloud_connection_id
}
	
data "huaweicloud_cc_network_instances" "filter_by_cloud_connection_id" {
  cloud_connection_id = local.cloud_connection_id
  
  depends_on = [
    huaweicloud_cc_network_instance.test1,
    huaweicloud_cc_network_instance.test2,
  ]
}
	  
output "is_cloud_connection_id_filter_useful" {
  value = length(data.huaweicloud_cc_network_instances.filter_by_cloud_connection_id.network_instances) >= 1 && alltrue([
    for v in data.huaweicloud_cc_network_instances.filter_by_cloud_connection_id.network_instances[*].cloud_connection_id : 
      v == local.cloud_connection_id
  ])
}
  
locals {
  description = huaweicloud_cc_network_instance.test2.description
}
	
data "huaweicloud_cc_network_instances" "filter_by_description" {
  description = local.description
  
  depends_on = [
    huaweicloud_cc_network_instance.test1,
    huaweicloud_cc_network_instance.test2,
  ]
}
	  
output "is_description_filter_useful" {
  value = length(data.huaweicloud_cc_network_instances.filter_by_description.network_instances) >= 1 && alltrue(
    [for v in data.huaweicloud_cc_network_instances.filter_by_description.network_instances[*].description : v == local.description]
  )
}
`, baseConfig, name)
}
