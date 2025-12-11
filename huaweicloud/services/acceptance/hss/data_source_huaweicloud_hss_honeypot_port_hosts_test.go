package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceHoneypotPortHosts_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_honeypot_port_hosts.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running test, prepare a host with configuration honeypot port policy.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHoneypotPortHosts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),

					resource.TestCheckOutput("host_name_filter_useful", "true"),
					resource.TestCheckOutput("private_ip_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceHoneypotPortHosts_basic = `
data "huaweicloud_hss_honeypot_port_hosts" "test" {}

locals {
  host_name = data.huaweicloud_hss_honeypot_port_hosts.test.data_list[0].host_name
}

data "huaweicloud_hss_honeypot_port_hosts" "host_name_filter" {
  host_name = local.host_name
}

output "host_name_filter_useful" {
  value = length(data.huaweicloud_hss_honeypot_port_hosts.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_honeypot_port_hosts.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

locals {
  private_ip = data.huaweicloud_hss_honeypot_port_hosts.test.data_list[0].private_ip
}

data "huaweicloud_hss_honeypot_port_hosts" "private_ip_filter" {	
  private_ip = local.private_ip
}

output "private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_honeypot_port_hosts.private_ip_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_honeypot_port_hosts.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

data "huaweicloud_hss_honeypot_port_hosts" "eps_filter" {
  enterprise_project_id = "0"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_honeypot_port_hosts.eps_filter.data_list) > 0
}
`
