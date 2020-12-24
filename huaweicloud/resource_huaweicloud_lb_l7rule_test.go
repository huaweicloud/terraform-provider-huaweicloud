package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	l7rules "github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/l7policies"
)

func TestAccLBV2L7Rule_basic(t *testing.T) {
	var l7rule l7rules.Rule
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_lb_l7rule.l7rule_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2L7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLBV2L7RuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists(resourceName, &l7rule),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(resourceName, "value", "/api"),
					resource.TestMatchResourceAttr(resourceName, "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(resourceName, "l7policy_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: testAccCheckLBV2L7RuleConfig_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists(resourceName, &l7rule),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "value", "/images"),
				),
			},
		},
	})
}

func testAccCheckLBV2L7RuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := config.ElbV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lb_l7rule" {
			continue
		}

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("Unable to find l7policy_id")
		}

		_, err := l7rules.GetRule(lbClient, l7policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Rule still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2L7RuleExists(n string, l7rule *l7rules.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.ElbV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud load balancing client: %s", err)
		}

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("Unable to find l7policy_id")
		}

		found, err := l7rules.GetRule(lbClient, l7policyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Policy not found")
		}

		*l7rule = *found

		return nil
	}
}

func testAccCheckLBV2L7RuleConfig(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_l7policy" "l7policy_1" {
  name         = "%s"
  action       = "REDIRECT_TO_POOL"
  description  = "test description"
  position     = 1
  listener_id  = huaweicloud_lb_listener.listener_1.id
  redirect_pool_id = huaweicloud_lb_pool.pool_1.id
}
`, rName, rName, rName, rName)
}

func testAccCheckLBV2L7RuleConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_l7rule" "l7rule_1" {
  l7policy_id  = huaweicloud_lb_l7policy.l7policy_1.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
`, testAccCheckLBV2L7RuleConfig(rName))
}

func testAccCheckLBV2L7RuleConfig_update2(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_l7rule" "l7rule_1" {
  l7policy_id  = huaweicloud_lb_l7policy.l7policy_1.id
  type         = "PATH"
  compare_type = "STARTS_WITH"
  value        = "/images"
}
`, testAccCheckLBV2L7RuleConfig(rName))
}
