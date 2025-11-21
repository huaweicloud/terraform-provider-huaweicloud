package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRansomwareProtectionServers_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_ransomware_protection_servers.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Enable enterprise edition host protection and enable ransomware protection.
			acceptance.TestAccPreCheckHSSHostProtectionHostId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRansomwareProtectionServers_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.ransom_protection_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protect_policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_version"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_source"),

					resource.TestCheckOutput("host_name_filter_useful", "true"),
					resource.TestCheckOutput("host_id_filter_useful", "true"),
					resource.TestCheckOutput("os_type_filter_useful", "true"),
					resource.TestCheckOutput("private_ip_filter_useful", "true"),
					resource.TestCheckOutput("host_status_filter_useful", "true"),
					resource.TestCheckOutput("ransom_protection_status_filter_useful", "true"),
					resource.TestCheckOutput("protect_policy_name_filter_useful", "true"),
					resource.TestCheckOutput("policy_id_filter_useful", "true"),
					resource.TestCheckOutput("agent_status_filter_useful", "true"),
					resource.TestCheckOutput("eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceRansomwareProtectionServers_basic = `
data "huaweicloud_hss_ransomware_protection_servers" "test" {}

locals {
  host_name = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].host_name
  host_id   = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].host_id
  os_type   = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].os_type
}

data "huaweicloud_hss_ransomware_protection_servers" "host_name_filter" {	
  host_name = local.host_name
}

output "host_name_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_servers.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

data "huaweicloud_hss_ransomware_protection_servers" "host_id_filter" {
  host_id = local.host_id
}

output "host_id_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.host_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_servers.host_id_filter.data_list[*].host_id : v == local.host_id]
  )
}

data "huaweicloud_hss_ransomware_protection_servers" "os_type_filter" {	
  os_type = local.os_type
}

output "os_type_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.os_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_servers.os_type_filter.data_list[*].os_type : v == local.os_type]
  )
}

locals {
  private_ip    = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].private_ip
	host_status   = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].host_status
  ransom_status = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].ransom_protection_status
}

data "huaweicloud_hss_ransomware_protection_servers" "private_ip_filter" {
  private_ip = local.private_ip
}

output "private_ip_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.private_ip_filter.data_list) > 0 && alltrue(	
    [for v in data.huaweicloud_hss_ransomware_protection_servers.private_ip_filter.data_list[*].private_ip : v == local.private_ip]
  )
}

data "huaweicloud_hss_ransomware_protection_servers" "host_status_filter" {
  host_status = local.host_status
}

output "host_status_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.host_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_servers.host_status_filter.data_list[*].host_status : v == local.host_status]
  )
}

data "huaweicloud_hss_ransomware_protection_servers" "ransom_protection_status_filter" {
  ransom_protection_status = local.ransom_status
}

output "ransom_protection_status_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.ransom_protection_status_filter.data_list) > 0 && alltrue(	
    [for v in data.huaweicloud_hss_ransomware_protection_servers.ransom_protection_status_filter.data_list[*].ransom_protection_status : v
      == local.ransom_status]
  )
}

locals {
  policy_name  = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].protect_policy_name
  policy_id    = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].protect_policy_id
  agent_status = data.huaweicloud_hss_ransomware_protection_servers.test.data_list[0].agent_status
}

data "huaweicloud_hss_ransomware_protection_servers" "protect_policy_name_filter" {
  protect_policy_name = local.policy_name
}

output "protect_policy_name_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.protect_policy_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_servers.protect_policy_name_filter.data_list[*].protect_policy_name : v
      == local.policy_name]
  )
}

data "huaweicloud_hss_ransomware_protection_servers" "policy_id_filter" {	
  policy_id = local.policy_id
}

output "policy_id_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.policy_id_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_servers.policy_id_filter.data_list[*].protect_policy_id : v == local.policy_id]
  )
}

data "huaweicloud_hss_ransomware_protection_servers" "agent_status_filter" {
  agent_status = local.agent_status
}

output "agent_status_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.agent_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_servers.agent_status_filter.data_list[*].agent_status : v == local.agent_status]
  )
}

data "huaweicloud_hss_ransomware_protection_servers" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "eps_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_servers.eps_filter.data_list) > 0
}
`
