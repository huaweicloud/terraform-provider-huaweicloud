package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/iec/v1/firewalls"
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
					testAccCheckIecNetworkACLRulesCount(resourceKey, &fwGroup, 1, 1),
				),
			},
			{
				Config: testAccIecNetworkACL_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceKey, "description", "Updated by terraform test acc"),
					testAccCheckIecNetworkACLRulesCount(resourceKey, &fwGroup, 1, 1),
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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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

func TestAccIecNetworkACLResource_remove(t *testing.T) {
	rName := fmt.Sprintf("iec-acl-%s", acctest.RandString(5))
	resourceKey := "huaweicloud_iec_network_acl.acl_demo"
	var fwGroup firewalls.Firewall

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecNetworkACLDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIecNetworkACL_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName+"-update"),
					resource.TestCheckResourceAttr(resourceKey, "description", "Updated by terraform test acc"),
					testAccCheckIecNetworkACLRulesCount(resourceKey, &fwGroup, 1, 1),
				),
			},
			{
				Config: testAccIecNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "description", "Created by terraform test acc"),
					testAccCheckIecNetworkACLRulesCount(resourceKey, &fwGroup, 1, 1),
				),
			},
			{
				Config: testAccIecNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecNetworkACLExists(resourceKey, &fwGroup),
					resource.TestCheckResourceAttr(resourceKey, "status", "INACTIVE"),
				),
			},
		},
	})
}

func testAccCheckIecNetworkACLDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	iecV1Client, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_network_acl" {
			continue
		}

		_, err := firewalls.Get(iecV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IEC network acl still exists")
		}
	}

	return nil
}

func testAccCheckIecNetworkACLExists(n string, resource *firewalls.Firewall) resource.TestCheckFunc {
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

		found, err := firewalls.Get(iecV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("IEC Network ACL not found")
		}

		*resource = *found

		return nil
	}
}

func testAccCheckIecNetworkACLRulesCount(n string, resource *firewalls.Firewall, inboundCount int, outboundCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, _ := s.RootModule().Resources[n]
		config := testAccProvider.Meta().(*Config)
		iecV1Client, _ := config.IECV1Client(HW_REGION_NAME)
		found, _ := firewalls.Get(iecV1Client, rs.Primary.ID).Extract()
		if len(found.IngressFWPolicy.FirewallRules) != inboundCount {
			return fmt.Errorf("IEC Network ACL not found")
		}
		if len(found.EgressFWPolicy.FirewallRules) != outboundCount {
			return fmt.Errorf("IEC Network ACL not found")
		}
		*resource = *found

		return nil
	}
}

const testAccIecNetworkACLRules = `
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

resource "huaweicloud_iec_network_acl_rule" "rule_1" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_demo.id
  direction              = "ingress"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "132.156.0.0/16"
  destination_ip_address = "192.168.128.0/18"
  destination_port       = "445"
  enabled                = true
}

resource "huaweicloud_iec_network_acl_rule" "rule_2" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_demo.id
  direction              = "egress"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "192.168.128.0/18"
  destination_ip_address = "152.16.30.0/24"
  destination_port       = "45"
  enabled                = true
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
