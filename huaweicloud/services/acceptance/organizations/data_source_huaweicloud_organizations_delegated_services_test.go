package organizations

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataDelegatedServices_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_organizations_delegated_services.test"
		dc  = acceptance.InitDataSourceCheck(all)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckOrganizationsInviteAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataDelegatedServices_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(all, "delegated_services.#"),
					resource.TestCheckResourceAttrSet(all, "delegated_services.0.service_principal"),
				),
			},
		},
	})
}

func testAccDataDelegatedServices_basic() string {
	return fmt.Sprintf(`
# Without any filter parameters.
data "huaweicloud_organizations_delegated_services" "test" {
  account_id = "%s"
}
`, acceptance.HW_ORGANIZATIONS_INVITE_ACCOUNT_ID)
}
