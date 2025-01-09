package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccNodePoolScale_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccNodePoolScale_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAccNodePoolScale_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "zhangjishu-test-secgroup"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "rule1" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = huaweicloud_vpc.test.cidr
}

resource "huaweicloud_networking_secgroup_rule" "rule2" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "30000-32767"
  protocol          = "udp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "rule3" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = "30000-32767"
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "rule4" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "172.16.0.0/24"
}

resource "huaweicloud_networking_secgroup_rule" "rule5" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  action            = "allow"
  direction         = "ingress"
  ethertype         = "IPv4"
  remote_group_id   = huaweicloud_networking_secgroup.test.id
}

resource "huaweicloud_networking_secgroup_rule" "rule6" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  action            = "allow"
  direction         = "egress"
  ethertype         = "IPv4"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id               = huaweicloud_cce_cluster.test.id
  name                     = "%[2]s"
  os                       = "EulerOS 2.9"
  flavor_id                = data.huaweicloud_compute_flavors.test.ids[0]
  initial_node_count       = 0
  availability_zone        = data.huaweicloud_availability_zones.test.names[0]
  password                 = "Test@123"
  scall_enable             = false
  min_node_count           = 0
  max_node_count           = 0
  scale_down_cooldown_time = 0
  priority                 = 0
  type                     = "vm"
  security_groups          = [huaweicloud_networking_secgroup.test.id]

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}
`, common.TestVpc(name), name)
}

func testAccNodePoolScale_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool_scale" "test" {
  cluster_id         = huaweicloud_cce_cluster.test.id
  nodepool_id        = huaweicloud_cce_node_pool.test.id
  scale_groups       = ["default"]
  desired_node_count = 2
}
`, testAccNodePoolScale_base(name), name)
}
