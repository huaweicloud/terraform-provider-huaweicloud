package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceEndpoints_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_ga_endpoints.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byEndpointId   = "data.huaweicloud_ga_endpoints.filter_by_endpoint_id"
		dcByEndpointId = acceptance.InitDataSourceCheck(byEndpointId)

		byStatus   = "data.huaweicloud_ga_endpoints.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byResourceId   = "data.huaweicloud_ga_endpoints.filter_by_resource_id"
		dcByResourceId = acceptance.InitDataSourceCheck(byResourceId)

		byResourceType   = "data.huaweicloud_ga_endpoints.filter_by_resource_type"
		dcByResourceType = acceptance.InitDataSourceCheck(byResourceType)

		byHealthState   = "data.huaweicloud_ga_endpoints.filter_by_health_state"
		dcByHealthState = acceptance.InitDataSourceCheck(byHealthState)

		byIpAddress   = "data.huaweicloud_ga_endpoints.filter_by_ip_address"
		dcByIpAddress = acceptance.InitDataSourceCheck(byIpAddress)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEndpoints_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.endpoint_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.resource_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.resource_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.weight"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.health_state"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.ip_address"),
					resource.TestCheckResourceAttrSet(dataSourceName, "endpoints.0.created_at"),
					dcByEndpointId.CheckResourceExists(),
					resource.TestCheckOutput("endpoint_id_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByResourceId.CheckResourceExists(),
					resource.TestCheckOutput("resource_id_filter_is_useful", "true"),

					dcByResourceType.CheckResourceExists(),
					resource.TestCheckOutput("resource_type_filter_is_useful", "true"),

					dcByHealthState.CheckResourceExists(),
					resource.TestCheckOutput("health_state_filter_is_useful", "true"),

					dcByIpAddress.CheckResourceExists(),
					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceEndpoints_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ga_endpoints" "test" {
  depends_on = [
    huaweicloud_ga_endpoint.test
  ]

  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
}

locals {
  endpoint_id = data.huaweicloud_ga_endpoints.test.endpoints[0].id
}

data "huaweicloud_ga_endpoints" "filter_by_endpoint_id" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  endpoint_id       = local.endpoint_id
}

locals {
  endpoint_id_filter_result = [
    for v in data.huaweicloud_ga_endpoints.filter_by_endpoint_id.endpoints[*].id : v == local.endpoint_id
  ]
}

output "endpoint_id_filter_is_useful" {
  value = alltrue(local.endpoint_id_filter_result) && length(local.endpoint_id_filter_result) > 0
}

locals {
  status = data.huaweicloud_ga_endpoints.test.endpoints[0].status
}

data "huaweicloud_ga_endpoints" "filter_by_status" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  status            = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ga_endpoints.filter_by_status.endpoints[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  resource_id = data.huaweicloud_ga_endpoints.test.endpoints[0].resource_id
}

data "huaweicloud_ga_endpoints" "filter_by_resource_id" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  resource_id       = local.resource_id
}

locals {
  resource_id_filter_result = [
    for v in data.huaweicloud_ga_endpoints.filter_by_resource_id.endpoints[*].resource_id : v == local.resource_id
  ]
}

output "resource_id_filter_is_useful" {
  value = alltrue(local.resource_id_filter_result) && length(local.resource_id_filter_result) > 0
}

locals {
  resource_type = data.huaweicloud_ga_endpoints.test.endpoints[0].resource_type
}

data "huaweicloud_ga_endpoints" "filter_by_resource_type" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  resource_type     = local.resource_type
}

locals {
  resource_type_filter_result = [
    for v in data.huaweicloud_ga_endpoints.filter_by_resource_type.endpoints[*].resource_type : 
    v == local.resource_type
  ]
}

output "resource_type_filter_is_useful" {
  value = alltrue(local.resource_type_filter_result) && length(local.resource_type_filter_result) > 0
}

locals {
  health_state = data.huaweicloud_ga_endpoints.test.endpoints[0].health_state
}

data "huaweicloud_ga_endpoints" "filter_by_health_state" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  health_state      = local.health_state
}

locals {
  health_state_filter_result = [
    for v in data.huaweicloud_ga_endpoints.filter_by_health_state.endpoints[*].health_state : v == local.health_state
  ]
}

output "health_state_filter_is_useful" {
  value = alltrue(local.health_state_filter_result) && length(local.health_state_filter_result) > 0
}

locals {
  ip_address = data.huaweicloud_ga_endpoints.test.endpoints[0].ip_address
}

data "huaweicloud_ga_endpoints" "filter_by_ip_address" {
  endpoint_group_id = huaweicloud_ga_endpoint_group.test.id
  ip_address        = local.ip_address
}

locals {
  ip_address_filter_result = [
    for v in data.huaweicloud_ga_endpoints.filter_by_ip_address.endpoints[*].ip_address : v == local.ip_address
  ]
}

output "ip_address_filter_is_useful" {
  value = alltrue(local.ip_address_filter_result) && length(local.ip_address_filter_result) > 0
}
`, testEndpoint_basic(name))
}
