package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func TestAccLBV2LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckULB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "name", "loadbalancer_1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestMatchResourceAttr(resourceName, "vip_port_id",
						regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
			{
				Config: testAccLBV2LoadBalancerConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "loadbalancer_1_updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
				),
			},
		},
	})
}

func TestAccLBV2LoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var sg_1, sg_2 groups.SecGroup
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckULB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancer_secGroup,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_1", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_1", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_1),
				),
			},
			{
				Config: testAccLBV2LoadBalancer_secGroup_update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "2"),
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_1),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
			{
				Config: testAccLBV2LoadBalancer_secGroup_update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_1),
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
		},
	})
}

func testAccCheckLBV2LoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lb_loadbalancer" {
			continue
		}

		_, err := loadbalancers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LoadBalancer still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2LoadBalancerExists(
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
		networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
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

func testAccCheckLBV2LoadBalancerHasSecGroup(
	lb *loadbalancers.LoadBalancer, sg *groups.SecGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		port, err := ports.Get(networkingClient, lb.VipPortID).Extract()
		if err != nil {
			return err
		}

		for _, p := range port.SecurityGroups {
			if p == sg.ID {
				return nil
			}
		}

		return fmt.Errorf("LoadBalancer does not have the security group")
	}
}

var testAccLBV2LoadBalancerConfig_basic = fmt.Sprintf(`
resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_SUBNET_ID)

var testAccLBV2LoadBalancerConfig_update = fmt.Sprintf(`
resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name = "loadbalancer_1_updated"
  admin_state_up = "true"
  vip_subnet_id = "%s"

  tags = {
    key1  = "value1"
    owner = "terraform_update"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, OS_SUBNET_ID)

var testAccLBV2LoadBalancer_secGroup = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "huaweicloud_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "%s"
    security_group_ids = [
      huaweicloud_networking_secgroup.secgroup_1.id
    ]
}
`, OS_SUBNET_ID)

var testAccLBV2LoadBalancer_secGroup_update1 = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "huaweicloud_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "%s"
    security_group_ids = [
      huaweicloud_networking_secgroup.secgroup_1.id,
      huaweicloud_networking_secgroup.secgroup_2.id
    ]
}
`, OS_SUBNET_ID)

var testAccLBV2LoadBalancer_secGroup_update2 = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "huaweicloud_networking_secgroup" "secgroup_2" {
  name = "secgroup_2"
  description = "secgroup_2"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
    name = "loadbalancer_1"
    vip_subnet_id = "%s"
    security_group_ids = [
      huaweicloud_networking_secgroup.secgroup_2.id
    ]
}
`, OS_SUBNET_ID)
