package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceEndpointConnections_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_apig_endpoint_connections.test"
		rName      = acceptance.RandomAccResourceName()
		dc         = acceptance.InitDataSourceCheck(dataSource)

		byEndpointId   = "data.huaweicloud_apig_endpoint_connections.by_endpoint_id_filter"
		dcByEndpointId = acceptance.InitDataSourceCheck(byEndpointId)

		byPacketId   = "data.huaweicloud_apig_endpoint_connections.by_packet_id_filter"
		dcByPacketId = acceptance.InitDataSourceCheck(byPacketId)

		byStatus   = "data.huaweicloud_apig_endpoint_connections.by_status_filter"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceEndpointConnections_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dataSource, "connections.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dataSource, "connections.0.domain_id"),
					// An example of time format is "yyyy-MM-ddTHH:mm:ss+08:00".
					resource.TestMatchResourceAttr(dataSource, "connections.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(dataSource, "connections.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					dcByEndpointId.CheckResourceExists(),
					resource.TestCheckOutput("is_endpoint_id_filter_useful", "true"),
					dcByPacketId.CheckResourceExists(),
					resource.TestCheckOutput("is_packet_id_filter_useful", "true"),
					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceEndpointConnections_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_endpoint_connections" "test" {
  depends_on = [
    huaweicloud_apig_endpoint_connection_management.test
  ]

  instance_id = huaweicloud_apig_instance.test.id
}

# By endpoint ID filter
locals {
  endpoint_id = huaweicloud_apig_endpoint_connection_management.test.id
}

data "huaweicloud_apig_endpoint_connections" "by_endpoint_id_filter" {
  instance_id = huaweicloud_apig_instance.test.id
  endpoint_id = local.endpoint_id
}

locals {
  endpoint_id_filter_result = [
    for v in data.huaweicloud_apig_endpoint_connections.by_endpoint_id_filter.connections[*].id : v == local.endpoint_id
  ]
}

output "is_endpoint_id_filter_useful" {
  value = length(local.endpoint_id_filter_result) > 0 && alltrue(local.endpoint_id_filter_result)
}

# By packet ID filter
# There is no "packet_id" field in the corresponding resource.
locals {
  packet_id = data.huaweicloud_apig_endpoint_connections.by_endpoint_id_filter.connections[0].packet_id
}

data "huaweicloud_apig_endpoint_connections" "by_packet_id_filter" {
  instance_id = huaweicloud_apig_instance.test.id
  packet_id   = local.packet_id
}

locals {
  packet_id_filter_result = [
    for v in data.huaweicloud_apig_endpoint_connections.by_packet_id_filter.connections[*].packet_id : v == local.packet_id
  ]
}

output "is_packet_id_filter_useful" {
  value = length(local.packet_id_filter_result) > 0 && alltrue(local.packet_id_filter_result)
}

# By status filter
locals {
  status = huaweicloud_apig_endpoint_connection_management.test.status
}

data "huaweicloud_apig_endpoint_connections" "by_status_filter" {
  instance_id = huaweicloud_apig_instance.test.id
  status      = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_apig_endpoint_connections.by_status_filter.connections[*].status : v == local.status
  ]
}

output "is_status_filter_useful" {
  value = length(local.status_filter_result) > 0 && alltrue(local.status_filter_result)
}
`, testAccEndpointConnectionManagement_basic_step1(name, acceptance.RandomAccResourceName()))
}
