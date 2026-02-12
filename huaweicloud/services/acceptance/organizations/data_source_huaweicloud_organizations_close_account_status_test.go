package organizations

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Please make sure to have at least one account in the organization that is either pending closure or closed.
func TestAccDataCloseAccountStatus_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_organizations_close_account_status.test"
		dc  = acceptance.InitDataSourceCheck(all)

		byStates   = "data.huaweicloud_organizations_close_account_status.filter_by_states"
		dcByStates = acceptance.InitDataSourceCheck(byStates)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			// Due to quota limit, use environment variable to ensure that there is at least one account that is either
			// pending closure or closed.
			acceptance.TestAccPreCheckOrganizationsCloseAccountId(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataCloseAccountStatus_basic(),
				Check: resource.ComposeTestCheckFunc(
					// Without any filter parameters.
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "close_account_statuses.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "close_account_statuses.0.account_id"),
					resource.TestCheckResourceAttrSet(all, "close_account_statuses.0.state"),
					resource.TestCheckResourceAttrSet(all, "close_account_statuses.0.organization_id"),
					resource.TestCheckResourceAttrSet(all, "close_account_statuses.0.created_at"),
					resource.TestCheckResourceAttrSet(all, "close_account_statuses.0.updated_at"),
					// Filter by 'states' parameter.
					dcByStates.CheckResourceExists(),
					resource.TestCheckOutput("is_states_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataCloseAccountStatus_basic() string {
	return fmt.Sprintf(`
# Without any filter parameters.
data "huaweicloud_organizations_close_account_status" "test" {}

output "is_account_id_exists" {
  value = contains(data.huaweicloud_organizations_close_account_status.test.close_account_statuses[*].account_id, "%[1]s")
}

locals {
  state = try(data.huaweicloud_organizations_close_account_status.test.close_account_statuses[0].state, null)
}

# Filter by 'states' parameter.
data "huaweicloud_organizations_close_account_status" "filter_by_states" {
  states = [local.state]
}

locals {
  states_filter_result = [for v in data.huaweicloud_organizations_close_account_status.filter_by_states.close_account_statuses[*].state :
  v == local.state]
}

output "is_states_filter_useful" {
  value = length(local.states_filter_result) > 0 && alltrue(local.states_filter_result)
}
`, acceptance.HW_ORGANIZATIONS_CLOSE_ACCOUNT_ID)
}
