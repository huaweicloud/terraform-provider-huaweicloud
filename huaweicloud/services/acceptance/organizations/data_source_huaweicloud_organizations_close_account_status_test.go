package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsCloseAccountStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_close_account_status.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationsCloseAccountStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "close_account_statuses.#"),
					resource.TestCheckResourceAttrSet(dataSource, "close_account_statuses.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "close_account_statuses.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "close_account_statuses.0.organization_id"),
					resource.TestCheckResourceAttrSet(dataSource, "close_account_statuses.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "close_account_statuses.0.updated_at"),
					resource.TestCheckOutput("states_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsCloseAccountStatus_basic() string {
	return `
data "huaweicloud_organizations_close_account_status" "test" {}

locals {
  state = "pending_closure"
}
data "huaweicloud_organizations_close_account_status" "states_filter" {
  states = ["pending_closure"]
}
output "states_filter_is_useful" {
  value = length(data.huaweicloud_organizations_close_account_status.states_filter.close_account_statuses) > 0 && alltrue(
  [for v in data.huaweicloud_organizations_close_account_status.states_filter.close_account_statuses[*].state : v == local.state]
  )  
}
`
}
