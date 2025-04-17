package organizations

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceOrganizationsCreateAccountStatus_basic(t *testing.T) {
	dataSource := "data.huaweicloud_organizations_create_account_status.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceOrganizationsCreateAccountStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "create_account_statuses.#"),
					resource.TestCheckResourceAttrSet(dataSource, "create_account_statuses.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "create_account_statuses.0.state"),
					resource.TestCheckResourceAttrSet(dataSource, "create_account_statuses.0.account_id"),
					resource.TestCheckResourceAttrSet(dataSource, "create_account_statuses.0.account_name"),
					resource.TestCheckResourceAttrSet(dataSource, "create_account_statuses.0.completed_at"),
					resource.TestCheckResourceAttrSet(dataSource, "create_account_statuses.0.created_at"),
					resource.TestCheckOutput("states_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceOrganizationsCreateAccountStatus_basic() string {
	return `
data "huaweicloud_organizations_create_account_status" "test" {}

locals {
  state = "succeeded"
}
data "huaweicloud_organizations_create_account_status" "states_filter" {
  states = [local.state]
}
output "states_filter_is_useful" {
  value = length(data.huaweicloud_organizations_create_account_status.states_filter.create_account_statuses) > 0 && alltrue(
  [for v in data.huaweicloud_organizations_create_account_status.states_filter.create_account_statuses[*].state : v == local.state]
  )  
}
`
}
