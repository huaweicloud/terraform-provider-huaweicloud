package huaweicloud

import (
	"reflect"
	"testing"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccFWRuleV2_basic(t *testing.T) {
	rule1 := &rules.Rule{
		Name:      "rule_1",
		Protocol:  "udp",
		Action:    "deny",
		IPVersion: 4,
		Enabled:   true,
	}

	rule2 := &rules.Rule{
		Name:                 "rule_1",
		Protocol:             "udp",
		Action:               "deny",
		Description:          "Terraform accept test",
		IPVersion:            4,
		SourceIPAddress:      "1.2.3.4",
		DestinationIPAddress: "4.3.2.0/24",
		SourcePort:           "444",
		DestinationPort:      "555",
		Enabled:              true,
	}

	rule3 := &rules.Rule{
		Name:                 "rule_1",
		Protocol:             "tcp",
		Action:               "allow",
		Description:          "Terraform accept test updated",
		IPVersion:            4,
		SourceIPAddress:      "1.2.3.0/24",
		DestinationIPAddress: "4.3.2.8",
		SourcePort:           "666",
		DestinationPort:      "777",
		Enabled:              false,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2_basic_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("huaweicloud_fw_rule_v2.rule_1", rule1),
				),
			},
			{
				Config: testAccFWRuleV2_basic_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("huaweicloud_fw_rule_v2.rule_1", rule2),
				),
			},
			{
				ResourceName:      "huaweicloud_fw_rule_v2.rule_1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFWRuleV2_basic_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("huaweicloud_fw_rule_v2.rule_1", rule3),
				),
			},
		},
	})
}

func TestAccFWRuleV2_anyProtocol(t *testing.T) {
	rule := &rules.Rule{
		Name:            "rule_1",
		Description:     "Allow any protocol",
		Protocol:        "",
		Action:          "allow",
		IPVersion:       4,
		SourceIPAddress: "192.168.199.0/24",
		Enabled:         true,
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2_anyProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("huaweicloud_fw_rule_v2.rule_1", rule),
				),
			},
		},
	})
}

func testAccCheckFWRuleV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	fwClient, err := config.FwV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_fw_rule_v2" {
			continue
		}
		_, err = rules.Get(fwClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Firewall rule (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWRuleV2Exists(n string, expected *rules.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		fwClient, err := config.FwV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud fw client: %s", err)
		}

		var found *rules.Rule
		for i := 0; i < 5; i++ {
			// Firewall rule creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = rules.Get(fwClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					//lintignore:R018
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		expected.ID = found.ID
		// Erase the tenant id because we don't want to compare
		// it as long it is not present in the expected
		found.TenantID = ""

		if !reflect.DeepEqual(expected, found) {
			return fmtp.Errorf("Expected:\n%#v\nFound:\n%#v", expected, found)
		}

		return nil
	}
}

const testAccFWRuleV2_basic_1 = `
resource "huaweicloud_fw_rule_v2" "rule_1" {
	name = "rule_1"
	protocol = "udp"
	action = "deny"
}
`

const testAccFWRuleV2_basic_2 = `
resource "huaweicloud_fw_rule_v2" "rule_1" {
	name = "rule_1"
	description = "Terraform accept test"
	protocol = "udp"
	action = "deny"
	ip_version = 4
	source_ip_address = "1.2.3.4"
	destination_ip_address = "4.3.2.0/24"
	source_port = "444"
	destination_port = "555"
	enabled = true
}
`

const testAccFWRuleV2_basic_3 = `
resource "huaweicloud_fw_rule_v2" "rule_1" {
	name = "rule_1"
	description = "Terraform accept test updated"
	protocol = "tcp"
	action = "allow"
	ip_version = 4
	source_ip_address = "1.2.3.0/24"
	destination_ip_address = "4.3.2.8"
	source_port = "666"
	destination_port = "777"
	enabled = false
}
`

const testAccFWRuleV2_anyProtocol = `
resource "huaweicloud_fw_rule_v2" "rule_1" {
	name = "rule_1"
	description = "Allow any protocol"
	protocol = "any"
	action = "allow"
	ip_version = 4
	source_ip_address = "192.168.199.0/24"
	enabled = true
}
`
