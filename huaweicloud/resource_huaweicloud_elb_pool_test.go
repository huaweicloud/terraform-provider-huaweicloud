package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/elb/v3/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccElbV3Pool_basic(t *testing.T) {
	var pool pools.Pool
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	rNameUpdate := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_elb_pool.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckElbV3PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccElbV3PoolConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckElbV3PoolExists(resourceName, &pool),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "ROUND_ROBIN"),
				),
			},
			{
				Config: testAccElbV3PoolConfig_update(rName, rNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rNameUpdate),
					resource.TestCheckResourceAttr(resourceName, "lb_method", "LEAST_CONNECTIONS"),
				),
			},
		},
	})
}

func testAccCheckElbV3PoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	elbClient, err := config.ElbV3Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_elb_pool" {
			continue
		}

		_, err := pools.Get(elbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmtp.Errorf("Pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckElbV3PoolExists(n string, pool *pools.Pool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		elbClient, err := config.ElbV3Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
		}

		found, err := pools.Get(elbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Member not found")
		}

		*pool = *found

		return nil
	}
}

func testAccElbV3PoolConfig_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.subnet_id
  ipv6_network_id = data.huaweicloud_vpc_subnet.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id

  forward_eip = true

  idle_timeout = 60
  request_timeout = 60
  response_timeout = 60
}

resource "huaweicloud_elb_pool" "test" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_elb_listener.test.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rName)
}

func testAccElbV3PoolConfig_update(rName, rNameUpdate string) string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_elb_loadbalancer" "test" {
  name            = "%s"
  ipv4_subnet_id  = data.huaweicloud_vpc_subnet.test.subnet_id
  ipv6_network_id = data.huaweicloud_vpc_subnet.test.id

  availability_zone = [
    data.huaweicloud_availability_zones.test.names[0]
  ]
}

resource "huaweicloud_elb_listener" "test" {
  name            = "%s"
  description     = "test description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_elb_loadbalancer.test.id

  forward_eip = true

  idle_timeout = 60
  request_timeout = 60
  response_timeout = 60
}

resource "huaweicloud_elb_pool" "test" {
  name           = "%s"
  protocol       = "HTTP"
  lb_method      = "LEAST_CONNECTIONS"
  listener_id    = huaweicloud_elb_listener.test.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`, rName, rName, rNameUpdate)
}
