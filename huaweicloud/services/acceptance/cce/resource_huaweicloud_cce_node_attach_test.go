package cce

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func getAttachedNodeFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.CceV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating CCE v3 client: %s", err)
	}
	return nodes.Get(client, state.Primary.Attributes["cluster_id"], state.Primary.ID).Extract()
}

func TestAccNodeAttach_basic(t *testing.T) {
	var (
		node nodes.Nodes

		name         = acceptance.RandomAccResourceNameWithDash()
		updateName   = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_attach.test"

		baseConfig = testAccNodeAttach_base(name)

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getAttachedNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodeAttach_basic_step1(baseConfig, name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.9"),
				),
			},
			{
				Config: testAccNodeAttach_basic_step2(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
				),
			},
			{
				Config: testAccNodeAttach_basic_step3(baseConfig, updateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "os", "CentOS 7.6"),
				),
			},
		},
	})
}

func testAccNodeAttach_base(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 24.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[2]s"
}

resource "huaweicloud_compute_instance" "test" {
  name                        = "%[2]s"
  image_id                    = data.huaweicloud_images_image.test.id
  flavor_id                   = "sn3.large.2"
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  key_pair                    = huaweicloud_kps_keypair.test.name
  delete_disks_on_termination = true
  
  system_disk_type = "SAS"
  system_disk_size = 40
  
  data_disks {
	type = "SAS"
	size = "100"
  }
  
  network {
	uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [
      image_id, tags, name
    ]
  }
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, common.TestVpc(name), name)
}

func testAccNodeAttach_basic_step1(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  key_pair   = huaweicloud_kps_keypair.test.name
  os         = "EulerOS 2.9"
  name       = "%[2]s"

  max_pods         = 20
  docker_base_size = 10
  lvm_config       = "dockerThinpool=vgpaas/90%%VG;kubernetesLV=vgpaas/10%%VG"

  labels = {
    test_key = "test_value"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, baseConfig, name)
}

func testAccNodeAttach_basic_step2(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  key_pair   = huaweicloud_kps_keypair.test.name
  os         = "EulerOS 2.9"
  name       = "%[2]s"

  max_pods         = 20
  docker_base_size = 10
  lvm_config       = "dockerThinpool=vgpaas/90%%VG;kubernetesLV=vgpaas/10%%VG"

  labels = {
    test_key = "test_value"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, baseConfig, name)
}

func testAccNodeAttach_basic_step3(baseConfig, name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  key_pair   = huaweicloud_kps_keypair.test.name
  os         = "CentOS 7.6"
  name       = "%[2]s"

  max_pods         = 20
  docker_base_size = 10
  lvm_config       = "dockerThinpool=vgpaas/90%%VG;kubernetesLV=vgpaas/10%%VG"

  labels = {
    test_key = "test_value"
  }

  taints {
    key    = "test_key"
    value  = "test_value"
    effect = "NoSchedule"
  }

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, baseConfig, name)
}

func TestAccNodeAttach_prePaid(t *testing.T) {
	var (
		node nodes.Nodes

		name         = acceptance.RandomAccResourceNameWithDash()
		resourceName = "huaweicloud_cce_node_attach.test"

		rc = acceptance.InitResourceCheck(
			resourceName,
			&node,
			getAttachedNodeFunc,
		)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNodeAttach_prePaid(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.9"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
				),
			},
		},
	})
}

func testAccNodeAttach_prePaidBase(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 24.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "test" {
  name                        = "%[2]s"
  image_id                    = data.huaweicloud_images_image.test.id
  flavor_id                   = "sn3.large.2"
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  admin_pass                  = "Test@123"
  delete_disks_on_termination = true

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
  
  system_disk_type = "SAS"
  system_disk_size = 40
  
  data_disks {
	type = "SAS"
	size = "100"
  }
  
  network {
	uuid = huaweicloud_vpc_subnet.test.id
  }

  lifecycle {
    ignore_changes = [
      image_id, tags, name
    ]
  }
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%[2]s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, common.TestVpc(name), name)
}

func testAccNodeAttach_prePaid(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  password   = "Test@123"
  os         = "EulerOS 2.9"
  name       = "%[2]s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccNodeAttach_prePaidBase(name), name)
}
