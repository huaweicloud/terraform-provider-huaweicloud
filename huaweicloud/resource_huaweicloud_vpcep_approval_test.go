package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/endpoints"
	"github.com/huaweicloud/golangsdk/openstack/vpcep/v1/services"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func TestAccVPCEndpointApproval_Basic(t *testing.T) {
	var service services.Service
	var endpoint endpoints.Endpoint

	rName := fmtp.Sprintf("acc-test-%s", acctest.RandString(4))
	resourceName := "huaweicloud_vpcep_approval.approval"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVPCEPServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVPCEndpointApproval_Basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVPCEPServiceExists("huaweicloud_vpcep_service.test", &service),
					testAccCheckVPCEndpointExists("huaweicloud_vpcep_endpoint.test", &endpoint),
					resource.TestCheckResourceAttrPtr(resourceName, "id", &service.ID),
					resource.TestCheckResourceAttrPtr(resourceName, "connections.0.endpoint_id", &endpoint.ID),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "accepted"),
				),
			},
			{
				Config: testAccVPCEndpointApproval_Update(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "connections.0.endpoint_id", &endpoint.ID),
					resource.TestCheckResourceAttr(resourceName, "connections.0.status", "rejected"),
				),
			},
		},
	})
}

func testAccVPCEndpointApproval_Basic(rName string) string {
	return fmtp.Sprintf(`
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
	return fmtp.Sprintf(`
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
