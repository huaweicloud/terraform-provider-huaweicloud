package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.huaweicloud_networking_port_v2.port_1", "id",
						"huaweicloud_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.huaweicloud_networking_port_v2.port_2", "id",
						"huaweicloud_networking_port_v2.port_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.huaweicloud_networking_port_v2.port_3", "id",
						"huaweicloud_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_networking_port_v2.port_3", "all_fixed_ips.#", "1"),
				),
			},
		},
	})
}

const testAccNetworkingV2PortDataSource_basic = `
resource "huaweicloud_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  cidr       = "10.0.0.0/24"
  ip_version = 4
}

data "huaweicloud_networking_secgroup_v2" "default" {
  name = "default"
}

resource "huaweicloud_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = "${huaweicloud_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  security_group_ids = [
    "${data.huaweicloud_networking_secgroup_v2.default.id}",
  ]

  fixed_ip {
    subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
  }
}

resource "huaweicloud_networking_port_v2" "port_2" {
  name               = "port_2"
  network_id         = "${huaweicloud_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  security_group_ids = [
    "${data.huaweicloud_networking_secgroup_v2.default.id}",
  ]
}

data "huaweicloud_networking_port_v2" "port_1" {
  name           = "${huaweicloud_networking_port_v2.port_1.name}"
  admin_state_up = "${huaweicloud_networking_port_v2.port_1.admin_state_up}"

  security_group_ids = [
    "${data.huaweicloud_networking_secgroup_v2.default.id}",
  ]
}

data "huaweicloud_networking_port_v2" "port_2" {
  name           = "${huaweicloud_networking_port_v2.port_2.name}"
  admin_state_up = "${huaweicloud_networking_port_v2.port_2.admin_state_up}"
}

data "huaweicloud_networking_port_v2" "port_3" {
  fixed_ip = "${huaweicloud_networking_port_v2.port_1.all_fixed_ips.0}"
}
`
