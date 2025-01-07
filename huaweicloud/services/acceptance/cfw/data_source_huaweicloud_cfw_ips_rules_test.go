package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCfwIpsRules_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_ips_rules.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckCfw(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCfwIpsRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.affected_application"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.default_status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_id"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_level"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_name"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_group"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_status"),
					resource.TestCheckResourceAttrSet(dataSource, "records.0.ips_rules_type"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_like_filter_useful", "true"),
					resource.TestCheckOutput("is_status_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCfwIpsRules_basic() string {
	return fmt.Sprintf(`
%[1]s

locals {
  ips_id        = "340710"
  ips_name_like = "web"
  ips_status    = "OBSERVE"
}

data "huaweicloud_cfw_ips_rules" "test" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_rules.test.records) >= 1
}

data "huaweicloud_cfw_ips_rules" "filter_by_id" {
  object_id = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  ips_id    = local.ips_id
}

output "is_id_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_rules.filter_by_id.records) >= 1 && alltrue(
    [for r in data.huaweicloud_cfw_ips_rules.filter_by_id.records[*] : r.ips_id == local.ips_id]
  )
}

data "huaweicloud_cfw_ips_rules" "filter_by_name_like" {
  object_id     = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  ips_name_like = local.ips_name_like
}

output "is_name_like_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_rules.filter_by_name_like.records) >= 1 && alltrue(
    [for r in data.huaweicloud_cfw_ips_rules.filter_by_name_like.records[*] : can(regex("(?i).*web.*", r.ips_name))]
  )
}

data "huaweicloud_cfw_ips_rules" "filter_by_status" {
  object_id  = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  ips_status = local.ips_status
}

output "is_status_filter_useful" {
  value = length(data.huaweicloud_cfw_ips_rules.filter_by_status.records) >= 1 && alltrue(
    [for r in data.huaweicloud_cfw_ips_rules.filter_by_status.records[*] : r.ips_status == local.ips_status]
  )
}
`, testAccDatasourceFirewalls_basic())
}
