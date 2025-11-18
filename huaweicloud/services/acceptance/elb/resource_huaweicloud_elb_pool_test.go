package elb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getELBPoolResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, err
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_id}", state.Primary.ID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(getResp)
}

func TestAccElbV3Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "description", "test pool description"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.type", "APP_COOKIE"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.cookie_name", "testCookie"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_member_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "public_border_group", "center"),
					resource.TestCheckResourceAttr(resourceName, "az_affinity.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_percentage", "20"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_count", "-1"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_unhealthy_fallback_strategy", "forward_to_all_healthy_member"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "description", "test pool description update"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "consoleProtection"),
					resource.TestCheckResourceAttr(resourceName, "protection_reason", "test protection reason"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.type", "APP_COOKIE"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.cookie_name", "testCookie"),
					resource.TestCheckResourceAttr(resourceName, "minimum_healthy_member_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "az_affinity.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_percentage", "-1"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_count", "200"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_unhealthy_fallback_strategy", "forward_to_healthy_member_of_remote_az"),
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

func TestAccElbV3Pool_basic_with_loadbalancer(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_loadbalancer(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id",
						"huaweicloud_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_loadbalancer(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id",
						"huaweicloud_elb_loadbalancer.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "type", "instance"),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id",
						"huaweicloud_vpc.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
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

func TestAccElbV3Pool_basic_with_listener(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_listener(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_listener(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttrPair(resourceName, "listener_id",
						"huaweicloud_elb_listener.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
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

func TestAccElbV3Pool_basic_with_type_ip(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_type_ip(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "protection_status", "nonProtection"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_type_ip(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "slow_start_duration", "100"),
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

func TestAccElbV3Pool_basic_with_protocol_tcp(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_protocol_tcp(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "dualstack"),
					resource.TestCheckResourceAttr(resourceName, "any_port_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "az_affinity.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_percentage", "20"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_count", "-1"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_unhealthy_fallback_strategy", "forward_to_all_healthy_member"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_protocol_tcp(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
					resource.TestCheckResourceAttr(resourceName, "type", "ip"),
					resource.TestCheckResourceAttr(resourceName, "ip_version", "dualstack"),
					resource.TestCheckResourceAttr(resourceName, "any_port_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "deletion_protection_enable", "true"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.type", "SOURCE_IP"),
					resource.TestCheckResourceAttr(resourceName, "persistence.0.timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "az_affinity.0.enable", "true"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_percentage", "-1"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_minimum_healthy_member_count", "200"),
					resource.TestCheckResourceAttr(resourceName,
						"az_affinity.0.az_unhealthy_fallback_strategy", "forward_to_healthy_member_of_remote_az"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "updated_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"deletion_protection_enable"},
			},
		},
	})
}

func TestAccElbV3Pool_basic_with_connection_drain(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_connection_drain(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "connection_drain_enabled", "true"),
					resource.TestCheckResourceAttr(resourceName, "connection_drain_timeout", "80"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update_with_connection_drain_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "connection_drain_enabled", "false"),
					resource.TestCheckResourceAttr(resourceName, "connection_drain_timeout", "60"),
				),
			},
		},
	})
}

func TestAccElbV3Pool_basic_with_ip_protocol(t *testing.T) {
	var pool pools.Pool
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_elb_pool.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&pool,
		getELBPoolResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckElbGatewayType(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic_with_ip_protocol(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "protocol", "IP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "2_TUPLE_HASH"),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id",
						"huaweicloud_elb_loadbalancer.test", "id"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_basic_with_ip_protocol_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "protocol", "IP"),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "5_TUPLE_HASH"),
					resource.TestCheckResourceAttrPair(resourceName, "loadbalancer_id",
						"huaweicloud_elb_loadbalancer.test", "id"),
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

func testAccElbV3PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  type        = "instance"
  vpc_id      = huaweicloud_vpc.test.id
  description = "test pool description"

  public_border_group          = "center"
  minimum_healthy_member_count = 1

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3PoolConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "LEAST_CONNECTIONS"
  type        = "instance"
  vpc_id      = huaweicloud_vpc.test.id
  description = "test pool description update"

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"

  public_border_group          = "center"
  minimum_healthy_member_count = 0

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
`, common.TestVpc(rName), rNameUpdate)
}

func testAccElbV3PoolConfig_basic_with_loadbalancer(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_pool" "test" {
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}
`, testAccElbV3LoadBalancerConfig_basic(rName), rName)
}

func testAccElbV3PoolConfig_update_with_loadbalancer(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_pool" "test" {
  name            = "%s"
  protocol        = "HTTP"
  lb_method       = "LEAST_CONNECTIONS"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
  type            = "instance"
  vpc_id          = huaweicloud_vpc.test.id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"
}
`, testAccElbV3LoadBalancerConfig_basic(rName), rNameUpdate)
}

func testAccElbV3PoolConfig_basic_with_listener_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%[1]s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]

  tags = {
    key   = "value"
    owner = "terraform"
  }
}

resource "huaweicloud_elb_listener" "test" {
  name            = "%[1]s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id

  idle_timeout     = 62
  request_timeout  = 63
  response_timeout = 64

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName)
}

func testAccElbV3PoolConfig_basic_with_listener(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_elb_listener.test.id
}
`, testAccElbV3PoolConfig_basic_with_listener_base(rName), rName)
}

