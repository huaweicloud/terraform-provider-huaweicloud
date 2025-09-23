package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesAgentDimensions_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_agent_dimensions.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesAgentDimensions_mountPoint(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSource, "dimensions.0.name", "mount_point"),
					resource.TestCheckResourceAttrSet(dataSource, "dimensions.0.value"),
					resource.TestCheckResourceAttrSet(dataSource, "dimensions.0.origin_value"),
				),
			},
		},
	})
}

func testDataSourceCesAgentDimensions_mountPoint() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_agent_dimensions" "test" {
  instance_id = huaweicloud_compute_instance.test.id
  dim_name    = "mount_point"
}
`, testDataSourceCesAgentDimensions_base())
}

func testDataSourceCesAgentDimensions_base() string {
	name := acceptance.RandomAccResourceName()
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
  cpu_core_count    = 8
  memory_size       = 16
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
	  
resource "huaweicloud_compute_instance" "test" {
  name               = "%[1]s"
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
}`, name)
}
