package vpc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/networking/v3/security/rules"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getSecGroupRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.NetworkingV3Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return rules.Get(c, state.Primary.ID)
}

func TestAccNetworkingSecGroupRule_basic(t *testing.T) {
	var secgroupRule rules.SecurityGroupRule
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_networking_secgroup_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getSecGroupRuleResourceFunc,
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
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_networking_secgroup_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getSecGroupRuleResourceFunc,
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
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_networking_secgroup_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getSecGroupRuleResourceFunc,
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
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_networking_secgroup_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&secgroupRule,
		getSecGroupRuleResourceFunc,
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

func testAccNetworkingSecGroupRule_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "test" {
  name        = "%s-secgroup"
  description = "terraform security group rule acceptance test"
}
`, rName)
}

func testAccNetworkingSecGroupRule_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "test" {
  direction         = "ingress"
  description       = "This is a basic acc test"
  ethertype         = "IPv4"
  ports             = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_oldPorts(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_min    = 80
  port_range_max    = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_remoteGroup(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "test" {
  direction         = "ingress"
  ethertype         = "IPv4"
  ports             = 80
  protocol          = "tcp"
  remote_group_id   = huaweicloud_networking_secgroup.test.id
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}

func testAccNetworkingSecGroupRule_lowerCaseCIDR(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_networking_secgroup_rule" "test" {
  direction         = "ingress"
  ethertype         = "IPv6"
  ports             = 80
  protocol          = "tcp"
  remote_ip_prefix  = "2001:558:FC00::/39"
  security_group_id = huaweicloud_networking_secgroup.test.id
}
`, testAccNetworkingSecGroupRule_base(rName))
}
