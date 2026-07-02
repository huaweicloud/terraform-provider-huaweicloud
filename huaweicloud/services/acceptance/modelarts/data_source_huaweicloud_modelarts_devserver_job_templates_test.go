package modelarts

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDevServerJobTemplates_basic(t *testing.T) {
	var (
		all = "data.huaweicloud_modelarts_devserver_job_templates.all"
		dc  = acceptance.InitDataSourceCheck(all)

		byId   = "data.huaweicloud_modelarts_devserver_job_templates.filter_by_template_id"
		dcById = acceptance.InitDataSourceCheck(byId)

		byName   = "data.huaweicloud_modelarts_devserver_job_templates.filter_by_name"
		dcByName = acceptance.InitDataSourceCheck(byName)

		byType   = "data.huaweicloud_modelarts_devserver_job_templates.filter_by_type"
		dcByType = acceptance.InitDataSourceCheck(byType)
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDevServerJobTemplates_basic,
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestMatchResourceAttr(all, "templates.#", regexp.MustCompile(`^[1-9]([0-9]*)?$`)),
					resource.TestCheckResourceAttrSet(all, "templates.0.id"),
					resource.TestCheckResourceAttrSet(all, "templates.0.name"),
					resource.TestCheckResourceAttrSet(all, "templates.0.description"),
					resource.TestCheckResourceAttrSet(all, "templates.0.flavor_type"),
					resource.TestCheckResourceAttrSet(all, "templates.0.type"),

					dcById.CheckResourceExists(),
					resource.TestCheckOutput("is_id_filter_useful", "true"),

					dcByName.CheckResourceExists(),
					resource.TestCheckOutput("is_name_filter_useful", "true"),

					dcByType.CheckResourceExists(),
					resource.TestCheckOutput("is_type_filter_useful", "true"),
				),
			},
		},
	})
}

const testAccDevServerJobTemplates_basic string = `
# Query all DevServer job templates without any filter.
data "huaweicloud_modelarts_devserver_job_templates" "all" {}

# Filter by 'template_id'.
locals {
  template_id = try(data.huaweicloud_modelarts_devserver_job_templates.all.templates[0].id, "")
}

data "huaweicloud_modelarts_devserver_job_templates" "filter_by_template_id" {
  template_id = local.template_id
}

locals {
  filter_by_template_id_result = [
    for v in data.huaweicloud_modelarts_devserver_job_templates.filter_by_template_id.templates[*].id : v == local.template_id
  ]
}

output "is_id_filter_useful" {
  value = length(local.filter_by_template_id_result) > 0 && alltrue(local.filter_by_template_id_result)
}

# Filter by 'name'.
locals {
  template_name = try(data.huaweicloud_modelarts_devserver_job_templates.all.templates[0].name, "")
}

data "huaweicloud_modelarts_devserver_job_templates" "filter_by_name" {
  name = local.template_name
}

locals {
  filter_by_name_result = [
    for v in data.huaweicloud_modelarts_devserver_job_templates.filter_by_name.templates[*].name : v == local.template_name
  ]
}

output "is_name_filter_useful" {
  value = length(local.filter_by_name_result) > 0 && alltrue(local.filter_by_name_result)
}

# Filter by 'type'.
locals {
  template_type = try(data.huaweicloud_modelarts_devserver_job_templates.all.templates[0].type, "")
}

data "huaweicloud_modelarts_devserver_job_templates" "filter_by_type" {
  type = local.template_type
}

locals {
  filter_by_type_result = [
    for v in data.huaweicloud_modelarts_devserver_job_templates.filter_by_type.templates[*].type : v == local.template_type
  ]
}

output "is_type_filter_useful" {
  value = length(local.filter_by_type_result) > 0 && alltrue(local.filter_by_type_result)
}
`
