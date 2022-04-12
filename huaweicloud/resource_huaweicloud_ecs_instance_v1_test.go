package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccEcsV1Instance_basic(t *testing.T) {
	var instance cloudservers.CloudServer
	rName := fmt.Sprintf("tf-acc-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDeprecated(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEcsV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEcsV1Instance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists("huaweicloud_ecs_instance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "availability_zone", HW_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "auto_recovery", "true"),
				),
			},
			{
				ResourceName:      "huaweicloud_ecs_instance_v1.instance_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
			{
				Config: testAccEcsV1Instance_update(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEcsV1InstanceExists("huaweicloud_ecs_instance_v1.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "availability_zone", HW_AVAILABILITY_ZONE),
					resource.TestCheckResourceAttr(
						"huaweicloud_ecs_instance_v1.instance_1", "auto_recovery", "false"),
				),
			},
		},
	})
}

func testAccCheckEcsV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	computeClient, err := config.ComputeV1Client(HW_REGION_NAME)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_ecs_instance_v1" {
			continue
		}

		server, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "DELETED" {
				return fmtp.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckEcsV1InstanceExists(n string, instance *cloudservers.CloudServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		computeClient, err := config.ComputeV1Client(HW_REGION_NAME)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		found, err := cloudservers.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmtp.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccEcsV1Instance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_ecs_instance_v1" "instance_1" {
  name     = "%s"
  image_id = data.huaweicloud_images_image.test.id
  flavor   = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id   = data.huaweicloud_vpc_subnet.test.vpc_id

  nics {
    network_id = data.huaweicloud_vpc_subnet.test.id
  }

  password          = "Password@123"
  security_groups   = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  auto_recovery     = true

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName)
}

func testAccEcsV1Instance_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_secgroup_v2" "secgroup_1" {
  name        = "secgroup_ecs"
  description = "a security group"
}

resource "huaweicloud_ecs_instance_v1" "instance_1" {
  name     = "%s_updated"
  image_id = data.huaweicloud_images_image.test.id
  flavor   = data.huaweicloud_compute_flavors.test.ids[0]
  vpc_id   = data.huaweicloud_vpc_subnet.test.vpc_id

  nics {
    network_id = data.huaweicloud_vpc_subnet.test.id
  }

  password                    = "Password@123"
  security_groups             = [data.huaweicloud_networking_secgroup.test.id, huaweicloud_compute_secgroup_v2.secgroup_1.id]
  availability_zone           = data.huaweicloud_availability_zones.test.names[0]
  auto_recovery               = false
  delete_disks_on_termination = true

  tags = {
    foo  = "bar1"
    key1 = "value"
  }
}
`, testAccCompute_data, rName)
}
