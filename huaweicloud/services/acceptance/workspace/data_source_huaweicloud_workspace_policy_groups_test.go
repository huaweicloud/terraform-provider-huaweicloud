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

		filterByUpdateTime   = "data.huaweicloud_workspace_policy_groups.filter_by_update_time"
		dcFilterByUpdateTime = acceptance.InitDataSourceCheck(filterByUpdateTime)

		filterByDescription   = "data.huaweicloud_workspace_policy_groups.filter_by_description"
		dcFilterByDescription = acceptance.InitDataSourceCheck(filterByDescription)

		filterByNameWithAccurate   = "data.huaweicloud_workspace_policy_groups.filter_by_name_accurate"
		dcFilterByNameWithAccurate = acceptance.InitDataSourceCheck(filterByNameWithAccurate)
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
					// all
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "policy_groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.policy_group_id"),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.policy_group_name"),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.priority"),
					resource.TestCheckResourceAttrSet(dcName, "policy_groups.0.update_time"),
					// filter by id
					dcFilterById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					// filter by name
					dcFilterByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					// filter by priority
					dcFilterByPriority.CheckResourceExists(),
					resource.TestCheckOutput("is_priority_filter_useful", "true"),
					// filter by update time
					dcFilterByUpdateTime.CheckResourceExists(),
					resource.TestCheckOutput("is_update_time_filter_useful", "true"),
					// filter by description
					dcFilterByDescription.CheckResourceExists(),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					// filter by name with accurate
					dcFilterByNameWithAccurate.CheckResourceExists(),
					resource.TestCheckOutput("is_accurate_name_filter_useful", "true"),
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
    huaweicloud_workspace_app_policy_group.test,
    huaweicloud_workspace_app_policy_group.nontest,
  ]
}

locals {
  policy_group_id         = data.huaweicloud_workspace_policy_groups.all.policy_groups[0].policy_group_id
  policy_group_name       = data.huaweicloud_workspace_policy_groups.all.policy_groups[0].policy_group_name
  priority                = data.huaweicloud_workspace_policy_groups.all.policy_groups[0].priority
  update_time             = data.huaweicloud_workspace_policy_groups.all.policy_groups[0].update_time
  description             = data.huaweicloud_workspace_policy_groups.all.policy_groups[0].description
}

# filter by policy group id
data "huaweicloud_workspace_policy_groups" "filter_by_policy_group_id" {
  policy_group_id = local.policy_group_id

  depends_on = [
    huaweicloud_workspace_app_policy_group.test,
    huaweicloud_workspace_app_policy_group.nontest,
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

# filter by policy group name
data "huaweicloud_workspace_policy_groups" "filter_by_policy_group_name" {
  policy_group_name = local.policy_group_name

  depends_on = [
    huaweicloud_workspace_app_policy_group.test,
    huaweicloud_workspace_app_policy_group.nontest,
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_workspace_policy_groups.filter_by_policy_group_name.policy_groups) > 0
}

# filter by priority
data "huaweicloud_workspace_policy_groups" "filter_by_priority" {
  priority = local.priority

  depends_on = [
    huaweicloud_workspace_app_policy_group.test,
    huaweicloud_workspace_app_policy_group.nontest,
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

# filter by update time
data "huaweicloud_workspace_policy_groups" "filter_by_update_time" {
  update_time = local.update_time

  depends_on = [
    huaweicloud_workspace_app_policy_group.test,
    huaweicloud_workspace_app_policy_group.nontest,
  ]
}

output "is_update_time_filter_useful" {
  value = length(data.huaweicloud_workspace_policy_groups.filter_by_update_time.policy_groups) > 0
}

# filter by description
data "huaweicloud_workspace_policy_groups" "filter_by_description" {
  description = local.description

  depends_on = [
    huaweicloud_workspace_app_policy_group.test,
    huaweicloud_workspace_app_policy_group.nontest,
  ]
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_workspace_policy_groups.filter_by_description.policy_groups[*].description : 
      v == local.description
  ]
}

output "is_description_filter_useful" {
  value = length(local.description_filter_result) > 0 && alltrue(local.description_filter_result)
}

# filter by description with accurate condition
data "huaweicloud_workspace_policy_groups" "filter_by_name_accurate" {
  policy_group_name = local.policy_group_name
  is_group_name_accurate = true

  depends_on = [
    huaweicloud_workspace_app_policy_group.test,
    huaweicloud_workspace_app_policy_group.nontest,
  ]
}

locals {
  accurate_name_filter_result = [
    for v in data.huaweicloud_workspace_policy_groups.filter_by_name_accurate.policy_groups[*].policy_group_name : 
      v == local.policy_group_name
  ]
}

output "is_accurate_name_filter_useful" {
  value = length(local.accurate_name_filter_result) < 2 && alltrue(local.accurate_name_filter_result)
}
`, testAccDataSourcePolicyGroups_base(name))
}

func testAccDataSourcePolicyGroups_base(name string) string {
	return fmt.Sprintf(`
locals {
  user_name           = "terraform"
  email_address       = "www.terraform@163.com"
  policy_name         = "%[1]s"
  nontest_policy_name = "non_%[1]s"
}

resource "huaweicloud_workspace_user" "test" {
  name  = local.user_name
  email = local.email_address
}

resource "huaweicloud_workspace_app_policy_group" "test" {
  name     = local.policy_name

  targets {
    type = "USER"
    id   = huaweicloud_workspace_user.test.id
    name = huaweicloud_workspace_user.test.name
  }

  policies = jsonencode({
    "client": {
      "automatic_reconnection_interval" : 5,
      "session_persistence_time" : 180,
      "forbid_screen_capture" : false
    }
  })
}

resource "huaweicloud_workspace_app_policy_group" "nontest" {
  name     = local.nontest_policy_name

  targets {
    type = "USER"
    id   = huaweicloud_workspace_user.test.id
    name = huaweicloud_workspace_user.test.name
  }

  policies = jsonencode({
    "client": {
      "automatic_reconnection_interval" : 5,
      "session_persistence_time" : 180,
      "forbid_screen_capture" : false
    }
  })
}
`, name)
}
