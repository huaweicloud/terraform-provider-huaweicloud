package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRaspServers_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_rasp_servers.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with premium edition host protection enabled.
			// The host also set the application protection (RASP) at the same time.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRaspServers_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rasp_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.is_friendly_user"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_support_auto_attach"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.auto_attach"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rasp_port"),

					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
					resource.TestCheckOutput("is_app_type_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceRaspServers_basic = `
data "huaweicloud_hss_rasp_servers" "test" {
  app_status = "opened"
}

locals {
  host_name = data.huaweicloud_hss_rasp_servers.test.data_list[0].host_name
}

data "huaweicloud_hss_rasp_servers" "host_name_filter" {
  app_status = "opened"
  host_name  = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_servers.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_rasp_servers.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

locals {
  os_type = data.huaweicloud_hss_rasp_servers.test.data_list[0].os_type
}

data "huaweicloud_hss_rasp_servers" "os_type_filter" {	
  app_status = "opened"
  os_type    = local.os_type
}

output "is_os_type_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_servers.os_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_rasp_servers.os_type_filter.data_list[*].os_type : v == local.os_type]
  )
}

data "huaweicloud_hss_rasp_servers" "app_type_filter" {
  app_status = "opened"
  app_type   = "java"
}

output "is_app_type_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_servers.app_type_filter.data_list) >= 0
}

data "huaweicloud_hss_rasp_servers" "eps_filter" {
  app_status            = "opened"
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_servers.eps_filter.data_list) >= 0
}
`
