package fgs

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataApplicationTemplates_basic(t *testing.T) {
	var (
		all               = "data.huaweicloud_fgs_application_templates.all"
		dcForAllTemplates = acceptance.InitDataSourceCheck(all)

		byRuntime           = "data.huaweicloud_fgs_application_templates.filter_by_runtime"
		dcByRuntime         = acceptance.InitDataSourceCheck(byRuntime)
		byNotFoundRuntime   = "data.huaweicloud_fgs_application_templates.filter_by_not_found_runtime"
		dcByNotFoundRuntime = acceptance.InitDataSourceCheck(byNotFoundRuntime)

		byCategory           = "data.huaweicloud_fgs_application_templates.filter_by_category"
		dcByCategory         = acceptance.InitDataSourceCheck(byCategory)
		byNotFoundCategory   = "data.huaweicloud_fgs_application_templates.filter_by_not_found_category"
		dcByNotFoundCategory = acceptance.InitDataSourceCheck(byNotFoundCategory)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataApplicationTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					// Without filter parameters.
					dcForAllTemplates.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "templates.#", regexp.MustCompile(`[1-9][0-9]*`)),
					// Filter by application runtime.
					dcByRuntime.CheckResourceExists(),
					resource.TestCheckOutput("is_runtime_filter_useful", "true"),
					dcByNotFoundRuntime.CheckResourceExists(),
					resource.TestCheckOutput("runtime_not_found_validation_pass", "true"),
					// Filter by application category.
					dcByCategory.CheckResourceExists(),
					resource.TestCheckOutput("is_category_filter_useful", "true"),
					dcByNotFoundCategory.CheckResourceExists(),
					resource.TestCheckOutput("category_not_found_validation_pass", "true"),
					// Check the attributes.
					resource.TestCheckResourceAttrSet(all, "templates.0.id"),
					resource.TestCheckResourceAttrSet(all, "templates.0.name"),
					resource.TestCheckResourceAttrSet(all, "templates.0.runtime"),
					resource.TestCheckResourceAttrSet(all, "templates.0.category"),
					resource.TestCheckResourceAttrSet(all, "templates.0.description"),
					resource.TestCheckResourceAttrSet(all, "templates.0.type"),
				),
			},
		},
	})
}

const testAccDataApplicationTemplates_basic = `
# Without any filter parameter.
data "huaweicloud_fgs_application_templates" "all" {}

# Filter by application runtime.
locals {
  application_runtime = data.huaweicloud_fgs_application_templates.all.templates[0].runtime
}

data "huaweicloud_fgs_application_templates" "filter_by_runtime" {
  runtime = local.application_runtime
}

data "huaweicloud_fgs_application_templates" "filter_by_not_found_runtime" {
  runtime = "runtime_not_found"
}

locals {
  runtime_filter_result = [for v in data.huaweicloud_fgs_application_templates.filter_by_runtime.templates[*].runtime :
    v == local.application_runtime]
}

output "is_runtime_filter_useful" {
  value = length(local.runtime_filter_result) > 0 && alltrue(local.runtime_filter_result)
}

output "runtime_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_application_templates.filter_by_not_found_runtime.templates) == 0
}

# Filter by application category.
locals {
  application_category = data.huaweicloud_fgs_application_templates.all.templates[0].category
}

data "huaweicloud_fgs_application_templates" "filter_by_category" {
  category = local.application_category
}

data "huaweicloud_fgs_application_templates" "filter_by_not_found_category" {
  category = "category_not_found"
}

locals {
  category_filter_result = [for v in data.huaweicloud_fgs_application_templates.filter_by_category.templates[*].category :
    v == local.application_category]
}

output "is_category_filter_useful" {
  value = length(local.category_filter_result) > 0 && alltrue(local.category_filter_result)
}

output "category_not_found_validation_pass" {
  value = length(data.huaweicloud_fgs_application_templates.filter_by_not_found_category.templates) == 0
}
`
