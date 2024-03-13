package ga

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceHealthChecks_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceNameWithDash()
		dataSourceName = "data.huaweicloud_ga_health_checks.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byHealthCheckId   = "data.huaweicloud_ga_health_checks.filter_by_health_check_id"
		dcByHealthCheckId = acceptance.InitDataSourceCheck(byHealthCheckId)

		byEndpointGroupId   = "data.huaweicloud_ga_health_checks.filter_by_endpoint_group_id"
		dcByEndpointGroupId = acceptance.InitDataSourceCheck(byEndpointGroupId)

		byStatus   = "data.huaweicloud_ga_health_checks.filter_by_status"
		dcByStatus = acceptance.InitDataSourceCheck(byStatus)

		byProtocol   = "data.huaweicloud_ga_health_checks.filter_by_protocol"
		dcByProtocol = acceptance.InitDataSourceCheck(byProtocol)

		byEnabled   = "data.huaweicloud_ga_health_checks.filter_by_enabled"
		dcByEnabled = acceptance.InitDataSourceCheck(byEnabled)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHealthChecks_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.endpoint_group_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.protocol"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.port"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.interval"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.timeout"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.max_retries"),
					resource.TestCheckResourceAttrSet(dataSourceName, "health_checks.0.created_at"),
					dcByHealthCheckId.CheckResourceExists(),
					resource.TestCheckOutput("health_check_id_filter_is_useful", "true"),

					dcByEndpointGroupId.CheckResourceExists(),
					resource.TestCheckOutput("endpoint_group_id_filter_is_useful", "true"),

					dcByProtocol.CheckResourceExists(),
					resource.TestCheckOutput("protocol_filter_is_useful", "true"),

					dcByStatus.CheckResourceExists(),
					resource.TestCheckOutput("status_filter_is_useful", "true"),

					dcByEnabled.CheckResourceExists(),
					resource.TestCheckOutput("enabled_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceHealthChecks_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_ga_health_checks" "test" {
  depends_on = [
    huaweicloud_ga_health_check.test
  ]
}

locals {
  health_check_id = data.huaweicloud_ga_health_checks.test.health_checks[0].id
}

data "huaweicloud_ga_health_checks" "filter_by_health_check_id" {
  health_check_id = local.health_check_id
}

locals {
  health_check_id_filter_result = [
    for v in data.huaweicloud_ga_health_checks.filter_by_health_check_id.health_checks[*].id : 
    v == local.health_check_id
  ]
}

output "health_check_id_filter_is_useful" {
  value = alltrue(local.health_check_id_filter_result) && length(local.health_check_id_filter_result) > 0
}

locals {
  endpoint_group_id = data.huaweicloud_ga_health_checks.test.health_checks[0].endpoint_group_id
}

data "huaweicloud_ga_health_checks" "filter_by_endpoint_group_id" {
  endpoint_group_id = local.endpoint_group_id
}

locals {
  endpoint_group_id_filter_result = [
    for v in data.huaweicloud_ga_health_checks.filter_by_endpoint_group_id.health_checks[*].endpoint_group_id : 
    v == local.endpoint_group_id
  ]
}

output "endpoint_group_id_filter_is_useful" {
  value = alltrue(local.endpoint_group_id_filter_result) && length(local.endpoint_group_id_filter_result) > 0
}

locals {
  protocol = data.huaweicloud_ga_health_checks.test.health_checks[0].protocol
}

data "huaweicloud_ga_health_checks" "filter_by_protocol" {
  protocol = local.protocol
}

locals {
  protocol_filter_result = [
    for v in data.huaweicloud_ga_health_checks.filter_by_protocol.health_checks[*].protocol : v == local.protocol
  ]
}

output "protocol_filter_is_useful" {
  value = alltrue(local.protocol_filter_result) && length(local.protocol_filter_result) > 0
}

locals {
  status = data.huaweicloud_ga_health_checks.test.health_checks[0].status
}

data "huaweicloud_ga_health_checks" "filter_by_status" {
  status = local.status
}

locals {
  status_filter_result = [
    for v in data.huaweicloud_ga_health_checks.filter_by_status.health_checks[*].status : v == local.status
  ]
}

output "status_filter_is_useful" {
  value = alltrue(local.status_filter_result) && length(local.status_filter_result) > 0
}

locals {
  enabled = data.huaweicloud_ga_health_checks.test.health_checks[0].enabled
}

data "huaweicloud_ga_health_checks" "filter_by_enabled" {
  enabled = local.enabled
}

locals {
  enabled_filter_result = [
    for v in data.huaweicloud_ga_health_checks.filter_by_enabled.health_checks[*].enabled : v == local.enabled
  ]
}

output "enabled_filter_is_useful" {
  value = alltrue(local.enabled_filter_result) && length(local.enabled_filter_result) > 0
}
`, testHealthCheck_basic(name))
}
