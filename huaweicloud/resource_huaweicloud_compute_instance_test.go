package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/volumeattach"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/security/groups"
	"github.com/huaweicloud/golangsdk/pagination"
)

func TestAccComputeV2Instance_basic(t *testing.T) {
	var instance servers.Server

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"force_delete",
				},
			},
		},
	})
}

func TestAccComputeV2Instance_disks(t *testing.T) {
	var instance servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_disks,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					resource.TestCheckResourceAttr(
						"huaweicloud_compute_instance.instance_1", "availability_zone", HW_AVAILABILITY_ZONE),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_tags(t *testing.T) {
	var instance servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_tags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeV2InstanceTags(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTags(&instance, "key", "value"),
				),
			},
			{
				Config: testAccComputeV2Instance_tags2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeV2InstanceTags(&instance, "foo2", "bar2"),
					testAccCheckComputeV2InstanceTags(&instance, "key", "value2"),
				),
			},
			{
				Config: testAccComputeV2Instance_notags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeV2InstanceNoTags(&instance),
				),
			},
			{
				Config: testAccComputeV2Instance_tags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeV2InstanceTags(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTags(&instance, "key", "value"),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_secgroupMulti(t *testing.T) {
	rName := acctest.RandString(5)
	var instance_1 servers.Server
	var secgroup_1, secgroup_2 groups.SecGroup

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_secgroupMulti(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_1", &secgroup_1),
					testAccCheckComputeV2InstanceExists(
						"huaweicloud_compute_instance.instance_1", &instance_1),
				),
			},
			{
				Config: testAccComputeV2Instance_secgroupMultiUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"huaweicloud_networking_secgroup.secgroup_2", &secgroup_2),
					testAccCheckComputeV2InstanceExists(
						"huaweicloud_compute_instance.instance_1", &instance_1),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeImage(t *testing.T) {
	var instance servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeImage,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
			{
				ResourceName:      "huaweicloud_compute_instance.instance_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"force_delete",
				},
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeVolume(t *testing.T) {
	var instance servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeVolume,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
					testAccCheckComputeV2InstanceBootVolumeAttachment(&instance),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_bootFromVolumeForceNew(t *testing.T) {
	var instance1_1 servers.Server
	var instance1_2 servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_bootFromVolumeForceNew_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"huaweicloud_compute_instance.instance_1", &instance1_1),
				),
			},
			{
				Config: testAccComputeV2Instance_bootFromVolumeForceNew_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"huaweicloud_compute_instance.instance_1", &instance1_2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1_1, &instance1_2),
				),
			},
			{
				ResourceName:      "huaweicloud_compute_instance.instance_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"stop_before_destroy",
					"force_delete",
				},
			},
		},
	})
}

func TestAccComputeV2Instance_changeFixedIP(t *testing.T) {
	var instance1_1 servers.Server
	var instance1_2 servers.Server

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_changeFixedIP_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"huaweicloud_compute_instance.instance_1", &instance1_1),
				),
			},
			{
				Config: testAccComputeV2Instance_changeFixedIP_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(
						"huaweicloud_compute_instance.instance_1", &instance1_2),
					testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(&instance1_1, &instance1_2),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_stopBeforeDestroy(t *testing.T) {
	var instance servers.Server
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_stopBeforeDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("huaweicloud_compute_instance.instance_1", &instance),
				),
			},
		},
	})
}

func testAccCheckComputeV2InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_compute_instance" {
			continue
		}

		server, err := servers.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			if server.Status != "SOFT_DELETED" {
				return fmt.Errorf("Instance still exists")
			}
		}
	}

	return nil
}

func testAccCheckComputeV2InstanceExists(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		found, err := servers.Get(computeClient, rs.Primary.ID).Extract()
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

func testAccCheckComputeV2InstanceDoesNotExist(n string, instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		_, err = servers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil
			}
			return err
		}

		return fmt.Errorf("Instance still exists")
	}
}

func testAccCheckComputeV2InstanceTags(
	instance *servers.Server, k, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.computeV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
		}

		taglist, err := tags.Get(client, "cloudservers", instance.ID).Extract()
		for _, val := range taglist.Tags {
			if k != val.Key {
				continue
			}

			if v == val.Value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, val.Value)
		}

		return fmt.Errorf("Tag not found: %s", k)
	}
}

func testAccCheckComputeV2InstanceNoTags(
	instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.computeV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud compute v1 client: %s", err)
		}

		taglist, err := tags.Get(client, "cloudservers", instance.ID).Extract()

		if taglist.Tags == nil {
			return nil
		}
		if len(taglist.Tags) == 0 {
			return nil
		}

		return fmt.Errorf("Expected no tags, but found %v", taglist.Tags)
	}
}

