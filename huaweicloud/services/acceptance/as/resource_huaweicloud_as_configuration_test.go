package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/configurations"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func getASConfigurationResourceFunc(cfg *config.Config, state *terraform.ResourceState) (interface{}, error) {
	asClient, err := cfg.AutoscalingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating autoscaling client: %s", err)
	}
	return configurations.Get(asClient, state.Primary.ID).Extract()
}

func TestAccASConfiguration_basic(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.acc_as_config"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.image",
						"data.huaweicloud_images_image.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.flavor",
						"data.huaweicloud_compute_flavors.test", "ids.0"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.key_name",
						"huaweicloud_kps_keypair.acc_key", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.security_group_ids.0",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.flavor_priority_policy", "PICK_FIRST"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.metadata.some_key", "some_value"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.public_ip.0.eip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_config.0.user_data"),
					resource.TestCheckResourceAttrSet(resourceName, "create_time"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_config.0.metadata",
					"instance_config.0.user_data",
				},
			},
		},
	})
}

func TestAccASConfiguration_spot_ecsPassword(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.acc_as_config"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_spot(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.image",
						"data.huaweicloud_images_image.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.flavor",
						"data.huaweicloud_compute_flavors.test", "ids.0"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.key_name", ""),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.security_group_ids.0",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.charging_mode", "spot"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.flavor_priority_policy", "COST_FIRST"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.ecs_group_id",
						"huaweicloud_compute_servergroup.test", "id"),

					resource.TestCheckResourceAttr(resourceName, "instance_config.0.personality.0.path", "/etc/foo.txt"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.personality.0.content", utils.Base64EncodeString("test content")),

					resource.TestCheckResourceAttr(resourceName, "instance_config.0.metadata.some_key", "some_value"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.public_ip.0.eip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_config.0.user_data"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_config.0.metadata", "instance_config.0.user_data"},
			},
		},
	})
}

func TestAccASConfiguration_windowsPassword(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.acc_as_config"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_windowsPassword(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.image",
						"data.huaweicloud_images_image.windows_test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.flavor",
						"data.huaweicloud_compute_flavors.test", "ids.0"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.security_group_ids.0",
						"huaweicloud_networking_secgroup.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.key_name", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.user_data", ""),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.admin_pass", "testTT123!"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.personality.0.path", "fbbo"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.personality.0.content", utils.Base64EncodeString("test content")),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.#", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_config.0.admin_pass"},
			},
		},
	})
}

func TestAccASConfiguration_instance(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.acc_as_config"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_instance(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttrSet(resourceName, "instance_config.0.user_data"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.instance_id",
						"huaweicloud_compute_instance.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_config.0.instance_id",
					"instance_config.0.user_data",
				},
			},
		},
	})
}

func TestAccASConfiguration_DEH(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAsDedicatedHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_DEH(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.tenancy", "dedicated"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.dedicated_host_id", acceptance.HW_DEDICATED_HOST_ID),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.0.volume_type", "GPSSD"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.key_name",
						"huaweicloud_kps_keypair.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_config.0.instance_id",
				},
			},
		},
	})
}

func TestAccASConfiguration_bandwidth_new_disk(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_bandwidth_new_disk(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.0.volume_type", "GPSSD2"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.0.iops", "8000"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.0.throughput", "125"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.1.volume_type", "ESSD2"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.1.iops", "6000"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.public_ip.0.eip.0.bandwidth.0.id",
						"huaweicloud_vpc_bandwidth.test", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_config.0.instance_id",
				},
			},
		},
	})
}

func TestAccASConfiguration_snapshot(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_snapshot(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.0.volume_type", "GPSSD"),
					resource.TestCheckResourceAttrPair(resourceName, "instance_config.0.disk.0.snapshot_id",
						"data.huaweicloud_cbr_backup.test", "children.0.extend_info.0.snapshot_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_config.0.instance_id",
				},
			},
		},
	})
}

