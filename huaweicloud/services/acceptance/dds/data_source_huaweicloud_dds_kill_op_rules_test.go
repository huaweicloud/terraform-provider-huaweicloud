package dds

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceKillOpRules_basic(t *testing.T) {
	var (
		dataSource = "data.huaweicloud_dds_kill_op_rules.test"
		dc         = acceptance.InitDataSourceCheck(dataSource)
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDDSInstanceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceKillOpRules_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.operation_types"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.namespaces"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.client_ips"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.plan_summary"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.node_type"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.max_concurrency"),
					resource.TestCheckResourceAttrSet(dataSource, "rules.0.secs_running"),

					resource.TestCheckOutput("operation_types_filter_useful", "true"),
					resource.TestCheckOutput("namespaces_filter_useful", "true"),
					resource.TestCheckOutput("status_filter_useful", "true"),
					resource.TestCheckOutput("plan_summary_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataSourceKillOpRules_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dds_kill_op_rules" "test" {
  instance_id = "%[1]s"
}

locals {
  operation_types = data.huaweicloud_dds_kill_op_rules.test.rules[0].operation_types
  namespaces      = data.huaweicloud_dds_kill_op_rules.test.rules[0].namespaces
  status          = data.huaweicloud_dds_kill_op_rules.test.rules[0].status
  plan_summary    = data.huaweicloud_dds_kill_op_rules.test.rules[0].plan_summary
}

data "huaweicloud_dds_kill_op_rules" "operation_types_filter" {
  instance_id     = "%[1]s"
  operation_types = local.operation_types
}

output "operation_types_filter_useful" {
  value = length(data.huaweicloud_dds_kill_op_rules.operation_types_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_dds_kill_op_rules.operation_types_filter.rules[*].operation_types : v == local.operation_types]
  )
}

data "huaweicloud_dds_kill_op_rules" "namespaces_filter" {
  instance_id = "%[1]s"	
  namespaces  = local.namespaces
}

output "namespaces_filter_useful" {
  value = length(data.huaweicloud_dds_kill_op_rules.namespaces_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_dds_kill_op_rules.namespaces_filter.rules[*].namespaces : v == local.namespaces]
  )
}

data "huaweicloud_dds_kill_op_rules" "status_filter" {
  instance_id = "%[1]s"		
  status      = local.status
}

output "status_filter_useful" {
  value = length(data.huaweicloud_dds_kill_op_rules.status_filter.rules) > 0 && alltrue(
    [for v in data.huaweicloud_dds_kill_op_rules.status_filter.rules[*].status : v == local.status]
  )
}

data "huaweicloud_dds_kill_op_rules" "plan_summary_filter" {
  instance_id  = "%[1]s"	
  plan_summary = local.plan_summary
}

output "plan_summary_filter_useful" {
  value = length(data.huaweicloud_dds_kill_op_rules.plan_summary_filter.rules) > 0 && alltrue(	
    [for v in data.huaweicloud_dds_kill_op_rules.plan_summary_filter.rules[*].plan_summary : v == local.plan_summary]
  )
}
`, acceptance.HW_DDS_INSTANCE_ID)
}
