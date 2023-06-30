package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccAccountAssociate_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_organizations_account_associate.test"

	// lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheckOrganizations(t)
			acceptance.TestAccPreCheckOrganizationsOrganizationalUnitId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccountAssociate_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "account_id",
						"huaweicloud_organizations_account.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "parent_id",
						"huaweicloud_organizations_organizational_unit.test", "id"),
					resource.TestCheckResourceAttrSet(rName, "name"),
					resource.TestCheckResourceAttrSet(rName, "urn"),
					resource.TestCheckResourceAttrSet(rName, "joined_at"),
					resource.TestCheckResourceAttrSet(rName, "joined_method"),
				),
			},
			{
				Config: testAccountAssociate_basic_update(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(rName, "account_id",
						"huaweicloud_organizations_account.test", "id"),
					resource.TestCheckResourceAttr(rName, "parent_id",
						acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID),
				),
			},
		},
	})
}

func testAccountAssociate_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[2]s"
  parent_id = huaweicloud_organizations_organization.test.root_id
}

resource "huaweicloud_organizations_account_associate" "test" {
  account_id = huaweicloud_organizations_account.test.id
  parent_id  = huaweicloud_organizations_organizational_unit.test.id
}
`, testAccount_basic(name), name)
}

func testAccountAssociate_basic_update(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_organizations_account_associate" "test" {
  account_id = huaweicloud_organizations_account.test.id
  parent_id  = "%s"
}
`, testAccount_basic(name), name, acceptance.HW_ORGANIZATIONS_ORGANIZATIONAL_UNIT_ID)
}
