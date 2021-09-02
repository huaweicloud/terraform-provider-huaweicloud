package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/lbaas_v2/monitors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_lb_monitor.monitor_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
				),
			},
			{
				Config: testAccLBV2MonitorConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "3"),
					resource.TestCheckResourceAttr(resourceName, "port", "8888"),
				),
			},
		},
	})
}

func TestAccLBV2Monitor_udp(t *testing.T) {
	var monitor monitors.Monitor
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_lb_monitor.monitor_udp"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_udp(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "UDP_CONNECT"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
				),
			},
		},
	})
}

func TestAccLBV2Monitor_http(t *testing.T) {
	var monitor monitors.Monitor
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_lb_monitor.monitor_http"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2MonitorConfig_http(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", "HTTP"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
					resource.TestCheckResourceAttr(resourceName, "max_retries", "5"),
					resource.TestCheckResourceAttr(resourceName, "url_path", "/api"),
					resource.TestCheckResourceAttr(resourceName, "http_method", "GET"),
					resource.TestCheckResourceAttr(resourceName, "expected_codes", "200-202"),
				),
			},
		},
	})
}

func testAccCheckLBV2MonitorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV2Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lb_monitor" {
			continue
		}

		_, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Monitor still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MonitorExists(n string, monitor *monitors.Monitor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		elbClient, err := config.ElbV2Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
		}

		found, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Monitor not found")
		}

		*monitor = *found

		return nil
	}
}

func testAccLBV2MonitorConfig_base(rName string) string {
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
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}
`, rName, rName, rName)
}

func testAccLBV2MonitorConfig_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_monitor" "monitor_1" {
  pool_id     = huaweicloud_lb_pool.pool_1.id
  name        = "%s"
  type        = "TCP"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
`, testAccLBV2MonitorConfig_base(rName), rName)
}

func testAccLBV2MonitorConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_monitor" "monitor_1" {
  pool_id     = huaweicloud_lb_pool.pool_1.id
  name        = "%s"
  type        = "TCP"
  delay       = 30
  timeout     = 15
  max_retries = 3
  port        = 8888
}
`, testAccLBV2MonitorConfig_base(rName), rNameUpdate)
}

func testAccLBV2MonitorConfig_http(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_lb_monitor" "monitor_http" {
  pool_id        = huaweicloud_lb_pool.pool_1.id
  name           = "%s"
  type           = "HTTP"
  delay          = 20
  timeout        = 10
  max_retries    = 5
  url_path       = "/api"
  expected_codes = "200-202"
}
`, testAccLBV2MonitorConfig_base(rName), rName)
}

func testAccLBV2MonitorConfig_udp(rName string) string {
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
  protocol        = "UDP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "UDP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_lb_monitor" "monitor_udp" {
  pool_id     = huaweicloud_lb_pool.pool_1.id
  name        = "%s"
  type        = "UDP_CONNECT"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
`, rName, rName, rName, rName)
}
