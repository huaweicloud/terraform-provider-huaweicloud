package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getL7PolicyResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud LB v2 client: %s", err)
	}
	resp, err := l7policies.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the l7policy (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2L7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy
	rName := acceptance.RandomAccResourceNameWithDash()
	rUpdateName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_l7policy.l7policy_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&l7Policy,
		getL7PolicyResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLBV2L7PolicyConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(resourceName, "position", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_lb_listener.listener_1", "id"),
				),
			},
			{
				Config: testAccCheckLBV2L7PolicyConfig_update(rName, rUpdateName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rUpdateName),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(resourceName, "position", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_lb_listener.listener_1", "id"),
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

func testAccCheckLBV2L7PolicyConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_l7policy" "l7policy_1" {
  name             = "%[1]s"
  action           = "REDIRECT_TO_POOL"
  description      = "test description"
  position         = 1
  listener_id      = huaweicloud_lb_listener.listener_1.id
  redirect_pool_id = huaweicloud_lb_pool.pool_1.id
}
`, rName)
}

func testAccCheckLBV2L7PolicyConfig_update(rName, rUpdateName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%[1]s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name            = "%[1]s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_2" {
  name            = "%[1]s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_l7policy" "l7policy_1" {
  name             = "%[2]s"
  action           = "REDIRECT_TO_POOL"
  description      = "test description update"
  position         = 1
  listener_id      = huaweicloud_lb_listener.listener_1.id
  redirect_pool_id = huaweicloud_lb_pool.pool_2.id
}
`, rName, rUpdateName)
}
