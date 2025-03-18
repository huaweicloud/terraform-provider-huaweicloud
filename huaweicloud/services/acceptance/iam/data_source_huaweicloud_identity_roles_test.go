package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityRolesDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_roles.roles"
	dc := acceptance.InitDataSourceCheck(dataSourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRolesDataSource_by_name,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "roles.#"),
					resource.TestCheckResourceAttr(dataSourceName, "roles.0.name", "system_all_64"),
					resource.TestCheckResourceAttr(dataSourceName, "roles.0.display_name", "OBS ReadOnlyAccess"),
				),
			},
			{
				Config: testAccIdentityRolesDataSource_by_displayname,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "roles.#"),
					resource.TestCheckResourceAttr(dataSourceName, "roles.0.name", "kms_adm"),
					resource.TestCheckResourceAttr(dataSourceName, "roles.0.display_name", "KMS Administrator"),
				),
			},
			{
				Config: testAccIdentityRolesDataSource_empty_filter,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "roles.#", "0"),
				),
			},
			{
				Config: testAccIdentityRolesDataSource_multiple_roles,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "roles.#"),
					resource.TestCheckResourceAttr(dataSourceName, "roles.0.name", "system_all_64"),
					resource.TestCheckResourceAttr(dataSourceName, "roles.1.name", "kms_adm"),
				),
			},
		},
	})
}

const testAccIdentityRolesDataSource_by_name = `
data "huaweicloud_identity_roles" "roles" {
  name = "system_all_64"
}
`

const testAccIdentityRolesDataSource_by_displayname = `
data "huaweicloud_identity_roles" "roles" {
  display_name = "KMS Administrator"
}
`

const testAccIdentityRolesDataSource_empty_filter = `
data "huaweicloud_identity_roles" "roles" {
  name = "nonexistent_role"
}
`

const testAccIdentityRolesDataSource_multiple_roles = `
data "huaweicloud_identity_roles" "roles" {
}
`
