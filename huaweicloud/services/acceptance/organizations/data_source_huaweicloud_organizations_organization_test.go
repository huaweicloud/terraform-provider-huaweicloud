package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataOrganization_basic(t *testing.T) {
	var (
		name = acceptance.RandomAccResourceName()

		dcName = "data.huaweicloud_organizations_organization.all"
		dc     = acceptance.InitDataSourceCheck(dcName)
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
				Config: testAccDataOrganization_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dcName, "master_account_name"),
					resource.TestCheckResourceAttrSet(dcName, "master_account_id"),
					resource.TestCheckResourceAttrSet(dcName, "created_at"),
					resource.TestCheckResourceAttrSet(dcName, "root_id"),
					resource.TestCheckResourceAttrSet(dcName, "root_name"),
					resource.TestCheckResourceAttrSet(dcName, "root_urn"),
					resource.TestMatchResourceAttr(dcName, "root_tags.%", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestMatchResourceAttr(dcName, "enabled_policy_types.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
				),
			},
		},
	})
}

func testAccDataOrganization_basic(name string) string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_organization" "test" {}

resource "huaweicloud_tms_resource_tags" "test" {
  resources {
    resource_type = "organizations:roots"
    resource_id   = data.huaweicloud_organizations_organization.test.root_id
  }

  tags = {
    "%[1]s" = "%[1]s"
  }
}

data "huaweicloud_organizations_organization" "all" {
  depends_on = [huaweicloud_tms_resource_tags.test]
}
`, name)
}
