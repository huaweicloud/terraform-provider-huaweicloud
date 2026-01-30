package iam

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataPermissions_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_permissions.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byName   = "data.huaweicloud_identity_permissions.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byCatalog   = "data.huaweicloud_identity_permissions.filter_by_catalog"
		dcByCatalog = acceptance.InitDataSourceCheck(byCatalog)

		byType   = "data.huaweicloud_identity_permissions.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byCustom   = "data.huaweicloud_identity_permissions.filter_by_custom"
		dcByCustom = acceptance.InitDataSourceCheck(byCustom)

		byScopeType   = "data.huaweicloud_identity_permissions.filter_by_scope_type"
		dcByScopeType = acceptance.InitDataSourceCheck(byScopeType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataPermissions_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "permissions.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					dcByCatalog.CheckResourceExists(),
					resource.TestCheckOutput("catalog_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					dcByCustom.CheckResourceExists(),
					resource.TestCheckOutput("custom_filter_is_useful", "true"),
					dcByScopeType.CheckResourceExists(),
					resource.TestCheckOutput("scope_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

const testAccDataPermissions_basic = `
# All
data "huaweicloud_identity_permissions" "all" {
}

# Filter by name
locals {
  name = "KMS Administrator"
}

data "huaweicloud_identity_permissions" "filter_by_name" {
  name = local.name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_identity_permissions.filter_by_name.permissions[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by catalog
locals {
  catalog = "ELB"
}

data "huaweicloud_identity_permissions" "filter_by_catalog" {
  catalog = local.catalog
}

locals {
  catalog_filter_result = [
    for v in data.huaweicloud_identity_permissions.filter_by_catalog.permissions[*].catalog : v == local.catalog
  ]
}

output "catalog_filter_is_useful" {
  value = length(local.catalog_filter_result) > 0 && alltrue(local.catalog_filter_result)
}

# Filter by type
locals {
  type         = "system-policy"
  type_catalog = "OBS"
}

data "huaweicloud_identity_permissions" "filter_by_type" {
  type    = local.type
  catalog = local.type_catalog
}

locals {
  type_filter_result = [
    for v in data.huaweicloud_identity_permissions.filter_by_type.permissions[*].catalog : v == local.type_catalog
  ]
}

output "type_filter_is_useful" {
  value = length(local.type_filter_result) > 0 && alltrue(local.type_filter_result)
}

# Filter by custom
locals {
  custom         = "custom"
  custom_catalog = "CUSTOMED"
}

data "huaweicloud_identity_permissions" "filter_by_custom" {
  type = local.custom
}

locals {
  custom_filter_result = [
    for v in data.huaweicloud_identity_permissions.filter_by_custom.permissions[*].catalog : v == local.custom_catalog
  ]
}

output "custom_filter_is_useful" {
  value = alltrue(local.custom_filter_result)
}

# Filter by scope type
locals {
  scope_type      = "project"
  scope_type_name = "CCE FullAccess"
}

data "huaweicloud_identity_permissions" "filter_by_scope_type" {
  scope_type = local.scope_type
  name       = local.scope_type_name
}

output "scope_type_filter_is_useful" {
  value = alltrue([length(data.huaweicloud_identity_permissions.filter_by_scope_type.permissions) == 1])
}
`
