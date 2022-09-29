package vpcep

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/vpcep/v1/endpoints"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccVPCEndpointApproval_Basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	rName := acceptance.RandomAccResourceNameWithDash()
	resourceName := "huaweicloud_vpcep_approval.approval"

	rc := acceptance.InitResourceCheck(
		"huaweicloud_vpcep_endpoint.test",
		&endpoint,
		getVpcepEndpointResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointApproval_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttrPair(resourceName, "id", "huaweicloud_vpcep_service.test", "id"),
					resource.TestCheckResourceAttrPair(resourceName, "connections.0.endpoint_id",
						"huaweicloud_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "accepted"),
				),
			},
			{
				Config: testAccVPCEndpointApproval_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "connections.0.endpoint_id",
						"huaweicloud_vpcep_endpoint.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "rejected"),
				),
			},
		},
	})
}

func testAccVPCEndpointApproval_Basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = true

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
  lifecycle {
    ignore_changes = [enable_dns]
  }
}

resource "huaweicloud_vpcep_approval" "approval" {
  service_id = huaweicloud_vpcep_service.test.id
  endpoints  = [huaweicloud_vpcep_endpoint.test.id]
}
`, testAccVPCEndpoint_Precondition(rName), rName)
}

func testAccVPCEndpointApproval_Update(rName string) string {
	return fmt.Sprintf(`
%s

resource "huaweicloud_vpcep_service" "test" {
  name        = "%s"
  server_type = "VM"
  vpc_id      = data.huaweicloud_vpc.myvpc.id
  port_id     = huaweicloud_compute_instance.ecs.network[0].port
  approval    = true

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
  lifecycle {
    ignore_changes = [enable_dns]
  }
}

resource "huaweicloud_vpcep_approval" "approval" {
  service_id = huaweicloud_vpcep_service.test.id
  endpoints  = []
}
`, testAccVPCEndpoint_Precondition(rName), rName)
}
