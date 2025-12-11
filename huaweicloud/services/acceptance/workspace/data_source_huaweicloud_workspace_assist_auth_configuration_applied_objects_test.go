package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataAssistAuthConfigurationAppliedObjects_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_workspace_assist_auth_configuration_applied_objects.all"
		dc  = acceptance.InitDataSourceCheck(all)

		filterByObjectTypeUser        = "data.huaweicloud_workspace_assist_auth_configuration_applied_objects.filter_by_object_type_user"
		dcFilterByObjectTypeUser      = acceptance.InitDataSourceCheck(filterByObjectTypeUser)
		filterByObjectTypeUserGroup   = "data.huaweicloud_workspace_assist_auth_configuration_applied_objects.filter_by_object_type_user_group"
		dcFilterByObjectTypeUserGroup = acceptance.InitDataSourceCheck(filterByObjectTypeUserGroup)

		filterByObjectName   = "data.huaweicloud_workspace_assist_auth_configuration_applied_objects.filter_by_object_name"
		dcFilterByObjectName = acceptance.InitDataSourceCheck(filterByObjectName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAssistAuthConfigurationAppliedObjects_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameter.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "objects.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "objects.0.type"),
					resource.TestCheckResourceAttrSet(all, "objects.0.id"),
					resource.TestCheckResourceAttrSet(all, "objects.0.name"),
					// Filter by 'object_type' parameter and the value is 'USER'.
					dcFilterByObjectTypeUser.CheckResourceExists(),
					resource.TestCheckOutput("is_object_type_user_filter_useful", "true"),
					// Filter by 'object_type' parameter and the value is 'USER_GROUP'.
					dcFilterByObjectTypeUserGroup.CheckResourceExists(),
					resource.TestCheckOutput("is_object_type_user_group_filter_useful", "true"),
					// Filter by 'object_name' parameter.
					dcFilterByObjectName.CheckResourceExists(),
					resource.TestCheckOutput("is_object_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataAssistAuthConfigurationAppliedObjects_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "test" {
  count = 6

  name  = format("%[1]s-user-%%d", count.index)
  email = format("test%%d@example.com", count.index)
}

resource "huaweicloud_workspace_user_group" "test" {
  count = 2

  name = format("%[1]s-group-%%d", count.index)
  type = "LOCAL"

  dynamic "users" {
    for_each = slice(huaweicloud_workspace_user.test[*], count.index*2+2, count.index*2+4)

    content {
      id = users.value.id
    }
  }
}

resource "huaweicloud_workspace_assist_auth_configuration_object_management" "test" {
  dynamic "objects" {
    for_each = huaweicloud_workspace_user.test[*]

    content {
      type = "USER"
      id   = objects.value.id
      name = objects.value.name
    }
  }

  dynamic "objects" {
    for_each = huaweicloud_workspace_user_group.test[*]

    content {
      type = "USER_GROUP"
      id   = objects.value.id
      name = objects.value.name
    }
  }
}
`, name)
}

func testAccDataAssistAuthConfigurationAppliedObjects_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameter.
data "huaweicloud_workspace_assist_auth_configuration_applied_objects" "all" {
  depends_on = [
    huaweicloud_workspace_assist_auth_configuration_object_management.test,
  ]
}

# Filter by 'object_type' parameter and the value is 'USER'.
data "huaweicloud_workspace_assist_auth_configuration_applied_objects" "filter_by_object_type_user" {
  depends_on = [
    huaweicloud_workspace_assist_auth_configuration_object_management.test,
  ]

  object_type = "USER"
}

locals {
  object_type_user_filter_result = [
    for v in data.huaweicloud_workspace_assist_auth_configuration_applied_objects.filter_by_object_type_user.objects[*].type :
    v == "USER"
  ]
}

output "is_object_type_user_filter_useful" {
  value = length(local.object_type_user_filter_result) > 0 && alltrue(local.object_type_user_filter_result)
}

# Filter by 'object_type' parameter and the value is 'USER_GROUP'.
data "huaweicloud_workspace_assist_auth_configuration_applied_objects" "filter_by_object_type_user_group" {
  depends_on = [
    huaweicloud_workspace_assist_auth_configuration_object_management.test,
  ]

  object_type = "USER_GROUP"
}

locals {
  object_type_user_group_filter_result = [
    for v in data.huaweicloud_workspace_assist_auth_configuration_applied_objects.filter_by_object_type_user_group.objects[*].type :
    v == "USER_GROUP"
  ]
}

output "is_object_type_user_group_filter_useful" {
  value = length(local.object_type_user_group_filter_result) > 0 && alltrue(local.object_type_user_group_filter_result)
}

# Filter by 'object_name' parameter.
locals {
  object_name = try(data.huaweicloud_workspace_assist_auth_configuration_applied_objects.all.objects[0].name, "NOT_FOUND")
}

data "huaweicloud_workspace_assist_auth_configuration_applied_objects" "filter_by_object_name" {
  depends_on = [
    huaweicloud_workspace_assist_auth_configuration_object_management.test,
  ]

  object_name = local.object_name
}

locals {
  object_name_filter_result = [
    for v in data.huaweicloud_workspace_assist_auth_configuration_applied_objects.filter_by_object_name.objects[*].name :
    v == local.object_name
  ]
}

output "is_object_name_filter_useful" {
  value = length(local.object_name_filter_result) > 0 && alltrue(local.object_name_filter_result)
}
`, testAccDataAssistAuthConfigurationAppliedObjects_base(name), name)
}
