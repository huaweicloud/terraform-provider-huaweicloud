package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Currently, only account-level effective policy queries are supported.
func TestAccDataEffectivePolicies_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_organizations_effective_policies.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
			acceptance.TestAccPreCheckOrganizationsAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataEffectivePolicies_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "last_updated_at"),
					resource.TestCheckResourceAttrSet(all, "policy_content"),
				),
			},
		},
	})
}

func testAccDataEffectivePolicies_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test" {
  name        = "%[1]s"
  type        = "tag_policy"
  description = "Created by terraform script"
  content     = jsonencode({
    tags = {
      test_tag = {
        tag_key = {
          "@@assign" = "test_tag"
        }
      }
    }
  })
}

resource "huaweicloud_organizations_policy_attach" "test" {
  policy_id = huaweicloud_organizations_policy.test.id
  entity_id = "%[2]s"
}
`, rName, acceptance.HW_ORGANIZATIONS_ACCOUNT_ID)
}

func testAccDataEffectivePolicies_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

# Without any filter parameters.
data "huaweicloud_organizations_effective_policies" "test" {
  depends_on = [huaweicloud_organizations_policy_attach.test]

  entity_id   = "%[2]s"
  policy_type = "tag_policy"
}
`, testAccDataEffectivePolicies_base(rName), acceptance.HW_ORGANIZATIONS_ACCOUNT_ID)
}
