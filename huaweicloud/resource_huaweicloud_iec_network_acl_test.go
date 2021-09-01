package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIecNetworkACLResource_basic(t *testing.T) {
	rName := fmt.Sprintf("iec-acl-%s", acctest.RandString(5))
	resourceKey := "huaweicloud_iec_network_acl.acl_demo"
	var fwGroup firewalls.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "description", "Created by terraform test acc"),
					testAccCheckIecNetworkACLNetBlockExists(&fwGroup),
				),
			},
			{
				ResourceName:      resourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIecNetworkACL_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceKey, "description", "Updated by terraform test acc"),
					testAccCheckIecNetworkACLNetBlockExists(&fwGroup),
				),
			},
		},
	})
}

func TestAccIecNetworkACLResource_no_subnets(t *testing.T) {
	rName := fmt.Sprintf("acc-fw-%s", acctest.RandString(5))
	resourceKey := "huaweicloud_iec_network_acl.acl_demo"
	var fwGroup firewalls.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName+"-noSubnet"),
					resource.TestCheckResourceAttr(resourceKey, "description", "Iec network acl without subents"),
					resource.TestCheckResourceAttr(resourceKey, "status", "INACTIVE"),
				),
			},
		},
	})
}

func testAccCheckIecNetworkACLDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_network_acl" {
			continue
		}

		_, err := firewalls.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("IEC network acl still exists")
		}
	}

	return nil
}

func testAccCheckIecNetworkACLExists(n string, resource *firewalls.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := firewalls.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("IEC Network ACL not found")
		}

		*resource = *found

		return nil
	}
}

func testAccCheckIecNetworkACLNetBlockExists(resource *firewalls.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(resource.Subnets) == 0 {
			return fmtp.Errorf("The Subnet of IEC Network ACL is not set.")
		}
		return nil
	}
}

var testAccIecNetworkACLRules string = `
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "vpc_demo"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name        = "subnet_demo"
  cidr        = "192.168.128.0/18"
  vpc_id      = huaweicloud_iec_vpc.vpc_test.id
  site_id     = data.huaweicloud_iec_sites.sites_test.sites[0].id
  gateway_ip  = "192.168.128.3"
}
`

func testAccIecNetworkACL_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iec_network_acl" "acl_demo" {
  name = "%s"
  description = "Created by terraform test acc"
  networks {
    vpc_id = huaweicloud_iec_vpc.vpc_test.id
    subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
  }
}
`, testAccIecNetworkACLRules, rName)
}

func testAccIecNetworkACL_basic_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_iec_network_acl" "acl_demo" {
  name = "%s-update"
  description = "Updated by terraform test acc"
  networks {
    vpc_id = huaweicloud_iec_vpc.vpc_test.id
    subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
  }
}
`, testAccIecNetworkACLRules, rName)
}

func testAccIecNetworkACL_no_subnets(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_network_acl" "acl_demo" {
  name        = "%s-noSubnet"
  description = "Iec network acl without subents"
}
`, rName)
}
