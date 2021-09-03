package vpc

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccVpcSubnetV1_basic(t *testing.T) {
	var subnet subnets.Subnet

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_subnet.test"
	rNameUpdate := rName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_enable", "false"),
					resource.TestCheckResourceAttr(resourceName, "dhcp_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "dns_list.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcSubnetV1_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value_updated"),
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

func TestAccVpcSubnetV1_ipv6(t *testing.T) {
	var subnet subnets.Subnet

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_subnet.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetV1_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists(resourceName, &subnet),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
				),
			},
			{
				Config: testAccVpcSubnetV1_ipv6(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "ipv6_enable", "true"),
					resource.TestMatchResourceAttr(resourceName, "ipv6_cidr",
						regexp.MustCompile("([[:xdigit:]]*):([[:xdigit:]]*:){1,6}[[:xdigit:]]*/\\d{1,3}")),
					resource.TestMatchResourceAttr(resourceName, "ipv6_gateway",
						regexp.MustCompile("([[:xdigit:]]*):([[:xdigit:]]*:){1,6}([[:xdigit:]]){1,4}")),
					resource.TestCheckResourceAttrSet(resourceName, "ipv6_subnet_id"),
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

func testAccCheckVpcSubnetV1Destroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	subnetClient, err := config.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating huaweicloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_subnet" {
			continue
		}

		_, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Subnet still exists")
		}
	}

	return nil
}
func testAccCheckVpcSubnetV1Exists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		subnetClient, err := config.NetworkingV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating huaweicloud Vpc client: %s", err)
		}

		found, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

func testAccVpcSubnet_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}
`, rName)
}

func testAccVpcSubnetV1_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id

  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVpcSubnet_base(rName), rName)
}

func testAccVpcSubnetV1_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id

  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value_updated"
  }
}
`, testAccVpcSubnet_base(rName), rName)
}

func testAccVpcSubnetV1_ipv6(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpc_subnet" "test" {
  name              = "%s"
  cidr              = "192.168.0.0/16"
  gateway_ip        = "192.168.0.1"
  vpc_id            = huaweicloud_vpc.test.id
  ipv6_enable       = true
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccVpcSubnet_base(rName), rName)
}
