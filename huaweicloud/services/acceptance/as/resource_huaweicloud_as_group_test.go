package as

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccASGroup_basic(t *testing.T) {
	var asGroup groups.Group
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_group.acc_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "0"),
					resource.TestCheckResourceAttr(resourceName, "lbaas_listeners.0.protocol_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
					resource.TestCheckResourceAttr(resourceName, "tags.key", "value"),
					resource.TestCheckResourceAttr(resourceName, "multi_az_scaling_policy", "EQUILIBRIUM_DISTRIBUTE"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "300"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_time", "5"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "600"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zones.#"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"delete_instances",
				},
			},
			{
				Config: testASGroup_basic_disable(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "status", "PAUSED"),
				),
			},
			{
				Config: testASGroup_basic_enable(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "multi_az_scaling_policy", "PICK_FIRST"),
					resource.TestCheckResourceAttr(resourceName, "cool_down_time", "600"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_time", "15"),
					resource.TestCheckResourceAttr(resourceName, "health_periodic_audit_grace_period", "900"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
		},
	})
}

func TestAccASGroup_withEpsId(t *testing.T) {
	var asGroup groups.Group
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_group.acc_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheckEpsID(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_withEpsId(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.HW_ENTERPRISE_PROJECT_ID_TEST),
				),
			},
		},
	})
}

func TestAccASGroup_forceDelete(t *testing.T) {
	var asGroup groups.Group
	rName := acceptance.RandomAccResourceName()
	resourceName := "huaweicloud_as_group.acc_as_group"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckASGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testASGroup_forceDelete(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckASGroupExists(resourceName, &asGroup),
					resource.TestCheckResourceAttr(resourceName, "desire_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "min_instance_number", "2"),
					resource.TestCheckResourceAttr(resourceName, "max_instance_number", "5"),
					resource.TestCheckResourceAttr(resourceName, "instances.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "status", "INSERVICE"),
				),
			},
		},
	})
}

func testAccCheckASGroupDestroy(s *terraform.State) error {
	config := acceptance.TestAccProvider.Meta().(*config.Config)
	asClient, err := config.AutoscalingV1Client(acceptance.HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("error creating autoscaling client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_as_group" {
			continue
		}

		_, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("AS group still exists")
		}
	}

	return nil
}

func testAccCheckASGroupExists(n string, group *groups.Group) resource.TestCheckFunc {
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

		found, err := groups.Get(asClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Autoscaling Group not found")
		}

		group = &found
		return nil
	}
}

func testASGroup_Base(rName string) string {
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

resource "huaweicloud_lb_loadbalancer" "loadbalancer_1" {
  name          = "%s"
  vip_subnet_id = data.huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

resource "huaweicloud_lb_listener" "listener_1" {
  name            = "%s"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = huaweicloud_lb_loadbalancer.loadbalancer_1.id
}

resource "huaweicloud_lb_pool" "pool_1" {
  name        = "%s"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = huaweicloud_lb_listener.listener_1.id
}

resource "huaweicloud_as_configuration" "acc_as_config"{
  scaling_configuration_name = "%s"
  instance_config {
	image    = data.huaweicloud_images_image.test.id
	flavor   = data.huaweicloud_compute_flavors.test.ids[0]
    key_name = huaweicloud_compute_keypair.acc_key.id
    disk {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  }
}`, rName, rName, rName, rName, rName, rName)
}

func testASGroup_basic(rName string) string {
	return fmt.Sprintf(`
%s

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
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_basic_disable(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = data.huaweicloud_vpc.test.id
  enable                   = false

  networks {
    id = data.huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.secgroup.id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_basic_enable(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = data.huaweicloud_vpc.test.id
  enable                   = true

  multi_az_scaling_policy            = "PICK_FIRST"
  cool_down_time                     = 600
  health_periodic_audit_time         = 15
  health_periodic_audit_grace_period = 900

  networks {
    id = data.huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.secgroup.id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName)
}

func testASGroup_withEpsId(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  vpc_id                   = data.huaweicloud_vpc.test.id
  enterprise_project_id    = "%s"

  networks {
    id = data.huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.secgroup.id
  }
  lbaas_listeners {
    pool_id       = huaweicloud_lb_pool.pool_1.id
    protocol_port = huaweicloud_lb_listener.listener_1.protocol_port
  }
  tags = {
    foo = "bar"
    key = "value"
  }
}
`, testASGroup_Base(rName), rName, acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}

func testASGroup_forceDelete(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_as_group" "acc_as_group"{
  scaling_group_name       = "%s"
  scaling_configuration_id = huaweicloud_as_configuration.acc_as_config.id
  min_instance_number      = 2
  max_instance_number      = 5
  force_delete             = true
  vpc_id                   = data.huaweicloud_vpc.test.id

  networks {
    id = data.huaweicloud_vpc_subnet.test.id
  }
  security_groups {
    id = huaweicloud_networking_secgroup.secgroup.id
  }
}
`, testASGroup_Base(rName), rName)
}
