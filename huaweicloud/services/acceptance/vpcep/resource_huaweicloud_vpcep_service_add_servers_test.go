package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance/common"
)

func TestAccVpcepServiceAddServers_basic(t *testing.T) {
	rName := acceptance.RandomAccResourceNameWithDash()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testVpcepServiceAddServers_basic(rName),
			},
		},
	})
}

func testVpcepServiceAddServers_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_elb_flavors" "flavors" {
  type = "L4"
}

data "huaweicloud_availability_zones" "test" {}

resource "huaweicloud_lb_loadbalancer" "test" {
  vip_subnet_id = huaweicloud_vpc_subnet.test.ipv4_subnet_id
}

locals {
  availability_zone_id = data.huaweicloud_availability_zones.test.names[0]
}

resource "huaweicloud_elb_loadbalancer" "test" {
  name              = "%[2]s"
  description       = "created by terraform acc test"
  vpc_id            = huaweicloud_vpc.test.id
  ipv4_subnet_id    = huaweicloud_vpc_subnet.test.ipv4_subnet_id
  l4_flavor_id      = data.huaweicloud_elb_flavors.flavors.flavors[0].id
  availability_zone = [local.availability_zone_id]
}

resource "huaweicloud_vpcep_service" "test" {
  name        = "%[3]s"
  server_type = "LB"
  vpc_id      = huaweicloud_vpc.test.id
  port_id     = huaweicloud_lb_loadbalancer.test.vip_port_id
  approval    = true
  port_mapping {
    service_port  = 8080
    terminal_port = 80
  }
}

resource "huaweicloud_vpcep_service_add_servers" "test" {
  vpc_endpoint_service_id = huaweicloud_vpcep_service.test.id

  server_resources {
    resource_id          = huaweicloud_elb_loadbalancer.test.id
    availability_zone_id = local.availability_zone_id
  }
}
`, common.TestVpc(rName), rName, rName)
}
