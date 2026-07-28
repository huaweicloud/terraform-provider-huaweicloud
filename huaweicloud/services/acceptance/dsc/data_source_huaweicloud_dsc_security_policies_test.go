package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceSecurityPolicies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dsc_security_policies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSecurityPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.related_datasource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.related_datasource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.related_datasource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_list.0.related_instance_id"),

					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceSecurityPolicies_basic = `
data "huaweicloud_dsc_security_policies" "test" {}

# Filter by name.
locals {
  name = data.huaweicloud_dsc_security_policies.test.policy_list[0].name
}

data "huaweicloud_dsc_security_policies" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dsc_security_policies.filter_by_name.policy_list[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = alltrue(local.name_filter_result) && length(local.name_filter_result) > 0
}

# Filter by type.
locals {
  type = data.huaweicloud_dsc_security_policies.test.policy_list[0].type
}

data "huaweicloud_dsc_security_policies" "filter_by_type" {
  type = local.type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dsc_security_policies.filter_by_type.policy_list[*].type : v == local.type
  ]
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`
