package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
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
					resource.TestCheckResourceAttr(resourceName, "status", "ACTIVE"),
					resource.TestCheckResourceAttrSet(resourceName, "system_disk_id"),
					resource.TestCheckResourceAttrSet(resourceName, "security_groups.#"),
					resource.TestCheckResourceAttrSet(resourceName, "network.#"),
					resource.TestCheckResourceAttrSet(resourceName, "volume_attached.#"),
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

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_disks(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_prePaid(t *testing.T) {
	var instance servers.Server

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_prePaid(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
				),
			},
		},
	})
}

func TestAccComputeV2Instance_tags(t *testing.T) {
	var instance servers.Server

	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	resourceName := "huaweicloud_compute_instance.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Instance_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTags(&instance, "key", "value"),
				),
			},
			{
				Config: testAccComputeV2Instance_tags2(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(&instance, "foo2", "bar2"),
					testAccCheckComputeV2InstanceTags(&instance, "key", "value2"),
				),
			},
			{
				Config: testAccComputeV2Instance_notags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceNoTags(&instance),
				),
			},
			{
				Config: testAccComputeV2Instance_tags(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists(resourceName, &instance),
					testAccCheckComputeV2InstanceTags(&instance, "foo", "bar"),
					testAccCheckComputeV2InstanceTags(&instance, "key", "value"),
				),
			},
		},
	})
}

func testAccCheckComputeV2InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(HW_REGION_NAME)
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
		computeClient, err := config.ComputeV2Client(HW_REGION_NAME)
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

func testAccCheckComputeV2InstanceTags(
	instance *servers.Server, k, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		client, err := config.ComputeV1Client(HW_REGION_NAME)
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
		client, err := config.ComputeV1Client(HW_REGION_NAME)
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

const testAccCompute_data = `
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}
`

func testAccComputeV2Instance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_disks(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
 
  system_disk_type = "SAS"
  system_disk_size = 50

  data_disks {
    type = "SAS"
    size = "10"
  }

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_prePaid(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  charging_mode = "prePaid"
  period_unit   = "month"
  period        = 1
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_tags(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_tags2(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }

  tags = {
    foo2 = "bar2"
    key = "value2"
  }
}
`, testAccCompute_data, rName)
}

func testAccComputeV2Instance_notags(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name              = "%s"
  image_id          = data.huaweicloud_images_image.test.id
  flavor_id         = data.huaweicloud_compute_flavors.test.ids[0]
  security_groups   = ["default"]
  availability_zone = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}
