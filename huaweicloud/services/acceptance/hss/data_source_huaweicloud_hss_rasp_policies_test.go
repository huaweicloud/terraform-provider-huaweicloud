package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRaspPolicies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_rasp_policies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Before running the test case, prepare an HSS protection policy.
			acceptance.TestAccPreCheckHSSPolicyId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRaspPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_id"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.policy_name"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.os_type"),
					resource.TestCheckResourceAttrSet(dataSource, "data_list.0.rule_name"),

					resource.TestCheckOutput("is_policy_name_filter_useful", "true"),
					resource.TestCheckOutput("is_os_type_filter_useful", "true"),
					resource.TestCheckOutput("is_eps_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceRaspPolicies_basic = `
data "huaweicloud_hss_rasp_policies" "test" {}

locals {
  policy_name = data.huaweicloud_hss_rasp_policies.test.data_list[0].policy_name
}

data "huaweicloud_hss_rasp_policies" "policy_name_filter" {
  policy_name = local.policy_name
}

output "is_policy_name_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_policies.policy_name_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_rasp_policies.policy_name_filter.data_list[*].policy_name : v == local.policy_name]
  )
}

locals {
  os_type = data.huaweicloud_hss_rasp_policies.test.data_list[0].os_type
}

data "huaweicloud_hss_rasp_policies" "os_type_filter" {
  os_type = local.os_type
}

output "is_os_type_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_policies.os_type_filter.data_list) > 0 && alltrue(
    [for v in data.huaweicloud_hss_rasp_policies.os_type_filter.data_list[*].os_type : v == local.os_type]
  )
}

data "huaweicloud_hss_rasp_policies" "eps_filter" {
  enterprise_project_id = "all_granted_eps"
}

output "is_eps_filter_useful" {
  value = length(data.huaweicloud_hss_rasp_policies.eps_filter.data_list) > 0
}
`
