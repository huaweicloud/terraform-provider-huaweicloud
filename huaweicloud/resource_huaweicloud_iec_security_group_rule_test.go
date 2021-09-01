package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/iec/v1/security/groups"
	"github.com/chnsz/golangsdk/openstack/iec/v1/security/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccIecSecurityGroupRuleResource_Basic(t *testing.T) {

	groupName := "huaweicloud_iec_security_group.my_group"
	ruleName1 := "huaweicloud_iec_security_group_rule.rule_1"
	ruleName2 := "huaweicloud_iec_security_group_rule.rule_2"
	rName := fmt.Sprintf("iec-secgroup-%s", acctest.RandString(5))

	var group groups.RespSecurityGroupEntity
	var rule1, rule2 rules.RespSecurityGroupRule

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIecSecurityGroupRuleV1Destory,
		Steps: []resource.TestStep{
			{
				Config: testAccIecSecurityGroupRuleV1_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIecSecGroupV1Exists(groupName, &group),
					testAccCheckIecSecurityGroupRuleV1Exists(ruleName1, &rule1),
					resource.TestCheckResourceAttr(ruleName1, "direction", "egress"),
					resource.TestCheckResourceAttr(ruleName1, "protocol", "tcp"),
					resource.TestCheckResourceAttr(ruleName1, "port_range_min", "445"),
					resource.TestCheckResourceAttr(ruleName1, "port_range_max", "445"),
					testAccCheckIecSecurityGroupRuleV1Exists(ruleName2, &rule2),
					resource.TestCheckResourceAttr(ruleName2, "direction", "ingress"),
					resource.TestCheckResourceAttr(ruleName2, "protocol", "udp"),
					resource.TestCheckResourceAttr(ruleName2, "port_range_min", "20"),
					resource.TestCheckResourceAttr(ruleName2, "port_range_max", "20"),
				),
			},
		},
	})
}

func testAccCheckIecSecGroupV1Exists(n string, group *groups.RespSecurityGroupEntity) resource.TestCheckFunc {

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID has been seted")
		}

		config := testAccProvider.Meta().(*config.Config)
		iecClient, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := groups.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("IEC Security group not found")
		}
		*group = *found
		return nil
	}
}

func testAccCheckIecSecurityGroupRuleV1Exists(n string, rule *rules.RespSecurityGroupRule) resource.TestCheckFunc {

	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not fount: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID has been seted")
		}

		config := testAccProvider.Meta().(*config.Config)
		iecClient, err := config.IECV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
		}

		found, err := rules.Get(iecClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		if found.SecurityGroupRule.ID != rs.Primary.ID {
			return fmtp.Errorf("IEC Security group rule not found")
		}
		*rule = *found
		return nil
	}
}

func testAccCheckIecSecurityGroupRuleV1Destory(state *terraform.State) error {

	config := testAccProvider.Meta().(*config.Config)
	iecClient, err := config.IECV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	for _, rs := range state.RootModule().Resources {
		if rs.Type != "huaweicloud_iec_security_group_rule" {
			continue
		}

		_, err := rules.Get(iecClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("IEC Security group rule still exists")
		}
	}

	return nil
}

func testAccIecSecurityGroupRuleV1_Basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_iec_security_group" "my_group" {
  name = "%s"
}

resource "huaweicloud_iec_security_group_rule" "rule_1" {
  direction = "egress"
  port_range_min = 445
  port_range_max = 445
  protocol = "tcp" 
  security_group_id = huaweicloud_iec_security_group.my_group.id
  remote_ip_prefix = "0.0.0.0/0"
}

resource "huaweicloud_iec_security_group_rule" "rule_2" {
  direction = "ingress"
  port_range_min = "20"
  port_range_max = "20"
  protocol = "udp" 
  security_group_id = huaweicloud_iec_security_group.my_group.id
  remote_ip_prefix = "0.0.0.0/0"
}
`, rName)
}
