package aom

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_aom_organization_accounts.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckMultiAccountAggregationRuleEnable(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationAccounts_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "accounts.#"),
					resource.TestCheckResourceAttrSet(dataSource, "accounts.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "accounts.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "accounts.0.urn"),
					resource.TestCheckResourceAttrSet(dataSource, "accounts.0.join_method"),
					resource.TestCheckResourceAttrSet(dataSource, "accounts.0.joined_at"),
				),
			},
		},
	})
}

const testDataSourceOrganizationAccounts_basic string = `data "huaweicloud_aom_organization_accounts" "test" {}`
