package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/services"
)

func TestAccVPCEPServiceBasic(t *testing.T) {
	var service services.Service

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "huaweicloud_vpcep_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEPServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPServiceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEPServiceExists(resourceName, &service),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "false"),
					resource.TestCheckResourceAttr(resourceName, "server_type", "VM"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8080"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
				),
			},
			{
				Config: testAccVPCEPServiceUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-"+rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttr(resourceName, "approval", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.protocol", "TCP"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.service_port", "8088"),
					resource.TestCheckResourceAttr(resourceName, "port_mapping.0.terminal_port", "80"),
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

func TestAccVPCEPServicePermission(t *testing.T) {
	var service services.Service

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "huaweicloud_vpcep_service.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEPServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEPServicePermission(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEPServiceExists(resourceName, &service),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttrSet(resourceName, "permissions.#"),
				),
			},
			{
				Config: testAccVPCEPServicePermissionUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "status", "available"),
					resource.TestCheckResourceAttrSet(resourceName, "permissions.#"),
				),
			},
		},
	})
}

func testAccCheckVPCEPServiceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcepClient, err := config.VPCEPClient(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating VPC endpoint client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpcep_service" {
			continue
		}

		_, err := services.Get(vpcepClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VPC endpoint service still exists")
		}
	}

	return nil
}

func testAccCheckVPCEPServiceExists(n string, service *services.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		vpcepClient, err := config.VPCEPClient(HW_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating VPC endpoint client: %s", err)
		}

		found, err := services.Get(vpcepClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VPC endpoint service not found")
		}

		*service = *found

		return nil
	}
}

func testAccVPCEPServicePrecondition(rName string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

resource "huaweicloud_compute_instance" "ecs" {
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

func testAccVPCEPServiceBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEPServicePrecondition(rName), rName)
}

func testAccVPCEPServiceUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_service" "test" {
  name        = "tf-%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = true

  port_mapping {
    service_port  = 8088
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc-update"
  }
}
`, testAccVPCEPServicePrecondition(rName), rName)
}

func testAccVPCEPServicePermission(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = false
  permissions = ["iam:domain::1234", "iam:domain::5678"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, testAccVPCEPServicePrecondition(rName), rName)
}

func testAccVPCEPServicePermissionUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = false
  permissions = ["iam:domain::1234", "iam:domain::abcd"]

  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}
`, testAccVPCEPServicePrecondition(rName), rName)
}
