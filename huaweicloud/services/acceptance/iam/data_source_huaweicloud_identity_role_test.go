package iam

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityRoleDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identity_role.role_1"
	dc := acceptance.InitDataSourceCheck(resourceName)

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
					resource.TestCheckResourceAttr(resourceName, "name", "system_all_64"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "OBS ReadOnlyAccess"),
				),
			},
			{
				Config: testAccIdentityRoleDataSource_by_displayname,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", "kms_adm"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "KMS Administrator"),
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
