package rms

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDatasourceTemplates_basic(t *testing.T) {
	rName := "data.huaweicloud_rms_assignment_package_templates.test"
	dc := acceptance.InitDataSourceCheck(rName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(rName, "templates.0.id"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.template_key"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.description"),
					resource.TestCheckResourceAttrSet(rName, "templates.0.template_body"),
					resource.TestCheckOutput("template_key_filter_is_useful", "true"),
					resource.TestCheckOutput("description_filter_is_useful", "true"),
				),
			},
		},
	})
}

func testAccDatasourceTemplates_basic() string {
	return `
data "huaweicloud_rms_assignment_package_templates" "test" {
}

data "huaweicloud_rms_assignment_package_templates" "template_key_filter" {
  template_key = data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key
}

data "huaweicloud_rms_assignment_package_templates" "description_filter" {
  description = data.huaweicloud_rms_assignment_package_templates.test.templates.0.description
}

locals {
  template_key_filter_result = [for v in data.huaweicloud_rms_assignment_package_templates.template_key_filter.
  templates[*].template_key:v == data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key]
  description_filter_result = [for v in data.huaweicloud_rms_assignment_package_templates.description_filter.
  templates[*].description:v == data.huaweicloud_rms_assignment_package_templates.test.templates.0.description]
}

output "template_key_filter_is_useful" {
  value = alltrue(local.template_key_filter_result) && length(local.template_key_filter_result) > 0
}

output "description_filter_is_useful" {
  value = alltrue(local.description_filter_result) && length(local.description_filter_result) > 0
}`
}
