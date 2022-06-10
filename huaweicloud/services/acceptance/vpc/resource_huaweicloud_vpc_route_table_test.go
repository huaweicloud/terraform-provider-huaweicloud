package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/routetables"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getRouteTableResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network v1 client: %s", err)
	}
	return routetables.Get(c, state.Primary.ID).Extract()
}

func TestAccVpcRouteTable_basic(t *testing.T) {
	var route routetables.RouteTable

	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_route_table.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&route,
		getRouteTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTable_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
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

func TestAccVpcRouteTable_multiRoutes(t *testing.T) {
	var route routetables.RouteTable
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_vpc_route_table.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&route,
		getRouteTableResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVpcRouteTable_multiRoutes(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform with multi routes"),
					resource.TestCheckResourceAttr(resourceName, "route.#", "6"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "0"),
				),
			},
		},
	})
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

func testAccVpcRouteTable_multiRoutes(rName string) string {
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
  description = "created by terraform with multi routes"

  route {
    destination = "172.16.1.0/24"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering one rule"
  }
  route {
    destination = "172.16.2.0/24"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering two rule"
  }
  route {
    destination = "172.16.3.0/24"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering three rule"
  }
  route {
    destination = "172.16.4.0/24"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering four rule"
  }
  route {
    destination = "172.16.5.0/24"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering five rule"
  }
  route {
    destination = "172.16.6.0/24"
    type        = "peering"
    nexthop     = huaweicloud_vpc_peering_connection.test.id
    description = "peering six rule"
  }
}
`, testAccVpcRouteTable_base(rName), rName, rName)
}
