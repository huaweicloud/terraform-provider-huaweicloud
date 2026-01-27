package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCustomRole_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identity_custom_role.all"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckAdminOnly(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataCustomRole_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "id", regexp.MustCompile(`^[a-zA-Z0-9-]+$`)),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataCustomRole_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identity_role" "test" {
  name        = "%[1]s"
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
`, name)
}

func testAccDataCustomRole_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identity_custom_role" "all" {
  name = huaweicloud_identity_role.test.name
}

output "is_name_filter_useful" {
  value = data.huaweicloud_identity_custom_role.all.name == "%[2]s"
}
`, testAccDataCustomRole_base(name), name)
}
