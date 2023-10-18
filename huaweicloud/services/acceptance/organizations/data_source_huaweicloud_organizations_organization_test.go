package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceOrganization_basic(t *testing.T) {
	rName := "data.huaweicloud_organizations_organization.test"
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
				Config: testAccDatasourceOrganization_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "root_tags.key1", "value1"),
					resource.TestCheckResourceAttr(rName, "root_tags.key2", "value2"),
					resource.TestCheckResourceAttr(rName, "enabled_policy_types.0", "service_control_policy"),
				),
			},
		},
	})
}

func testAccDatasourceOrganization_basic() string {
	return `data "huaweicloud_organizations_organization" "test" {}`
}
