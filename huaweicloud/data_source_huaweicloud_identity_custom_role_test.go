package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIdentityCustomRoleDataSource_basic(t *testing.T) {
	var rName = fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCustomRoleDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityCustomDataSourceID("data.huaweicloud_identity_custom_role.role_1"),
					resource.TestCheckResourceAttr(
						"data.huaweicloud_identity_custom_role.role_1", "name", rName),
				),
			},
		},
	})
}

func testAccCheckIdentityCustomDataSourceID(n string) resource.TestCheckFunc {
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

func testAccIdentityCustomRoleDataSource_basic(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_role" test {
  name        = "%s"
  description = "created by terraform"
  type        = "AX"
  policy      = <<EOF
{
  "Version": "1.1",
  "Statement": [
    {
      "Action": [
        "obs:bucket:GetBucketAcl"
      ],
      "Effect": "Allow",
      "Resource": [
        "obs:*:*:bucket:*"
      ]
    }
  ]
}
EOF
}

data "huaweicloud_identity_custom_role" "role_1" {
  name = huaweicloud_identity_role.test.name
}
`, rName)
}
