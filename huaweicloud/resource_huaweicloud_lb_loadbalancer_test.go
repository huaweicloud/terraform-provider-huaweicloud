package huaweicloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/elb/v2/loadbalancers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func TestAccLBV2LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform"),
					resource.TestMatchResourceAttr(resourceName, "vip_port_id",
						regexp.MustCompile("^[a-f0-9-]+")),
				),
			},
			{
				Config: testAccLBV2LoadBalancerConfig_update(rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
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
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameSecg1 := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameSecg2 := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancer_secGroup(rName, rNameSecg1, rNameSecg2),
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
				Config: testAccLBV2LoadBalancer_secGroup_update1(rName, rNameSecg1, rNameSecg2),
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
				Config: testAccLBV2LoadBalancer_secGroup_update2(rName, rNameSecg1, rNameSecg2),
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

func TestAccLBV2LoadBalancer_withEpsId(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckEpsID(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2LoadBalancerExists(resourceName, &lb),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func testAccCheckLBV2LoadBalancerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	elbClient, err := config.LoadBalancerClient(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_lb_loadbalancer" {
			continue
		}

		_, err := loadbalancers.Get(elbClient, rs.Primary.ID).Extract()
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
		elbClient, err := config.LoadBalancerClient(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := loadbalancers.Get(elbClient, rs.Primary.ID).Extract()
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
		networkingClient, err := config.NetworkingV2Client(HW_REGION_NAME)
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

func testAccLBV2LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.subnet_id

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
`, rName)
}

func testAccLBV2LoadBalancerConfig_update(rNameUpdate string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name           = "%s"
  admin_state_up = "true"
  vip_subnet_id  = data.huaweicloud_vpc_subnet.test.subnet_id

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
`, rNameUpdate)
}

func testAccLBV2LoadBalancer_secGroup(rName, rNameSecg1, rNameSecg2 string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "secgroup_1"
}

resource "huaweicloud_networking_secgroup" "secgroup_2" {
  name        = "%s"
  description = "secgroup_2"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name               = "%s"
  vip_subnet_id      = data.huaweicloud_vpc_subnet.test.subnet_id
  security_group_ids = [
    huaweicloud_networking_secgroup.secgroup_1.id
  ]
}
`, rNameSecg1, rNameSecg2, rName)
}

func testAccLBV2LoadBalancer_secGroup_update1(rName, rNameSecg1, rNameSecg2 string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "secgroup_1"
}

resource "huaweicloud_networking_secgroup" "secgroup_2" {
  name        = "%s"
  description = "secgroup_2"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name               = "%s"
  vip_subnet_id      = data.huaweicloud_vpc_subnet.test.subnet_id
  security_group_ids = [
    huaweicloud_networking_secgroup.secgroup_1.id,
    huaweicloud_networking_secgroup.secgroup_2.id
  ]
}
`, rNameSecg1, rNameSecg2, rName)
}

func testAccLBV2LoadBalancer_secGroup_update2(rName, rNameSecg1, rNameSecg2 string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name        = "%s"
  description = "secgroup_1"
}

resource "huaweicloud_networking_secgroup" "secgroup_2" {
  name        = "%s"
  description = "secgroup_2"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name               = "%s"
  vip_subnet_id      = data.huaweicloud_vpc_subnet.test.subnet_id
  security_group_ids = [
    huaweicloud_networking_secgroup.secgroup_2.id
  ]
}
`, rNameSecg1, rNameSecg2, rName)
}

func testAccLBV2LoadBalancerConfig_withEpsId(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name                  = "%s"
  vip_subnet_id         = data.huaweicloud_vpc_subnet.test.subnet_id
  enterprise_project_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, HW_ENTERPRISE_PROJECT_ID_TEST)
}
