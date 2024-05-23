package css

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCssLogstashTemplates_basic(t *testing.T) {
	dataSource := "data.huaweicloud_css_logstash_templates.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDataSourceCssLogstashTemplates_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "system_templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "system_templates.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "custom_templates.#"),
					resource.TestCheckResourceAttrSet(dataSource, "custom_templates.0.name"),

					resource.TestCheckOutput("type_filter_is_useful", "true"),
					resource.TestCheckOutput("system_template_name_filter_is_useful", "true"),
					resource.TestCheckOutput("custom_template_name_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceDataSourceCssLogstashTemplates_basic(name string) string {
	return fmt.Sprintf(`
%s

data "huaweicloud_css_logstash_templates" "test" {
  depends_on = [huaweicloud_css_logstash_custom_template.test]
}

locals {
  system_template_name = data.huaweicloud_css_logstash_templates.test.system_templates[0].name
  custom_template_name = data.huaweicloud_css_logstash_templates.test.custom_templates[0].name
}

data "huaweicloud_css_logstash_templates" "filter_by_type" {
  depends_on = [huaweicloud_css_logstash_custom_template.test]

  type = "custom"
}

data "huaweicloud_css_logstash_templates" "filter_by_system_template_name" {
  name = local.system_template_name
}

data "huaweicloud_css_logstash_templates" "filter_by_custom_template_name" {
  depends_on = [huaweicloud_css_logstash_custom_template.test]

  name = local.custom_template_name
}

locals {
  system_list_by_type = data.huaweicloud_css_logstash_templates.filter_by_type.system_templates
  custom_list_by_type = data.huaweicloud_css_logstash_templates.filter_by_type.custom_templates
  system_list_by_name = data.huaweicloud_css_logstash_templates.filter_by_system_template_name.system_templates
  custom_list_by_name = data.huaweicloud_css_logstash_templates.filter_by_custom_template_name.custom_templates
}

output "type_filter_is_useful" {
  value = length(local.system_list_by_type) == 0 && length(local.custom_list_by_type) > 0
}

output "system_template_name_filter_is_useful" {
  value = length(local.system_list_by_name) > 0 && alltrue(
    [for v in local.system_list_by_name[*].name : v == local.system_template_name]
  )
}

output "custom_template_name_filter_is_useful" {
  value = length(local.custom_list_by_name) > 0 && alltrue(
	[for v in local.custom_list_by_name[*].name : v == local.custom_template_name]
  )
}
`, testLogstashCustomTemplate_basic(name))
}
