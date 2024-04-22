package dcs

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDcsAccounts_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dcs_accounts.all"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDcsAccounts_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "accounts.#"),

					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_role_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
					resource.TestCheckOutput("is_description_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDcsAccounts_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dcs_accounts" "all" {
  depends_on = [huaweicloud_dcs_account.test]

  instance_id = huaweicloud_dcs_instance.instance_1.id
}

// filter by name
data "huaweicloud_dcs_accounts" "filter_by_name" {
  depends_on = [huaweicloud_dcs_account.test]

  instance_id  = huaweicloud_dcs_instance.instance_1.id
  account_name = huaweicloud_dcs_account.test.account_name
}

locals {
  filter_result_by_name = [for v in data.huaweicloud_dcs_accounts.filter_by_name.accounts[*].account_name : 
    v == huaweicloud_dcs_account.test.account_name]
}

output "is_name_filter_useful" {
  value = length(local.filter_result_by_name) == 1 && alltrue(local.filter_result_by_name) 
}

// filter by type
data "huaweicloud_dcs_accounts" "filter_by_type" {
  depends_on = [huaweicloud_dcs_account.test]

  instance_id  = huaweicloud_dcs_instance.instance_1.id
  account_type = huaweicloud_dcs_account.test.account_type
}

locals {
  filter_result_by_type = [for v in data.huaweicloud_dcs_accounts.filter_by_type.accounts[*].account_type : 
    v == huaweicloud_dcs_account.test.account_type]
}

output "is_type_filter_useful" {
  value = length(local.filter_result_by_type) > 0 && alltrue(local.filter_result_by_type) 
}

// filter by role
data "huaweicloud_dcs_accounts" "filter_by_role" {
  depends_on = [huaweicloud_dcs_account.test]

  instance_id  = huaweicloud_dcs_instance.instance_1.id
  account_role = huaweicloud_dcs_account.test.account_role
}

locals {
  filter_result_by_role = [for v in data.huaweicloud_dcs_accounts.filter_by_role.accounts[*].account_role : 
    v == huaweicloud_dcs_account.test.account_role]
}

output "is_role_filter_useful" {
  value = length(local.filter_result_by_role) > 0 && alltrue(local.filter_result_by_role) 
}

// filter by status
data "huaweicloud_dcs_accounts" "filter_by_status" {
  instance_id = huaweicloud_dcs_instance.instance_1.id
  status      = huaweicloud_dcs_account.test.status
}

locals {
  filter_result_by_status = [for v in data.huaweicloud_dcs_accounts.filter_by_status.accounts[*].status : 
    v == huaweicloud_dcs_account.test.status]
}

output "is_status_filter_useful" {
  value = length(local.filter_result_by_status) > 0 && alltrue(local.filter_result_by_status) 
}

// filter by description
data "huaweicloud_dcs_accounts" "filter_by_description" {
  depends_on = [huaweicloud_dcs_account.test]

  instance_id = huaweicloud_dcs_instance.instance_1.id
  description = huaweicloud_dcs_account.test.description
}

locals {
  filter_result_by_description = [for v in data.huaweicloud_dcs_accounts.filter_by_description.accounts[*].description : 
    v == huaweicloud_dcs_account.test.description]
}

output "is_description_filter_useful" {
  value = length(local.filter_result_by_description) > 0 && alltrue(local.filter_result_by_description) 
}
`, testAccDcsAccount_basic(name))
}
