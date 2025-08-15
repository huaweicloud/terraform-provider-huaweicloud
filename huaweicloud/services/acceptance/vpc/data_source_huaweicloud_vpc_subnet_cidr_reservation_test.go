package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVpcSubnetCidrReservation_basic(t *testing.T) {
	var name = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_vpc_subnet_cidr_reservation.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVpcSubnetCidrReservationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcSubnetCidrReservation_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "4"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "mask"),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region", // Field not returned by API
				},
			},
			{
				Config: testAccVpcSubnetCidrReservation_update(name + "-update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name+"-update"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform"),
				),
			},
		},
	})
}

func testAccCheckVpcSubnetCidrReservationDestroy(s *terraform.State) error {
	// Implement actual check logic here
	return nil
}

func testAccVpcSubnetCidrReservation_basic(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_subnet_cidr_reservation" "test" {
  subnet_id = huaweicloud_vpc_subnet.test.id
  ip_version   = 4
  mask         = 26
  name         = "%s"
  description  = "created by terraform"
}
`, name, name, name)
}

func testAccVpcSubnetCidrReservation_update(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_subnet_cidr_reservation" "test" {
  subnet_id = huaweicloud_vpc_subnet.test.id
  ip_version   = 4
  mask         = 26
  name         = "%s"
  description  = "updated by terraform"
}
`, name, name, name)
}
