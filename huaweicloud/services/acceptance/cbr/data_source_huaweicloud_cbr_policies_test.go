package cbr

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccPoliciesDataSource_basic(t *testing.T) {
	var (
		name           = acceptance.RandomAccResourceName()
		dataSourceName = "data.huaweicloud_cbr_policies.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPoliciesDataSource_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.name"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.type"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.enabled"),
					resource.TestCheckResourceAttrSet(dataSourceName, "policies.0.associated_vaults.0.vault_id"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_enabled_filter_useful", "true"),
					resource.TestCheckOutput("is_vault_id_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccPoliciesDataSource_base(name string) string {
	return fmt.Sprintf(`

resource "huaweicloud_cbr_policy" "test" {
  name            = "%[1]s"
  type            = "backup"
  backup_quantity = 5

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}

resource "huaweicloud_cbr_vault" "test" {
  name            = "%[1]s"
  type            = "server"
  protection_type = "backup"
  size            = 200

  policy {
    id = huaweicloud_cbr_policy.test.id
  }
}
`, name)
}

func testAccPoliciesDataSource_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cbr_policies" "test" {
  depends_on = [
    huaweicloud_cbr_vault.test
  ]  
}

locals {
  policy_id = data.huaweicloud_cbr_policies.test.policies[0].id
}
data "huaweicloud_cbr_policies" "id_filter" {
  policy_id = local.policy_id
}
output "is_id_filter_useful" {
  value = length(data.huaweicloud_cbr_policies.id_filter.policies) > 0 && alltrue( 
    [for v in data.huaweicloud_cbr_policies.id_filter.policies[*].id : v == local.policy_id]
  )  
}

locals {
  name = data.huaweicloud_cbr_policies.test.policies[0].name
}
data "huaweicloud_cbr_policies" "name_filter" {
  name = local.name
}
output "is_name_filter_useful" {
  value = length(data.huaweicloud_cbr_policies.name_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_policies.name_filter.policies[*].name : v == local.name]
  )  
}

locals {
  type = data.huaweicloud_cbr_policies.test.policies[0].type
}
data "huaweicloud_cbr_policies" "type_filter" {
  type = local.type
}
output "is_type_filter_useful" {
  value = length(data.huaweicloud_cbr_policies.type_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_policies.type_filter.policies[*].type : v == local.type]
  )  
}

locals {
  enabled = data.huaweicloud_cbr_policies.test.policies[0].enabled
}
data "huaweicloud_cbr_policies" "enabled_filter" {
  enabled = local.enabled
}
output "is_enabled_filter_useful" {
  value = length(data.huaweicloud_cbr_policies.enabled_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_policies.enabled_filter.policies[*].enabled : v == local.enabled]
  )  
}

locals {
  vault_id = data.huaweicloud_cbr_policies.test.policies[0].associated_vaults[0].vault_id
}
data "huaweicloud_cbr_policies" "vault_id_filter" {
  vault_id = local.vault_id
}
output "is_vault_id_filter_useful" {
  value = length(data.huaweicloud_cbr_policies.vault_id_filter.policies) > 0 && alltrue(
    [for v in data.huaweicloud_cbr_policies.vault_id_filter.policies[*].associated_vaults[0].vault_id : v == local.vault_id]
  )  
}
`, testAccPoliciesDataSource_base(name))
}
