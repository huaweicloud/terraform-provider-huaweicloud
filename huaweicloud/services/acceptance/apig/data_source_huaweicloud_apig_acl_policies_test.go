package apig

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAclPolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_apig_acl_policies.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAclPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "policies.0.entity_type"),

					resource.TestCheckOutput("policy_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("entity_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAclPolicies_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_apig_acl_policies" "test" {
  depends_on = [
    huaweicloud_apig_acl_policy.test
  ]

  instance_id = huaweicloud_apig_instance.test.id
}

data "huaweicloud_apig_acl_policies" "policy_id_filter" {
  instance_id = huaweicloud_apig_instance.test.id
  policy_id   = local.policy_id
}
  
locals {
  policy_id = data.huaweicloud_apig_acl_policies.test.policies[0].id
}
  
output "policy_id_filter_is_useful" {
  value = length(data.huaweicloud_apig_acl_policies.policy_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_apig_acl_policies.policy_id_filter.policies[*].id : v == local.policy_id]
  )
}

data "huaweicloud_apig_acl_policies" "name_filter" {
  instance_id = huaweicloud_apig_instance.test.id
  name        = local.name
}
  
locals {
  name = data.huaweicloud_apig_acl_policies.test.policies[0].name
}
  
output "name_filter_is_useful" {
  value = length(data.huaweicloud_apig_acl_policies.name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_apig_acl_policies.name_filter.policies[*].name : v == local.name]
  )
}

data "huaweicloud_apig_acl_policies" "type_filter" {
  instance_id = huaweicloud_apig_instance.test.id
  type        = local.type
}
  
locals {
  type = data.huaweicloud_apig_acl_policies.test.policies[0].type
}
  
output "type_filter_is_useful" {
  value = length(data.huaweicloud_apig_acl_policies.type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_apig_acl_policies.type_filter.policies[*].type : v == local.type]
  )
}

data "huaweicloud_apig_acl_policies" "entity_type_filter" {
  instance_id = huaweicloud_apig_instance.test.id
  entity_type = local.entity_type
}
  
locals {
  entity_type = data.huaweicloud_apig_acl_policies.test.policies[0].entity_type
}
  
output "entity_type_filter_is_useful" {
  value = length(data.huaweicloud_apig_acl_policies.entity_type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_apig_acl_policies.entity_type_filter.policies[*].entity_type : v == local.entity_type]
  )
}
`, testAccDataSourceApiAssociatedAclPolicies_base())
}
