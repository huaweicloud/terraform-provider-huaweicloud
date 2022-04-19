package ecs

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/ecs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/block_devices"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func getVolumeAttachResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.ComputeV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating HuaweiCloud compute v1 client: %s", err)
	}

	instanceId, volumeId, err := ecs.ParseComputeVolumeAttachmentId(state.Primary.ID)
	if err != nil {
		return nil, err
	}

	found, err := block_devices.Get(c, instanceId, volumeId).Extract()
	if err != nil {
		return nil, err
	}

	if found.ServerId != instanceId || found.VolumeId != volumeId {
		return nil, fmt.Errorf("volumeAttach not found %s", state.Primary.ID)
	}

	return found, nil
}

func TestAccComputeVolumeAttach_basic(t *testing.T) {
	var va block_devices.VolumeAttachment
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_compute_volume_attach.va_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&va,
		getVolumeAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttach_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_compute_instance.instance_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_id", "huaweicloud_evs_volume.test", "id"),
					resource.TestCheckResourceAttrSet(resourceName, "pci_address"),
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

func TestAccComputeVolumeAttach_device(t *testing.T) {
	var va block_devices.VolumeAttachment
	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_compute_volume_attach.va_1"
	rc := acceptance.InitResourceCheck(
		resourceName,
		&va,
		getVolumeAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttach_device(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "instance_id",
						"huaweicloud_compute_instance.instance_1", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "volume_id", "huaweicloud_evs_volume.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "device", "/dev/vdb"),
					resource.TestCheckResourceAttrSet(resourceName, "pci_address"),
				),
			},
		},
	})
}

func TestAccComputeVolumeAttach_multiple(t *testing.T) {
	var va block_devices.VolumeAttachment
	rName := acceptance.RandomAccResourceNameWithDash()
	rc := acceptance.InitResourceCheck(
		"huaweicloud_compute_volume_attach.test",
		&va,
		getVolumeAttachResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeVolumeAttach_multiple(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckMultiResourcesExists(2),
					resource.TestCheckResourceAttrPair("huaweicloud_compute_volume_attach.test.0", "instance_id",
						"huaweicloud_compute_instance.test.0", "id"),
					resource.TestCheckResourceAttrPair("huaweicloud_compute_volume_attach.test.0", "volume_id",
						"huaweicloud_evs_volume.test", "id"),
					resource.TestCheckResourceAttrPair("huaweicloud_compute_volume_attach.test.1", "instance_id",
						"huaweicloud_compute_instance.test.1", "id"),
					resource.TestCheckResourceAttrPair("huaweicloud_compute_volume_attach.test.1", "volume_id",
						"huaweicloud_evs_volume.test", "id"),
				),
			},
		},
	})
}

func testAccComputeVolumeAttach_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
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

resource "huaweicloud_compute_volume_attach" "va_1" {
  instance_id = huaweicloud_compute_instance.instance_1.id
  volume_id   = huaweicloud_evs_volume.test.id
}
`, testAccCompute_data, rName, rName)
}

func testAccComputeVolumeAttach_device(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_evs_volume" "test" {
  name              = "%s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
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

resource "huaweicloud_compute_volume_attach" "va_1" {
  instance_id = huaweicloud_compute_instance.instance_1.id
  volume_id   = huaweicloud_evs_volume.test.id
  device      = "/dev/vdb"
}
`, testAccCompute_data, rName, rName)
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
  name        = "Ubuntu 20.04 server 64bit"
  most_recent = true
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccComputeVolumeAttach_multiple(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_evs_volume" "test" {
  name              = "%[2]s"
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  volume_type       = "SAS"
  size              = 10
  
  multiattach = true
}

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

resource "huaweicloud_compute_volume_attach" "test" {
  count = 2

  instance_id = huaweicloud_compute_instance.test[count.index].id
  volume_id   = huaweicloud_evs_volume.test.id
}
`, testAccCompute_data, rName, rName)
}
