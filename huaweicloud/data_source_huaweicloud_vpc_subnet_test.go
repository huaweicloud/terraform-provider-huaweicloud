package huaweicloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcSubnetV1DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	tmp := strconv.Itoa(acctest.RandIntRange(1, 254))
	cidr := fmt.Sprintf("172.16.%s.0/24", string(tmp))
	gateway := fmt.Sprintf("172.16.%s.1", string(tmp))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVpcSubnetV1Config(rName, cidr, gateway),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceVpcSubnetV1Check("data.huaweicloud_vpc_subnet.by_id", rName, cidr, gateway),
					testAccDataSourceVpcSubnetV1Check("data.huaweicloud_vpc_subnet.by_cidr", rName, cidr, gateway),
					testAccDataSourceVpcSubnetV1Check("data.huaweicloud_vpc_subnet.by_name", rName, cidr, gateway),
					testAccDataSourceVpcSubnetV1Check("data.huaweicloud_vpc_subnet.by_vpc_id", rName, cidr, gateway),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_subnet.by_id", "status", "ACTIVE"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_vpc_subnet.by_id", "dhcp_enable", "true"),
				),
			},
		},
	})
}

func testAccDataSourceVpcSubnetV1Check(n, name, cidr, gateway_ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", n)
		}

		subnetRs, ok := s.RootModule().Resources["huaweicloud_vpc_subnet.test"]
		if !ok {
			return fmt.Errorf("can't find huaweicloud_vpc_subnet.test in state")
		}

		attr := rs.Primary.Attributes

		if attr["id"] != subnetRs.Primary.Attributes["id"] {
			return fmt.Errorf(
				"id is %s; want %s",
				attr["id"],
				subnetRs.Primary.Attributes["id"],
			)
		}

		if attr["cidr"] != cidr {
			return fmt.Errorf("bad subnet cidr %s, expected: %s", attr["cidr"], cidr)
		}
		if attr["name"] != name {
			return fmt.Errorf("bad subnet name %s", attr["name"])
		}
		if attr["gateway_ip"] != gateway_ip {
			return fmt.Errorf("bad subnet gateway_ip %s", attr["gateway_ip"])
		}

		return nil
	}
}

func testAccDataSourceVpcSubnetV1Config(rName, cidr, gateway string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "%s"
}

resource "huaweicloud_vpc_subnet" "test" {
  name              = "%s"
  cidr              = "%s"
  gateway_ip        = "%s"
  vpc_id            = huaweicloud_vpc.test.id

  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

data "huaweicloud_vpc_subnet" "by_id" {
  id = huaweicloud_vpc_subnet.test.id
}

data "huaweicloud_vpc_subnet" "by_cidr" {
  cidr = huaweicloud_vpc_subnet.test.cidr
}

data "huaweicloud_vpc_subnet" "by_name" {
  name = huaweicloud_vpc_subnet.test.name
}

data "huaweicloud_vpc_subnet" "by_vpc_id" {
  vpc_id = huaweicloud_vpc_subnet.test.vpc_id
}
`, rName, cidr, rName, cidr, gateway)
}
