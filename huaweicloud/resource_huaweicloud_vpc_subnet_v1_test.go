package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v1/subnets"
)

func TestAccVpcSubnetV1_basic(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcSubnetV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists("huaweicloud_vpc_subnet_v1.subnet_1", &subnet),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_subnet_v1.subnet_1", "name", "huaweicloud_subnet"),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_subnet_v1.subnet_1", "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_subnet_v1.subnet_1", "gateway_ip", "192.168.0.1"),
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_subnet_v1.subnet_1", "availability_zone", OS_AVAILABILITY_ZONE),
				),
			},
			resource.TestStep{
				Config: testAccVpcSubnetV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_vpc_subnet_v1.subnet_1", "name", "huaweicloud_subnet_1"),
				),
			},
		},
	})
}

func TestAccVpcSubnetV1_timeout(t *testing.T) {
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVpcSubnetV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVpcSubnetV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetV1Exists("huaweicloud_vpc_subnet_v1.subnet_1", &subnet),
				),
			},
		},
	})
}

func testAccCheckVpcSubnetV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	subnetClient, err := config.networkingV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud vpc client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpc_subnet_v1" {
			continue
		}

		_, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Subnet still exists")
		}
	}

	return nil
}
func testAccCheckVpcSubnetV1Exists(n string, subnet *subnets.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		subnetClient, err := config.networkingV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud Vpc client: %s", err)
		}

		found, err := subnets.Get(subnetClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Subnet not found")
		}

		*subnet = *found

		return nil
	}
}

var testAccVpcSubnetV1_basic = fmt.Sprintf(`
resource "huaweicloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet_v1" "subnet_1" {
  name = "huaweicloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${huaweicloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"

}
`, OS_AVAILABILITY_ZONE)
var testAccVpcSubnetV1_update = fmt.Sprintf(`
resource "huaweicloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet_v1" "subnet_1" {
  name = "huaweicloud_subnet_1"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${huaweicloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"
 }
`, OS_AVAILABILITY_ZONE)

var testAccVpcSubnetV1_timeout = fmt.Sprintf(`
resource "huaweicloud_vpc_v1" "vpc_1" {
  name = "vpc_test"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet_v1" "subnet_1" {
  name = "huaweicloud_subnet"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = "${huaweicloud_vpc_v1.vpc_1.id}"
  availability_zone = "%s"

 timeouts {
    create = "5m"
    delete = "5m"
  }

}
`, OS_AVAILABILITY_ZONE)
