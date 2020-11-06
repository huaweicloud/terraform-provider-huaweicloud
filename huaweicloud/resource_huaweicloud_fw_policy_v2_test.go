package huaweicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/policies"
)

func TestAccFWPolicyV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"huaweicloud_fw_policy_v2.policy_1", "", "", 0),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_addRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_addRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"huaweicloud_fw_policy_v2.policy_1", "policy_1", "terraform acceptance test", 2),
				),
			},
			{
				ResourceName:      "huaweicloud_fw_policy_v2.policy_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccFWPolicyV2_deleteRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_deleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"huaweicloud_fw_policy_v2.policy_1", "policy_1", "terraform acceptance test", 1),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_timeout(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"huaweicloud_fw_policy_v2.policy_1", "", "", 0),
				),
			},
		},
	})
}

func testAccCheckFWPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	fwClient, err := config.fwV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_fw_policy_v2" {
			continue
		}
		_, err = policies.Get(fwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall policy (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWPolicyV2Exists(n, name, description string, ruleCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		fwClient, err := config.fwV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
		}

		var found *policies.Policy
		for i := 0; i < 5; i++ {
			// Firewall policy creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = policies.Get(fwClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case name != found.Name:
			err = fmt.Errorf("Expected name <%s>, but found <%s>", name, found.Name)
		case description != found.Description:
			err = fmt.Errorf("Expected description <%s>, but found <%s>", description, found.Description)
		case ruleCount != len(found.Rules):
			err = fmt.Errorf("Expected rule count <%d>, but found <%d>", ruleCount, len(found.Rules))
		}

		if err != nil {
			return err
		}

		return nil
	}
}

const testAccFWPolicyV2_basic = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
}
`

const testAccFWPolicyV2_addRules = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
  description =  "terraform acceptance test"
  rules = [
    "${huaweicloud_fw_rule_v2.udp_deny.id}",
    "${huaweicloud_fw_rule_v2.tcp_allow.id}"
  ]
}

resource "huaweicloud_fw_rule_v2" "tcp_allow" {
  protocol = "tcp"
  action = "allow"
}

resource "huaweicloud_fw_rule_v2" "udp_deny" {
  protocol = "udp"
  action = "deny"
}
`

const testAccFWPolicyV2_deleteRules = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
  name = "policy_1"
  description =  "terraform acceptance test"
  rules = [
    "${huaweicloud_fw_rule_v2.udp_deny.id}"
  ]
}

resource "huaweicloud_fw_rule_v2" "udp_deny" {
  protocol = "udp"
  action = "deny"
}
`

const testAccFWPolicyV2_timeout = `
resource "huaweicloud_fw_policy_v2" "policy_1" {
  timeouts {
    create = "5m"
  }
}
`
