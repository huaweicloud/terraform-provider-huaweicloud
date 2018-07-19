package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"
)

func TestAccVpcPeeringConnectionV2_basic(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcPeeringConnectionV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2Exists("huaweicloud_vpc_peering_connection_v2.peering_1", &peering),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_peering_connection_v2.peering_1", "name", "huaweicloud_peering"),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_peering_connection_v2.peering_1", "status", "ACTIVE"),
				),
			},
			resource.TestStep{
				Config: testAccVpcPeeringConnectionV2_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_peering_connection_v2.peering_1", "name", "huaweicloud_peering_1"),
				),
			},
		},
	})
}

func TestAccVpcPeeringConnectionV2_timeout(t *testing.T) {
	var peering peerings.Peering

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringConnectionV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcPeeringConnectionV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2Exists("huaweicloud_vpc_peering_connection_v2.peering_1", &peering),
				),
			},
		},
	})
}

func testAccCheckVpcPeeringConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	peeringClient, err := config.networkingHwV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud Peering client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_peering_connection_v2" {
			continue
		}

		_, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Vpc Peering Connection still exists")
		}
	}

	return nil
}

func testAccCheckVpcPeeringConnectionV2Exists(n string, peering *peerings.Peering) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		peeringClient, err := config.networkingHwV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud Peering client: %s", err)
		}

		found, err := peerings.Get(peeringClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Vpc peering Connection not found")
		}

		*peering = *found

		return nil
	}
}

const testAccVpcPeeringConnectionV2_basic = `
resource "huaweicloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_peering_connection_v2" "peering_1" {
  name = "huaweicloud_peering"
  vpc_id = "${huaweicloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${huaweicloud_vpc_v1.vpc_2.id}"
}
`
const testAccVpcPeeringConnectionV2_update = `
resource "huaweicloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_peering_connection_v2" "peering_1" {
  name = "huaweicloud_peering_1"
  vpc_id = "${huaweicloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${huaweicloud_vpc_v1.vpc_2.id}"
}
`
const testAccVpcPeeringConnectionV2_timeout = `
resource "huaweicloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_v1" "vpc_2" {
  name = "vpc_test1"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_peering_connection_v2" "peering_1" {
  name = "huaweicloud_peering"
  vpc_id = "${huaweicloud_vpc_v1.vpc_1.id}"
  peer_vpc_id = "${huaweicloud_vpc_v1.vpc_2.id}"

 timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
