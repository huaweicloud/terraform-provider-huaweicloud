package vpcep

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceVPCEPServiceConnections_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_vpcep_service_connections.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceVpcepServiceConnections_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "connections.0.endpoint_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.marker_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.domain_id"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.status"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.created_at"),
					resource.TestCheckResourceAttrSet(rName, "connections.0.updated_at"),

					resource.TestCheckOutput("endpoint_id_filter_is_useful", "true"),

					resource.TestCheckOutput("marker_id_filter_is_useful", "true"),

					resource.TestCheckOutput("status_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceVpcepServiceConnections_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_vpcep_service_connections" "test" {
  service_id = huaweicloud_vpcep_endpoint.test.service_id
}

data "huaweicloud_vpcep_service_connections" "endpoint_id_filter" {
  service_id  = huaweicloud_vpcep_endpoint.test.service_id
  endpoint_id = data.huaweicloud_vpcep_service_connections.test.connections.0.endpoint_id
}

locals {
  endpoint_id = data.huaweicloud_vpcep_service_connections.test.connections.0.endpoint_id
}

output "endpoint_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_service_connections.endpoint_id_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_service_connections.endpoint_id_filter.connections[*].endpoint_id : v == local.endpoint_id]
  )  
}

data "huaweicloud_vpcep_service_connections" "marker_id_filter" {
  service_id = huaweicloud_vpcep_endpoint.test.service_id
  marker_id  = data.huaweicloud_vpcep_service_connections.test.connections.0.marker_id
}
  
locals {
  marker_id = data.huaweicloud_vpcep_service_connections.test.connections.0.marker_id
}
  
output "marker_id_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_service_connections.marker_id_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_service_connections.marker_id_filter.connections[*].marker_id : v == local.marker_id]
  )  
}

data "huaweicloud_vpcep_service_connections" "status_filter" {
  service_id = huaweicloud_vpcep_endpoint.test.service_id
  status     = data.huaweicloud_vpcep_service_connections.test.connections.0.status
}
  
locals {
  status = data.huaweicloud_vpcep_service_connections.test.connections.0.status
}
  
output "status_filter_is_useful" {
  value = length(data.huaweicloud_vpcep_service_connections.status_filter.connections) > 0 && alltrue(
    [for v in data.huaweicloud_vpcep_service_connections.status_filter.connections[*].status : v == local.status]
  )  
}
`, testAccVPCEndpoint_Basic(name))
}
