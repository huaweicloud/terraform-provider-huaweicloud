package cce

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cce/v3/nodes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccCCENodeAttachV3_basic(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_attach.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeAttachV3_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.5"),
				),
			},
			{
				Config: testAccCCENodeAttachV3_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.5"),
				),
			},
			{
				Config: testAccCCENodeAttachV3_reset(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar_update"),
					resource.TestCheckResourceAttr(resourceName, "tags.key_update", "value_update"),
					resource.TestCheckResourceAttr(resourceName, "os", "CentOS 7.6"),
				),
			},
		},
	})
}

func TestAccCCENodeAttachV3_prePaid(t *testing.T) {
	var node nodes.Nodes

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_cce_node_attach.test"
	//clusterName here is used to provide the cluster id to fetch cce node.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckChargingMode(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodeAttachV3_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists(resourceName, clusterName, &node),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "os", "EulerOS 2.5"),
					resource.TestCheckResourceAttr(resourceName, "charging_mode", "prePaid"),
				),
			},
		},
	})
}

func testAccCCENodeAttachV3_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "EulerOS 2.5 64bit"
  most_recent = true
}

resource "huaweicloud_compute_keypair" "test" {
  name = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_compute_instance" "test" {
  name                        = "%s"
  image_id                    = data.huaweicloud_images_image.test.id
  flavor_id                   = "sn3.large.2"
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  key_pair                    = huaweicloud_compute_keypair.test.name
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
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName, rName, rName)
}

func testAccCCENodeAttachV3_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  key_pair   = huaweicloud_compute_keypair.test.name
  os         = "EulerOS 2.5"
  name       = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCENodeAttachV3_Base(rName), rName)
}

func testAccCCENodeAttachV3_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  key_pair   = huaweicloud_compute_keypair.test.name
  os         = "EulerOS 2.5"
  name       = "%s"

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, testAccCCENodeAttachV3_Base(rName), rNameUpdate)
}

func testAccCCENodeAttachV3_reset(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  key_pair   = huaweicloud_compute_keypair.test.name
  os         = "CentOS 7.6"
  name       = "%s"

  tags = {
    foo        = "bar_update"
    key_update = "value_update"
  }
}
`, testAccCCENodeAttachV3_Base(rName), rNameUpdate)
}

func testAccCCENodeAttachV3_prePaidBase(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "EulerOS 2.5 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "test" {
  name                        = "%s"
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
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc.test.id
  subnet_id              = huaweicloud_vpc_subnet.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName, rName)
}

func testAccCCENodeAttachV3_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_attach" "test" {
  cluster_id = huaweicloud_cce_cluster.test.id
  server_id  = huaweicloud_compute_instance.test.id
  password   = "Test@123"
  os         = "EulerOS 2.5"
  name       = "%s"

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCCENodeAttachV3_prePaidBase(rName), rName)
}
