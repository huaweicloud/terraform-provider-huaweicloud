package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccVpcPeeringConnectionV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcPeeringConnectionV2Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.huaweicloud_vpc_peering_connection.by_id"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_id", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_id", "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.huaweicloud_vpc_peering_connection.by_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_vpc_id", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_vpc_id", "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.huaweicloud_vpc_peering_connection.by_peer_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_peer_vpc_id", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_peer_vpc_id", "status", "ACTIVE"),
					testAccCheckVpcPeeringConnectionV2DataSourceID("data.huaweicloud_vpc_peering_connection.by_vpc_ids"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_vpc_ids", "name", rName),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_peering_connection.by_vpc_ids", "status", "ACTIVE"),
				),
			},
		},
	})
}

func testAccCheckVpcPeeringConnectionV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find vpc peering connection data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("vpc peering connection data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceVpcPeeringConnectionV2Config(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "vpc_1" {
  name = "%s_1"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "vpc_2" {
  name = "%s_2"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_peering_connection" "peering_1" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.vpc_1.id
  peer_vpc_id = huaweicloud_vpc.vpc_2.id
}

data "huaweicloud_vpc_peering_connection" "by_id" {
  id = huaweicloud_vpc_peering_connection.peering_1.id
}

data "huaweicloud_vpc_peering_connection" "by_vpc_id" {
  vpc_id = huaweicloud_vpc_peering_connection.peering_1.vpc_id
}

data "huaweicloud_vpc_peering_connection" "by_peer_vpc_id" {
  peer_vpc_id = huaweicloud_vpc_peering_connection.peering_1.peer_vpc_id
}

data "huaweicloud_vpc_peering_connection" "by_vpc_ids" {
  vpc_id      = huaweicloud_vpc_peering_connection.peering_1.vpc_id
  peer_vpc_id = huaweicloud_vpc_peering_connection.peering_1.peer_vpc_id
}
`, rName, rName, rName)
}
