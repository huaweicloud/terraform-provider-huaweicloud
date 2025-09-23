package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsEffectivePolicies_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_effective_policies.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

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
				Config: testDataSourceOrganizationsEffectivePolicies_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "last_updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "policy_content"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsEffectivePolicies_base(rName string) string {
	return fmt.Sprintf(`
resource "huaweicloud_organizations_policy" "test"{
  name        = "%[1]s"
  type        = "tag_policy"
  description = "test description"
  content     = jsonencode(
    {
      "tags":{
        "test_tag":{
          "tag_key":{
            "@@assign":"test_tag"
          }
        }
      }
    }
  )
}

resource "huaweicloud_organizations_policy_attach" "test" {
  policy_id = huaweicloud_organizations_policy.test.id
  entity_id = "%[2]s"
}
`, rName, acceptance.HW_ORGANIZATIONS_ACCOUNT_ID)
}

func testDataSourceOrganizationsEffectivePolicies_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_organizations_effective_policies" "test" {
  depends_on = [huaweicloud_organizations_policy_attach.test]

  entity_id   = "%[2]s"
  policy_type = "tag_policy"
}
`, testDataSourceOrganizationsEffectivePolicies_base(rName), acceptance.HW_ORGANIZATIONS_ACCOUNT_ID)
}
