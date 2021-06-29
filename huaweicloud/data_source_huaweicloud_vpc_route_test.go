package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcRouteV2DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRouteV2Config(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteV2DataSourceID("data.huaweicloud_vpc_route.by_id"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_route.by_id", "type", "peering"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_route.by_id", "destination", "192.168.0.0/16"),
					testAccCheckRouteV2DataSourceID("data.huaweicloud_vpc_route.by_vpc_id"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_route.by_vpc_id", "type", "peering"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_route.by_vpc_id", "destination", "192.168.0.0/16"),
				),
			},
		},
	})
}

func testAccCheckRouteV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find vpc route connection data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("vpc route connection data source ID not set")
		}

		return nil
	}
}

func testAccDataSourceRouteV2Config(rName string) string {
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

data "huaweicloud_vpc_route" "by_id" {
  id = huaweicloud_vpc_route.test.id
}

data "huaweicloud_vpc_route" "by_vpc_id" {
  vpc_id = huaweicloud_vpc_route.test.vpc_id
}
`, rName, rName+"2", rName)
}
