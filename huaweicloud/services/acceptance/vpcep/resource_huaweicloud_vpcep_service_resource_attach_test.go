package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccEndpointServiceResourceAttach_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testEndpointServiceResourceAttach_basic(rName),
			},
		},
	})
}

func testEndpointServiceResourceAttach_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_vpc" "test" {
  name = "%[1]s-vpc"
  cidr = "192.168.0.0/16"
}

resource "huaweicloud_vpc_subnet" "test" {
  name       = "%[1]s-subnet"
  vpc_id     = huaweicloud_vpc.test.id
  cidr       = cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0)
  gateway_ip = cidrhost(cidrsubnet(huaweicloud_vpc.test.cidr, 4, 0), 1)
}

data "huaweicloud_availability_zones" "test" {}

data "huaweicloud_images_image" "test" {
  name        = "Ubuntu 24.04 server 64bit"
  most_recent = true
}

data "huaweicloud_compute_flavors" "test" {
  availability_zone = data.huaweicloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

resource "huaweicloud_networking_secgroup" "test" {
  name                 = "%[1]s"
  delete_default_rules = true
}

resource "huaweicloud_compute_instance" "test" {
  count = 3

  name               = "%[1]s-${count.index}"
  image_id           = data.huaweicloud_images_image.test.id
  flavor_id          = data.huaweicloud_compute_flavors.test.ids[0]
  security_group_ids = [huaweicloud_networking_secgroup.test.id]
  availability_zone  = data.huaweicloud_availability_zones.test.names[0]

  network {
    uuid = huaweicloud_vpc_subnet.test.id
  }
}

resource "huaweicloud_vpcep_service" "test" {
  name        = "%[1]s"
  server_type = "VM"
  vpc_id      = huaweicloud_vpc.test.id
  port_id     = huaweicloud_compute_instance.test[0].network[0].port

  port_mapping {
    protocol      = "TCP"
    service_port  = 8080
    terminal_port = 80
  }
}
`, rName)
}

func testEndpointServiceResourceAttach_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_vpcep_service_resource_attach" "test" {
  service_id = huaweicloud_vpcep_service.test.id

  dynamic "server_resources" {
    for_each = slice(huaweicloud_compute_instance.test, 1, 3)
    content {
      resource_id          = server_resources.value.id
      availability_zone_id = data.huaweicloud_availability_zones.test.names[0]
    }
  }
}
`, testEndpointServiceResourceAttach_base(rName))
}
