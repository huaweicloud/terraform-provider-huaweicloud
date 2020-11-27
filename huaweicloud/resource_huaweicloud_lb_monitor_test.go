package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/monitors"
)

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor
	resourceName := "huaweicloud_lb_monitor.monitor_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckULB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2MonitorConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(resourceName, &monitor),
					resource.TestCheckResourceAttr(resourceName, "name", "monitor_1"),
					resource.TestCheckResourceAttr(resourceName, "delay", "20"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "10"),
				),
			},
			{
				Config: TestAccLBV2MonitorConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "monitor_1_updated"),
					resource.TestCheckResourceAttr(resourceName, "delay", "30"),
					resource.TestCheckResourceAttr(resourceName, "timeout", "15"),
				),
			},
		},
	})
}

func testAccCheckLBV2MonitorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	elbClient, err := config.elbV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lb_monitor" {
			continue
		}

		_, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Monitor still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MonitorExists(n string, monitor *monitors.Monitor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		elbClient, err := config.elbV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
		}

		found, err := monitors.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Monitor not found")
		}

		*monitor = *found

		return nil
	}
}

var TestAccLBV2MonitorConfig_basic = fmt.Sprintf(`
resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_lb_monitor" "monitor_1" {
  name        = "monitor_1"
  type        = "PING"
  delay       = 20
  timeout     = 10
  max_retries = 5
  pool_id     = huaweicloud_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, HW_SUBNET_ID)

var TestAccLBV2MonitorConfig_update = fmt.Sprintf(`
resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "loadbalancer_1"
  vip_subnet_id = "%s"
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_lb_monitor" "monitor_1" {
  name           = "monitor_1_updated"
  type           = "PING"
  delay          = 30
  timeout        = 15
  max_retries    = 10
  admin_state_up = "true"
  pool_id        = huaweicloud_lb_pool.pool_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, HW_SUBNET_ID)
