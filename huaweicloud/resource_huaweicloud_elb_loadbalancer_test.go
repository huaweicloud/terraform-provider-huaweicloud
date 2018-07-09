package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elb/loadbalancers"
)

func TestAccELBLoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckELB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccELBLoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists("huaweicloud_elb_loadbalancer.loadbalancer_1", &lb),
				),
			},
			resource.TestStep{
				Config: testAccELBLoadBalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_elb_loadbalancer.loadbalancer_1", "name", "loadbalancer_1_updated"),
					resource.TestCheckResourceAttr(
						"huaweicloud_elb_loadbalancer.loadbalancer_1", "admin_state_up", "0"),
				),
			},
		},
	})
}

func TestAccELBLoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancers.LoadBalancer

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckELB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckELBLoadBalancerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccELBLoadBalancer_internal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists(
						"huaweicloud_elb_loadbalancer.loadbalancer_1", &lb),
					resource.TestCheckResourceAttr(
						"huaweicloud_elb_loadbalancer.loadbalancer_1", "admin_state_up", "1"),
				),
			},
			resource.TestStep{
				Config: testAccELBLoadBalancer_internal_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckELBLoadBalancerExists(
						"huaweicloud_elb_loadbalancer.loadbalancer_1", &lb),
					resource.TestCheckResourceAttr(
						"huaweicloud_elb_loadbalancer.loadbalancer_1", "admin_state_up", "0"),
				),
			},
		},
	})
}

func testAccCheckELBLoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.loadElasticLoadBalancerClient(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_elb_loadbalancer" {
			continue
		}

		_, err := loadbalancers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckELBLoadBalancerExists(
	n string, lb *loadbalancers.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.loadElasticLoadBalancerClient(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}
		found, err := loadbalancers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*lb = *found

		return nil
	}
}

var testAccELBLoadBalancerConfig_basic = fmt.Sprintf(`
resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vpc_id = "%s"
  type = "External"
  bandwidth = "5"
  admin_state_up = 1

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)

var testAccELBLoadBalancerConfig_update = fmt.Sprintf(`
resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1_updated"
  vpc_id = "%s"
  type = "External"
  bandwidth = 3
  admin_state_up = 0

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_VPC_ID)

var testAccELBLoadBalancer_internal = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "huaweicloud_networking_router_interface_v2" "interface" {
  router_id = "%s"
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  type = "Internal"
  admin_state_up = 1
  vpc_id = "%s"
  az = "%s"
  tenantid = "%s"
  vip_subnet_id = "${huaweicloud_networking_network_v2.network_1.id}"
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
  depends_on = ["huaweicloud_networking_router_interface_v2.interface"]
}
`, OS_VPC_ID, OS_VPC_ID, OS_AVAILABILITY_ZONE, OS_TENANT_ID)

var testAccELBLoadBalancer_internal_update = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "huaweicloud_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "huaweicloud_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${huaweicloud_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
}

resource "huaweicloud_networking_router_interface_v2" "interface" {
  router_id = "%s"
  subnet_id = "${huaweicloud_networking_subnet_v2.subnet_1.id}"
}

resource "huaweicloud_elb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  type = "Internal"
  admin_state_up = 0
  vpc_id = "%s"
  az = "%s"
  tenantid = "%s"
  vip_subnet_id = "${huaweicloud_networking_network_v2.network_1.id}"
  security_group_id = "${huaweicloud_networking_secgroup_v2.secgroup_1.id}"
  depends_on = ["huaweicloud_networking_router_interface_v2.interface"]
}
`, OS_VPC_ID, OS_VPC_ID, OS_AVAILABILITY_ZONE, OS_TENANT_ID)
