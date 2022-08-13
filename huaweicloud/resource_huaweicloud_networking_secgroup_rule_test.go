package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/security/rules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkingSecGroupRule_basic(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "description", "This is a basic acc test"),
					resource.TestCheckResourceAttr(resourceRuleName, "ports", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceRuleName, "priority", "1"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_oldPorts(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_oldPorts(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_min", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "port_range_max", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "0.0.0.0/0"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_remoteGroup(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_remoteGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "ports", "80"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(resourceRuleName, "remote_group_id"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_lowerCaseCIDR(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_lowerCaseCIDR(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "2001:558:fc00::/39"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_noPorts(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_noPorts(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceRuleName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceRuleName, "protocol", "icmp"),
					resource.TestCheckResourceAttr(resourceRuleName, "remote_ip_prefix", "0.0.0.0/0"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_remoteAddressGroup(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_remoteAddressGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttrPair(resourceRuleName, "remote_address_group_id",
						"huaweicloud_vpc_address_group.test", "id"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_action(t *testing.T) {
	var (
		secgroupRule rules.SecurityGroupRule
		allowResName string = "huaweicloud_networking_secgroup_rule.allow"
		denyResName  string = "huaweicloud_networking_secgroup_rule.deny"
	)

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_action(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(allowResName, &secgroupRule),
					resource.TestCheckResourceAttr(allowResName, "action", "allow"),
					testAccCheckNetworkingSecGroupRuleExists(denyResName, &secgroupRule),
					resource.TestCheckResourceAttr(denyResName, "action", "deny"),
				),
			},
			{
				ResourceName:      allowResName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      denyResName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingSecGroupRule_priority(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	var resourceRuleName string = "huaweicloud_networking_secgroup_rule.test"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingSecGroupRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_priority(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupRuleExists(resourceRuleName, &secgroupRule),
					resource.TestCheckResourceAttr(resourceRuleName, "priority", "50"),
				),
			},
			{
				ResourceName:      resourceRuleName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckNetworkingSecGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_networking_secgroup_rule" {
			continue
		}

		_, err := rules.Get(networkingClient, rs.Primary.ID)
		if err == nil {
			return fmtp.Errorf("Security group rule still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingSecGroupRuleExists(n string, secGroupRule *rules.SecurityGroupRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := rules.Get(networkingClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Security group rule not found")
		}

		*secGroupRule = *found

		return nil
	}
}

func testAccNetworkingSecGroupRule_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_test" {
  name        = "%s-secgroup"
  description = "terraform security group rule acceptance test"
}
`, rName)
}

func testAccNetworkingSecGroupRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  description       = "This is a basic acc test"
  ethertype         = "IPv4"
  ports             = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_oldPorts(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 80
  port_range_max    = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_remoteGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 80
  protocol          = "tcp"
  remote_group_id   = huaweicloud_networking_secgroup.secgroup_test.id
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_lowerCaseCIDR(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "ingress"
  ethertype         = "IPv6"
  ports             = 80
  protocol          = "tcp"
  remote_ip_prefix  = "2001:558:FC00::/39"
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_noPorts(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_remoteAddressGroup(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpc_address_group" "test" {
  name = "%[2]s"

  addresses = [
    "192.168.10.12",
    "192.168.11.0-192.168.11.240",
  ]
}

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id       = huaweicloud_networking_secgroup.secgroup_test.id
  direction               = "ingress"
  ethertype               = "IPv4"
  ports                   = 80
  protocol                = "tcp"
  remote_address_group_id = huaweicloud_vpc_address_group.test.id
}
`, testAccNetworkingSecGroupRule_base(rName), rName)
}

func testAccNetworkingSecGroupRule_action(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "allow" {
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 80
  protocol          = "tcp"
  action            = "allow"
  remote_ip_prefix  = "0.0.0.0/0"
}

resource "huaweicloud_networking_secgroup_rule" "deny" {
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 8080
  protocol          = "tcp"
  action            = "deny"
  remote_ip_prefix  = "0.0.0.0/0"
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_priority(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.secgroup_test.id
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  priority          = 50
}
`, testAccNetworkingSecGroupRule_base(rName))
}
