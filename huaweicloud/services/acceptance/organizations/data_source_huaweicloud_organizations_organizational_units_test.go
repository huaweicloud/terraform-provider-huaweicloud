package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataOrganizationalUnits_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		all = "data.huaweicloud_organizations_organizational_units.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataOrganizationalUnits_basic(name),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "children.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "children.0.id"),
					resource.TestCheckResourceAttrSet(all, "children.0.name"),
					resource.TestCheckResourceAttrSet(all, "children.0.urn"),
					resource.TestCheckResourceAttrSet(all, "children.0.created_at"),
				),
			},
		},
	})
}

func testAccDataOrganizationalUnits_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_organizations_organizational_unit" "test" {
  name      = "%[1]s"
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}

# Without any filter parameters.
data "huaweicloud_organizations_organizational_units" "test" {
  parent_id = data.huaweicloud_organizations_organization.test.root_id

  depends_on = [huaweicloud_organizations_organizational_unit.test]
}
`, name)
}
