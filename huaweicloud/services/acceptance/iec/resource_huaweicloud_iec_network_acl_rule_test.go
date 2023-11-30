package iec

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getNetworkACLRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	iecClient, err := conf.IECV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating IEC client: %s", err)
	}

	fwGroup, err := firewalls.Get(iecClient, state.Primary.ID).Extract()
	if err != nil {
		return nil, err
	}

	if len(fwGroup.IngressFWPolicy.FirewallRules) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return fwGroup, nil
}

func TestAccNetworkACLRuleResource_basic(t *testing.T) {
	var fwGroup firewalls.RespFirewallRulesEntity

	aclResourceName := "huaweicloud_iec_network_acl.acl_test"
	aclRuleResourceName := "huaweicloud_iec_network_acl_rule.rule_test"

	rc := acceptance.InitResourceCheck(
		aclResourceName,
		&fwGroup,
		getNetworkACLRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_basic(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclRuleResourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(aclRuleResourceName, "action", "allow"),
					resource.TestCheckResourceAttr(aclRuleResourceName, "destination_port", "445"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_update(),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(aclRuleResourceName, "protocol", "udp"),
					resource.TestCheckResourceAttr(aclRuleResourceName, "action", "deny"),
					resource.TestCheckResourceAttr(aclRuleResourceName, "destination_port", "23-30"),
				),
			},
		},
	})
}

func testAccNetworkACLRule_basic() string {
	return `
resource "huaweicloud_iec_network_acl" "acl_test" {
  name = "iec-acl-basic"
}

resource "huaweicloud_iec_network_acl_rule" "rule_test" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_test.id
  direction              = "ingress"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "0.0.0.0/0"
  destination_ip_address = "0.0.0.0/0"
  destination_port       = "445"
  enabled                = true
}
`
}

func testAccNetworkACLRule_basic_update() string {
	return `
resource "huaweicloud_iec_network_acl" "acl_test" {
  name = "iec-acl-update"
}

resource "huaweicloud_iec_network_acl_rule" "rule_test" {
  network_acl_id         = huaweicloud_iec_network_acl.acl_test.id
  direction              = "ingress"
  protocol               = "udp"
  action                 = "deny"
  source_ip_address      = "0.0.0.0/0"
  destination_ip_address = "0.0.0.0/0"
  destination_port       = "23-30"
  enabled                = true
}
`
}
