package huaweicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/rules"
)

func TestAccNetworkACLRule_basic(t *testing.T) {
	resourceKey := "huaweicloud_network_acl_rule.rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_basic_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", "rule_1"),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceKey, "action", "deny"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "true"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", "rule_1"),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceKey, "action", "deny"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceKey, "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceKey, "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(resourceKey, "source_port", "444"),
					resource.TestCheckResourceAttr(resourceKey, "destination_port", "555"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", "rule_1"),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceKey, "action", "allow"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceKey, "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(resourceKey, "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(resourceKey, "source_port", "666"),
					resource.TestCheckResourceAttr(resourceKey, "destination_port", "777"),
				),
			},
			{
				ResourceName:      resourceKey,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkACLRule_anyProtocol(t *testing.T) {
	resourceKey := "huaweicloud_network_acl_rule.rule_any"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_anyProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", "rule_any"),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "any"),
					resource.TestCheckResourceAttr(resourceKey, "action", "allow"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceKey, "source_ip_address", "192.168.199.0/24"),
				),
			},
		},
	})
}

func testAccCheckNetworkACLRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	fwClient, err := config.fwV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_network_acl_rule" {
			continue
		}
		_, err = rules.Get(fwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Network ACL rule (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckNetworkACLRuleExists(key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[key]
		if !ok {
			return fmt.Errorf("Not found: %s", key)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set in %s", key)
		}

		config := testAccProvider.Meta().(*Config)
		fwClient, err := config.fwV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
		}

		var found *rules.Rule
		for i := 0; i < 5; i++ {
			// Network ACL rule creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = rules.Get(fwClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		if found == nil {
			return fmt.Errorf("Network ACL rule (%s) is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccNetworkACLRule_basic_1 = `
resource "huaweicloud_network_acl_rule" "rule_1" {
  name = "rule_1"
  protocol = "udp"
  action = "deny"
}
`

const testAccNetworkACLRule_basic_2 = `
resource "huaweicloud_network_acl_rule" "rule_1" {
  name = "rule_1"
  description = "Terraform accept test"
  protocol = "udp"
  action = "deny"
  source_ip_address = "1.2.3.4"
  destination_ip_address = "4.3.2.0/24"
  source_port = "444"
  destination_port = "555"
  enabled = true
}
`

const testAccNetworkACLRule_basic_3 = `
resource "huaweicloud_network_acl_rule" "rule_1" {
  name = "rule_1"
  description = "Terraform accept test updated"
  protocol = "tcp"
  action = "allow"
  source_ip_address = "1.2.3.0/24"
  destination_ip_address = "4.3.2.8"
  source_port = "666"
  destination_port = "777"
  enabled = false
}
`

const testAccNetworkACLRule_anyProtocol = `
resource "huaweicloud_network_acl_rule" "rule_any" {
  name = "rule_any"
  description = "Allow any protocol"
  protocol = "any"
  action = "allow"
  source_ip_address = "192.168.199.0/24"
  enabled = true
}
`
