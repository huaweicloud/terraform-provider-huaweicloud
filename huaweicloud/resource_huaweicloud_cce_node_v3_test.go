package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
)

func TestAccCCENodesV3_basic(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCCENode(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCCENodeV3_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists("huaweicloud_cce_node_v3.node_1", "huaweicloud_cce_cluster_v3.cluster_1", &node),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_node_v3.node_1", "name", "test-node"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_node_v3.node_1", "flavor", "s1.medium"),
				),
			},
			resource.TestStep{
				Config: testAccCCENodeV3_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_node_v3.node_1", "name", "test-node2"),
				),
			},
		},
	})
}

func TestAccCCENodesV3_timeout(t *testing.T) {
	var node nodes.Nodes

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCENodeV3Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCCENodeV3_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3Exists("huaweicloud_cce_node_v3.node_1", "huaweicloud_cce_cluster_v3.cluster_1", &node),
				),
			},
		},
	})
}

func testAccCheckCCENodeV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	var clusterId string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "huaweicloud_cce_cluster_v3" {
			clusterId = rs.Primary.ID
		}

		if rs.Type != "huaweicloud_cce_node_v3" {
			continue
		}

		_, err := nodes.Get(cceClient, clusterId, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Node still exists")
		}
	}

	return nil
}

func testAccCheckCCENodeV3Exists(n string, cluster string, node *nodes.Nodes) resource.TestCheckFunc {
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

		found, err := nodes.Get(cceClient, c.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Node not found")
		}

		*node = *found

		return nil
	}
}

var testAccCCENodeV3_basic = fmt.Sprintf(`
resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  name = "huaweicloud-cce"
  cluster_type="VirtualMachine"
  flavor="cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}

resource "huaweicloud_cce_node_v3" "node_1" {
cluster_id = "${huaweicloud_cce_cluster_v3.cluster_1.id}"
  name = "test-node"
  flavor="s1.medium"
  iptype="5_bgp"
  billing_mode=0
  az= "%s"
  sshkey="%s"
  root_volume = {
    size= 40,
    volumetype= "SATA"
  }
  bandwidth_charge_mode="traffic"
  sharetype= "PER"
  bandwidth_size= 100,
  data_volumes = [
    {
      size= 100,
      volumetype= "SATA"
    },
  ]
}`, OS_VPC_ID, OS_SUBNET_ID, OS_AVAILABILITY_ZONE, OS_SSH_KEY)

var testAccCCENodeV3_update = fmt.Sprintf(`
resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  name = "huaweicloud-cce"
  cluster_type="VirtualMachine"
  flavor="cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}

resource "huaweicloud_cce_node_v3" "node_1" {
cluster_id = "${huaweicloud_cce_cluster_v3.cluster_1.id}"
  name = "test-node2"
  flavor="s1.medium"
  iptype="5_bgp"
  billing_mode=0
  az= "%s"
  sshkey="%s"
  root_volume = {
    size= 40,
    volumetype= "SATA"
  }
  bandwidth_charge_mode="traffic"
  sharetype= "PER"
  bandwidth_size= 100,
  data_volumes = [
    {
      size= 100,
      volumetype= "SATA"
    },
  ]
}`, OS_VPC_ID, OS_SUBNET_ID, OS_AVAILABILITY_ZONE, OS_SSH_KEY)

var testAccCCENodeV3_timeout = fmt.Sprintf(`
resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  name = "huaweicloud-cce"
  cluster_type="VirtualMachine"
  flavor="cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}

resource "huaweicloud_cce_node_v3" "node_1" {
  cluster_id = "${huaweicloud_cce_cluster_v3.cluster_1.id}"
  name = "test-node2"
  flavor="s1.medium"
  iptype="5_bgp"
  billing_mode=0
  az= "%s"
  sshkey="%s"
  root_volume = {
    size= 40,
    volumetype= "SATA"
  }
  bandwidth_charge_mode="traffic"
  sharetype= "PER"
  bandwidth_size= 100,
  data_volumes = [
    {
      size= 100,
      volumetype= "SATA"
    },
  ]
timeouts {
create = "10m"
delete = "10m"
} 
}
`, OS_VPC_ID, OS_SUBNET_ID, OS_AVAILABILITY_ZONE, OS_SSH_KEY)
