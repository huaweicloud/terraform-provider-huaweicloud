package vpc

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getACLRuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.FwV2Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud Network client: %s", err)
	}
	return rules.Get(c, state.Primary.ID).Extract()
}

func TestAccNetworkACLRule_basic(t *testing.T) {
	var rule rules.Rule
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_network_acl_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&rule,
		getACLRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_basic_1(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "udp"),
					resource.TestCheckResourceAttr(resourceName, "action", "deny"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(resourceName, "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(resourceName, "source_port", "444"),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "555"),
				),
			},
			{
				Config: testAccNetworkACLRule_basic_3(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(resourceName, "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(resourceName, "source_port", "666"),
					resource.TestCheckResourceAttr(resourceName, "destination_port", "777"),
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

func TestAccNetworkACLRule_anyProtocol(t *testing.T) {
	var rule rules.Rule
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_network_acl_rule.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&rule,
		getACLRuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkACLRule_anyProtocol(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "any"),
					resource.TestCheckResourceAttr(resourceName, "action", "allow"),
					resource.TestCheckResourceAttr(resourceName, "enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "source_ip_address", "192.168.199.0/24"),
				),
			},
		},
	})
}

func testAccNetworkACLRule_basic_1(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "test" {
  name     = "%s"
  protocol = "udp"
  action   = "deny"
}
`, rName)
}

func testAccNetworkACLRule_basic_2(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "test" {
  name                   = "%s"
  description            = "Terraform accept test"
  protocol               = "udp"
  action                 = "deny"
  source_ip_address      = "1.2.3.4"
  destination_ip_address = "4.3.2.0/24"
  source_port            = "444"
  destination_port       = "555"
  enabled                = true
}
`, rName)
}

func testAccNetworkACLRule_basic_3(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "test" {
  name                   = "%s"
  description            = "Terraform accept test updated"
  protocol               = "tcp"
  action                 = "allow"
  source_ip_address      = "1.2.3.0/24"
  destination_ip_address = "4.3.2.8"
  source_port            = "666"
  destination_port       = "777"
  enabled                = false
}
`, rName)
}

func testAccNetworkACLRule_anyProtocol(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_network_acl_rule" "test" {
  name              = "%s"
  description       = "Allow any protocol"
  protocol          = "any"
  action            = "allow"
  source_ip_address = "192.168.199.0/24"
  enabled           = true
}
`, rName)
}
