package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSettingPlugins_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_setting_plugins.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSettingPlugins_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),

					resource.TestCheckOutput("host_name_filter_useful", "true"),
					resource.TestCheckOutput("host_id_filter_useful", "true"),
					resource.TestCheckOutput("private_ip_filter_useful", "true"),
					resource.TestCheckOutput("os_type_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSettingPlugins_basic = `
data "huaweicloud_hss_setting_plugins" "test" {
  name = "opa-docker-authz"
}

locals {
  host_name = data.huaweicloud_hss_setting_plugins.test.data_list[0].host_name
  host_id   = data.huaweicloud_hss_setting_plugins.test.data_list[0].host_id
}

data "huaweicloud_hss_setting_plugins" "host_name_filter" {
  name      = "opa-docker-authz"
  host_name = local.host_name
}

output "host_name_filter_useful" {
  value = length(data.huaweicloud_hss_setting_plugins.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_plugins.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

data "huaweicloud_hss_setting_plugins" "host_id_filter" {
  name    = "opa-docker-authz"
  host_id = local.host_id
}

output "host_id_filter_useful" {
  value = length(data.huaweicloud_hss_setting_plugins.host_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_plugins.host_id_filter.data_list[*].host_id : v == local.host_id]
  )
}

locals {
  private_ip = data.huaweicloud_hss_setting_plugins.test.data_list[0].private_ip
  os_type    = data.huaweicloud_hss_setting_plugins.test.data_list[0].os_type
}

data "huaweicloud_hss_setting_plugins" "private_ip_filter" {
  name       = "opa-docker-authz"
  private_ip = local.private_ip
}

output "private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_setting_plugins.private_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_setting_plugins.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

data "huaweicloud_hss_setting_plugins" "os_type_filter" {
  name    = "opa-docker-authz"
  os_type = local.os_type
}

output "os_type_filter_useful" {
  value = length(data.huaweicloud_hss_setting_plugins.os_type_filter.data_list) > 0 && alltrue(	
    [for v in data.huaweicloud_hss_setting_plugins.os_type_filter.data_list[*].os_type : v == local.os_type]
  )
}

data "huaweicloud_hss_setting_plugins" "eps_filter" {
  name                  = "opa-docker-authz"
  enterprise_project_id = "0"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_setting_plugins.eps_filter.data_list) > 0
}
`
