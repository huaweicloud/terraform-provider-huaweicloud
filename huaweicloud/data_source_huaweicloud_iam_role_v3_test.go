package huaweicloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIAMRoleV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIAMRoleV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.huaweicloud_iam_role_v3.roles", "id"),
				),
			},
		},
	})
}

const testAccIAMRoleV3DataSource_basic = `
data "huaweicloud_iam_role_v3" "roles" {
}
`
