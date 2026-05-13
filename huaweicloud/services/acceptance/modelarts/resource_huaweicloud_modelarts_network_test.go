package modelarts

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/modelarts"
)

func getNetworkResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NewServiceClient("modelarts", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	return modelarts.GetNetworkById(client, state.Primary.ID)
}

func TestAccNetwork_basic(t *testing.T) {
	var (
		obj interface{}

		rName = "huaweicloud_modelarts_network.test"
		rc    = acceptance.InitResourceCheck(rName, &obj, getNetworkResourceFunc)

		name = acceptance.RandomAccResourceNameWithDash()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				// Create a network and with two SFS Turbo connections.
				Config: testAccNetwork_basic_step1(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "10.168.0.0/16"),
					resource.TestCheckResourceAttr(rName, "workspace_id", "0"),
					resource.TestCheckResourceAttr(rName, "peer_connections.#", "0"),
					resource.TestCheckResourceAttr(rName, "sfs_turbos.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.0.name",
						"huaweicloud_sfs_turbo.test.0", "name"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.0.id",
						"huaweicloud_sfs_turbo.test.0", "id"),
					resource.TestCheckResourceAttrSet(rName, "sfs_turbos.0.uri"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.1.name",
						"huaweicloud_sfs_turbo.test.1", "name"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.1.id",
						"huaweicloud_sfs_turbo.test.1", "id"),
					resource.TestCheckResourceAttrSet(rName, "sfs_turbos.1.uri"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttrSet(rName, "subnet_id"),
				),
			},
			{
				// Associate two peering connections and change the first associated SFS Turbo connection to another.
				Config: testAccNetwork_basic_step2(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "10.168.0.0/16"),
					resource.TestCheckResourceAttr(rName, "workspace_id", "0"),
					resource.TestCheckResourceAttr(rName, "peer_connections.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.vpc_id",
						"huaweicloud_vpc.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.subnet_id",
						"huaweicloud_vpc_subnet.test.0", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.1.vpc_id",
						"huaweicloud_vpc.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.1.subnet_id",
						"huaweicloud_vpc_subnet.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "sfs_turbos.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.0.name",
						"huaweicloud_sfs_turbo.test.2", "name"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.0.id",
						"huaweicloud_sfs_turbo.test.2", "id"),
					resource.TestCheckResourceAttrSet(rName, "sfs_turbos.0.uri"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.1.name",
						"huaweicloud_sfs_turbo.test.1", "name"),
					resource.TestCheckResourceAttrPair(rName, "sfs_turbos.1.id",
						"huaweicloud_sfs_turbo.test.1", "id"),
					resource.TestCheckResourceAttrSet(rName, "sfs_turbos.1.uri"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
					resource.TestCheckResourceAttrSet(rName, "subnet_id"),
				),
			},
			{
				// Remove all SFS Turbo connections and change the first associated peering connection to another.
				Config: testAccNetwork_basic_step3(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "cidr", "10.168.0.0/16"),
					resource.TestCheckResourceAttr(rName, "workspace_id", "0"),
					resource.TestCheckResourceAttr(rName, "peer_connections.#", "2"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.vpc_id",
						"huaweicloud_vpc.test.2", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.0.subnet_id",
						"huaweicloud_vpc_subnet.test.2", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.1.vpc_id",
						"huaweicloud_vpc.test.1", "id"),
					resource.TestCheckResourceAttrPair(rName, "peer_connections.1.subnet_id",
						"huaweicloud_vpc_subnet.test.1", "id"),
					resource.TestCheckResourceAttr(rName, "sfs_turbos.#", "0"),
					resource.TestCheckResourceAttr(rName, "status", "Active"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccNetwork_basic_base(name string) string {
	return fmt.Sprintf(`
variable "enterprise_project_id" {
  default = "%[1]s"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  count = 3

  name                  = format("%[2]s-%%d", count.index)
  cidr                  = cidrsubnet("192.168.0.0/16", 2, count.index) # 192.168.0.0/18, 192.168.128.0/18, 192.168.192.0/18
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}

resource "huaweicloud_vpc_subnet" "test" {
  count = 3

  name       = format("%[2]s-%%d", count.index)
  vpc_id     = huaweicloud_vpc.test[count.index].id
  cidr       = cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 6, 0)              # 192.168.0.0/24, 192.168.128.0/24, 192.168.192.0/24
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test[count.index].cidr, 6, 0), 1) # 192.168.0.1, 192.168.128.1, 192.168.192.1
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[2]s"
  delete_default_rules = true
}

# Make sure open the full ingress access for 111, 2048, 2049, 2051, 2052 and 20048 ports and about TCP and UDP protocols.
resource "huaweicloud_networking_secgroup_rule" "tcp_ingress_access" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  ports             = "111,2048,2049,2051,2052,20048"
}

resource "huaweicloud_networking_secgroup_rule" "udp_ingress_access" {
  security_group_id = huaweicloud_networking_secgroup.test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "udp"
  ports             = "111,2048,2049,2051,2052,20048"
}

resource "huaweicloud_sfs_turbo" "test" {
  depends_on = [
    huaweicloud_networking_secgroup_rule.tcp_ingress_access,
    huaweicloud_networking_secgroup_rule.udp_ingress_access,
  ]

  count = 3

  name                  = format("%[2]s-%%d", count.index)
  size                  = 1228
  share_proto           = "NFS"
  share_type            = "HPC"
  hpc_bandwidth         = "40M"
  vpc_id                = huaweicloud_vpc.test[count.index].id
  subnet_id             = huaweicloud_vpc_subnet.test[count.index].id
  security_group_id     = huaweicloud_networking_secgroup.test.id
  availability_zone     = data.huaweicloud_availability_zones.test.names[0]
  enterprise_project_id = var.enterprise_project_id != "" ? var.enterprise_project_id : null
}
`, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST, name)
}

func testAccNetwork_basic_step1(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_network" "test" {
  name = "%[2]s"
  cidr = "10.168.0.0/16" # The recommended connecting CIDR about SFS Turbo.

  dynamic "sfs_turbos" {
    for_each = slice(huaweicloud_sfs_turbo.test, 0, 2)

    content {
      name = sfs_turbos.value.name
      id   = sfs_turbos.value.id
    }
  }
}
`, testAccNetwork_basic_base(name), name)
}

func testAccNetwork_basic_step2(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_network" "test" {
  name = "%[2]s"
  cidr = "10.168.0.0/16" # The recommended connecting CIDR about SFS Turbo.

  dynamic "peer_connections" {
    for_each = slice(huaweicloud_vpc_subnet.test, 0, 2)

    content {
      vpc_id    = peer_connections.value.vpc_id
      subnet_id = peer_connections.value.id
    }
  }

  dynamic "sfs_turbos" {
    for_each = [element(huaweicloud_sfs_turbo.test, 2), element(huaweicloud_sfs_turbo.test, 1)] # Reverse traversal.

    content {
      name = sfs_turbos.value.name
      id   = sfs_turbos.value.id
    }
  }
}
`, testAccNetwork_basic_base(name), name)
}

func testAccNetwork_basic_step3(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_modelarts_network" "test" {
  depends_on = [huaweicloud_vpc_subnet.test]

  name = "%[2]s"
  cidr = "10.168.0.0/16"

  dynamic "peer_connections" {
    for_each = [element(huaweicloud_vpc_subnet.test, 2), element(huaweicloud_vpc_subnet.test, 1)] # Reverse traversal.

    content {
      vpc_id    = peer_connections.value.vpc_id
      subnet_id = peer_connections.value.id
    }
  }
}
`, testAccNetwork_basic_base(name), name)
}
