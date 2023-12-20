package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v1/security/rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNetworkSecGroupRuleResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := cfg.NetworkingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating VPC network v1 client: %s", err)
	}

	return rules.Get(client, state.Primary.ID)
}

func TestAccNetworkingSecGroupRule_basic(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a basic acc test"),
					resource.TestCheckResourceAttr(resourceName, "ports", "80"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
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

func TestAccNetworkingSecGroupRule_Egress(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_egress(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "direction", "egress"),
					resource.TestCheckResourceAttr(resourceName, "description", "This is a basic acc test"),
					resource.TestCheckResourceAttr(resourceName, "ports", "80"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
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

func TestAccNetworkingSecGroupRule_oldPorts(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_oldPorts(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "port_range_min", "80"),
					resource.TestCheckResourceAttr(resourceName, "port_range_max", "80"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", "0.0.0.0/0"),
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

func TestAccNetworkingSecGroupRule_remoteGroup(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_remoteGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "ports", "80"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttrSet(resourceName, "remote_group_id"),
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

func TestAccNetworkingSecGroupRule_lowerCaseCIDR(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.secgroup_rule_test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_lowerCaseCIDR(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", "2001:558:fc00::/39"),
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

func TestAccNetworkingSecGroupRule_noPorts(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_noPorts(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "direction", "ingress"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "icmp"),
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", "0.0.0.0/0"),
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

func TestAccNetworkingSecGroupRule_remoteAddressGroup(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_remoteAddressGroup(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "remote_address_group_id",
						"huaweicloud_vpc_address_group.test", "id"),
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

func TestAccNetworkingSecGroupRule_noRemote(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	resourceName := "huaweicloud_networking_secgroup_rule.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_noRemote(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "direction", "egress"),
					resource.TestCheckResourceAttr(resourceName, "ethertype", "IPv4"),
					resource.TestCheckResourceAttr(resourceName, "priority", "1"),
					// the IP address is shown as 0.0.0.0/0 in console, but it is empty from API response
					resource.TestCheckResourceAttr(resourceName, "remote_ip_prefix", ""),
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

func TestAccNetworkingSecGroupRule_action(t *testing.T) {
	var (
		secgroupRule rules.SecurityGroupRule
		allowResName string = "huaweicloud_networking_secgroup_rule.allow"
		denyResName  string = "huaweicloud_networking_secgroup_rule.deny"
	)

	rName := acceptance.RandomAccResourceNameWithDash()

	rc1 := acceptance.InitResourceCheck(
		allowResName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	rc2 := acceptance.InitResourceCheck(
		denyResName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc1.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_action(rName),
				Check: resource.ComposeTestCheckFunc(
					rc1.CheckResourceExists(),
					resource.TestCheckResourceAttr(allowResName, "action", "allow"),
					rc2.CheckResourceExists(),
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
	resourceName := "huaweicloud_networking_secgroup_rule.test"
	rName := acceptance.RandomAccResourceNameWithDash()

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getNetworkSecGroupRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingSecGroupRule_priority(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "priority", "50"),
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

func testAccNetworkingSecGroupRule_egress(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_test" {
  direction         = "egress"
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

func testAccNetworkingSecGroupRule_noRemote(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "test" {
  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
  direction         = "egress"
  ethertype         = "IPv4"
}
`, testAccSecGroup_noDefaultRules(rName))
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
