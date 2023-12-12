package waf

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceWAFAddressGroups_basic(t *testing.T) {
	name := acceptance.RandomAccResourceName()
	rName := "data.huaweicloud_waf_address_groups.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPrecheckWafInstance(t)
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

					resource.TestCheckOutput("name_filter_is_useful", "true"),

					resource.TestCheckOutput("ip_address_filter_is_useful", "true"),

					resource.TestCheckOutput("enterprise_project_id_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceAddressGroups_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_waf_address_groups" "test" {
  depends_on = [huaweicloud_waf_address_group.test]
}

data "huaweicloud_waf_address_groups" "name_filter" {
  name = data.huaweicloud_waf_address_groups.test.groups.0.name
}

locals {
  name = data.huaweicloud_waf_address_groups.test.groups.0.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_waf_address_groups.name_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_waf_address_groups.name_filter.groups[*].name : v == local.name]
  )  
}

data "huaweicloud_waf_address_groups" "ip_address_filter" {
  ip_address  = data.huaweicloud_waf_address_groups.test.groups.0.ip_addresses
}

locals {
  ip_address = data.huaweicloud_waf_address_groups.test.groups.0.ip_addresses
}

output "ip_address_filter_is_useful" {
  value = length(data.huaweicloud_waf_address_groups.ip_address_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_waf_address_groups.ip_address_filter.groups[*].ip_addresses : v == local.ip_address]
  )  
}

data "huaweicloud_waf_address_groups" "enterprise_project_id_filter" {
  enterprise_project_id  = data.huaweicloud_waf_address_groups.test.groups.0.enterprise_project_id
}
  
locals {
  enterprise_project_id = data.huaweicloud_waf_address_groups.test.groups.0.enterprise_project_id
}
  
output "enterprise_project_id_filter_is_useful" {
  value = length(data.huaweicloud_waf_address_groups.enterprise_project_id_filter.groups) > 0 && alltrue(
    [for v in data.huaweicloud_waf_address_groups.enterprise_project_id_filter.groups[*].enterprise_project_id : v == local.enterprise_project_id]
  )  
}
`, testAddressGroup_basic(name))
}
