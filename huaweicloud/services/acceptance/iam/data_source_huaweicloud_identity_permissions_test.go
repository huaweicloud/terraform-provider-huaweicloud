package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityPermissionsDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identity_permissions.test"
	dc := acceptance.InitDataSourceCheck(resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityPermissionsDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(resourceName, "permissions.#"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("catalog_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("custom_filter_is_useful", "true"),
					resource.TestCheckOutput("scope_type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccIdentityPermissionsDataSource_basic() string {
	return `
data "huaweicloud_identity_permissions" "test" {
}

data "huaweicloud_identity_permissions" "by_name" {
  name = "KMS Administrator"
}

data "huaweicloud_identity_permissions" "by_catalog" {
  catalog = "ELB"
}

data "huaweicloud_identity_permissions" "obs_policy" {
  type    = "system-policy"
  catalog = "OBS"
}

data "huaweicloud_identity_permissions" "custom" {
  type = "custom"
}

data "huaweicloud_identity_permissions" "by_scope_type" {
  scope_type = "project"
  name = "CCE FullAccess"
}

output "name_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.by_name.permissions[*].name : v == "KMS Administrator"])
}

output "catalog_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.by_catalog.permissions[*].catalog : v == "ELB"])
}

output "type_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.obs_policy.permissions[*].catalog : v == "OBS"])
}

output "custom_filter_is_useful" {
  value = alltrue([for v in data.huaweicloud_identity_permissions.custom.permissions[*].catalog : v == "CUSTOMED"])
}

output "scope_type_filter_is_useful" {
  value = alltrue([length(data.huaweicloud_identity_permissions.by_scope_type.permissions) == 1])
}
`
}
