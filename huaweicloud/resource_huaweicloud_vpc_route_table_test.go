package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/routables"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccVpcRouteTable_basic(t *testing.T) {
	var route routables.RouteTable

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_route_table.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcRouteTableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTable_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcRouteTableExists(resourceName, &route),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "0"),
				),
			},
			{
				Config: testAccVpcRouteTable_route(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform with routes"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "route.0.destination", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "route.0.type", "peering"),
				),
			},
			{
				Config: testAccVpcRouteTable_associate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform with subnets"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "route.0.destination", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "route.0.type", "peering"),
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

func testAccCheckVpcRouteTableDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	vpcClient, err := config.NetworkingV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating huaweicloud VPC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_route_table" {
			continue
		}

		_, err := routables.Get(vpcClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Route table still exists")
		}
	}

	return nil
}

func testAccCheckVpcRouteTableExists(n string, routeTable *routables.RouteTable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		vpcClient, err := config.NetworkingV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating huaweicloud VPC client: %s", err)
		}

		found, err := routables.Get(vpcClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("route table not found")
		}

		*routeTable = *found

		return nil
	}
}

func testAccVpcRouteTable_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test1" {
  name = "%s-1"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test1-1" {
  name       = "%s-1-1"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test1.id
}

resource "huaweicloud_vpc_subnet" "test1-2" {
  name       = "%s-1-2"
  cidr       = "192.168.10.0/24"
  gateway_ip = "192.168.10.1"
  vpc_id     = huaweicloud_vpc.test1.id
}

resource "huaweicloud_vpc" "test2" {
  name = "%s-2"
  cidr = "172.16.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test2-1" {
  name       = "%s-2-1"
  cidr       = "172.16.10.0/24"
  gateway_ip = "172.16.10.1"
  vpc_id     = huaweicloud_vpc.test2.id
}
`, rName, rName, rName, rName, rName)
}

func testAccVpcRouteTable_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_route_table" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test1.id
  description = "created by terraform"
}
`, testAccVpcRouteTable_base(rName), rName)
}

func testAccVpcRouteTable_route(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test1.id
  peer_vpc_id = huaweicloud_vpc.test2.id
}

resource "huaweicloud_vpc_route_table" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test1.id
  description = "created by terraform with routes"

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering rule"
  }
}
`, testAccVpcRouteTable_base(rName), rName, rName)
}

func testAccVpcRouteTable_associate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_peering_connection" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test1.id
  peer_vpc_id = huaweicloud_vpc.test2.id
}

resource "huaweicloud_vpc_route_table" "test" {
  name        = "%s"
  vpc_id      = huaweicloud_vpc.test1.id
  description = "created by terraform with subnets"

  subnets     = [
    huaweicloud_vpc_subnet.test1-1.id,
    huaweicloud_vpc_subnet.test1-2.id,
  ]

  route {
    destination = "172.16.0.0/16"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering rule"
  }
}
`, testAccVpcRouteTable_base(rName), rName, rName)
}