func TestAccASConfiguration_dataDiskImage(t *testing.T) {
	var (
		obj          interface{}
		rName        = acceptance.RandomAccResourceName()
		resourceName = "huaweicloud_as_configuration.test"
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&obj,
		getASConfigurationResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAsDataDiskImageId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_dataDiskImage(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.0.volume_type", "GPSSD"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.1.volume_type", "SAS"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.1.data_disk_image_id", acceptance.HW_IMS_DATA_DISK_IMAGE_ID),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"instance_config.0.instance_id",
				},
			},
		},
	})
}

//nolint:revive
func testAccASConfiguration_base(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  visibility  = "public"
  most_recent = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "acc_key" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}
`, common.TestBaseNetwork(rName), rName)
}

func testAccASConfiguration_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
    image              = data.huaweicloud_images_image.test.id
    flavor             = data.huaweicloud_compute_flavors.test.ids[0]
    key_name           = huaweicloud_kps_keypair.acc_key.id
    security_group_ids = [huaweicloud_networking_secgroup.test.id]

    metadata = {
      some_key = "some_value"
    }
    user_data = <<EOT
#!/bin/sh
echo "Hello World! The time is now $(date -R)!" | tee /root/output.txt
EOT

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    public_ip {
      eip {
        ip_type = "5_bgp"
        bandwidth {
          size          = 10
          share_type    = "PER"
          charging_mode = "traffic"
        }
      }
    }
  }
}
`, testAccASConfiguration_base(rName), rName)
}

func testAccASConfiguration_spot(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_servergroup" "test" {
  name     = "%[2]s"
  policies = ["anti-affinity"]
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%[2]s"
  instance_config {
    image                  = data.huaweicloud_images_image.test.id
    flavor                 = data.huaweicloud_compute_flavors.test.ids[0]
    security_group_ids     = [huaweicloud_networking_secgroup.test.id]
    charging_mode          = "spot"
    flavor_priority_policy = "COST_FIRST"
    ecs_group_id           = huaweicloud_compute_servergroup.test.id

    metadata = {
      some_key = "some_value"
    }

# The data injected by user_data is the password of user root for logging in to the ECS by default.
    user_data = <<EOT
#! /bin/bash
echo 'root:$6$V6azyeLwcD3CHlpY$BN3VVq18fmCkj66B4zdHLWevqcxlig' | chpasswd -e
EOT

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    personality {
      path    = "/etc/foo.txt"
      content = base64encode("test content")
    }

    public_ip {
      eip {
        ip_type = "5_bgp"
        bandwidth {
          size          = 10
          share_type    = "PER"
          charging_mode = "traffic"
        }
      }
    }
  }
}
`, testAccASConfiguration_base(rName), rName)
}

func testAccASConfiguration_windowsPassword(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_images_image" "windows_test" {
  name        = "Windows Server 2019 Datacenter 64bit English"
  visibility  = "public"
  most_recent = true
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%[2]s"
  instance_config {
    image              = data.huaweicloud_images_image.windows_test.id
    flavor             = data.huaweicloud_compute_flavors.test.ids[0]
    security_group_ids = [huaweicloud_networking_secgroup.test.id]
    admin_pass         = "testTT123!"

    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }

    personality {
      path    = "fbbo"
      content = base64encode("test content")
    }
  }
}
`, testAccASConfiguration_base(rName), rName)
}

func testAccASConfiguration_instance(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
    instance_id = huaweicloud_compute_instance.test.id
    key_name    = huaweicloud_kps_keypair.acc_key.id
    user_data   = "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkISBUaGUgdGltZSBpcyBub3cgJChkYXRlIC1SKSEiIHwgdGVlIC9yb290L291dHB1dC50eHQK"
  }
}
`, testAccASConfiguration_base(rName), rName, rName)
}

func testAccASConfiguration_newBase(name string) string {
	publicKeyValue := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A" +
		"/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmq" +
		"kr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9Co" +
		"WWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"

	return fmt.Sprintf(`
%[1]s

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  architecture = "x86"
  os           = "CentOS"
  visibility   = "public"
  most_recent  = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[3]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_kps_keypair" "test" {
  name       = "%[2]s"
  public_key = "%[3]s"
}

