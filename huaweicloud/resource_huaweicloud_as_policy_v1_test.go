package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccASV1Policy_basic(t *testing.T) {
	var asPolicy policies.Policy
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Policy_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1PolicyExists("huaweicloud_as_policy.acc_as_policy", &asPolicy),
				),
			},
		},
	})
}

func testAccCheckASV1PolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating huaweicloud autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_as_policy" {
			continue
		}

		_, err := policies.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS policy still exists")
		}
	}

	log.Printf("[DEBUG] testCheckASV1PolicyDestroy success!")

	return nil
}

func testAccCheckASV1PolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		asClient, err := config.AutoscalingV1Client(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating huaweicloud autoscaling client: %s", err)
		}

		found, err := policies.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] test found is: %#v", found)
		policy = &found

		return nil
	}
}

func testASV1Policy_basic(rName string) string {
	return fmt.Sprintf(`
data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_vpc" "test" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "test" {
  name = "subnet-default"
}

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

resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "%s"
  description = "This is a terraform test security group"
}

resource "huaweicloud_compute_keypair" "acc_key" {
  name       = "%s"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
	image    = data.huaweicloud_images_image.test.id
	flavor   = data.huaweicloud_compute_flavors.test.ids[0]
	key_name = huaweicloud_compute_keypair.acc_key.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
  }
}

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = data.huaweicloud_vpc.test.id
  networks {
    id = data.huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.secgroup.id
  }
}

resource "huaweicloud_as_policy" "acc_as_policy"{
  scaling_policy_name = "%s"
  scaling_policy_type = "SCHEDULED"
  scaling_group_id    = huaweicloud_as_group.acc_as_group.id

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time = "2021-12-22T12:00Z"
  }
}
`, rName, rName, rName, rName, rName)
}
