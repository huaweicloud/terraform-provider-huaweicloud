package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/configurations"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
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
				),
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
				),
			},
		},
	})
}

func testAccCheckASConfigurationDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(acceptance.HW_REGION_NAME)
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

func testAccASConfiguration_base(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

data "huaweicloud_networking_secgroup" "test" {
  name = "default"
}

resource "huaweicloud_compute_keypair" "acc_key" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}
`, rName)
}

func testAccASConfiguration_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
    image  = data.huaweicloud_images_image.test.id
    flavor = data.huaweicloud_compute_flavors.test.ids[0]
    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
    key_name = huaweicloud_compute_keypair.acc_key.id
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
  security_group_ids = [data.huaweicloud_networking_secgroup.test.id]

  network {
    uuid = data.huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
    instance_id = huaweicloud_compute_instance.test.id
    key_name    = huaweicloud_compute_keypair.acc_key.id
  }
}
`, testAccASConfiguration_base(rName), rName, rName)
}
