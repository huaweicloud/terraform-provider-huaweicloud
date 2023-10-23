package lts

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceCustomTemplates_basic(t *testing.T) {
	rName := "data.huaweicloud_lts_structuring_custom_templates.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckLtsEnableFlag(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceCustomTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "templates.0.id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.name"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.type"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.demo_log"),

					resource.TestCheckOutput("template_id_filter_is_useful", "true"),
					resource.TestCheckOutput("name_filter_is_useful", "true"),
					resource.TestCheckOutput("type_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceCustomTemplates_basic() string {
	return `
data "huaweicloud_lts_structuring_custom_templates" "test" {
}

data "huaweicloud_lts_structuring_custom_templates" "template_id_filter" {
  template_id = data.huaweicloud_lts_structuring_custom_templates.test.templates.0.id
}

output "template_id_filter_is_useful" {
  value = length(data.huaweicloud_lts_structuring_custom_templates.template_id_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_lts_structuring_custom_templates.template_id_filter.templates[*].id :
    v == data.huaweicloud_lts_structuring_custom_templates.template_id_filter.template_id]
  )
}

data "huaweicloud_lts_structuring_custom_templates" "name_filter" {
  name = data.huaweicloud_lts_structuring_custom_templates.test.templates.0.name
}

output "name_filter_is_useful" {
  value = length(data.huaweicloud_lts_structuring_custom_templates.name_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_lts_structuring_custom_templates.name_filter.templates[*].name :
    v == data.huaweicloud_lts_structuring_custom_templates.name_filter.name]
  )
}

data "huaweicloud_lts_structuring_custom_templates" "type_filter" {
  type = data.huaweicloud_lts_structuring_custom_templates.test.templates.0.type
}

output "type_filter_is_useful" {
  value = length(data.huaweicloud_lts_structuring_custom_templates.type_filter.templates) > 0 && alltrue(
    [for v in data.huaweicloud_lts_structuring_custom_templates.type_filter.templates[*].type :
    v == data.huaweicloud_lts_structuring_custom_templates.type_filter.type]
  )
}
`
}
