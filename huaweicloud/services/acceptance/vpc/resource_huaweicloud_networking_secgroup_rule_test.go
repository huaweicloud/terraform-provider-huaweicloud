package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"
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

func testAccCheckNetworkingSecGroupRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV3Client(HW_REGION_NAME)
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
		networkingClient, err := config.NetworkingV3Client(HW_REGION_NAME)
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
