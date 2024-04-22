package dli

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDliFlinkTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_dli_flink_templates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceDliFlinkTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.id"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "templates.0.type"),

					resource.TestCheckOutput("template_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceDliFlinkTemplates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_dli_flink_templates" "test" {
  depends_on = [
    huaweicloud_dli_flink_template.test
  ]
}

data "huaweicloud_dli_flink_templates" "template_id_filter" {
  template_id = local.template_id
}
  
locals {
  template_id = data.huaweicloud_dli_flink_templates.test.templates[0].id
}
  
output "template_id_filter_is_useful" {
  value = length(data.huaweicloud_dli_flink_templates.template_id_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flink_templates.template_id_filter.templates[*].id : v == local.template_id]
  )
}

data "huaweicloud_dli_flink_templates" "name_filter" {
  name = local.name
}
  
locals {
  name = data.huaweicloud_dli_flink_templates.test.templates[0].name
}
  
output "name_filter_is_useful" {
  value = length(data.huaweicloud_dli_flink_templates.name_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flink_templates.name_filter.templates[*].name : v == local.name]
  )
}

data "huaweicloud_dli_flink_templates" "type_filter" {
  type = local.type
}
  
locals {
  type = data.huaweicloud_dli_flink_templates.test.templates[0].type
}
  
output "type_filter_is_useful" {
  value = length(data.huaweicloud_dli_flink_templates.type_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_dli_flink_templates.type_filter.templates[*].type : v == local.type]
  )
}
`, testFlinkTemplate_basic(name))
}
