package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsTagPolicyServices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_tag_policy_services.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccount(t)
			acceptance.TestAccPreCheckOrganizationsOpen(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationsTagPolicyServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "services.#"),
					resource.TestCheckResourceAttrSet(dataSource, "services.0.service_name"),
					resource.TestCheckResourceAttrSet(dataSource, "services.0.resource_types.#"),
					resource.TestCheckResourceAttrSet(dataSource, "services.0.resource_types.0"),
					resource.TestCheckResourceAttrSet(dataSource, "services.0.support_all"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsTagPolicyServices_basic() string {
	return `
data "huaweicloud_organizations_tag_policy_services" "test" {}
`
}
