package huaweicloud

import (
	"testing"

	"regexp"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVpcPeeringConnectionAccepterV2_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcPeeringConnectionAccepterDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testAccVpcPeeringConnectionAccepterV2_basic, //TODO: Research why normal scenario with peer tenant id is not working in acceptance tests
				ExpectError: regexp.MustCompile(`VPC peering action not permitted: Can not accept/reject peering request not in PENDING_ACCEPTANCE state.`),
			},
		},
	})
}

func testAccCheckVpcPeeringConnectionAccepterDestroy(s *terraform.State) error {
	// We don't destroy the underlying VPC Peering Connection.
	return nil
}

const testAccVpcPeeringConnectionAccepterV2_basic = `
resource "huaweicloud_vpc_v1" "vpc_1" {
  name = "huawei_vpc_1"
  cidr = "192.168.0.0/16"
}
resource "huaweicloud_vpc_v1" "vpc_2" {
  name = "huawei_vpc_2"
  cidr = "192.168.0.0/16"
}
resource "huaweicloud_vpc_peering_connection_v2" "peering_1" {
    name = "huaweicloud"
    vpc_id = "${huaweicloud_vpc_v1.vpc_1.id}"
    peer_vpc_id = "${huaweicloud_vpc_v1.vpc_2.id}"
  }
resource "huaweicloud_vpc_peering_connection_accepter_v2" "peer" {
  vpc_peering_connection_id = "${huaweicloud_vpc_peering_connection_v2.peering_1.id}"
  accept = true

}
`
