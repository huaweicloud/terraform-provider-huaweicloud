package iam

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccIdentityRoleDataSource_basic(t *testing.T) {
	dataSourceName := "data.huaweicloud_identity_role.role_1"

	dc := acceptance.InitDataSourceCheck(dataSourceName)

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
					resource.TestCheckResourceAttr(dataSourceName, "name", "system_all_64"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "OBS ReadOnlyAccess"),
				),
			},
			{
				Config: testAccIdentityRoleDataSource_by_displayname,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "kms_adm"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "KMS Administrator"),
				),
			},
			{
				Config: testAccIdentityRoleDataSource_by_roleid,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(dataSourceName, "name", "kms_adm"),
					resource.TestCheckResourceAttr(dataSourceName, "display_name", "KMS Administrator"),
				),
			},
		},
	})
}

const testAccIdentityRoleDataSource_by_name = `
data "huaweicloud_identity_role" "role_1" {
  # OBS ReadOnlyAccess
  name = "system_all_64"
}
`
const testAccIdentityRoleDataSource_by_displayname = `
data "huaweicloud_identity_role" "role_1" {
  display_name = "KMS Administrator"
}
`

const testAccIdentityRoleDataSource_by_roleid = `
data "huaweicloud_identity_role" "role_3" {
  display_name = "KMS Administrator"
}

data "huaweicloud_identity_role" "role_1" {
  role_id = data.huaweicloud_identity_role.role_3.id
}
`
