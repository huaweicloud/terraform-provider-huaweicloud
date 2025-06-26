package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourcePolicyGroups_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_workspace_policy_groups.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterById   = "data.huaweicloud_workspace_policy_groups.filter_by_policy_group_id"
		dcFilterById = acceptance.InitDataSourceCheck(filterById)

		filterByName   = "data.huaweicloud_workspace_policy_groups.filter_by_policy_group_name"
		dcFilterByName = acceptance.InitDataSourceCheck(filterByName)

		filterByPriority   = "data.huaweicloud_workspace_policy_groups.filter_by_priority"
		dcFilterByPriority = acceptance.InitDataSourceCheck(filterByPriority)

		filterByDescription   = "data.huaweicloud_workspace_policy_groups.filter_by_description"
		dcFilterByDescription = acceptance.InitDataSourceCheck(filterByDescription)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourcePolicyGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Query policy groups without any filter parameter
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "policy_groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.policy_group_id"),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.policy_group_name"),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.priority"),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.update_time"),
					// Filter by ID
					dcFilterById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					// Filter by name
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// Filter by priority
					dcFilterByPriority.CheckResourceExists(),
					resource.TestCheckOutput("is_priority_filter_useful", "true"),
					// Filter by description
					dcFilterByDescription.CheckResourceExists(),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourcePolicyGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_policy_groups" "all" {
  depends_on = [
    huaweicloud_workspace_policy_group.test,
    huaweicloud_workspace_policy_group.nontest,
  ]
}

locals {
  policy_group_id   = try(data.huaweicloud_workspace_policy_groups.all.policy_groups[0].policy_group_id, "NOT_FOUND")
  policy_group_name = try(data.huaweicloud_workspace_policy_groups.all.policy_groups[0].policy_group_name, "NOT_FOUND")
  priority          = try(data.huaweicloud_workspace_policy_groups.all.policy_groups[0].priority, -1)
  update_time       = try(data.huaweicloud_workspace_policy_groups.all.policy_groups[0].update_time, "1900-01-01T01:01:01Z")
  description       = try(data.huaweicloud_workspace_policy_groups.all.policy_groups[0].description, "NOT_FOUND")
}

# Filter by policy group id
data "huaweicloud_workspace_policy_groups" "filter_by_policy_group_id" {
  policy_group_id = local.policy_group_id

  depends_on = [
    huaweicloud_workspace_policy_group.test,
    huaweicloud_workspace_policy_group.nontest,
  ]
}

locals {
  id_filter_result = [
    for v in data.huaweicloud_workspace_policy_groups.filter_by_policy_group_id.policy_groups[*].policy_group_id :
    v == local.policy_group_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.id_filter_result) < 2 && alltrue(local.id_filter_result)
}

# Filter by policy group name
data "huaweicloud_workspace_policy_groups" "filter_by_policy_group_name" {
  policy_group_name = local.policy_group_name

  depends_on = [
    huaweicloud_workspace_policy_group.test,
    huaweicloud_workspace_policy_group.nontest,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_workspace_policy_groups.filter_by_policy_group_name.policy_groups) > 0
}

# Filter by priority
data "huaweicloud_workspace_policy_groups" "filter_by_priority" {
  priority = local.priority

  depends_on = [
    huaweicloud_workspace_policy_group.test,
    huaweicloud_workspace_policy_group.nontest,
  ]
}

locals {
  priority_filter_result = [
    for v in data.huaweicloud_workspace_policy_groups.filter_by_priority.policy_groups[*].priority :
    v == local.priority
  ]
}

output "is_priority_filter_useful" {
  value = length(local.priority_filter_result) > 0 && alltrue(local.priority_filter_result)
}

# Filter by description
data "huaweicloud_workspace_policy_groups" "filter_by_description" {
  description = local.description

  depends_on = [
    huaweicloud_workspace_policy_group.test,
    huaweicloud_workspace_policy_group.nontest,
  ]
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_workspace_policy_groups.filter_by_description.policy_groups[*].description :
    strcontains(v, local.description)
  ]
}

output "is_description_filter_useful" {
  value = length(local.description_filter_result) > 0 && alltrue(local.description_filter_result)
}
`, testAccDataSourcePolicyGroups_base(name))
}

func testAccDataSourcePolicyGroups_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "test" {
  name  = "%[1]s"
  email = "www.%[1]s@example.com"
}

// The priority will automatically increment with the creation of the resource, no need to specify it manually.
resource "huaweicloud_workspace_policy_group" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"

  targets {
    type = "USER"
    id   = huaweicloud_workspace_user.test.id
    name = huaweicloud_workspace_user.test.name
  }

  policy {
    access_control {
      ip_access_control = "112.20.53.2|255.255.240.0;112.20.53.3|255.255.240.0"
    }
  }
}

resource "huaweicloud_workspace_policy_group" "nontest" {
  name = "non_%[1]s"

  targets {
    type = "USER"
    id   = huaweicloud_workspace_user.test.id
    name = huaweicloud_workspace_user.test.name
  }

  policy {
    access_control {
      ip_access_control = "112.20.53.2|255.255.240.0;112.20.53.3|255.255.240.0"
    }
  }
}
`, name)
}
