package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePluginAttachments_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_plugin_attachments.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePluginAttachments_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.host_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.agent_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.asset_value"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.plugin_status"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.protect_status"),

					resource.TestCheckOutput("is_host_ids_filter_useful", "true"),
					resource.TestCheckOutput("is_host_status_filter_useful", "true"),
					resource.TestCheckOutput("is_plugin_version_filter_useful", "true"),
					resource.TestCheckOutput("is_plugin_status_filter_useful", "true"),
					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_agent_status_filter_useful", "true"),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
					resource.TestCheckOutput("is_host_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePluginAttachments_base() string {
	return `
data "huaweicloud_hss_plugins" "test" {}
`
}

func testDataSourcePluginAttachments_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_plugin_attachments" "test" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
}

# Filter using host_ids.
locals {
  host_id = data.huaweicloud_hss_plugin_attachments.test.data_list[0].host_id
}

data "huaweicloud_hss_plugin_attachments" "host_ids_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  host_ids              = [local.host_id]
}

output "is_host_ids_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.host_ids_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.host_ids_filter.data_list[*].host_id : v == local.host_id]
  )
}

# Filter using host_status.
locals {
  host_status = data.huaweicloud_hss_plugin_attachments.test.data_list[0].host_status
}

data "huaweicloud_hss_plugin_attachments" "host_status_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  host_status           = [local.host_status]
}

output "is_host_status_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.host_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.host_status_filter.data_list[*].host_status : v == local.host_status]
  )
}

# Filter using plugin_version.
locals {
  plugin_version = data.huaweicloud_hss_plugin_attachments.test.data_list[0].plugin_version
}

data "huaweicloud_hss_plugin_attachments" "plugin_version_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  plugin_version        = local.plugin_version
}

output "is_plugin_version_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.plugin_version_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.plugin_version_filter.data_list[*].plugin_version : v == local.plugin_version]
  )
}

# Filter using plugin_status.
locals {
  plugin_status = data.huaweicloud_hss_plugin_attachments.test.data_list[0].plugin_status
}

data "huaweicloud_hss_plugin_attachments" "plugin_status_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  plugin_status         = local.plugin_status
}

output "is_plugin_status_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.plugin_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.plugin_status_filter.data_list[*].plugin_status : v == local.plugin_status]
  )
}

# Filter using host_name.
locals {
  host_name = data.huaweicloud_hss_plugin_attachments.test.data_list[0].host_name
}

data "huaweicloud_hss_plugin_attachments" "host_name_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  host_name             = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

# Filter using agent_status.
locals {
  agent_status = data.huaweicloud_hss_plugin_attachments.test.data_list[0].agent_status
}

data "huaweicloud_hss_plugin_attachments" "agent_status_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  agent_status          = local.agent_status
}

output "is_agent_status_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.agent_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.agent_status_filter.data_list[*].agent_status : v == local.agent_status]
  )
}

# Filter using os_type.
locals {
  os_type = data.huaweicloud_hss_plugin_attachments.test.data_list[0].os_type
}

data "huaweicloud_hss_plugin_attachments" "os_type_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  os_type               = local.os_type
}

output "is_os_type_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.os_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.os_type_filter.data_list[*].os_type : v == local.os_type]
  )
}

# Filter using host_type.
locals {
  host_type = data.huaweicloud_hss_plugin_attachments.test.data_list[0].host_type
}

data "huaweicloud_hss_plugin_attachments" "host_type_filter" {
  plugin_code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
  host_type             = local.host_type
}

output "is_host_type_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_attachments.host_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_attachments.host_type_filter.data_list[*].host_type : v == local.host_type]
  )
}
`, testDataSourcePluginAttachments_base())
}
