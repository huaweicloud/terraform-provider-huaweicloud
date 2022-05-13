package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/servergroups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccComputeServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
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

func TestAccComputeServerGroup_scheduler(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_scheduler(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_members(t *testing.T) {
	var instance cloudservers.CloudServer
	var sg servergroups.ServerGroup
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_members(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func testAccCheckComputeServerGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	ecsClient, err := config.ComputeV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_servergroup" {
			continue
		}

		_, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ServerGroup still exists")
		}
	}

	return nil
}

func testAccCheckComputeServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		ecsClient, err := config.ComputeV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		found, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ServerGroup not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckComputeInstanceInServerGroup(instance *cloudservers.CloudServer, sg *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(sg.Members) > 0 {
			for _, m := range sg.Members {
				if m == instance.ID {
					return nil
				}
			}
		}

		return fmt.Errorf("Instance %s is not part of Server Group %s", instance.ID, sg.ID)
	}
}

func testAccComputeServerGroup_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}
`, rName)
}

func testAccComputeServerGroup_scheduler(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
}

resource "huaweicloud_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  scheduler_hints {
    group = huaweicloud_compute_servergroup.sg_1.id
  }
  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeServerGroup_members(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_servergroup" "sg_1" {
  name     = "%s"
  policies = ["anti-affinity"]
  members  = [huaweicloud_compute_instance.instance_1.id]
}

resource "huaweicloud_compute_instance" "instance_1" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName, rName)
}
