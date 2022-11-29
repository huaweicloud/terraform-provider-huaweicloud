package lb

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/elb/v2/loadbalancers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v3/security/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getLoadBalancerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.LoadBalancerClient(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud LB v2 client: %s", err)
	}
	resp, err := loadbalancers.Get(c, state.Primary.ID).Extract()
	if resp == nil && err == nil {
		return resp, fmt.Errorf("Unable to find the loadbalancer (%s)", state.Primary.ID)
	}
	return resp, err
}

func TestAccLBV2LoadBalancer_basic(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameUpdate := rName + "-update"
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getLoadBalancerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "created by acceptance test"),
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
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "terraform_update"),
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

func TestAccLBV2LoadBalancer_secGroup(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	var sg_1, sg_2 groups.SecurityGroup
	rName := acceptance.RandomAccResourceNameWithDash()
	rNameSecg1 := acceptance.RandomAccResourceNameWithDash()
	rNameSecg2 := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getLoadBalancerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancer_secGroup(rName, rNameSecg1, rNameSecg2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					testAccCheckNetworkingV3SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_1", &sg_1),
					testAccCheckNetworkingV3SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_1", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_1),
				),
			},
			{
				Config: testAccLBV2LoadBalancer_secGroup_update1(rName, rNameSecg1, rNameSecg2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "2"),
					testAccCheckNetworkingV3SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_1),
					testAccCheckNetworkingV3SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_1),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_2),
				),
			},
			{
				Config: testAccLBV2LoadBalancer_secGroup_update2(rName, rNameSecg1, rNameSecg2),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "security_group_ids.#", "1"),
					testAccCheckNetworkingV3SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_1),
					testAccCheckNetworkingV3SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &sg_2),
					testAccCheckLBV2LoadBalancerHasSecGroup(&lb, &sg_2),
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

func TestAccLBV2LoadBalancer_withEpsId(t *testing.T) {
	var lb loadbalancers.LoadBalancer
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_lb_loadbalancer.loadbalancer_1"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&lb,
		getLoadBalancerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
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

func testAccCheckLBV2LoadBalancerHasSecGroup(
	lb *loadbalancers.LoadBalancer, sg *groups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := acceptance.TestAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
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

		return fmtp.Errorf("LoadBalancer does not have the security group")
	}
}

func testAccCheckNetworkingV3SecGroupExists(n string, secGroup *groups.SecurityGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV3Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}

		found, err := groups.Get(networkingClient, rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Security group not found")
		}

		*secGroup = *found

		return nil
	}
}

func testAccLBV2LoadBalancerConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  description   = "created by acceptance test"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  tags = {
    key   = "value"
    owner = "terraform"
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
  vip_subnet_id  = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id

  tags = {
    key1  = "value1"
    owner = "terraform_update"
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
  vip_subnet_id      = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
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
  vip_subnet_id      = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
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
  vip_subnet_id      = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
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
  vip_subnet_id         = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
  enterprise_project_id = "%s"

  tags = {
    key   = "value"
    owner = "terraform"
  }
}
`, rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
