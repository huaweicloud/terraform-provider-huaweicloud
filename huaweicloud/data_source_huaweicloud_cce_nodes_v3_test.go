package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCCENodesV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCCENode(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCCENodeV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3DataSourceID("data.huaweicloud_cce_node_v3.nodes"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_node_v3.nodes", "name", "test-node"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_node_v3.nodes", "flavor", "s1.medium"),
				),
			},
		},
	})
}

func testAccCheckCCENodeV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nodes data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Node data source ID not set ")
		}

		return nil
	}
}

var testAccCCENodeV3DataSource_basic = fmt.Sprintf(`
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
}
data "huaweicloud_cce_node_v3" "nodes" {
		cluster_id = "${huaweicloud_cce_cluster_v3.cluster_1.id}"
		name = "${huaweicloud_cce_node_v3.node_1.name}"
}
`, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE, OS_SSH_KEY)
