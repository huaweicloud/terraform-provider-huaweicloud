package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityRoleDataSource_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_identity_role.test"

		dc = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleDataSource_by_name,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "name", "system_all_64"),
					resource.TestCheckResourceAttr(all, "display_name", "OBS ReadOnlyAccess"),
					resource.TestCheckResourceAttrSet(all, "catalog"),
					resource.TestCheckResourceAttrSet(all, "type"),
					resource.TestCheckResourceAttrSet(all, "policy"),
				),
			},
			{
				Config: testAccIdentityRoleDataSource_by_displayname,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "name", "kms_adm"),
					resource.TestCheckResourceAttr(all, "display_name", "KMS Administrator"),
					resource.TestCheckResourceAttrSet(all, "catalog"),
					resource.TestCheckResourceAttrSet(all, "type"),
					resource.TestCheckResourceAttrSet(all, "policy"),
				),
			},
			{
				Config: testAccIdentityRoleDataSource_by_roleid,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(all, "name", "kms_adm"),
					resource.TestCheckResourceAttr(all, "display_name", "KMS Administrator"),
					resource.TestCheckResourceAttrSet(all, "catalog"),
					resource.TestCheckResourceAttrSet(all, "type"),
					resource.TestCheckResourceAttrSet(all, "policy"),
				),
			},
		},
	})
}

const testAccIdentityRoleDataSource_by_name = `
data "huaweicloud_identity_role" "test" {
  # OBS ReadOnlyAccess
  name = "system_all_64"
}
`
const testAccIdentityRoleDataSource_by_displayname = `
data "huaweicloud_identity_role" "test" {
  display_name = "KMS Administrator"
}
`

const testAccIdentityRoleDataSource_by_roleid = `
data "huaweicloud_identity_role" "test_by_roleid" {
  display_name = "KMS Administrator"
}

data "huaweicloud_identity_role" "test" {
  role_id = data.huaweicloud_identity_role.test_by_roleid.id
}
`
