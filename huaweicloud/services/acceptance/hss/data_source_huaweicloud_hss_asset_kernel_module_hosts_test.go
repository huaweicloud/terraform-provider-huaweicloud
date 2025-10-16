package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAssetKernelModuleHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_asset_kernel_module_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// This test case requires setting a host ID with host protection enabled.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAssetKernelModuleHosts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.kernel_module_info.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.kernel_module_info.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.kernel_module_info.0.srcversion"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.kernel_module_info.0.path"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.kernel_module_info.0.mode"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.kernel_module_info.0.desc"),

					resource.TestCheckOutput("host_name_filter_useful", "true"),
					resource.TestCheckOutput("host_ip_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceAssetKernelModuleHosts_basic = `
data "huaweicloud_hss_asset_kernel_module_hosts" "test" {
  name = "aesni_intel"
}

locals {
  host_name = data.huaweicloud_hss_asset_kernel_module_hosts.test.data_list[0].host_name
}

data "huaweicloud_hss_asset_kernel_module_hosts" "host_name_filter" {
  name      = "aesni_intel"
  host_name = local.host_name
}

output "host_name_filter_useful" {
  value = length(data.huaweicloud_hss_asset_kernel_module_hosts.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_kernel_module_hosts.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

locals {
  host_ip = data.huaweicloud_hss_asset_kernel_module_hosts.test.data_list[0].host_ip
}

data "huaweicloud_hss_asset_kernel_module_hosts" "host_ip_filter" {	
  name    = "aesni_intel"
  host_ip = local.host_ip
}

output "host_ip_filter_useful" {
  value = length(data.huaweicloud_hss_asset_kernel_module_hosts.host_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_asset_kernel_module_hosts.host_ip_filter.data_list[*].host_ip : v == local.host_ip]
  )
}

data "huaweicloud_hss_asset_kernel_module_hosts" "eps_filter" {
  name                  = "aesni_intel"
  enterprise_project_id = "0"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_asset_kernel_module_hosts.eps_filter.data_list) > 0
}
`
