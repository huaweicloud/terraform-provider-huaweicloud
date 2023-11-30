package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcSubnet_basic(t *testing.T) {
	var iecSubnet ieccommon.Subnet

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_iec_vpc_subnet.subnet_test"
	rNameUpdate := rName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnet_customer(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists(resourceName, &iecSubnet),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-subnet", rName)),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.128.0/18"),
					resource.TestCheckResourceAttr(resourceName, "gateway_ip", "192.168.128.1"),
					resource.TestCheckResourceAttr(resourceName, "dns_list.#", "2"),
				),
			},
			{
				Config: testAccVpcSubnet_customer_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcSubnetExists(resourceName, &iecSubnet),
					resource.TestCheckResourceAttr(resourceName, "name", fmt.Sprintf("%s-subnet", rNameUpdate)),
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

func testAccCheckVpcSubnetDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	iecV1Client, err := conf.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_vpc_subnet" {
			continue
		}

		_, err := subnets.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC VPC still exists")
		}
	}

	return nil
}

func testAccCheckVpcSubnetExists(n string, subnetResource *ieccommon.Subnet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		iecV1Client, err := config.IECV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating Huaweicloud IEC client: %s", err)
		}

		found, err := subnets.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC VPC not found")
		}

		*subnetResource = *found

		return nil
	}
}

func testAccVpcSubnet_customer(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "%s-vpc"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name       = "%s-subnet"
  cidr       = "192.168.128.0/18"
  vpc_id     = huaweicloud_iec_vpc.vpc_test.id
  site_id    = data.huaweicloud_iec_sites.sites_test.sites[0].id
  gateway_ip = "192.168.128.1"
}
`, rName, rName)
}

func testAccVpcSubnet_customer_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "%s-vpc"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name       = "%s-subnet"
  cidr       = "192.168.128.0/18"
  vpc_id     = huaweicloud_iec_vpc.vpc_test.id
  site_id    = data.huaweicloud_iec_sites.sites_test.sites[0].id
  gateway_ip = "192.168.128.1"
}
`, rName, rNameUpdate)
}
