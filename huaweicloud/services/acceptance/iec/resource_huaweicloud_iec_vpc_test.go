package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/vpcs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIecVpcV1_basic(t *testing.T) {
	var iecVPC ieccommon.VPC

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_iec_vpc.test"
	rNameUpdate := rName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIecVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecVpcV1_system(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecVpcV1Exists(resourceName, &iecVPC),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "192.168.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "mode", "SYSTEM"),
				),
			},
			{
				Config: testAccIecVpcV1_system_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecVpcV1Exists(resourceName, &iecVPC),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
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

func TestAccIecVpcV1_customer(t *testing.T) {
	var iecVPC ieccommon.VPC

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_iec_vpc.customer"
	rNameUpdate := rName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckIecVpcV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecVpcV1_customer(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecVpcV1Exists(resourceName, &iecVPC),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "cidr", "172.16.0.0/16"),
					resource.TestCheckResourceAttr(resourceName, "mode", "CUSTOMER"),
				),
			},
			{
				Config: testAccIecVpcV1_customer_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecVpcV1Exists(resourceName, &iecVPC),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "cidr", "172.30.0.0/16"),
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

func testAccCheckIecVpcV1Destroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	iecV1Client, err := conf.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_vpc" {
			continue
		}

		_, err := vpcs.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC VPC still exists")
		}
	}

	return nil
}

func testAccCheckIecVpcV1Exists(n string, vpcResource *ieccommon.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		conf := acceptance.TestAccProvider.Meta().(*config.Config)
		iecV1Client, err := conf.IECV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating IEC client: %s", err)
		}

		found, err := vpcs.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC VPC not found")
		}

		*vpcResource = *found

		return nil
	}
}

func testAccIecVpcV1_system(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}
`, rName)
}

func testAccIecVpcV1_system_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}
`, rName)
}

func testAccIecVpcV1_customer(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_vpc" "customer" {
  name = "%s"
  cidr = "172.16.0.0/16"
  mode = "CUSTOMER"
}
`, rName)
}

func testAccIecVpcV1_customer_update(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_vpc" "customer" {
  name = "%s"
  cidr = "172.30.0.0/16"
  mode = "CUSTOMER"
}
`, rName)
}
