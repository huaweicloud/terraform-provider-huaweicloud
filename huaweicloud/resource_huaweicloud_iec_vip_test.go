package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	iec_common "github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/ports"
)

func TestAccIecVipResource_basic(t *testing.T) {
	var iecPort iec_common.Port
	resourceName := "huaweicloud_iec_vip.vip_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecVip_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecVipExists(resourceName, &iecPort),
					resource.TestMatchResourceAttr(resourceName, "mac_address", regexp.MustCompile("^[0-9A-Fa-f]{2}\\:[0-9A-Fa-f]{2}\\:[0-9A-Fa-f]{2}\\:[0-9A-Fa-f]{2}\\:[0-9A-Fa-f]{2}\\:[0-9A-Fa-f]{2}$")),
					resource.TestCheckResourceAttr(resourceName, "fixed_ips.#", "1"),
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

func testAccCheckIecVipDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_vip" {
			continue
		}

		_, err := ports.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC Vip still exists")
		}
	}

	return nil
}

func testAccCheckIecVipExists(n string, resource *iec_common.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := ports.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC Vip not found")
		}

		*resource = *found

		return nil
	}
}

func testAccIecVip_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "iec-vpc-demo"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name = "iec-subnet-demo"
  cidr = "192.168.0.0/16"
  gateway_ip = "192.168.0.1"
  vpc_id = huaweicloud_iec_vpc.vpc_test.id
  site_id = data.huaweicloud_iec_sites.sites_test.sites[0].id
}

resource "huaweicloud_iec_vip" "vip_test" {
  subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
}
`)
}
