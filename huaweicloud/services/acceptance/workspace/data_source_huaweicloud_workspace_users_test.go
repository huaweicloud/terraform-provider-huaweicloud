package workspace

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataUsers_basic(t *testing.T) {
	var (
		name     = acceptance.RandomAccResourceName()
		password = acceptance.RandomPassword()

		dcName = "data.huaweicloud_workspace_users.all"
		dc     = acceptance.InitDataSourceCheck(dcName)

		filterByUserName   = "data.huaweicloud_workspace_users.filter_by_user_name"
		dcFilterByUserName = acceptance.InitDataSourceCheck(filterByUserName)

		filterByDescription   = "data.huaweicloud_workspace_users.filter_by_description"
		dcFilterByDescription = acceptance.InitDataSourceCheck(filterByDescription)

		filterByActiveType   = "data.huaweicloud_workspace_users.filter_by_active_type"
		dcFilterByActiveType = acceptance.InitDataSourceCheck(filterByActiveType)

		filterByGroupName   = "data.huaweicloud_workspace_users.filter_by_group_name"
		dcFilterByGroupName = acceptance.InitDataSourceCheck(filterByGroupName)

		computeByIsQueryTotalDesktops   = "data.huaweicloud_workspace_users.computor_by_is_query_total_desktops"
		dcComputeByIsQueryTotalDesktops = acceptance.InitDataSourceCheck(computeByIsQueryTotalDesktops)

		filterByEnterpriseProjectId   = "data.huaweicloud_workspace_users.filter_by_enterprise_project_id"
		dcFilterByEnterpriseProjectId = acceptance.InitDataSourceCheck(filterByEnterpriseProjectId)
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUsers_basic(name, password),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(dcName, "users.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(dcName, "users.0.id"),
					resource.TestCheckResourceAttrSet(dcName, "users.0.sid"),
					resource.TestCheckResourceAttrSet(dcName, "users.0.user_name"),
					resource.TestCheckResourceAttrSet(dcName, "users.0.total_desktops"),
					resource.TestCheckResourceAttrSet(dcName, "users.0.active_type"),
					dcFilterByUserName.CheckResourceExists(),
					resource.TestCheckOutput("is_user_name_filter_useful", "true"),
					dcFilterByDescription.CheckResourceExists(),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
					dcFilterByActiveType.CheckResourceExists(),
					resource.TestCheckOutput("is_active_type_filter_useful", "true"),
					dcFilterByGroupName.CheckResourceExists(),
					resource.TestCheckOutput("is_group_name_filter_useful", "true"),
					dcComputeByIsQueryTotalDesktops.CheckResourceExists(),
					resource.TestCheckOutput("is_query_total_desktops_computor_useful", "true"),
					dcFilterByEnterpriseProjectId.CheckResourceExists(),
					resource.TestCheckOutput("is_enterprise_project_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataUsers_base(name, password string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_user" "with_user_activate" {
  name        = "%[1]s_with_user"
  description = "Created user_with_user_activate by terraform script"
  active_type = "USER_ACTIVATE"
  email       = "terraform@test.com"
  password    = "%[2]s"

  account_expires            = "0"
  password_never_expires     = false
  enable_change_password     = true
  next_login_change_password = true
  disabled                   = false
}

resource "huaweicloud_workspace_user" "with_admin_activate" {
  name        = "%[1]s_with_admin"
  description = "Created user_with_admin_activate by terraform script"
  active_type = "ADMIN_ACTIVATE"
  email       = "terraform@test.com"
  password    = "%[2]s"

  account_expires            = "0"
  password_never_expires     = false
  enable_change_password     = true
  next_login_change_password = true
  disabled                   = false
}

resource "huaweicloud_workspace_user_group" "user_group_with_local" {
  name        = "%[1]s_with_local"
  type        = "LOCAL"
  description = "Created by terraform script"

  users {
    id = huaweicloud_workspace_user.with_user_activate.id
  }
}`, name, password)
}

func testAccDataUsers_basic(name, password string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_workspace_users" "all" {
  depends_on = [
    huaweicloud_workspace_user.with_user_activate,
    huaweicloud_workspace_user.with_admin_activate,
  ]
}

locals {
  user_name             = try(data.huaweicloud_workspace_users.all.users[0].user_name, "NOT_FOUND")
  description           = try(data.huaweicloud_workspace_users.all.users[0].description, "NOT_FOUND")
  active_type           = try(data.huaweicloud_workspace_users.all.users[0].active_type, "NOT_FOUND")
  group_name            = huaweicloud_workspace_user_group.user_group_with_local.name
  enterprise_project_id = try(data.huaweicloud_workspace_users.all.users[0].enterprise_project_id, "NOT_FOUND")
}

# Filter by user name
data "huaweicloud_workspace_users" "filter_by_user_name" {
  user_name = local.user_name

  depends_on = [
    huaweicloud_workspace_user.with_user_activate,
    huaweicloud_workspace_user.with_admin_activate,
  ]
}

locals {
  user_name_filter_result = [
    for v in data.huaweicloud_workspace_users.filter_by_user_name.users[*].user_name :
    strcontains(v, local.user_name)
  ]
}

output "is_user_name_filter_useful" {
  value = length(local.user_name_filter_result) > 0 && alltrue(local.user_name_filter_result)
}

# Filter by description
data "huaweicloud_workspace_users" "filter_by_description" {
  description = local.description

  depends_on = [
    huaweicloud_workspace_user.with_user_activate,
    huaweicloud_workspace_user.with_admin_activate,
  ]
}

locals {
  description_filter_result = [
    for v in data.huaweicloud_workspace_users.filter_by_description.users[*].description :
    strcontains(v, local.description)
  ]
}

output "is_description_filter_useful" {
  value = length(local.description_filter_result) > 0 && alltrue(local.description_filter_result)
}

# Filter by active type
data "huaweicloud_workspace_users" "filter_by_active_type" {
  active_type = local.active_type

  depends_on = [
    huaweicloud_workspace_user.with_user_activate,
    huaweicloud_workspace_user.with_admin_activate,
  ]
}

locals {
  active_type_filter_result = [
    for v in data.huaweicloud_workspace_users.filter_by_active_type.users[*].active_type :
    v == local.active_type
  ]
}

output "is_active_type_filter_useful" {
  value = length(local.active_type_filter_result) > 0 && alltrue(local.active_type_filter_result)
}

# Filter by group name
data "huaweicloud_workspace_users" "filter_by_group_name" {
  group_name = local.group_name

  depends_on = [
    huaweicloud_workspace_user.with_user_activate,
    huaweicloud_workspace_user.with_admin_activate,
  ]
}

locals {
  group_name_filter_result = [
    for v in data.huaweicloud_workspace_users.filter_by_group_name.users[*].group_names :
    contains(v, local.group_name)
  ]
}

output "is_group_name_filter_useful" {
  value = length(local.group_name_filter_result) > 0 && alltrue(local.group_name_filter_result)
}

# Computor by is query total desktops
data "huaweicloud_workspace_users" "computor_by_is_query_total_desktops" {
  is_query_total_desktops = false

  depends_on = [
    huaweicloud_workspace_user.with_user_activate,
    huaweicloud_workspace_user.with_admin_activate,
  ]
}

locals {
  is_query_total_desktops_computor_result = [
    for v in data.huaweicloud_workspace_users.computor_by_is_query_total_desktops.users[*].total_desktops :
    v == 0
  ]
}

output "is_query_total_desktops_computor_useful" {
  value = length(local.is_query_total_desktops_computor_result) > 0 && alltrue(local.is_query_total_desktops_computor_result)
}

# Filter by enterprise project ID
data "huaweicloud_workspace_users" "filter_by_enterprise_project_id" {
  enterprise_project_id = local.enterprise_project_id

  depends_on = [
    huaweicloud_workspace_user.with_user_activate,
    huaweicloud_workspace_user.with_admin_activate,
  ]
}

locals {
  is_enterprise_project_id_filter_result = [
    for v in data.huaweicloud_workspace_users.filter_by_enterprise_project_id.users[*].enterprise_project_id :
    v == local.enterprise_project_id
  ]
}

output "is_enterprise_project_id_filter_useful" {
  value = length(local.is_enterprise_project_id_filter_result) > 0 && alltrue(local.is_enterprise_project_id_filter_result)
}`, testAccDataUsers_base(name, password))
}