func testAccCheckComputeV2InstanceBootVolumeAttachment(
	instance *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var attachments []volumeattach.VolumeAttachment

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(HW_REGION_NAME)
		if err != nil {
			return err
		}

		err = volumeattach.List(computeClient, instance.ID).EachPage(
			func(page pagination.Page) (bool, error) {

				actual, err := volumeattach.ExtractVolumeAttachments(page)
				if err != nil {
					return false, fmt.Errorf("Unable to lookup attachment: %s", err)
				}

				attachments = actual
				return true, nil
			})

		if len(attachments) == 1 {
			return nil
		}

		return fmt.Errorf("No attached volume found.")
	}
}

func testAccCheckComputeV2InstanceInstanceIDsDoNotMatch(
	instance1, instance2 *servers.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if instance1.ID == instance2.ID {
			return fmt.Errorf("Instance was not recreated.")
		}

		return nil
	}
}

func testAccComputeV2Instance_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc_subnet_v1" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image_v2" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image_v2.test.id
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet_v1.test.id
  }
}
`, rName)
}

var testAccComputeV2Instance_disks = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }

  charging_mode = "prePaid"
  period_unit = "month"
  period = 1
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)

func testAccComputeV2Instance_secgroupMulti(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name                 = "secgroup1-%s"
  description          = "a security group by acceptance test"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
}

resource "huaweicloud_compute_instance" "instance_1" {
  name              = "ecs-%s"
  security_groups   = ["default", huaweicloud_networking_secgroup.secgroup_1.name]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}
`, rName, rName, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)
}

func testAccComputeV2Instance_secgroupMultiUpdate(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup_1" {
  name                 = "secgroup1-%s"
  description          = "a security group by acceptance test"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_1" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 22
  port_range_min    = 22
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_1.id
}

resource "huaweicloud_networking_secgroup" "secgroup_2" {
  name                 = "secgroup2-%s"
  description          = "another security group by acceptance test"
  delete_default_rules = true
}

resource "huaweicloud_networking_secgroup_rule" "secgroup_rule_2" {
  direction         = "ingress"
  ethertype         = "IPv4"
  port_range_max    = 80
  port_range_min    = 80
  protocol          = "tcp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = huaweicloud_networking_secgroup.secgroup_2.id
}

resource "huaweicloud_compute_instance" "instance_1" {
  name              = "ecs-%s"
  security_groups   = ["default", huaweicloud_networking_secgroup.secgroup_1.name, huaweicloud_networking_secgroup.secgroup_2.name]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}
`, rName, rName, rName, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)
}

var testAccComputeV2Instance_bootFromVolumeImage = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 40
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, HW_AVAILABILITY_ZONE, HW_IMAGE_ID, HW_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeVolume = fmt.Sprintf(`
resource "huaweicloud_blockstorage_volume_v2" "vol_1" {
  name = "vol_1"
  size = 40
  image_id = "%s"
  availability_zone = "%s"
}

resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  block_device {
    uuid = huaweicloud_blockstorage_volume_v2.vol_1.id
    source_type = "volume"
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, HW_IMAGE_ID, HW_AVAILABILITY_ZONE, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeForceNew_1 = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 40
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, HW_AVAILABILITY_ZONE, HW_IMAGE_ID, HW_NETWORK_ID)

var testAccComputeV2Instance_bootFromVolumeForceNew_2 = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  block_device {
    uuid = "%s"
    source_type = "image"
    volume_size = 41
    boot_index = 0
    destination_type = "volume"
    delete_on_termination = true
  }
  network {
    uuid = "%s"
  }
}
`, HW_AVAILABILITY_ZONE, HW_IMAGE_ID, HW_NETWORK_ID)

var testAccComputeV2Instance_changeFixedIP_1 = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
    fixed_ip_v4 = "192.168.0.24"
  }
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)

var testAccComputeV2Instance_changeFixedIP_2 = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
    fixed_ip_v4 = "192.168.0.25"
  }
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)

var testAccComputeV2Instance_stopBeforeDestroy = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  stop_before_destroy = true
  network {
    uuid = "%s"
  }
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)

var testAccComputeV2Instance_tags = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)

var testAccComputeV2Instance_tags2 = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
  tags = {
    foo2 = "bar2"
    key = "value2"
  }
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)

var testAccComputeV2Instance_notags = fmt.Sprintf(`
resource "huaweicloud_compute_instance" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  availability_zone = "%s"
  network {
    uuid = "%s"
  }
}
`, HW_AVAILABILITY_ZONE, HW_NETWORK_ID)
