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
)

func TestAccASConfiguration_basic(t *testing.T) {
	var asConfig configurations.Configuration
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_configuration.acc_as_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASConfigurationExists(resourceName, &asConfig),
					resource.TestCheckResourceAttr(resourceName, "scaling_configuration_name", rName),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.metadata.some_key", "some_value"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.disk.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.public_ip.0.eip.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "instance_config.0.security_group_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "instance_config.0.user_data"),
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

func TestAccASConfiguration_instance(t *testing.T) {
	var asConfig configurations.Configuration
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_configuration.acc_as_config"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_instance(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASConfigurationExists(resourceName, &asConfig),
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
				},
			},
		},
	})
}

func TestAccASConfiguration_DEH(t *testing.T) {
	var asConfig configurations.Configuration
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAsDedicatedHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_DEH(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASConfigurationExists(resourceName, &asConfig),
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
	var asConfig configurations.Configuration
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccASConfiguration_bandwidth_new_disk(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASConfigurationExists(resourceName, &asConfig),
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

func testAccCheckASConfigurationDestroy(s *terraform.State) error {
	conf := acceptance.TestAccProvider.Meta().(*config.Config)
	asClient, err := conf.AutoscalingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_as_configuration" {
			continue
		}

		_, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS configuration still exists")
		}
	}

	return nil
}

func testAccCheckASConfigurationExists(n string, configuration *configurations.Configuration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := acceptance.TestAccProvider.Meta().(*config.Config)
		asClient, err := config.AutoscalingV1Client(acceptance.HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("error creating autoscaling client: %s", err)
		}

		found, err := configurations.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Configuration not found")
		}

		configuration = &found
		return nil
	}
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
