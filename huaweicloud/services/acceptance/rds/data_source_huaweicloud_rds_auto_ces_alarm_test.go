package rds

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceRdsAutoCesAlarm_basic(t *testing.T) {
	dataSource := "data.huaweicloud_rds_auto_ces_alarm.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRdsAutoCesAlarm_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "entities.#"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.domain_id"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.domain_name"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.project_id"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.project_name"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.engine_name"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.new_instance_default"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.switch_status"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.topic_urn"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.created_at"),
					resource.TestCheckResourceAttrSet(dataSource, "entities.0.updated_at"),
				),
			},
		},
	})
}

func testAccDataSourceRdsAutoCesAlarm_basic() string {
	return `
data "huaweicloud_rds_auto_ces_alarm" "test" {}

data "huaweicloud_rds_auto_ces_alarm" "engine_filter" {
  engine = "mysql"
}
output "engine_filter_is_useful" {
  value = length(data.huaweicloud_rds_auto_ces_alarm.engine_filter.entities) > 0 && alltrue(
    [for v in data.huaweicloud_rds_auto_ces_alarm.engine_filter.entities[*].engine_name : v == "mysql"]
  )
}
`
}