resource "huaweicloud_vpc_bandwidth" "test" {
  name = "%[2]s"
  size = 5
}
`, common.TestBaseNetwork(name), name, publicKeyValue)
}

func testAccASConfiguration_DEH(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_configuration" "test"{
  scaling_configuration_name = "%[2]s"

  instance_config {
    image              = data.huaweicloud_images_image.test.id
    flavor             = data.huaweicloud_compute_flavors.test.ids[0]
    security_group_ids = [huaweicloud_networking_secgroup.test.id]
    key_name           = huaweicloud_kps_keypair.test.id
    tenancy            = "dedicated"
    dedicated_host_id  = "%[3]s"

    disk {
      size        = 40
      volume_type = "GPSSD"
      disk_type   = "SYS"
    }
  }
}
`, testAccASConfiguration_newBase(name), name, acceptance.HW_DEDICATED_HOST_ID)
}

func testAccASConfiguration_bandwidth_new_disk(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_configuration" "test"{
  scaling_configuration_name = "%[2]s"

  instance_config {
    image              = data.huaweicloud_images_image.test.id
    flavor             = data.huaweicloud_compute_flavors.test.ids[0]
    security_group_ids = [huaweicloud_networking_secgroup.test.id]
    key_name           = huaweicloud_kps_keypair.test.id

    disk {
      size        = 40
      volume_type = "GPSSD2"
      disk_type   = "SYS"
      iops        = 8000
      throughput  = 125
    }

    disk {
      size        = 40
      volume_type = "ESSD2"
      disk_type   = "DATA"
      iops        = 6000
    }

    public_ip {
      eip {
        ip_type = "5_bgp"
        bandwidth {
          share_type = "WHOLE"
          id         = huaweicloud_vpc_bandwidth.test.id
        }
      }
    }
  }
}
`, testAccASConfiguration_newBase(name), name)
}

func testAccASConfiguration_snapshot(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_compute_instance" "test" {
  name               = "%s"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name             = "%[2]s"
  type             = "server"
  consistent_level = "app_consistent"
  protection_type  = "backup"
  size             = 200
}

resource "huaweicloud_images_image" "test" {
  name        = "%[2]s"
  instance_id = huaweicloud_compute_instance.test.id
  description = "Terraform test"
  vault_id    = huaweicloud_cbr_vault.test.id
}

data "huaweicloud_cbr_backup" "test" {
  id = huaweicloud_images_image.test.backup_id
}

resource "huaweicloud_as_configuration" "test"{
  scaling_configuration_name = "%[2]s"

  instance_config {
    image              = huaweicloud_images_image.test.id
    flavor             = huaweicloud_compute_instance.test.flavor_name
    security_group_ids = [huaweicloud_networking_secgroup.test.id]
    key_name           = huaweicloud_kps_keypair.acc_key.id

    disk {
      size        = 40
      volume_type = "GPSSD"
      disk_type   = "SYS"
      snapshot_id = data.huaweicloud_cbr_backup.test.children.0.extend_info.0.snapshot_id
    }
  }
}
`, testAccASConfiguration_base(name), name)
}

func testAccASConfiguration_dataDiskImage(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_as_configuration" "test"{
  scaling_configuration_name = "%[2]s"

  instance_config {
    image              = data.huaweicloud_images_image.test.id
    flavor             = data.huaweicloud_compute_flavors.test.ids[0]
    security_group_ids = [huaweicloud_networking_secgroup.test.id]
    key_name           = huaweicloud_kps_keypair.acc_key.id

    disk {
      size        = 40
      volume_type = "GPSSD"
      disk_type   = "SYS"
    }

    disk {
      size               = 100
      volume_type        = "SAS"
      disk_type          = "DATA"
      data_disk_image_id = "%[3]s"
    }
  }
}
`, testAccASConfiguration_base(name), name, acceptance.HW_IMS_DATA_DISK_IMAGE_ID)
}
