package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesAlarmTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_alarm_templates.test"
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesAlarmTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_templates.0.template_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_templates.0.type"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_templates.0.created_at"),
					resource.TestCheckOutput("is_namespace_filter_useful", "true"),
					resource.TestCheckOutput("is_dimension_name_filter_useful", "true"),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesAlarmTemplates_basic() string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_ces_alarm_templates" "test" {
  depends_on = [
    huaweicloud_ces_alarm_template.test1,
    huaweicloud_ces_alarm_template.test2
  ]
}

locals {
  name           = huaweicloud_ces_alarm_template.test1.name
  type           = "custom"
  namespace      = "SYS.APIG"
  dimension_name = "mongodb_instance_id"
}

data "huaweicloud_ces_alarm_templates" "filter_by_namespace" {
  namespace = local.namespace
  
  depends_on = [
    huaweicloud_ces_alarm_template.test1,
    huaweicloud_ces_alarm_template.test2
  ]
}

output "is_namespace_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_templates.filter_by_namespace.alarm_templates) >= 1 
}

data "huaweicloud_ces_alarm_templates" "filter_by_dimension_name" {
  dimension_name = local.dimension_name
	
  depends_on = [
    huaweicloud_ces_alarm_template.test1,
    huaweicloud_ces_alarm_template.test2
  ]
}

output "is_dimension_name_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_templates.filter_by_dimension_name.alarm_templates) >= 1
}

data "huaweicloud_ces_alarm_templates" "filter_by_type" {
  type = local.type
	  
  depends_on = [
    huaweicloud_ces_alarm_template.test1,
    huaweicloud_ces_alarm_template.test2
  ]
}

output "is_type_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_templates.filter_by_type.alarm_templates) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_alarm_templates.filter_by_type.alarm_templates[*] : item.type == local.type]
  )
}

data "huaweicloud_ces_alarm_templates" "filter_by_name" {
  name = local.name

  depends_on = [
    huaweicloud_ces_alarm_template.test1,
    huaweicloud_ces_alarm_template.test2
  ]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_templates.filter_by_name.alarm_templates) >= 1 && alltrue(
    [for item in data.huaweicloud_ces_alarm_templates.filter_by_name.alarm_templates[*] : item.name == local.name]
  )
}
`, testDataSourceCesAlarmTemplates_base())
}

func testDataSourceCesAlarmTemplates_base() string {
	name := acceptance.RandomAccResourceName()
	return fmt.Sprintf(`
resource "huaweicloud_ces_alarm_template" "test1" {
  name        = "%[1]s1" 
  description = "It is template"
	  
  policies {
    namespace           = "SYS.APIG"
    dimension_name      = "api_id"
    metric_name         = "req_count_2xx"
    period              = 1
    filter              = "average"
    comparison_operator = "="
    value               = "10"
    unit                = "times/minute"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 300
  }
}

resource "huaweicloud_ces_alarm_template" "test2" {
  name        = "%[1]s2"
  description = "It is template"
  
  policies {
    namespace           = "SYS.DDS"
    dimension_name      = "mongodb_instance_id"
    metric_name         = "mongo003_insert_ps"
    period              = 300
    filter              = "max"
    comparison_operator = "<"
    value               = "300"
    unit                = "times/second"
    count               = 5
    alarm_level         = 3
    suppress_duration   = 3600
  }
}
`, name)
}
