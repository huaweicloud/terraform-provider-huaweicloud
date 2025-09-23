package hss

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// There will always be a default policy that will not be deleted,
// So there is no need to create resources in this acc test.
func TestAccDataSourceRansomwareProtectionPolicies_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_hss_ransomware_protection_policies.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRansomwareProtectionPolicies_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.#"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.protection_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.bait_protection_status"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.deploy_mode"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.protection_directory"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.protection_type"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.runtime_detection_status"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.count_associated_server"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.operating_system"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.default_policy"),

					resource.TestCheckOutput("is_policy_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("not_found_validation_pass", "true"),
					resource.TestCheckOutput("is_operating_system_filter_useful", "true"),
				),
			},
		},
	})
}

const testDataSourceRansomwareProtectionPolicies_basic string = `
data "huaweicloud_hss_ransomware_protection_policies" "test" {}

# Filter using policy ID.
locals {
  policy_id = data.huaweicloud_hss_ransomware_protection_policies.test.policies[0].id
}

data "huaweicloud_hss_ransomware_protection_policies" "policy_id_filter" {
  policy_id = local.policy_id
}

output "is_policy_id_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_policies.policy_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_policies.policy_id_filter.policies[*].id : v == local.policy_id]
  )
}

# Filter using name.
locals {
  name = data.huaweicloud_hss_ransomware_protection_policies.test.policies[0].name
}

data "huaweicloud_hss_ransomware_protection_policies" "name_filter" {
  name = local.name
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_policies.name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_policies.name_filter.policies[*].name : v == local.name]
  )
}

# Filter using non existent name.
data "huaweicloud_hss_ransomware_protection_policies" "not_found" {
  name = "resource_not_found"
}

output "not_found_validation_pass" {
  value = length(data.huaweicloud_hss_ransomware_protection_policies.not_found.policies) == 0
}

# Filter using operating_system.
locals {
  operating_system = data.huaweicloud_hss_ransomware_protection_policies.test.policies[0].operating_system
}

data "huaweicloud_hss_ransomware_protection_policies" "operating_system_filter" {
  operating_system = local.operating_system
}

output "is_operating_system_filter_useful" {
  value = length(data.huaweicloud_hss_ransomware_protection_policies.operating_system_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_hss_ransomware_protection_policies.operating_system_filter.policies[*].operating_system : v == local.operating_system]
  )
}
`
