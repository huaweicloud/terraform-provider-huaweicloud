package lb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func getPoolResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud LB v2 client: %s", err)
	}
	resp, err := pools.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the pool (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_pool.pool_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.type", "APP_COOKIE"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.cookie_name", "testCookie"),
				),
			},
			{
				Config: testAccLBV2PoolConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttr(resourceName, "description", "test description update"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
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

func testAccLBV2PoolConfig_base(rName string) string {
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
`, rName)
}

func testAccLBV2PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "%[2]s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
  description = "test description"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, testAccLBV2PoolConfig_base(rName), rName)
}

func testAccLBV2PoolConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_lb_pool" "pool_1" {
  name              = "%[2]s"
  protocol          = "HTTP"
  lb_method         = "LEAST_CONNECTIONS"
  listener_id       = huaweicloud_lb_listener.listener_1.id
  description       = "test description update"
  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, testAccLBV2PoolConfig_base(rName), rNameUpdate)
}
