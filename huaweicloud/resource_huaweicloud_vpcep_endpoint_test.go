package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/endpoints"
)

func TestAccVPCEndpointBasic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := fmt.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "huaweicloud_vpcep_endpoint.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEndpointExists(resourceName, &endpoint),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
				),
			},
			{
				Config: testAccVPCEndpointUpdate(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "tags.owner", "tf-acc-update"),
					resource.TestCheckResourceAttr(resourceName, "tags.foo", "bar"),
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

func TestAccVPCEndpointPublic(t *testing.T) {
	var endpoint endpoints.Endpoint
	resourceName := "huaweicloud_vpcep_endpoint.myendpoint"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointPublic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEndpointExists(resourceName, &endpoint),
					resource.TestCheckResourceAttr(resourceName, "status", "accepted"),
					resource.TestCheckResourceAttr(resourceName, "enable_dns", "true"),
					resource.TestCheckResourceAttr(resourceName, "enable_whitelist", "true"),
					resource.TestCheckResourceAttr(resourceName, "service_type", "interface"),
					resource.TestCheckResourceAttr(resourceName, "whitelist.#", "2"),
					resource.TestCheckResourceAttrSet(resourceName, "service_name"),
					resource.TestCheckResourceAttrSet(resourceName, "domain_name"),
					resource.TestCheckResourceAttrSet(resourceName, "ip_address"),
				),
			},
		},
	})
}

func testAccCheckVPCEndpointDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	vpcepClient, err := config.VPCEPClient(HW_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating VPC endpoint client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_vpcep_endpoint" {
			continue
		}

		_, err := endpoints.Get(vpcepClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VPC endpoint still exists")
		}
	}

	return nil
}

func testAccCheckVPCEndpointExists(n string, endpoint *endpoints.Endpoint) resource.TestCheckFunc {
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

		found, err := endpoints.Get(vpcepClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VPC endpoint not found")
		}

		*endpoint = *found

		return nil
	}
}

func testAccVPCEndpointPrecondition(rName string) string {
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

func testAccVPCEndpointBasic(rName string) string {
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

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id  = huaweicloud_vpcep_service.test.id
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  network_id  = data.huaweicloud_vpc_subnet.test.id
  enable_dns  = true

  tags = {
    owner = "tf-acc"
  }
}
`, testAccVPCEndpointPrecondition(rName), rName)
}

func testAccVPCEndpointUpdate(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_service" "test" {
  name        = "tf-%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = false

  port_mapping {
    service_port  = 8088
    terminal_port = 80
  }
  tags = {
    owner = "tf-acc"
  }
}

resource "huaweicloud_vpcep_endpoint" "test" {
  service_id  = huaweicloud_vpcep_service.test.id
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  network_id  = data.huaweicloud_vpc_subnet.test.id
  enable_dns  = true

  tags = {
    owner = "tf-acc-update"
    foo   = "bar"
  }
}
`, testAccVPCEndpointPrecondition(rName), rName)
}

var testAccVPCEndpointPublic string = `
data "huaweicloud_vpc" "myvpc" {
  name = "vpc-default"
}

data "huaweicloud_vpc_subnet" "mynet" {
  vpc_id = data.huaweicloud_vpc.myvpc.id
  name   = "subnet-default"
}

data "huaweicloud_vpcep_public_services" "cloud_service" {
  service_name = "dis"
}

resource "huaweicloud_vpcep_endpoint" "myendpoint" {
  service_id       = data.huaweicloud_vpcep_public_services.cloud_service.services[0].id
  vpc_id           = data.huaweicloud_vpc.myvpc.id
  network_id       = data.huaweicloud_vpc_subnet.mynet.id
  enable_dns       = true
  enable_whitelist = true
  whitelist        = ["192.168.0.0/24", "10.10.10.10"]
}
`
