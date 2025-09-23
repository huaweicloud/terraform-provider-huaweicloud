package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceHosts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.agent_status"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.protect_version"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.protect_charging_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.quota_id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.detect_result"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.policy_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.asset_value"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.open_time"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "hosts.0.enterprise_project_id"),

					resource.TestCheckOutput("is_host_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
					resource.TestCheckOutput("is_agent_status_filter_useful", "true"),
					resource.TestCheckOutput("is_protect_status_status_filter_useful", "true"),
					resource.TestCheckOutput("is_protect_version_filter_useful", "true"),
					resource.TestCheckOutput("is_protect_charging_mode_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceHosts_basic string = `
data "huaweicloud_hss_hosts" "test" {}

# Filter using host ID.
locals {
  host_id = data.huaweicloud_hss_hosts.test.hosts[0].id
}

data "huaweicloud_hss_hosts" "host_id_filter" {
  host_id = local.host_id
}

output "is_host_id_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.host_id_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_hosts.host_id_filter.hosts[*].id : v == local.host_id]
  )
}

# Filter using name.
locals {
  name = data.huaweicloud_hss_hosts.test.hosts[0].name
}

data "huaweicloud_hss_hosts" "name_filter" {
  name = local.name
}

# The name parameter is a fuzzy match, to prevent test case error, only check the length of the response list.
output "is_name_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.name_filter.hosts) > 0
}

# Filter using status
locals {
  status = data.huaweicloud_hss_hosts.test.hosts[0].status
}

data "huaweicloud_hss_hosts" "status_filter" {
  status = local.status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.status_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_hosts.status_filter.hosts[*].status : v == local.status]
  )
}

# Filter using os_type
locals {
  os_type = data.huaweicloud_hss_hosts.test.hosts[0].os_type
}

data "huaweicloud_hss_hosts" "os_type_filter" {
  os_type = local.os_type
}

output "is_os_type_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.os_type_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_hosts.os_type_filter.hosts[*].os_type : v == local.os_type]
  )
}

# Filter using agent_status
locals {
  agent_status = data.huaweicloud_hss_hosts.test.hosts[0].agent_status
}

data "huaweicloud_hss_hosts" "agent_status_filter" {
  agent_status = local.agent_status
}

output "is_agent_status_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.agent_status_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_hosts.agent_status_filter.hosts[*].agent_status : v == local.agent_status]
  )
}

# Filter using protect_status
locals {
  protect_status = data.huaweicloud_hss_hosts.test.hosts[0].protect_status
}

data "huaweicloud_hss_hosts" "protect_status_filter" {
  protect_status = local.protect_status
}

output "is_protect_status_status_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.protect_status_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_hosts.protect_status_filter.hosts[*].protect_status : v == local.protect_status]
  )
}

# Filter using protect_version
locals {
  protect_version = data.huaweicloud_hss_hosts.test.hosts[0].protect_version
}

data "huaweicloud_hss_hosts" "protect_version_filter" {
  protect_version = local.protect_version
}

output "is_protect_version_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.protect_version_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_hosts.protect_version_filter.hosts[*].protect_version : v == local.protect_version]
  )
}

# Filter using protect_charging_mode
locals {
  protect_charging_mode = data.huaweicloud_hss_hosts.test.hosts[0].protect_charging_mode
}

data "huaweicloud_hss_hosts" "protect_charging_mode_filter" {
  protect_charging_mode = local.protect_charging_mode
}

output "is_protect_charging_mode_filter_useful" {
  value = length(data.huaweicloud_hss_hosts.protect_charging_mode_filter.hosts) > 0 && alltrue(
    [for v in data.huaweicloud_hss_hosts.protect_charging_mode_filter.hosts[*].protect_charging_mode : v == local.protect_charging_mode]
  )
}

# Filter using non existent name.
data "huaweicloud_hss_hosts" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_hosts.not_found.hosts) == 0
}
`
