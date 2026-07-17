package dsc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDscDeviceSecurityPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dsc_device_security_policies.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDscDeviceSecurityPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.update_time"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.related_datasource_id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.related_datasource_name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.related_datasource_type"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.related_device_id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.ddm_config.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.ddm_policy_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.gde_config.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.gde_policy.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.sdm_config.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.sdm_policy_list.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.resource.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.target_resource.#"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testAccDataSourceDscDeviceSecurityPolicies_basic = `
data "huaweicloud_dsc_device_security_policies" "test" {}

# Filter by type
locals {
  policy_type = data.huaweicloud_dsc_device_security_policies.test.policies[0].type
}

data "huaweicloud_dsc_device_security_policies" "filter_by_type" {
  type = local.policy_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_dsc_device_security_policies.filter_by_type.policies[*].type : v == local.policy_type
  ]
}

output "type_filter_is_useful" {
  value = alltrue(local.type_filter_result) && length(local.type_filter_result) > 0
}
`