func testAccElbV3PoolConfig_update_with_listener(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "LEAST_CONNECTIONS"
  listener_id = huaweicloud_elb_listener.test.id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"
}
`, testAccElbV3PoolConfig_basic_with_listener_base(rName), rNameUpdate)
}

func testAccElbV3PoolConfig_basic_with_type_ip(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_pool" "test" {
  name      = "%s"
  protocol  = "HTTP"
  lb_method = "ROUND_ROBIN"
  type      = "ip"
}
`, rName)
}

func testAccElbV3PoolConfig_update_with_type_ip(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_pool" "test" {
  name      = "%s"
  protocol  = "HTTP"
  lb_method = "LEAST_CONNECTIONS"
  type      = "ip"

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test protection reason"
}
`, rNameUpdate)
}

func testAccElbV3PoolConfig_basic_with_protocol_tcp(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_pool" "test" {
  name            = "%s"
  protocol        = "TCP"
  lb_method       = "ROUND_ROBIN"
  type            = "ip"
  ip_version      = "dualstack"
  any_port_enable = true

  az_affinity {
    enable                               = true
    az_minimum_healthy_member_percentage = 20
    az_minimum_healthy_member_count      = -1
    az_unhealthy_fallback_strategy       = "forward_to_all_healthy_member"
  }
}
`, rName)
}

func testAccElbV3PoolConfig_update_with_protocol_tcp(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_elb_pool" "test" {
  name                       = "%s"
  protocol                   = "TCP"
  lb_method                  = "ROUND_ROBIN"
  type                       = "ip"
  ip_version                 = "dualstack"
  any_port_enable            = true
  deletion_protection_enable = true

  persistence {
    type    = "SOURCE_IP"
    timeout = 10
  }

  az_affinity {
    enable                               = true
    az_minimum_healthy_member_percentage = -1
    az_minimum_healthy_member_count      = 200
    az_unhealthy_fallback_strategy       = "forward_to_healthy_member_of_remote_az"
  }
}
`, rName)
}

func testAccElbV3PoolConfig_basic_with_connection_drain(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  name                     = "%[2]s"
  protocol                 = "TCP"
  lb_method                = "ROUND_ROBIN"
  type                     = "instance"
  vpc_id                   = huaweicloud_vpc.test.id
  connection_drain_enabled = true
  connection_drain_timeout = 80
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3PoolConfig_update_with_connection_drain_update(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  name                     = "%[2]s"
  protocol                 = "TCP"
  lb_method                = "ROUND_ROBIN"
  type                     = "instance"
  vpc_id                   = huaweicloud_vpc.test.id
  connection_drain_enabled = false
  connection_drain_timeout = 60
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3PoolConfig_basic_with_ip_protocol_base(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  loadbalancer_type = "gateway"
  description       = "test gateway description"

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}
`, common.TestVpc(rName), rName)
}

func testAccElbV3PoolConfig_basic_with_ip_protocol(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_elb_pool" "test" {
  name            = "%[2]s"
  protocol        = "IP"
  lb_method       = "2_TUPLE_HASH"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}
`, testAccElbV3PoolConfig_basic_with_ip_protocol_base(rName), rName)
}

func testAccElbV3PoolConfig_basic_with_ip_protocol_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_elb_pool" "test" {
  name            = "%[2]s"
  protocol        = "IP"
  lb_method       = "5_TUPLE_HASH"
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id
}
`, testAccElbV3PoolConfig_basic_with_ip_protocol_base(rName), rNameUpdate)
}
