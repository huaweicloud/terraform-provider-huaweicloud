package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	l7rules "github.com/chnsz/golangsdk/openstack/elb/v2/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getL7RuleResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud LB v2 client: %s", err)
	}
	resp, err := l7rules.GetRule(c, state.Primary.Attributes["l7policy_id"], state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the l7rule (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2L7Rule_basic(t *testing.T) {
	var l7rule l7rules.Rule
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_l7rule.l7rule_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&l7rule,
		getL7RuleResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLBV2L7RuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(resourceName, "value", "/api"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_lb_listener.listener_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "l7policy_id",
						"huaweicloud_lb_l7policy.l7policy_1", "id"),
				),
			},
			{
				Config: testAccCheckLBV2L7RuleConfig_update2(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "type", "PATH"),
					resource.TestCheckResourceAttr(resourceName, "compare_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(resourceName, "key", ""),
					resource.TestCheckResourceAttr(resourceName, "value", "/images"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLBL7RuleImportStateIdFunc(),
			},
		},
	})
}

func testAccLBL7RuleImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		policy, ok := s.RootModule().Resources["huaweicloud_lb_l7policy.l7policy_1"]
		if !ok {
			return "", fmt.Errorf("policy not found: %s", policy)
		}
		rule, ok := s.RootModule().Resources["huaweicloud_lb_l7rule.l7rule_1"]
		if !ok {
			return "", fmt.Errorf("rule not found: %s", rule)
		}
		if policy.Primary.ID == "" || rule.Primary.ID == "" {
			return "", fmt.Errorf("resource not found: %s/%s", policy.Primary.ID, rule.Primary.ID)
		}
		return fmt.Sprintf("%s/%s", policy.Primary.ID, rule.Primary.ID), nil
	}
}

func testAccCheckLBV2L7RuleConfig(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
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
