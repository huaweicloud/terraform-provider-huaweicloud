package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesAgentStatuses_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_agent_statuses.test"
	name := acceptance.RandomAccResourceName()
	baseConfig := testDataSourceCesAgentStatuses_base(name)
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesAgentStatuses_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "agent_status.0.instance_id"),
					resource.TestCheckResourceAttrSet(dataSource, "agent_status.0.uniagent_status"),
					resource.TestCheckResourceAttrSet(dataSource, "agent_status.0.extensions.#"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesAgentStatuses_basic(baseConfig string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_agent_statuses" "test" {
  instance_ids = [
    huaweicloud_compute_instance.test1.id,
    huaweicloud_compute_instance.test2.id,
  ]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_agent_statuses.test.agent_status) >= 2 
}

data "huaweicloud_ces_agent_statuses" "filter_by_status" {
  instance_ids = [
    huaweicloud_compute_instance.test1.id,
    huaweicloud_compute_instance.test2.id,
  ]

  extension_status = "running"
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_ces_agent_statuses.filter_by_status.agent_status) >= 1 && alltrue([
    for agent in data.huaweicloud_ces_agent_statuses.filter_by_status.agent_status[*] : 
      contains([
        for extension in agent.extensions : 
          extension.status if extension.name == "telescope"
      ], "running")
  ])
}
`, baseConfig)
}

func testDataSourceCesAgentStatuses_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.0.0/16"
}
	  
resource "huaweicloud_vpc_subnet" "test" { 
  name        = "%[1]s"
  cidr        = "192.168.0.0/24"
  gateway_ip  = "192.168.0.1"
  vpc_id      = huaweicloud_vpc.test.id
  ipv6_enable = true
}
	  
data "huaweicloud_availability_zones" "test" {}
	  
data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 8
}
	  
data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 22.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}
	  
resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "egress"
  ethertype         = "IPv4"
  action            = "allow"
  priority          = 1
  remote_ip_prefix  = "0.0.0.0/0"
}
	  
resource "huaweicloud_compute_instance" "test1" {
  name               = "%[1]s1"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  key_pair           = huaweicloud_kps_keypair.test.name
	  
  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_compute_instance" "test2" {
  name               = "%[1]s2"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  key_pair           = huaweicloud_kps_keypair.test.name
  agent_list         = "ces"
	  
  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
	  
  provisioner "local-exec" {
    command = "sleep 300"
  }
}
`, rName)
}
