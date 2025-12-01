package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAppPolicyGroups_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_workspace_app_policy_groups.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByPolicyGroupName   = "data.huaweicloud_workspace_app_policy_groups.filter_by_policy_group_name"
		dcFilterByPolicyGroupName = acceptance.InitDataSourceCheck(filterByPolicyGroupName)

		filterByPolicyGroupType   = "data.huaweicloud_workspace_app_policy_groups.filter_by_policy_group_type"
		dcFilterByPolicyGroupType = acceptance.InitDataSourceCheck(filterByPolicyGroupType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckWorkspaceAppServerGroup(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAppPolicyGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "policy_groups.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					//  Filter by 'policy_group_name' parameter.
					dcFilterByPolicyGroupName.CheckResourceExists(),
					resource.TestCheckOutput("is_policy_group_name_filter_useful", "true"),
					// Filter by 'policy_group_type' parameter.
					dcFilterByPolicyGroupType.CheckResourceExists(),
					resource.TestCheckOutput("is_policy_group_type_filter_useful", "true"),
					// Check attributes.
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.id"),
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.name"),
					resource.TestMatchResourceAttr(filterByPolicyGroupName, "policy_groups.0.targets.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.targets.0.id"),
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.targets.0.name"),
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.targets.0.type"),
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.policies"),
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.priority"),
					resource.TestCheckResourceAttrSet(filterByPolicyGroupName, "policy_groups.0.description"),
					resource.TestMatchResourceAttr(filterByPolicyGroupName, "policy_groups.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestMatchResourceAttr(filterByPolicyGroupName, "policy_groups.0.updated_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
				),
			},
		},
	})
}

func testAccDataAppPolicyGroups_base(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_workspace_service" "test" {}

resource "huaweicloud_workspace_app_server_group" "test" {
  name             = "%[1]s"
  os_type          = "Windows"
  app_type         = "SESSION_DESKTOP_APP"
  flavor_id        = "%[2]s"
  vpc_id           = data.huaweicloud_workspace_service.test.vpc_id
  subnet_id        = data.huaweicloud_workspace_service.test.network_ids[0]
  system_disk_type = "SAS"
  system_disk_size = 90
  is_vdi           = true
  image_id         = "%[3]s"
  image_type       = "gold"
  image_product_id = "%[4]s"
}

resource "huaweicloud_workspace_app_group" "test" {
  name            = "%[1]s"
  type            = "SESSION_DESKTOP_APP"
  server_group_id = huaweicloud_workspace_app_server_group.test.id
}

resource "huaweicloud_workspace_app_policy_group" "test" {
  name        = "%[1]s"
  description = "Created by terraform script"

  targets {
    id   = huaweicloud_workspace_app_group.test.id
    name = huaweicloud_workspace_app_group.test.name
    type = "APPGROUP"
  }

  policies = jsonencode({
    client = {
      automatic_reconnection_interval = 10
      session_persistence_time        = 120
      forbid_screen_capture           = true
    }
  })
}

resource "huaweicloud_workspace_app_policy_template" "test" {
  name = "%[1]s"

  policies = jsonencode({
    client = {
      automatic_reconnection_interval = 10
      session_persistence_time        = 120
      forbid_screen_capture           = true
    }
  })
}
`, name, acceptance.HW_WORKSPACE_APP_SERVER_GROUP_FLAVOR_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_ID,
		acceptance.HW_WORKSPACE_APP_SERVER_GROUP_IMAGE_PRODUCT_ID)
}

func testAccDataAppPolicyGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_app_policy_groups" "all" {
  depends_on = [huaweicloud_workspace_app_policy_group.test]
}

# Filter by 'policy_group_name' parameter.
locals {
  policy_group_name = huaweicloud_workspace_app_policy_group.test.name
}

data "huaweicloud_workspace_app_policy_groups" "filter_by_policy_group_name" {
  policy_group_name = local.policy_group_name

  depends_on = [huaweicloud_workspace_app_policy_group.test]
}

locals {
  policy_group_name_filter_result = [
    for v in data.huaweicloud_workspace_app_policy_groups.filter_by_policy_group_name.policy_groups[*].name :
    strcontains(v, local.policy_group_name)
  ]
}

output "is_policy_group_name_filter_useful" {
  value = length(local.policy_group_name_filter_result) > 0 && alltrue(local.policy_group_name_filter_result)
}

# Filter by 'policy_group_type' parameter.
data "huaweicloud_workspace_app_policy_groups" "filter_by_policy_group_type" {
  policy_group_type = 4

  depends_on = [huaweicloud_workspace_app_policy_template.test]
}

locals {
  policy_group_type_filter_result = data.huaweicloud_workspace_app_policy_groups.filter_by_policy_group_type.policy_groups[*].name
}

output "is_policy_group_type_filter_useful" {
  value = (
    length(local.policy_group_type_filter_result) > 0 &&
    contains(local.policy_group_type_filter_result, huaweicloud_workspace_app_policy_template.test.name)
  )
}
`, testAccDataAppPolicyGroups_base(name))
}
