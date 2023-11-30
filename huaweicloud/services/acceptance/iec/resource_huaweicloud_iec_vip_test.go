package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ieccommon "github.com/chnsz/golangsdk/openstack/iec/v1/common"
	"github.com/chnsz/golangsdk/openstack/iec/v1/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVipResource_basic(t *testing.T) {
	var iecPort ieccommon.Port
	rName := fmt.Sprintf("iec-%s", acctest.RandString(5))
	resourceName := "huaweicloud_iec_vip.vip_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVip_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVipExists(resourceName, &iecPort),
					resource.TestCheckResourceAttr(resourceName, "ip_address", "192.168.0.100"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
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

func TestAccVipResource_associate(t *testing.T) {
	var iecPort ieccommon.Port
	rName := fmt.Sprintf("iec-%s", acctest.RandString(5))
	resourceName := "huaweicloud_iec_vip.vip_test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckVipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVip_associate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVipExists(resourceName, &iecPort),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttr(resourceName, "allowed_addresses.#", "1"),
				),
			},
			{
				Config: testAccVip_disassociate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttr(resourceName, "allowed_addresses.#", "0"),
				),
			},
		},
	})
}

func testAccCheckVipDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	iecV1Client, err := conf.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating IEC client: %s", err)
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

func testAccCheckVipExists(n string, vipResource *ieccommon.Port) resource.TestCheckFunc {
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
			return fmt.Errorf("error creating IEC client: %s", err)
		}

		found, err := ports.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC Vip not found")
		}

		*vipResource = *found

		return nil
	}
}

func testAccVip_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "vpc-%s"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name       = "subnet-%s"
  vpc_id     = huaweicloud_iec_vpc.vpc_test.id
  site_id    = data.huaweicloud_iec_sites.sites_test.sites[0].id
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
}
`, rName, rName)
}

func testAccVip_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iec_vip" "vip_test" {
  subnet_id  = huaweicloud_iec_vpc_subnet.subnet_test.id
  ip_address = "192.168.0.100"
}
`, testAccVip_base(rName))
}

func testAccVip_associate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iec_vip" "vip_test" {
  subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
  port_ids  = [huaweicloud_iec_server.server_test.nics[0].port]
}
`, testAccServer_basic(rName))
}

func testAccVip_disassociate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iec_vip" "vip_test" {
  subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
}
`, testAccServer_basic(rName))
}
