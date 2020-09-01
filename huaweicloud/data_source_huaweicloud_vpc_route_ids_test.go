package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcRouteIdsV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteIdV2DataSource_vpcroute(rName),
			},
			{
				Config: testAccRouteIdV2DataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccRouteIdV2DataSourceID("data.huaweicloud_vpc_route_ids.route_ids"),
					resource.TestCheckResourceAttr("data.huaweicloud_vpc_route_ids.route_ids", "ids.#", "1"),
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

func testAccRouteIdV2DataSource_vpcroute(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc" "test2" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test.id
  peer_vpc_id = huaweicloud_vpc.test2.id
}

resource "huaweicloud_vpc_route" "test" {
  type        = "peering"
  nexthop     = huaweicloud_vpc_peering_connection.test.id
  destination = "192.168.0.0/16"
  vpc_id      = huaweicloud_vpc.test.id
}
`, rName, rName+"2", rName)
}

func testAccRouteIdV2DataSource_basic(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc_route_ids" "route_ids" {
  vpc_id = huaweicloud_vpc_route.test.vpc_id
}
`, testAccRouteIdV2DataSource_vpcroute(rName))
}
