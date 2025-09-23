package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsTrustedServices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_trusted_services.test"
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
				Config: testDataSourceOrganizationsTrustedServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "trusted_services.#"),
					resource.TestCheckResourceAttrSet(dataSource, "trusted_services.0.service_principal"),
					resource.TestCheckResourceAttrSet(dataSource, "trusted_services.0.enabled_at"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsTrustedServices_basic() string {
	return `
resource "huaweicloud_organizations_trusted_service" "test" {
  service = "service.SecMaster"
}

data "huaweicloud_organizations_trusted_services" "test" {}
`
}
