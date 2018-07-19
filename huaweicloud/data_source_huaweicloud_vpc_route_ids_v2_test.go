package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccVpcRouteIdsV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRouteIdV2DataSource_vpcroute,
			},
			resource.TestStep{
				Config: testAccRouteIdV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccRouteIdV2DataSourceID("data.huaweicloud_vpc_route_ids_v2.route_ids"),
					resource.TestCheckResourceAttr("data.huaweicloud_vpc_route_ids_v2.route_ids", "ids.#", "1"),
				),
			},
		},
	})
}
func testAccRouteIdV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find vpc route data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Vpc Route data source ID not set")
		}

		return nil
	}
}

const testAccRouteIdV2DataSource_vpcroute = `
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

resource "huaweicloud_vpc_route_v2" "route_1" {
   type = "peering"
  nexthop = "${huaweicloud_vpc_peering_connection_v2.peering_1.id}"
  destination = "192.168.0.0/16"
  vpc_id ="${huaweicloud_vpc_v1.vpc_1.id}"
}
`

var testAccRouteIdV2DataSource_basic = fmt.Sprintf(`
%s
data "huaweicloud_vpc_route_ids_v2" "route_ids" {
  vpc_id = "${huaweicloud_vpc_route_v2.route_1.vpc_id}"
}
`, testAccRouteIdV2DataSource_vpcroute)
