package huaweicloud

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccVpcSubnetV1DataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dNameByID := "data.huaweicloud_vpc_subnet.by_id"
	dNameByCIDR := "data.huaweicloud_vpc_subnet.by_cidr"
	dNameByName := "data.huaweicloud_vpc_subnet.by_name"
	dNameByVpcID := "data.huaweicloud_vpc_subnet.by_vpc_id"
	tmp := strconv.Itoa(acctest.RandIntRange(1, 254))
	cidr := fmt.Sprintf("172.16.%s.0/24", string(tmp))
	gateway := fmt.Sprintf("172.16.%s.1", string(tmp))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1DataSource_basic(rName, cidr, gateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1DataSourceID(dNameByID),
					resource.TestCheckResourceAttr(dNameByID, "name", rName),
					resource.TestCheckResourceAttr(dNameByID, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByID, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByID, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByID, "dhcp_enable", "true"),
					testAccCheckVpcSubnetV1DataSourceID(dNameByCIDR),
					resource.TestCheckResourceAttr(dNameByCIDR, "name", rName),
					resource.TestCheckResourceAttr(dNameByCIDR, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByCIDR, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByCIDR, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByCIDR, "dhcp_enable", "true"),
					testAccCheckVpcSubnetV1DataSourceID(dNameByName),
					resource.TestCheckResourceAttr(dNameByName, "name", rName),
					resource.TestCheckResourceAttr(dNameByName, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByName, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByName, "dhcp_enable", "true"),
					testAccCheckVpcSubnetV1DataSourceID(dNameByVpcID),
					resource.TestCheckResourceAttr(dNameByVpcID, "name", rName),
					resource.TestCheckResourceAttr(dNameByVpcID, "cidr", cidr),
					resource.TestCheckResourceAttr(dNameByVpcID, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dNameByVpcID, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dNameByVpcID, "dhcp_enable", "true"),
				),
			},
		},
	})
}

func TestAccVpcSubnetV1DataSource_ipv6(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	tmp := strconv.Itoa(acctest.RandIntRange(1, 254))
	cidr := fmt.Sprintf("172.16.%s.0/24", string(tmp))
	gateway := fmt.Sprintf("172.16.%s.1", string(tmp))
	dName := "data.huaweicloud_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1DataSource_ipv6(rName, cidr, gateway),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1DataSourceID(dName),
					resource.TestCheckResourceAttr(dName, "name", rName),
					resource.TestCheckResourceAttr(dName, "cidr", cidr),
					resource.TestCheckResourceAttr(dName, "gateway_ip", gateway),
					resource.TestCheckResourceAttr(dName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(dName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttr(dName, "ipv6_enable", "true"),
					resource.TestMatchResourceAttr(dName, "ipv6_cidr",
						regexp.MustCompile("([[:xdigit:]]*):([[:xdigit:]]*:){1,6}[[:xdigit:]]*/\\d{1,3}")),
					resource.TestMatchResourceAttr(dName, "ipv6_gateway",
						regexp.MustCompile("([[:xdigit:]]*):([[:xdigit:]]*:){1,6}([[:xdigit:]]){1,4}")),
					resource.TestCheckResourceAttrSet(dName, "ipv6_subnet_id"),
				),
			},
		},
	})
}

func testAccCheckVpcSubnetV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find %s in state", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Vpc Subnet data source ID not set")
		}

		return nil
	}
}

func testAccVpcSubnetV1DataSource_basic(rName, cidr, gateway string) string {
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

func testAccVpcSubnetV1DataSource_ipv6(rName, cidr, gateway string) string {
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
  ipv6_enable       = true
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
}

data "huaweicloud_vpc_subnet" "test" {
  id = huaweicloud_vpc_subnet.test.id
}
`, rName, cidr, rName, cidr, gateway)
}
