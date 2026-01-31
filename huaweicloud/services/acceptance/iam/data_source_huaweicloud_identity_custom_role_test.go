package iam

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataCustomRole_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		byName   = "data.huaweicloud_identity_custom_role.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byId   = "data.huaweicloud_identity_custom_role.filter_by_id"
		dcById = acceptance.InitDataSourceCheck(byId)
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
				Config: testAccDataCustomRole_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dcByName.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byName, "catalog"),
					dcById.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(byId, "catalog"),
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

# Filter by name
locals {
  name = huaweicloud_identity_role.test.name
}

data "huaweicloud_identity_custom_role" "filter_by_name" {
  name = local.name

  # Waiting for the role to be created.
  depends_on = [huaweicloud_identity_role.test]
}

# Filter by ID
locals {
  id = huaweicloud_identity_role.test.id
}

data "huaweicloud_identity_custom_role" "filter_by_id" {
  id = local.id

  # Waiting for the role to be created.
  depends_on = [huaweicloud_identity_role.test]
}
`, testAccDataCustomRole_base(name))
}
