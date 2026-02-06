package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataV3Roles_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_roles.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byDisplayName   = "data.huaweicloud_identity_roles.filter_by_display_name"
		dcByDisplayName = acceptance.InitDataSourceCheck(byDisplayName)

		byCatalog   = "data.huaweicloud_identity_roles.filter_by_catalog"
		dcByCatalog = acceptance.InitDataSourceCheck(byCatalog)

		byType   = "data.huaweicloud_identity_roles.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)

		byPermissionType   = "data.huaweicloud_identity_roles.filter_by_permission_type"
		dcByPermissionType = acceptance.InitDataSourceCheck(byPermissionType)

		byDomainId   = "data.huaweicloud_identity_roles.filter_by_domain_id"
		dcByDomainId = acceptance.InitDataSourceCheck(byDomainId)

		name = acceptance.RandomAccResourceName()
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
			acceptance.TestAccPrecheckDomainId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataV3Roles_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "roles.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					dcByDisplayName.CheckResourceExists(),
					resource.TestCheckOutput("display_name_filter_is_useful", "true"),
					dcByCatalog.CheckResourceExists(),
					resource.TestCheckOutput("catalog_filter_is_useful", "true"),
					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					dcByPermissionType.CheckResourceExists(),
					resource.TestCheckOutput("permission_type_filter_is_useful", "true"),
					dcByDomainId.CheckResourceExists(),
					resource.TestCheckOutput("domain_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDataV3Roles_base(name string) string {
	return fmt.Sprintf(`
variable "role_configs" {
  type = list(object({
    suffix      = string
    description = string
    action      = string
  }))
  default = [
    {
      suffix      = "1"
      description = "Created by terraform script for acc test 1"
      action      = "obs:bucket:GetBucketAcl"
    },
    {
      suffix      = "2"
      description = "Created by terraform script for acc test 2"
      action      = "ecs:servers:list"
    },
  ]
}

resource "huaweicloud_identity_role" "test" {
  count = length(var.role_configs)

  name        = "%[1]s_${var.role_configs[count.index].suffix}"
  description = var.role_configs[count.index].description
  type        = "AX"
  policy      = jsonencode({
    Version = "1.1"

    Statement = [
      {
        Action = [var.role_configs[count.index].action]
        Effect = "Allow"
      }
    ]
  })
}
`, name)
}

func testAccDataV3Roles_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

# All
data "huaweicloud_identity_roles" "all" {
  depends_on = [huaweicloud_identity_role.test]
}

# Filter by display_name
locals {
  display_name = try(data.huaweicloud_identity_roles.all.roles[0].display_name, "NOT_FOUND")
}

data "huaweicloud_identity_roles" "filter_by_display_name" {
  display_name = local.display_name
}

locals {
  display_name_filter_result = [
    for v in data.huaweicloud_identity_roles.filter_by_display_name.roles[*].display_name : v == local.display_name
  ]
}

output "display_name_filter_is_useful" {
  value = length(local.display_name_filter_result) > 0 && alltrue(local.display_name_filter_result)
}

# Filter by catalog
locals {
  catalog = try(data.huaweicloud_identity_roles.all.roles[0].catalog, "NOT_FOUND")
}

data "huaweicloud_identity_roles" "filter_by_catalog" {
  catalog = local.catalog
}

locals {
  catalog_filter_result = [
    for v in data.huaweicloud_identity_roles.filter_by_catalog.roles[*].catalog : v == local.catalog
  ]
}

output "catalog_filter_is_useful" {
  value = length(local.catalog_filter_result) > 0 && alltrue(local.catalog_filter_result)
}

# Filter by type
locals {
  type = "all"
}

data "huaweicloud_identity_roles" "filter_by_type" {
  type    = local.type
  catalog = "OBS"
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_identity_roles.filter_by_type.roles) > 0
}

# Filter by permission_type
locals {
  permission_type = "policy"
}

data "huaweicloud_identity_roles" "filter_by_permission_type" {
  permission_type = local.permission_type
  catalog         = "VPC"
}

output "permission_type_filter_is_useful" {
  value = length(data.huaweicloud_identity_roles.filter_by_permission_type.roles) > 0
}

# Filter by domain_id (custom policies)
data "huaweicloud_identity_roles" "filter_by_domain_id" {
  domain_id = "%[2]s"

  depends_on = [huaweicloud_identity_role.test]
}

output "domain_id_filter_is_useful" {
  value = length(data.huaweicloud_identity_roles.filter_by_domain_id.roles) >= 2
}
`, testAccDataV3Roles_base(name), acceptance.HW_DOMAIN_ID)
}
