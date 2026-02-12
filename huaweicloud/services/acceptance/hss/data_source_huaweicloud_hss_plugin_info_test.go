package hss

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePluginInfo_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_hss_plugin_info.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourcePluginInfo_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "total_num"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.agent_version"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.arch"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.size"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.cpu_limit"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.memory_limit"),
					resource.TestCheckResourceAttrSet(dataSourceName, "data_list.0.update_time"),

					resource.TestCheckOutput("is_agent_version_filter_useful", "true"),
					resource.TestCheckOutput("is_plugin_arch_filter_useful", "true"),
					resource.TestCheckOutput("is_plugin_os_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourcePluginInfo_base() string {
	return `
data "huaweicloud_hss_plugins" "test" {}
`
}

func testDataSourcePluginInfo_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_hss_plugin_info" "test" {
  code                  = data.huaweicloud_hss_plugins.test.data_list[0].code
  enterprise_project_id = "all_granted_eps"
}

# Filter using agent_version.
locals {
  agent_version = data.huaweicloud_hss_plugin_info.test.data_list[0].agent_version
}

data "huaweicloud_hss_plugin_info" "agent_version_filter" {
  code          = data.huaweicloud_hss_plugins.test.data_list[0].code
  agent_version = local.agent_version
}

output "is_agent_version_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_info.agent_version_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_info.agent_version_filter.data_list[*].agent_version : v == local.agent_version]
  )
}

# Filter using plugin_arch.
locals {
  plugin_arch = data.huaweicloud_hss_plugin_info.test.data_list[0].arch
}

data "huaweicloud_hss_plugin_info" "plugin_arch_filter" {
  code        = data.huaweicloud_hss_plugins.test.data_list[0].code
  plugin_arch = local.plugin_arch
}

output "is_plugin_arch_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_info.plugin_arch_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_info.plugin_arch_filter.data_list[*].arch : v == local.plugin_arch]
  )
}

# Filter using plugin_os_type.
locals {
  plugin_os_type = data.huaweicloud_hss_plugin_info.test.data_list[0].os_type
}

data "huaweicloud_hss_plugin_info" "plugin_os_type_filter" {
  code           = data.huaweicloud_hss_plugins.test.data_list[0].code
  plugin_os_type = local.plugin_os_type
}

output "is_plugin_os_type_filter_useful" {
  value = length(data.huaweicloud_hss_plugin_info.plugin_os_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_plugin_info.plugin_os_type_filter.data_list[*].os_type : v == local.plugin_os_type]
  )
}
`, testDataSourcePluginInfo_base())
}
