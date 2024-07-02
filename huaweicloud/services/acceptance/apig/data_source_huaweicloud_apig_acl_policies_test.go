package apig

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAclPolicies_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_apig_acl_policies.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byId   = "data.huaweicloud_apig_acl_policies.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_apig_acl_policies.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_apig_acl_policies.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byEntityType   = "data.huaweicloud_apig_acl_policies.filter_by_entity_type"
		dcByEntityType = acceptance.InitDataSourceCheck(byEntityType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckApigSubResourcesRelatedInfo(t)
			acceptance.TestAccPreCheckApigChannelRelatedInfo(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAclPolicies_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckResourceAttrPair(byId, "policies.0.id", "huaweicloud_apig_acl_policy.test", "id"),
					resource.TestCheckResourceAttrPair(byId, "policies.0.name", "huaweicloud_apig_acl_policy.test", "name"),
					resource.TestCheckResourceAttrPair(byId, "policies.0.type", "huaweicloud_apig_acl_policy.test", "type"),
					resource.TestCheckResourceAttrPair(byId, "policies.0.value", "huaweicloud_apig_acl_policy.test", "value"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.bind_num"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.entity_type"),
					resource.TestCheckResourceAttrSet(byId, "policies.0.updated_at"),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					dcByEntityType.CheckResourceExists(),
					resource.TestCheckOutput("is_entity_type_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceAclPolicies_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_apig_acl_policies" "test" {
  depends_on = [
    huaweicloud_apig_acl_policy.test
  ]

  instance_id = local.instance_id
}

# Filter by ID
locals {
  policy_id = huaweicloud_apig_acl_policy.test.id
}

data "huaweicloud_apig_acl_policies" "filter_by_id" {
  instance_id = local.instance_id
  policy_id   = local.policy_id
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_apig_acl_policies.filter_by_id.policies[*].id : v == local.policy_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) > 0 && alltrue(local.id_filter_result)
}

# Filter by name
locals {
  policy_name = huaweicloud_apig_acl_policy.test.name
}

data "huaweicloud_apig_acl_policies" "filter_by_name" {
  depends_on = [huaweicloud_apig_acl_policy.test]

  instance_id = local.instance_id
  name        = local.policy_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_apig_acl_policies.filter_by_name.policies[*].name : v == local.policy_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by type
locals {
  policy_type = huaweicloud_apig_acl_policy.test.type
}

data "huaweicloud_apig_acl_policies" "filter_by_type" {
  depends_on = [huaweicloud_apig_acl_policy.test]

  instance_id = local.instance_id
  type        = local.policy_type
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_apig_acl_policies.filter_by_type.policies[*].type : v == local.policy_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by entity type
locals {
  entity_type = huaweicloud_apig_acl_policy.test.entity_type
}

data "huaweicloud_apig_acl_policies" "filter_by_entity_type" {
  depends_on = [huaweicloud_apig_acl_policy.test]

  instance_id = local.instance_id
  entity_type = local.entity_type
}

locals {
  entity_type_filter_result = [
    for v in data.huaweicloud_apig_acl_policies.filter_by_entity_type.policies[*].entity_type : v == local.entity_type
  ]
}

output "is_entity_type_filter_useful" {
  value = length(local.entity_type_filter_result) > 0 && alltrue(local.entity_type_filter_result)
}

`, testAccDataSourceApiAssociatedAclPolicies_base())
}
