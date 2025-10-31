package rgc

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationalUnitAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_organizational_unit_accounts.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckRGCOrganizationID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceOrganizationalUnitAccountsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSource, "managed_organizational_unit_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.landing_zone_version"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.manage_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.account_name"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.account_type"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.owner"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.parent_organizational_unit_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.parent_organizational_unit_name"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.regions.#"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.regions.0.region"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.regions.0.region_status"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.updated_at"),
					resource.TestCheckResourceAttrSet(dataSource, "region"),
				),
			},
		},
	})
}

func testAccDataSourceOrganizationalUnitAccountsConfig() string {
	return fmt.Sprintf(`
data "huaweicloud_rgc_organizational_unit_accounts" "test" {
  managed_organizational_unit_id = "%[1]s"
}
`, acceptance.HW_RGC_ORGANIZATIONAL_UNIT_ID)
}
