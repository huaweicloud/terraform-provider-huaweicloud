package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSecurityResourcePermissionPolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_dataarts_security_resource_permission_policies.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byPolicyName   = "data.huaweicloud_dataarts_security_resource_permission_policies.filter_by_policy_name"
		dcByPolicyName = acceptance.InitDataSourceCheck(byPolicyName)

		byResourceName   = "data.huaweicloud_dataarts_security_resource_permission_policies.filter_by_resource_name"
		dcByResourceName = acceptance.InitDataSourceCheck(byResourceName)

		byMemberName   = "data.huaweicloud_dataarts_security_resource_permission_policies.filter_by_member_name"
		dcByMemberName = acceptance.InitDataSourceCheck(byMemberName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
			acceptance.TestAccPreCheckDataArtsConnectionID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSecurityResourcePermissionPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					// Filter by 'policy_name' parameter.
					dcByPolicyName.CheckResourceExists(),
					resource.TestCheckOutput("is_policy_name_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(byPolicyName, "policies.0.policy_id"),
					resource.TestCheckResourceAttrSet(byPolicyName, "policies.0.policy_name"),
					resource.TestMatchResourceAttr(byPolicyName, "policies.0.create_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(byPolicyName, "policies.0.create_user"),
					resource.TestMatchResourceAttr(byPolicyName, "policies.0.update_time",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					// Filter by 'resource_name' parameter.
					dcByResourceName.CheckResourceExists(),
					resource.TestCheckOutput("is_resource_name_filter_useful", "true"),
					// Filter by 'member_name' parameter.
					dcByMemberName.CheckResourceExists(),
					resource.TestCheckOutput("is_member_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSecurityResourcePermissionPolicies_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = "%[1]s"
}

data "huaweicloud_dataarts_studio_data_connections" "test" {
  workspace_id  = "%[1]s"
  connection_id = "%[2]s"
}

locals {
  connection_name = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].name, "NOT_FOUND")
  member_name     = try(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles[0].name, "NOT_FOUND")
}

resource "huaweicloud_dataarts_security_resource_permission_policy" "test" {
  workspace_id = "%[1]s"
  name         = "%[3]s"

  resources {
    resource_id   = try(data.huaweicloud_dataarts_studio_data_connections.test.connections[0].id, "NOT_FOUND")
    resource_name = local.connection_name
    resource_type = "DATA_CONNECTION"
  }

  members {
    member_id   = try(data.huaweicloud_dataarts_studio_workspace_user_roles.test.roles[0].id, "NOT_FOUND")
    member_name = local.member_name
    member_type = "WORKSPACE_ROLE"
  }
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID, acceptance.HW_DATAARTS_CONNECTION_ID, name)
}

func testAccDataSecurityResourcePermissionPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_dataarts_security_resource_permission_policies" "test" {
  workspace_id = "%[2]s"

  depends_on = [huaweicloud_dataarts_security_resource_permission_policy.test]
}

# Filter by 'policy_name' parameter.
locals {
  policy_name = huaweicloud_dataarts_security_resource_permission_policy.test.name
}

data "huaweicloud_dataarts_security_resource_permission_policies" "filter_by_policy_name" {
  workspace_id = "%[2]s"
  policy_name  = local.policy_name

  depends_on = [huaweicloud_dataarts_security_resource_permission_policy.test]
}

locals {
  policy_name_filter_result = [
    for v in data.huaweicloud_dataarts_security_resource_permission_policies.filter_by_policy_name.policies[*].policy_name
    : strcontains(v, local.policy_name)
  ]
}

output "is_policy_name_filter_useful" {
  value = length(local.policy_name_filter_result) > 0 && alltrue(local.policy_name_filter_result)
}

# Filter by 'resource_name' parameter.
data "huaweicloud_dataarts_security_resource_permission_policies" "filter_by_resource_name" {
  workspace_id  = "%[2]s"
  resource_name = local.connection_name

  depends_on = [huaweicloud_dataarts_security_resource_permission_policy.test]
}

locals {
  resource_name_filter_result = [for item in
    [for v in data.huaweicloud_dataarts_security_resource_permission_policies.filter_by_resource_name.policies : v.resources[*].resource_name] :
    contains(item, local.connection_name)
  ]
}

output "is_resource_name_filter_useful" {
  value = length(local.resource_name_filter_result) > 0 && alltrue(local.resource_name_filter_result)
}

# Filter by 'member_name' parameter.
data "huaweicloud_dataarts_security_resource_permission_policies" "filter_by_member_name" {
  workspace_id = "%[2]s"
  member_name  = local.member_name

  depends_on = [huaweicloud_dataarts_security_resource_permission_policy.test]
}

locals {
  member_name_filter_result = [
    for item in [for v in data.huaweicloud_dataarts_security_resource_permission_policies.filter_by_member_name.policies : v.members[*].member_name]
    : contains(item, local.member_name)
  ]
}

output "is_member_name_filter_useful" {
  value = length(local.member_name_filter_result) > 0 && alltrue(local.member_name_filter_result)
}
`, testAccDataSecurityResourcePermissionPolicies_base(name), acceptance.HW_DATAARTS_WORKSPACE_ID)
}
