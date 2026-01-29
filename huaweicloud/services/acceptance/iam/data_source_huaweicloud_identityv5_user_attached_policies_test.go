package iam

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please ensure that the user executing the acceptance test has 'admin' permission.
func TestAccDataV5UserAttachedPolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_identityv5_user_attached_policies.test"
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
				Config: testAccDataV5UserAttachedPolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "user_id"),
					resource.TestMatchResourceAttr(all, "attached_policies.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrPair(all, "attached_policies.0.policy_id", "huaweicloud_identity_policy.test", "id"),
					resource.TestCheckResourceAttrPair(all, "attached_policies.0.policy_name", "huaweicloud_identity_policy.test", "name"),
					resource.TestCheckResourceAttrPair(all, "attached_policies.0.urn", "huaweicloud_identity_policy.test", "urn"),
					resource.TestCheckResourceAttrSet(all, "attached_policies.0.attached_at"),
				),
			},
		},
	})
}

func testAccDataV5UserAttachedPolicies_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_identityv5_user" "test" {
  name = "%[1]s"
}

resource "huaweicloud_identity_policy" "test" {
  name            = "%[1]s"
  description     = "test policy for terraform"
  policy_document = jsonencode(
    {
      Statement = [
        {
          Action = ["*"]
          Effect = "Allow"
        }
      ]
      Version = "5.0"
    }
  )
}

resource "huaweicloud_identityv5_policy_user_attach" "test" {
  policy_id = huaweicloud_identity_policy.test.id
  user_id   = huaweicloud_identityv5_user.test.id
}
`, name)
}

func testAccDataV5UserAttachedPolicies_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_identityv5_user_attached_policies" "test" {
  user_id = huaweicloud_identityv5_user.test.id

  depends_on = [huaweicloud_identityv5_policy_user_attach.test]
}
`, testAccDataV5UserAttachedPolicies_base(name))
}
