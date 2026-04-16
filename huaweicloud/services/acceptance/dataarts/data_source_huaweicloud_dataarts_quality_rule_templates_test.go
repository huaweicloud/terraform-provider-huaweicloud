package dataarts

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataQualityRuleTemplates_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_dataarts_quality_rule_templates.all"
		dc  = acceptance.InitDataSourceCheck(all)

		// filter by name
		byName   = "data.huaweicloud_dataarts_quality_rule_templates.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		// filter by not found name
		byNotFoundName   = "data.huaweicloud_dataarts_quality_rule_templates.filter_by_not_found_name"
		dcByNotFoundName = acceptance.InitDataSourceCheck(byNotFoundName)

		// filter by category ID
		byCategoryId   = "data.huaweicloud_dataarts_quality_rule_templates.filter_by_category_id"
		dcByCategoryId = acceptance.InitDataSourceCheck(byCategoryId)

		// filter by system template
		bySystemTemplate   = "data.huaweicloud_dataarts_quality_rule_templates.filter_by_system_template"
		dcBySystemTemplate = acceptance.InitDataSourceCheck(bySystemTemplate)

		// filter by creator
		byCreator   = "data.huaweicloud_dataarts_quality_rule_templates.filter_by_creator"
		dcByCreator = acceptance.InitDataSourceCheck(byCreator)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckDataArtsWorkSpaceID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataQualityRuleTemplates_basic(),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "templates.#", regexp.MustCompile(`[1-9]([0-9]*)?`)),
					resource.TestCheckResourceAttrSet(all, "templates.0.id"),
					resource.TestCheckResourceAttrSet(all, "templates.0.name"),
					resource.TestCheckResourceAttrSet(all, "templates.0.category_id"),
					resource.TestCheckResourceAttrSet(all, "templates.0.dimension"),
					resource.TestCheckResourceAttrSet(all, "templates.0.type"),
					resource.TestCheckResourceAttrSet(all, "templates.0.system_template"),
					resource.TestCheckResourceAttrSet(all, "templates.0.sql_info"),
					resource.TestCheckResourceAttrSet(all, "templates.0.result_description"),
					resource.TestMatchResourceAttr(all, "templates.0.created_at",
						regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}?(Z|([+-]\d{2}:\d{2}))$`)),
					resource.TestCheckResourceAttrSet(all, "templates.0.creator"),

					// filter by name
					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					// filter by not found name
					dcByNotFoundName.CheckResourceExists(),
					resource.TestCheckOutput("is_not_found_name_filter_useful", "true"),

					// filter by category ID
					dcByCategoryId.CheckResourceExists(),
					resource.TestCheckOutput("is_category_id_filter_useful", "true"),

					// filter by system template
					dcBySystemTemplate.CheckResourceExists(),
					resource.TestCheckOutput("is_system_template_filter_useful", "true"),

					// filter by creator
					dcByCreator.CheckResourceExists(),
					resource.TestCheckOutput("is_creator_filter_useful", "true"),
				),
			},
		},
	})
}

func testAccDataQualityRuleTemplates_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_dataarts_quality_rule_templates" "all" {
  workspace_id = "%[1]s"
}

# Filter by name
locals {
  template_name = data.huaweicloud_dataarts_quality_rule_templates.all.templates[0].name
}

data "huaweicloud_dataarts_quality_rule_templates" "filter_by_name" {
  workspace_id = "%[1]s"

  name = local.template_name
}

locals {
  name_filter_result = [
    for v in data.huaweicloud_dataarts_quality_rule_templates.filter_by_name.templates[*].name : v == local.template_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.name_filter_result) > 0 && alltrue(local.name_filter_result)
}

# Filter by not found name
data "huaweicloud_dataarts_quality_rule_templates" "filter_by_not_found_name" {
  workspace_id = "%[1]s"

  name = "not_found_template"
}

output "is_not_found_name_filter_useful" {
  value = length(data.huaweicloud_dataarts_quality_rule_templates.filter_by_not_found_name.templates) == 0
}

# Filter by category ID
locals {
  category_id = data.huaweicloud_dataarts_quality_rule_templates.all.templates[0].category_id
}

data "huaweicloud_dataarts_quality_rule_templates" "filter_by_category_id" {
  workspace_id = "%[1]s"

  category_id = local.category_id
}

locals {
  category_id_filter_result = [
    for v in data.huaweicloud_dataarts_quality_rule_templates.filter_by_category_id.templates[*].category_id : v == local.category_id
  ]
}

output "is_category_id_filter_useful" {
  value = length(local.category_id_filter_result) > 0 && alltrue(local.category_id_filter_result)
}

# Filter by system template
locals {
  system_template = data.huaweicloud_dataarts_quality_rule_templates.all.templates[0].system_template
}

data "huaweicloud_dataarts_quality_rule_templates" "filter_by_system_template" {
  workspace_id = "%[1]s"

  system_template = local.system_template
}

locals {
  system_template_filter_result = [
    for v in data.huaweicloud_dataarts_quality_rule_templates.filter_by_system_template.templates[*].system_template : v == local.system_template
  ]
}

output "is_system_template_filter_useful" {
  value = length(local.system_template_filter_result) > 0 && alltrue(local.system_template_filter_result)
}

# Filter by creator
locals {
  template_creator = data.huaweicloud_dataarts_quality_rule_templates.all.templates[0].creator
}

data "huaweicloud_dataarts_quality_rule_templates" "filter_by_creator" {
  workspace_id = "%[1]s"

  creator = local.template_creator
}

locals {
  creator_filter_result = [
    for v in data.huaweicloud_dataarts_quality_rule_templates.filter_by_creator.templates[*].creator : v == local.template_creator
  ]
}

output "is_creator_filter_useful" {
  value = length(local.creator_filter_result) > 0 && alltrue(local.creator_filter_result)
}
`, acceptance.HW_DATAARTS_WORKSPACE_ID)
}
