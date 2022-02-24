package vpc

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/firewall_groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getACLResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FwV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return firewall_groups.Get(c, state.Primary.ID).Extract()
}

func TestAccNetworkACL_basic(t *testing.T) {
	var fwGroup firewall_groups.FirewallGroup
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_network_acl.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&fwGroup,
		getACLResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by terraform test acc"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "inbound_policy_id"),
					resource.TestCheckResourceAttr(resourceName, "inbound_rules.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "1"),
				),
			},
			{
				Config: testAccNetworkACL_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName+"_update"),
					resource.TestCheckResourceAttr(resourceName, "description", "updated by terraform test acc"),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "inbound_rules.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkACL_no_subnets(t *testing.T) {
	var fwGroup firewall_groups.FirewallGroup
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_network_acl.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&fwGroup,
		getACLResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "network acl without subents"),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "inbound_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "0"),
				),
			},
		},
	})
}

func TestAccNetworkACL_remove(t *testing.T) {
	var fwGroup firewall_groups.FirewallGroup
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_network_acl.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&fwGroup,
		getACLResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACL_basic_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "2"),
				),
			},
			{
				Config: testAccNetworkACL_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "1"),
				),
			},
			{
				Config: testAccNetworkACL_no_subnets(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "status", "INACTIVE"),
					resource.TestCheckResourceAttr(resourceName, "subnets.#", "0"),
				),
			},
		},
	})
}

func testAccNetworkACLRules(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "subnet_1" {
  name       = "%s_1"
  cidr       = "192.168.0.0/24"
  gateway_ip = "192.168.0.1"
  vpc_id     = huaweicloud_vpc.test.id
}

resource "huaweicloud_vpc_subnet" "subnet_2" {
	name       = "%s_2"
	cidr       = "192.168.10.0/24"
	gateway_ip = "192.168.10.1"
	vpc_id     = huaweicloud_vpc.test.id
  }

resource "huaweicloud_network_acl_rule" "rule_1" {
  name             = "%s_1"
  description      = "drop TELNET traffic"
  action           = "deny"
  protocol         = "tcp"
  destination_port = "23"
}

resource "huaweicloud_network_acl_rule" "rule_2" {
  name             = "%s_2"
  description      = "drop NTP traffic"
  action           = "deny"
  protocol         = "udp"
  destination_port = "123"
}
`, name, name, name, name, name)
}

func testAccNetworkACL_basic(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_network_acl" "test" {
  name        = "%s"
  description = "created by terraform test acc"

  inbound_rules = [huaweicloud_network_acl_rule.rule_1.id]
  subnets       = [huaweicloud_vpc_subnet.subnet_1.id]
}
`, testAccNetworkACLRules(name), name)
}

func testAccNetworkACL_basic_update(name string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_network_acl" "test" {
  name        = "%s_update"
  description = "updated by terraform test acc"

  inbound_rules = [
    huaweicloud_network_acl_rule.rule_1.id,
    huaweicloud_network_acl_rule.rule_2.id,
  ]
  subnets = [
    huaweicloud_vpc_subnet.subnet_1.id,
    huaweicloud_vpc_subnet.subnet_2.id,
  ]
}
`, testAccNetworkACLRules(name), name)
}

func testAccNetworkACL_no_subnets(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl" "test" {
  name        = "%s"
  description = "network acl without subents"
}
`, name)
}
