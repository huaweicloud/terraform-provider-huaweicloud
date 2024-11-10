package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

// Before running the test case, please ensure that there is at least one WAF instance in the current region.
func TestAccDataSourceAddressGroups_basic(t *testing.T) {
	var (
		name  = acceptance.RandomAccResourceName()
		rName = "data.huaweicloud_waf_address_groups.test"
		dc    = acceptance.InitDataSourceCheck(rName)

		byName   = "data.huaweicloud_waf_address_groups.name_filter"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byIpAddress   = "data.huaweicloud_waf_address_groups.ip_address_filter"
		dcByIpAddress = acceptance.InitDataSourceCheck(byIpAddress)

		byAllParameters   = "data.huaweicloud_waf_address_groups.all_parameters_filter"
		dcByAllParameters = acceptance.InitDataSourceCheck(byAllParameters)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAddressGroups_basic(name),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "groups.0.id"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.name"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.ip_addresses"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.enterprise_project_id"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.description"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.share_count"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.accept_count"),
					resource.TestCheckResourceAttrSet(rName, "groups.0.process_status"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("name_filter_is_useful", "true"),

					dcByIpAddress.CheckResourceExists(),
					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),

					dcByAllParameters.CheckResourceExists(),
					resource.TestCheckOutput("all_parameters_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceAddressGroups_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_waf_address_groups" "test" {
  enterprise_project_id = "%[2]s"

  depends_on = [huaweicloud_waf_address_group.test]
}

# Filter by name
locals {
  name = data.huaweicloud_waf_address_groups.test.groups.0.name
}

data "huaweicloud_waf_address_groups" "name_filter" {
  enterprise_project_id = "%[2]s"
  name                  = local.name

  depends_on = [huaweicloud_waf_address_group.test]
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_waf_address_groups.name_filter.groups[*].name : v == local.name
  ]
}

output "name_filter_is_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by ip_address
locals {
  ip_address = data.huaweicloud_waf_address_groups.test.groups.0.ip_addresses
}

data "huaweicloud_waf_address_groups" "ip_address_filter" {
  enterprise_project_id = "%[2]s"
  ip_address            = local.ip_address

  depends_on = [huaweicloud_waf_address_group.test]
}

locals {
  ip_address_filter_result = [
    for v in data.huaweicloud_waf_address_groups.ip_address_filter.groups[*].ip_addresses : v == local.ip_address
  ]
}

output "ip_address_filter_is_useful" {
  value = length(local.ip_address_filter_result) > 0 && alltrue(local.ip_address_filter_result)
}

# Filter by all parameters
data "huaweicloud_waf_address_groups" "all_parameters_filter" {
  enterprise_project_id = "%[2]s"
  name                  = local.name
  ip_address            = local.ip_address

  depends_on = [huaweicloud_waf_address_group.test]
}

locals {
  all_parameters_filter_result = [
    for v in data.huaweicloud_waf_address_groups.all_parameters_filter.groups[*] :
    v.name == local.name && v.ip_addresses == local.ip_address
  ]
}

output "all_parameters_filter_is_useful" {
  value = length(local.all_parameters_filter_result) > 0 && alltrue(local.all_parameters_filter_result)
}
`, testAddressGroup_basic(name), acceptance.HW_ENTERPRISE_PROJECT_ID_TEST)
}
