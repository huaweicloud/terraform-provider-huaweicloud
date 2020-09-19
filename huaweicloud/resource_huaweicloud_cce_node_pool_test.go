package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodepools"
)

func TestAccCCENodePool_basic(t *testing.T) {
	var nodePool nodepools.NodePool

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	updateName := rName + "update"
	resourceName := "huaweicloud_cce_node_pool.test"
	//clusterName here is used to provide the cluster id to fetch cce node pool.
	clusterName := "huaweicloud_cce_cluster.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodePoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCCENodePool_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodePoolExists(resourceName, clusterName, &nodePool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				Config: testAccCCENodePool_update(rName, updateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "initial_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "scall_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "min_node_count", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_node_count", "9"),
					resource.TestCheckResourceAttr(resourceName, "scale_down_cooldown_time", "100"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
				),
			},
		},
	})
}

func testAccCheckCCENodePoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	var clusterId string
	var nodepollId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_cce_cluster" {
			clusterId = rs.Primary.ID
		}

		if rs.Type == "huaweicloud_cce_node_pool" {
			nodepollId = rs.Primary.ID
		}

		if clusterId == "" || nodepollId == "" {
			continue
		}

		_, err := nodepools.Get(cceClient, clusterId, nodepollId).Extract()
		if err == nil {
			return fmt.Errorf("Node still exists")
		}
	}

	return nil
}

func testAccCheckCCENodePoolExists(n string, cluster string, nodePool *nodepools.NodePool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		c, ok := s.RootModule().Resources[cluster]
		if !ok {
			return fmt.Errorf("Cluster not found: %s", c)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		if c.Primary.ID == "" {
			return fmt.Errorf("Cluster id is not set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
		}

		found, err := nodepools.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Node Pool not found")
		}

		*nodePool = *found

		return nil
	}
}

func testAccCCENodePool_Base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_compute_keypair" "test" {
  name = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_cce_cluster" "test" {
  name                   = "%s"
  cluster_type           = "VirtualMachine"
  flavor_id              = "cce.s1.small"
  vpc_id                 = huaweicloud_vpc_v1.test.id
  subnet_id              = huaweicloud_vpc_subnet_v1.test.id
  container_network_type = "overlay_l2"
}
`, testAccCCEClusterV3_Base(rName), rName, rName)
}

func testAccCCENodePool_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id         = huaweicloud_cce_cluster.test.id
  name               = "%s"
  os                 = "EulerOS 2.5"
  flavor_id          = "s6.large.2"
  initial_node_count = 1
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  key_pair           = huaweicloud_compute_keypair.test.name
  scall_enable      = false
  min_node_count    = 0
  max_node_count    = 0
  scale_down_cooldown_time = 0
  priority          = 0
  type 				= "vm"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, testAccCCENodePool_Base(rName), rName)
}

func testAccCCENodePool_update(rName, updateName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_cce_node_pool" "test" {
  cluster_id         = huaweicloud_cce_cluster.test.id
  name               = "%s"
  os                 = "EulerOS 2.5"
  flavor_id          = "s6.large.2"
  initial_node_count = 2
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  key_pair           = huaweicloud_compute_keypair.test.name
  scall_enable      = true
  min_node_count    = 2
  max_node_count    = 9
  scale_down_cooldown_time = 100
  priority          = 1
  type 				= "vm"

  root_volume {
    size       = 40
    volumetype = "SATA"
  }
  data_volumes {
    size       = 100
    volumetype = "SATA"
  }
}
`, testAccCCENodePool_Base(rName), updateName)
}
