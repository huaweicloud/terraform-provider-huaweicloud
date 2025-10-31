package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsDelegatedServices_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_delegated_services.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationsDelegatedServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "delegated_services.#"),
					resource.TestCheckResourceAttrSet(dataSource, "delegated_services.0.service_principal"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsDelegatedServices_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_organizations_delegated_services" "test" {
  account_id = "%s"
}
`, acceptance.HW_ORGANIZATIONS_INVITE_ACCOUNT_ID)
}
