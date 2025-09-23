package deprecated

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccNetworkACLRule_basic(t *testing.T) {
	resourceKey := "huaweicloud_network_acl_rule.rule_1"
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_basic_1(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
					resource.TestCheckResourceAttr(resourceKey, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceKey, "action", "deny"),
					resource.TestCheckResourceAttr(resourceKey, "enabled", "true"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
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
				Config: testAccNetworkACLRule_basic_3(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
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
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckDeprecated(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckNetworkACLRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_anyProtocol(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkACLRuleExists(resourceKey),
					resource.TestCheckResourceAttr(resourceKey, "name", rName),
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
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	fwClient, err := config.FwV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_network_acl_rule" {
			continue
		}
		_, err = rules.Get(fwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Network ACL rule (%s) still exists.", rs.Primary.ID)
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
			return fmtp.Errorf("Not found: %s", key)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set in %s", key)
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		fwClient, err := config.FwV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud fw client: %s", err)
		}

		found, err := rules.Get(fwClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Network ACL rule not found")
		}

		return nil
	}
}

func testAccNetworkACLRule_basic_1(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "rule_1" {
  name = "%s"
  protocol = "udp"
  action = "deny"
}
`, rName)
}

func testAccNetworkACLRule_basic_2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "rule_1" {
  name = "%s"
  description = "Terraform accept test"
  protocol = "udp"
  action = "deny"
  source_ip_address = "1.2.3.4"
  destination_ip_address = "4.3.2.0/24"
  source_port = "444"
  destination_port = "555"
  enabled = true
}
`, rName)
}

func testAccNetworkACLRule_basic_3(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "rule_1" {
  name = "%s"
  description = "Terraform accept test updated"
  protocol = "tcp"
  action = "allow"
  source_ip_address = "1.2.3.0/24"
  destination_ip_address = "4.3.2.8"
  source_port = "666"
  destination_port = "777"
  enabled = false
}
`, rName)
}

func testAccNetworkACLRule_anyProtocol(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "rule_any" {
  name = "%s"
  description = "Allow any protocol"
  protocol = "any"
  action = "allow"
  source_ip_address = "192.168.199.0/24"
  enabled = true
}
`, rName)
}
