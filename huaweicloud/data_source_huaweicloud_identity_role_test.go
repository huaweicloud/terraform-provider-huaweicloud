package huaweicloud

import (
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIdentityRoleDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_identity_role.role_1"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityRoleDataSource_by_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "system_all_64"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "OBS ReadOnlyAccess"),
				),
			},
			{
				Config: testAccIdentityRoleDataSource_by_displayname,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "kms_adm"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "KMS Administrator"),
				),
			},
		},
	})
}

func testAccCheckIdentityDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmtp.Errorf("Can't find role data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmtp.Errorf("Role data source ID not set")
		}

		return nil
	}
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
