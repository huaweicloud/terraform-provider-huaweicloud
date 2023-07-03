package ecs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/servergroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccComputeServerGroup_basic(t *testing.T) {
	var sg servergroups.ServerGroup
	var instance cloudservers.CloudServer

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					testAccCheckComputeServerGroupExists(resourceName, &sg),
				),
			},
			{
				Config: testAccComputeServerGroup_members(rName, 0),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("huaweicloud_compute_instance.test.0", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
				),
			},
			{
				Config: testAccComputeServerGroup_members(rName, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("huaweicloud_compute_instance.test.1", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
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
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_compute_servergroup.sg_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_scheduler(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupExists(resourceName, &sg),
					testAccCheckComputeInstanceExists("huaweicloud_compute_instance.test", &instance),
					testAccCheckComputeInstanceInServerGroup(&instance, &sg),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeServerGroup_concurrency(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckComputeServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroup_concurrency(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckOutput("members_attached", "true"),
					resource.TestCheckOutput("volumes_attached", "true"),
				),
			},
		},
	})
}

func testAccCheckComputeServerGroupDestroy(s *terraform.State) error {
	cfg := acceptance.TestAccProvider.Meta().(*config.Config)
	ecsClient, err := cfg.ComputeV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_servergroup" {
			continue
		}

		_, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("server group still exists")
		}
	}

	return nil
}

func testAccCheckComputeServerGroupExists(n string, kp *servergroups.ServerGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		cfg := acceptance.TestAccProvider.Meta().(*config.Config)
		ecsClient, err := cfg.ComputeV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating compute client: %s", err)
		}

		found, err := servergroups.Get(ecsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("server group not found")
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

		return fmt.Errorf("instance %s does not belong to server group %s", instance.ID, sg.ID)
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

func testAccComputeServerGroup_members(rName string, idx int) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = "%[2]s-${count.index}"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_compute_servergroup" "sg_1" {
  name     = "%[2]s"
  policies = ["anti-affinity"]
  members  = [huaweicloud_compute_instance.test.%d.id]
}
`, testAccCompute_data, rName, idx)
}

func testAccComputeServerGroup_scheduler(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_servergroup" "sg_1" {
  name     = "%[2]s"
  policies = ["anti-affinity"]
}

resource "huaweicloud_compute_instance" "test" {
  name               = "%[2]s"
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
`, testAccCompute_data, rName)
}

func testAccComputeServerGroup_concurrency(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_images_images" "test" {
  flavor_id  = data.huaweicloud_compute_flavors.test.ids[0]
  os         = "Ubuntu"
  visibility = "public"
}

resource "huaweicloud_vpc" "test" {
  name = "%[1]s"
  cidr = "192.168.192.0/20"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 1), 1)
}

resource "huaweicloud_networking_secgroup" "test" {
  name = "%[1]s"
}

resource "huaweicloud_kps_keypair" "test" {
  name = "%[1]s"
}

resource "huaweicloud_compute_instance" "test" {
  count = 2

  name               = "%[1]s-${count.index}"
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  image_id           = data.huaweicloud_images_images.test.images[0].id
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]
  key_pair           = huaweicloud_kps_keypair.test.name

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_compute_servergroup" "test" {
  count = 2

  name     = "%[1]s-${count.index}"
  policies = ["anti-affinity"]

  members = [
    huaweicloud_compute_instance.test[count.index].id,
  ]

  # make sure the resource can be applied with "huaweicloud_compute_volume_attach" at the same time
  depends_on = [huaweicloud_evs_volume.test]
}

resource "huaweicloud_evs_volume" "test" {
  count = 4

  name              = "%[1]s-${count.index}"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  device_type = "SCSI"
  volume_type = "SAS"
  size        = 40
  multiattach = true
}

resource "huaweicloud_compute_volume_attach" "attach_volumes_to_compute_test_1" {
  count = 4

  instance_id = huaweicloud_compute_instance.test[0].id
  volume_id   = huaweicloud_evs_volume.test[count.index].id
}

resource "huaweicloud_compute_volume_attach" "attach_volumes_to_compute_test_2" {
  count = 4

  instance_id = huaweicloud_compute_instance.test[1].id
  volume_id   = huaweicloud_evs_volume.test[count.index].id
}

locals {
  attach_members_1 = huaweicloud_compute_servergroup.test[0].members
  attach_members_2 = huaweicloud_compute_servergroup.test[1].members

  attach_devices_1 = [for d in huaweicloud_compute_volume_attach.attach_volumes_to_compute_test_1[*].device : d != ""]
  attach_devices_2 = [for d in huaweicloud_compute_volume_attach.attach_volumes_to_compute_test_2[*].device : d != ""]
}

output "members_attached" {
  value = length(local.attach_members_1) == 1 && length(local.attach_members_2) == 1
}

output "volumes_attached" {
  value = length(local.attach_devices_1) == 4 && length(local.attach_devices_2) == 4
}
`, name)
}
