package rgc

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rgc_accounts.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceAccounts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.account_type"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.landing_zone_version"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.manage_account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.owner"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.parent_organizational_unit_id"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.parent_organizational_unit_name"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "managed_accounts.0.updated_at"),
				),
			},
		},
	})
}

const testDataSourceDataSourceAccounts_basic = `
data "huaweicloud_rgc_accounts" "test" {}
`
