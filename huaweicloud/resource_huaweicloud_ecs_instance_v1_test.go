package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
)

func TestAccEcsV1Instance_basic(t *testing.T) {
	var instance cloudservers.CloudServer

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1Instance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists("huaweicloud_ecs_instance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "availability_zone", OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "auto_recovery", "true"),
				),
			},
			{
				Config: testAccEcsV1Instance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists("huaweicloud_ecs_instance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "availability_zone", OS_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "auto_recovery", "false"),
				),
			},
		},
	})
}

func testAccCheckEcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_ecs_instance_v1" {
			continue
		}

		server, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckEcsV1InstanceExists(n string, instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		found, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

var testAccEcsV1Instance_basic = fmt.Sprintf(`
resource "huaweicloud_ecs_instance_v1" "instance_1" {
  name     = "server_1"
  image_id = "%s"
  flavor   = "%s"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  password          = "Password@123"
  security_groups   = ["default"]
  availability_zone = "%s"
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, OS_IMAGE_ID, OS_FLAVOR_NAME, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)

var testAccEcsV1Instance_update = fmt.Sprintf(`
resource "huaweicloud_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_ecs"
  description = "a security group"
}

resource "huaweicloud_ecs_instance_v1" "instance_1" {
  name     = "server_updated"
  image_id = "%s"
  flavor   = "%s"
  vpc_id   = "%s"

  nics {
    network_id = "%s"
  }

  password                    = "Password@123"
  security_groups             = ["default", "${huaweicloud_compute_secgroup_v2.secgroup_1.name}"]
  availability_zone           = "%s"
  auto_recovery               = false
  delete_disks_on_termination = true

  tags = {
    foo = "bar1"
    key1 = "value"
  }
}
`, OS_IMAGE_ID, OS_FLAVOR_NAME, OS_VPC_ID, OS_NETWORK_ID, OS_AVAILABILITY_ZONE)
