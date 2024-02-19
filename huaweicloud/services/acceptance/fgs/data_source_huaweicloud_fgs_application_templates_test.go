package fgs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccFunctionGraphApplicationTemplates_basic(t *testing.T) {
	var (
		dataSourceName = "data.huaweicloud_fgs_application_templates.test"
		dc             = acceptance.InitDataSourceCheck(dataSourceName)

		byRuntime   = "data.huaweicloud_fgs_application_templates.filter_by_runtime"
		dcByRuntime = acceptance.InitDataSourceCheck(byRuntime)

		byCategory   = "data.huaweicloud_fgs_application_templates.filter_by_category"
		dcByCategory = acceptance.InitDataSourceCheck(byCategory)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFunctionGraphApplicationTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),

					dcByRuntime.CheckResourceExists(),
					resource.TestCheckOutput("is_runtime_filter_useful", "true"),
					resource.TestCheckResourceAttrSet(dataSourceName, "templates.#"),
					resource.TestCheckResourceAttrSet(byRuntime, "templates.0.id"),
					resource.TestCheckResourceAttrSet(byRuntime, "templates.0.name"),

					dcByCategory.CheckResourceExists(),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccFunctionGraphApplicationTemplates_basic = `
data "huaweicloud_fgs_application_templates" "test" {}

// By runtime filter
locals {
  runtime = data.huaweicloud_fgs_application_templates.test.templates[0].runtime
}

data "huaweicloud_fgs_application_templates" "filter_by_runtime" {
  runtime = local.runtime
}

output "is_runtime_filter_useful" {
  value = length(data.huaweicloud_fgs_application_templates.filter_by_runtime.templates) >= 1 && alltrue(
    [for v in data.huaweicloud_fgs_application_templates.filter_by_runtime.templates[*].runtime : v == local.runtime]
  )
}

// By category filter
locals {
  category = data.huaweicloud_fgs_application_templates.test.templates[0].category
}

data "huaweicloud_fgs_application_templates" "filter_by_category" {
  category = local.category
}

output "is_category_filter_useful" {
  value = length(data.huaweicloud_fgs_application_templates.filter_by_category.templates) > 0 && alltrue(
    [for v in data.huaweicloud_fgs_application_templates.filter_by_category.templates[*].category : v == local.category]
  )
}
`
