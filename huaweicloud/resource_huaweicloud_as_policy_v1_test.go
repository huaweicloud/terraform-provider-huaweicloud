package huaweicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/policies"
)

func TestAccASV1Policy_basic(t *testing.T) {
	var asPolicy policies.Policy

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccAsConfigPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckASV1PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASV1Policy_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASV1PolicyExists("huaweicloud_as_policy.acc_as_policy", &asPolicy),
				),
			},
		},
	})
}

func testAccCheckASV1PolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
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

		config := testAccProvider.Meta().(*Config)
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

var testASV1Policy_basic = fmt.Sprintf(`
resource "huaweicloud_networking_secgroup" "secgroup" {
  name        = "terraform"
  description = "This is a terraform test security group"
}

resource "huaweicloud_compute_keypair" "acc_key" {
  name       = "acc_key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDAjpC1hwiOCCmKEWxJ4qzTTsJbKzndLo1BCz5PcwtUnflmU+gHJtWMZKpuEGVi29h0A/+ydKek1O18k10Ff+4tyFjiHDQAT9+OfgWf7+b1yK+qDip3X1C0UPMbwHlTfSGWLGZquwhvEFx9k3h/M+VtMvwR1lJ9LUyTAImnNjWG7TAIPmui30HvM2UiFEmqkr4ijq45MyX2+fLIePLRIFuu1p4whjHAQYufqyno3BS48icQb4p6iVEZPo4AE2o9oIyQvj2mx4dk5Y8CgSETOZTYDOR3rU2fZTRDRgPJDH9FWvQjF5tA0p3d9CoWWd2s6GKKbfoUIi8R/Db1BSPJwkqB jrp-hp-pc"
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "acc_as_config"
  instance_config {
	image    = "%s"
	key_name = huaweicloud_compute_keypair.acc_key.id
    disk {
      size        = 40
      volume_type = "SATA"
      disk_type   = "SYS"
    }
  }
}

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "acc_as_group"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = "%s"
  networks {
    id = "%s"
  }
  security_groups {
    id = huaweicloud_networking_secgroup.secgroup.id
  }
}

resource "huaweicloud_as_policy" "acc_as_policy"{
  scaling_policy_name = "terraform"
  scaling_policy_type = "SCHEDULED"
  scaling_group_id    = huaweicloud_as_group.acc_as_group.id

  scaling_policy_action {
    operation       = "ADD"
    instance_number = 1
  }
  scheduled_policy {
    launch_time = "2020-12-22T12:00Z"
  }
}
`, HW_IMAGE_ID, HW_VPC_ID, HW_NETWORK_ID)
