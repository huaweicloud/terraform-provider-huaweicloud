package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/routes"
)

func TestAccVpcRouteV2_basic(t *testing.T) {
	var route routes.Route

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_route.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRouteV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteV2_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRouteV2Exists(resourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "destination", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "type", "peering"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckRouteV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	routeClient, err := config.NetworkingV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud route client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_route" {
			continue
		}

		_, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Route still exists")
		}
	}

	return nil
}

func testAccCheckRouteV2Exists(n string, route *routes.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		routeClient, err := config.NetworkingV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud route client: %s", err)
		}

		found, err := routes.Get(routeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.RouteID != rs.Primary.ID {
			return fmt.Errorf("route not found")
		}

		*route = *found

		return nil
	}
}

func testAccRouteV2_basic(rName string) string {
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
