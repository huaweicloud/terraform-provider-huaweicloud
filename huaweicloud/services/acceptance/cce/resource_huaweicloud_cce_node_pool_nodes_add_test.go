package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccPoolNodesAdd_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolNodesAdd_basic(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func TestAccPoolNodesAdd_removeNodes(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
			acceptance.TestAccPreCheckUserId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccPoolNodesAdd_removeNodes(name),
				// there is nothing to check, if no error occurred, that means the test is successful
			},
		},
	})
}

func testAccPoolNodesAdd_base(name string) string {
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

  root_volume {
    size       = 40
    volumetype = "SSD"
  }
  data_volumes {
    size       = 100
    volumetype = "SSD"
  }
}

resource "huaweicloud_networking_secgroup" "secgroup" {
  name = "%[2]s"
}

data "huaweicloud_images_image" "myimage" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "basic1" {
  name               = "%[2]s-1"
  image_id           = data.huaweicloud_images_image.myimage.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.secgroup.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  key_pair           = huaweicloud_kps_keypair.test.name
  user_id            = "%[3]s"
  system_disk_type   = "SSD"
  system_disk_size   = "40"

  data_disks {
    type = "SSD"
    size = 100
  }

  charging_mode = "prePaid"
  period        = 1
  period_unit   = "month"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [ image_id, security_group_ids, tags ]
  }
}

resource "huaweicloud_compute_instance" "basic2" {
  name               = "%[2]s-2"
  image_id           = data.huaweicloud_images_image.myimage.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.secgroup.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  key_pair           = huaweicloud_kps_keypair.test.name
  user_id            = "%[3]s"
  system_disk_type   = "SSD"
  system_disk_size   = "40"

  data_disks {
    type = "SSD"
    size = 100
  }

  charging_mode = "prePaid"
  period        = 1
  period_unit   = "month"

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [ image_id, security_group_ids, tags ]
  }
}
`, common.TestVpc(name), name, acceptance.HW_USER_ID)
}

func testAccPoolNodesAdd_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool_nodes_add" "test" {
  cluster_id  = huaweicloud_cce_cluster.test.id
  nodepool_id = huaweicloud_cce_node_pool.test.id

  node_list {
    server_id = huaweicloud_compute_instance.basic1.id
  }
  node_list {
    server_id = huaweicloud_compute_instance.basic2.id
  }
}
`, testAccPoolNodesAdd_base(name), name)
}

func testAccPoolNodesAdd_removeNodes(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_pool_nodes_add" "test" {
  cluster_id  = huaweicloud_cce_cluster.test.id
  nodepool_id = huaweicloud_cce_node_pool.test.id

  remove_nodes_on_delete = true

  node_list {
    server_id = huaweicloud_compute_instance.basic1.id
  }
  node_list {
    server_id = huaweicloud_compute_instance.basic2.id
  }
}
`, testAccPoolNodesAdd_base(name), name)
}
