package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceContainerNodes_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_container_nodes.test"
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
				Config: testDataSourceContainerNodes_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.host_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.agent_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protect_status"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protect_interrupt"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.protect_degradation"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.container_tags"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.private_ip"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.enterprise_project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.detect_result"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.asset"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.vulnerability"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.intrusion"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_group_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_group_name"),

					resource.TestCheckOutput("is_host_name_filter_useful", "true"),
					resource.TestCheckOutput("is_agent_status_filter_useful", "true"),
					resource.TestCheckOutput("is_protect_status_filter_useful", "true"),
					resource.TestCheckOutput("is_container_tags_filter_useful", "true"),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
				),
			},
		},
	})
}

const testDataSourceContainerNodes_basic string = `
data "huaweicloud_hss_container_nodes" "test" {}

# Filter using host_name.
locals {
  host_name = data.huaweicloud_hss_container_nodes.test.data_list[0].host_name
}

data "huaweicloud_hss_container_nodes" "host_name_filter" {
  host_name = local.host_name
}

output "is_host_name_filter_useful" {
  value = length(data.huaweicloud_hss_container_nodes.host_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_nodes.host_name_filter.data_list[*].host_name : v == local.host_name]
  )
}

# Filter using agent_status.
locals {
  agent_status = data.huaweicloud_hss_container_nodes.test.data_list[0].agent_status
}

data "huaweicloud_hss_container_nodes" "agent_status_filter" {
  agent_status = local.agent_status
}

output "is_agent_status_filter_useful" {
  value = length(data.huaweicloud_hss_container_nodes.agent_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_nodes.agent_status_filter.data_list[*].agent_status : v == local.agent_status]
  )
}

# Filter using protect_status.
locals {
  protect_status = data.huaweicloud_hss_container_nodes.test.data_list[0].protect_status
}

data "huaweicloud_hss_container_nodes" "protect_status_filter" {
  protect_status = local.protect_status
}

output "is_protect_status_filter_useful" {
  value = length(data.huaweicloud_hss_container_nodes.protect_status_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_nodes.protect_status_filter.data_list[*].protect_status : v == local.protect_status]
  )
}

# Filter using container_tags.
locals {
  container_tags = data.huaweicloud_hss_container_nodes.test.data_list[0].container_tags
}

data "huaweicloud_hss_container_nodes" "container_tags_filter" {
  container_tags = local.container_tags
}

output "is_container_tags_filter_useful" {
  value = length(data.huaweicloud_hss_container_nodes.container_tags_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_container_nodes.container_tags_filter.data_list[*].container_tags : v == local.container_tags]
  )
}

# Filter using enterprise_project_id.
data "huaweicloud_hss_container_nodes" "enterprise_project_id_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_enterprise_project_id_filter_useful" {
  value = length(data.huaweicloud_hss_container_nodes.container_tags_filter.data_list) > 0
}

# Filter using non existent host_name
data "huaweicloud_hss_container_nodes" "not_found" {
  host_name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_container_nodes.not_found.data_list) == 0
}
`
