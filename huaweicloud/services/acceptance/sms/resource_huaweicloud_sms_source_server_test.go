package sms

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/sms"
)

func getSourceServerResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("sms", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SMS client: %s", err)
	}

	return sms.GetSourceServer(client, state.Primary.ID)
}

func TestAccResourceSourceServer_basic(t *testing.T) {
	var obj interface{}
	name := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_sms_source_server.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getSourceServerResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceServer_basic(name),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disks.0.name", "Disk 0"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
					resource.TestCheckResourceAttrSet(resourceName, "os_type"),
					resource.TestCheckResourceAttrSet(resourceName, "os_version"),
					resource.TestCheckResourceAttrSet(resourceName, "firmware"),
					resource.TestCheckResourceAttrSet(resourceName, "cpu_quantity"),
					resource.TestCheckResourceAttrSet(resourceName, "memory"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.partition_style"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.device_use"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.used_size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.os_disk"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.device_use"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.file_system"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.index"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.mount_point"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.used_size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.size_per_cluster"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.ip"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.netmask"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.mtu"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.mac"),
					resource.TestCheckResourceAttrSet(resourceName, "agent_version"),
					resource.TestCheckResourceAttrSet(resourceName, "migration_cycle"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "oem_system"),
					resource.TestCheckResourceAttrSet(resourceName, "has_tc"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"virtualization_type", "linux_block_check", "domain_id", "has_rsync",
					"paravirtualization", "raw_devices", "driver_files", "system_services", "account_rights", "boot_loader",
					"system_dir", "kernel_version", "start_type", "io_read_wait", "platform", "migprojectid", "copystate"},
			},
			{
				Config: testAccSourceServer_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "disks.0.name", "Disk 2"),
					resource.TestCheckResourceAttrSet(resourceName, "ip"),
					resource.TestCheckResourceAttrSet(resourceName, "os_type"),
					resource.TestCheckResourceAttrSet(resourceName, "os_version"),
					resource.TestCheckResourceAttrSet(resourceName, "firmware"),
					resource.TestCheckResourceAttrSet(resourceName, "cpu_quantity"),
					resource.TestCheckResourceAttrSet(resourceName, "memory"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.partition_style"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.device_use"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.used_size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.os_disk"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.#"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.device_use"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.file_system"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.index"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.mount_point"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.used_size"),
					resource.TestCheckResourceAttrSet(resourceName, "disks.0.physical_volumes.0.size_per_cluster"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.#"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.name"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.ip"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.netmask"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.gateway"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.mtu"),
					resource.TestCheckResourceAttrSet(resourceName, "networks.0.mac"),
					resource.TestCheckResourceAttrSet(resourceName, "agent_version"),
					resource.TestCheckResourceAttrSet(resourceName, "migration_cycle"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "oem_system"),
					resource.TestCheckResourceAttrSet(resourceName, "has_tc"),
				),
			},
		},
	})
}

func testAccSourceServer_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sms_source_server" "test" {
  ip                 = huaweicloud_compute_instance.test.access_ip_v4
  name               = "%[2]s"
  os_type            = "LINUX"
  os_version         = "Ubuntu 18.04 server 64bit"
  firmware           = "BIOS"
  boot_loader        = "GRUB"
  raw_devices        = ""
  has_rsync          = true
  paravirtualization = true
  cpu_quantity       = 2
  memory             = 4018196480
  agent_version      = "25.2.0"

  disks {
    name            = "Disk 0"
    partition_style = "MBR"
    device_use      = "BOOT"
    size            = 85897247744
    used_size       = 42943137792
    physical_volumes {
      device_use  = "OS"
      file_system = "ext4"
      mount_point = "/"
      name        = "/dev/vda1"
      size        = 42943137792
      used_size   = 1071640576
    }
  }

  networks {
    name    = "eth0"
    ip      = huaweicloud_compute_instance.test.access_ip_v4
    netmask = "255.255.255.0"
    gateway = "192.168.0.1"
    mtu     = 1
    mac     = huaweicloud_compute_instance.test.network[0].mac
  }
}
`, testAccComputeInstance_basic(name), name)
}

func testAccSourceServer_updated(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_sms_source_server" "test" {
  ip                 = huaweicloud_compute_instance.test.access_ip_v4
  os_type            = "LINUX"
  os_version         = "Ubuntu 18.04 server 64bit"
  firmware           = "BIOS"
  boot_loader        = "GRUB"
  raw_devices        = ""
  has_rsync          = true
  paravirtualization = true
  cpu_quantity       = 2
  memory             = 4018196480
  agent_version      = "25.2.0"
  copystate          = "unavailable"
  migration_cycle    = "checking"

  disks {
    name            = "Disk 2"
    partition_style = "MBR"
    device_use      = "BOOT"
    size            = 85897247744
    used_size       = 42943137792
    physical_volumes {
      device_use  = "OS"
      file_system = "ext4"
      mount_point = "/"
      name        = "/dev/vda1"
      size        = 42943137792
      used_size   = 1071640576
    }
  }

  networks {
    name    = "eth0"
    ip      = huaweicloud_compute_instance.test.access_ip_v4
    netmask = "255.255.255.0"
    gateway = "192.168.0.1"
    mtu     = 1
    mac     = huaweicloud_compute_instance.test.network[0].mac
  }
}
`, testAccComputeInstance_basic(name), name)
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

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}
`

func testAccComputeInstance_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}
`, testAccCompute_data, rName)
}
