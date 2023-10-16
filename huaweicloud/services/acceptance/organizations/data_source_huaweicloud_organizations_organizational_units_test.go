package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceOrganizationalUnits_basic(t *testing.T) {
	rName := "data.huaweicloud_organizations_organizational_units.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceOrganizationalUnits_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "children.#"),
					resource.TestCheckResourceAttrSet(rName, "children.0.id"),
					resource.TestCheckResourceAttrSet(rName, "children.0.name"),
					resource.TestCheckResourceAttrSet(rName, "children.0.urn"),
					resource.TestCheckResourceAttrSet(rName, "children.0.created_at"),
				),
			},
		},
	})
}

func testAccDatasourceOrganizationalUnits_basic() string {
	return fmt.Sprintf(`
%s

data "huaweicloud_organizations_organizational_units" "test" {
  parent_id = data.huaweicloud_organizations_organization.test.root_id
}
`, testAccDatasourceOrganization_basic())
}
